package config

import (
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/copystructure"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/waypoint/pkg/config/funcs"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// App represents a single application.
type App struct {
	Name   string            `hcl:",label"`
	Path   string            `hcl:"path,optional"`
	Labels map[string]string `hcl:"labels,optional"`
	URL    *AppURL           `hcl:"url,block" default:"{}"`
	Config *genericConfig    `hcl:"config,block"`

	Runner *Runner `hcl:"runner,block"`

	BuildRaw   *hclBuild `hcl:"build,block"`
	DeployRaw  *hclStage `hcl:"deploy,block"`
	ReleaseRaw *hclStage `hcl:"release,block"`

	Body hcl.Body `hcl:",body"`

	ctx    *hcl.EvalContext
	config *Config
}

// AppURL configures the App-specific URL settings.
type AppURL struct {
	AutoHostname *bool `hcl:"auto_hostname,optional"`
}

type hclApp struct {
	Name string `hcl:",label"`
	Path string `hcl:"path,optional"`

	// We need these raw values to determine the plugins need to be used.
	BuildRaw   *hclBuild `hcl:"build,block"`
	DeployRaw  *hclStage `hcl:"deploy,block"`
	ReleaseRaw *hclStage `hcl:"release,block"`

	Runner *Runner `hcl:"runner,block"`

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}

// hclLabeled is used to partially decode only the labels from a
// structure that supports it.
type hclLabeled struct {
	Labels map[string]string `hcl:"labels,optional"`
	Remain hcl.Body          `hcl:",remain"`
}

// Apps returns the names of all the apps.
func (c *Config) Apps() []string {
	var result []string
	for _, app := range c.hclConfig.Apps {
		result = append(result, app.Name)
	}

	return result
}

// App returns the configured app named n. If the app doesn't exist, this
// will return (nil, nil).
func (c *Config) App(n string, ctx *hcl.EvalContext) (*App, error) {
	ctx = appendContext(c.ctx, ctx)

	// Find the app by progressively decoding
	var rawApp *hclApp
	for _, app := range c.hclConfig.Apps {
		if app.Name == n {
			rawApp = app
			break
		}
	}
	if rawApp == nil {
		return nil, nil
	}

	// Determine the app path
	appPath := rawApp.Path
	if !filepath.IsAbs(appPath) {
		appPath = filepath.Join(c.pathData["project"], appPath)
	}

	// Update our path data to contain the app path
	pathData := copystructure.Must(copystructure.Copy(c.pathData)).(map[string]string)
	pathData["app"] = appPath

	// Build a new context with our app-scoped values
	ctx = ctx.NewChild()
	addPathValue(ctx, pathData)
	addMapVariable(ctx, "app", map[string]string{
		"name": rawApp.Name,
	})

	// Full decode
	var app App
	if diag := gohcl.DecodeBody(rawApp.Body, finalizeContext(ctx), &app); diag.HasErrors() {
		return nil, diag
	}
	app.Name = rawApp.Name
	app.Path = appPath
	app.Runner = rawApp.Runner
	app.ctx = ctx
	app.config = c
	if app.Config != nil {
		app.Config.ctx = ctx
		app.Config.scopeFunc = func(cv *pb.ConfigVar) {
			cv.Target.AppScope = &pb.ConfigVar_Target_Application{
				Application: app.Ref(),
			}
		}
	}

	return &app, nil
}

// Ref returns the ref for this app.
func (c *App) Ref() *pb.Ref_Application {
	return &pb.Ref_Application{
		Application: c.Name,
		Project:     c.config.Project,
	}
}

// ConfigVars returns the configuration variables for the app, including
// merging the configuration variables from the project level.
//
// For access to only the app-level config vars, use the Config attribute directly.
func (c *App) ConfigVars() ([]*pb.ConfigVar, error) {
	vars, err := c.config.Config.ConfigVars()
	if err != nil {
		return nil, err
	}

	appVars, err := c.Config.ConfigVars()
	if err != nil {
		return nil, err
	}

	return append(vars, appVars...), nil
}

// ConfigMetadata holds information about the configuration variables or process
// themselves.
type ConfigMetadata struct {
	FileChangeSignal string
}

// ConfigMetadata returns any configuration metadata about the project and app.
func (c *App) ConfigMetadata() (*ConfigMetadata, *ConfigMetadata) {
	var app, proj *ConfigMetadata

	if c.config.Config != nil {
		proj = &ConfigMetadata{
			FileChangeSignal: c.config.Config.FileChangeSignal,
		}
	}

	if c.Config != nil {
		app = &ConfigMetadata{
			FileChangeSignal: c.Config.FileChangeSignal,
		}
	}

	return proj, app
}

// Build loads the Build section of the configuration.
func (c *App) Build(ctx *hcl.EvalContext) (*Build, error) {
	ctx = appendContext(c.ctx, ctx)

	body := c.BuildRaw.Body
	scope, err := scopeMatchStage(ctx, c.BuildRaw.WorkspaceScoped, c.BuildRaw.LabelScoped)
	if err != nil {
		return nil, err
	}
	if scope != nil {
		body = scope.Body
	}

	var b Build
	if diag := gohcl.DecodeBody(body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// Registry loads the Registry section of the configuration.
func (c *App) Registry(ctx *hcl.EvalContext) (*Registry, error) {
	// Registry is optional
	if c.BuildRaw == nil || c.BuildRaw.Registry == nil {
		return nil, nil
	}

	body := c.BuildRaw.Registry.Body
	scope, err := scopeMatchStage(ctx,
		c.BuildRaw.Registry.WorkspaceScoped,
		c.BuildRaw.Registry.LabelScoped)
	if err != nil {
		return nil, err
	}
	if scope != nil {
		body = scope.Body
	}

	var b Registry
	ctx = appendContext(c.ctx, ctx)
	if diag := gohcl.DecodeBody(body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// Deploy loads the associated section of the configuration.
func (c *App) Deploy(ctx *hcl.EvalContext) (*Deploy, error) {
	ctx = appendContext(c.ctx, ctx)

	body := c.DeployRaw.Body
	scope, err := scopeMatchStage(ctx, c.DeployRaw.WorkspaceScoped, c.DeployRaw.LabelScoped)
	if err != nil {
		return nil, err
	}
	if scope != nil {
		body = scope.Body
	}

	var b Deploy
	if diag := gohcl.DecodeBody(body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// Release loads the associated section of the configuration.
func (c *App) Release(ctx *hcl.EvalContext) (*Release, error) {
	if c.ReleaseRaw == nil {
		return nil, nil
	}

	body := c.ReleaseRaw.Body
	scope, err := scopeMatchStage(ctx, c.ReleaseRaw.WorkspaceScoped, c.ReleaseRaw.LabelScoped)
	if err != nil {
		return nil, err
	}
	if scope != nil {
		body = scope.Body
	}

	var b Release
	ctx = appendContext(c.ctx, ctx)
	if diag := gohcl.DecodeBody(body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// BuildUse returns the plugin "use" value.
func (c *App) BuildUse(ctx *hcl.EvalContext) (string, error) {
	if c.BuildRaw == nil {
		return "", nil
	} else if c.BuildRaw.Use == nil {
		return "", nil
	}

	useType := c.BuildRaw.Use.Type
	stage, err := scopeMatchStage(ctx, c.BuildRaw.WorkspaceScoped, c.BuildRaw.LabelScoped)
	if err != nil {
		return "", err
	}
	if stage != nil {
		useType = stage.Use.Type
	}

	return useType, nil
}

// RegistryUse returns the plugin "use" value.
func (c *App) RegistryUse(ctx *hcl.EvalContext) (string, error) {
	if c.BuildRaw == nil || c.BuildRaw.Registry == nil {
		return "", nil
	} else if c.BuildRaw.Registry.Use == nil {
		return "", nil
	}

	useType := c.BuildRaw.Registry.Use.Type
	stage, err := scopeMatchStage(ctx, c.BuildRaw.Registry.WorkspaceScoped, c.BuildRaw.Registry.LabelScoped)
	if err != nil {
		return "", err
	}
	if stage != nil {
		useType = stage.Use.Type
	}

	return useType, nil
}

// DeployUse returns the plugin "use" value.
func (c *App) DeployUse(ctx *hcl.EvalContext) (string, error) {
	if c.DeployRaw == nil {
		return "", nil
	} else if c.DeployRaw.Use == nil {
		return "", nil
	}

	useType := c.DeployRaw.Use.Type
	stage, err := scopeMatchStage(ctx, c.DeployRaw.WorkspaceScoped, c.DeployRaw.LabelScoped)
	if err != nil {
		return "", err
	}
	if stage != nil {
		useType = stage.Use.Type
	}

	return useType, nil
}

// ReleaseUse returns the plugin "use" value.
func (c *App) ReleaseUse(ctx *hcl.EvalContext) (string, error) {
	if c.ReleaseRaw == nil {
		return "", nil
	} else if c.ReleaseRaw.Use == nil {
		return "", nil
	}

	useType := c.ReleaseRaw.Use.Type
	stage, err := scopeMatchStage(ctx, c.ReleaseRaw.WorkspaceScoped, c.ReleaseRaw.LabelScoped)
	if err != nil {
		return "", err
	}
	if stage != nil {
		useType = stage.Use.Type
	}

	return useType, nil
}

// BuildLabels returns the labels for this stage.
func (c *App) BuildLabels(ctx *hcl.EvalContext) (map[string]string, error) {
	if c.BuildRaw == nil {
		return nil, nil
	}

	ctx = appendContext(c.ctx, ctx)
	return labels(ctx, c.BuildRaw.Body)
}

// RegistryLabels returns the labels for this stage.
func (c *App) RegistryLabels(ctx *hcl.EvalContext) (map[string]string, error) {
	if c.BuildRaw == nil || c.BuildRaw.Registry == nil {
		return nil, nil
	}

	ctx = appendContext(c.ctx, ctx)

	// Get both build and registry labels
	allLabels, err := labels(ctx, c.BuildRaw.Body)
	if err != nil {
		return nil, err
	}
	registryLabels, err := labels(ctx, c.BuildRaw.Registry.Body)
	if err != nil {
		return nil, err
	}

	// Merge em
	for k, v := range registryLabels {
		allLabels[k] = v
	}

	return allLabels, nil
}

// DeployLabels returns the labels for this stage.
func (c *App) DeployLabels(ctx *hcl.EvalContext) (map[string]string, error) {
	if c.DeployRaw == nil {
		return nil, nil
	}

	ctx = appendContext(c.ctx, ctx)
	return labels(ctx, c.DeployRaw.Body)
}

// ReleaseLabels returns the labels for this stage.
func (c *App) ReleaseLabels(ctx *hcl.EvalContext) (map[string]string, error) {
	if c.ReleaseRaw == nil {
		return nil, nil
	}

	ctx = appendContext(c.ctx, ctx)
	return labels(ctx, c.ReleaseRaw.Body)
}

// labels reads the labels from the given body (if they are available),
// merges them using lm, and returns the final merged set of labels. This
// also returns a new EvalContext that has the `labels` HCL variable set.
func labels(ctx *hcl.EvalContext, body hcl.Body) (map[string]string, error) {
	// First decode into our structure that only reads labels.
	var labeled hclLabeled
	if diag := gohcl.DecodeBody(body, finalizeContext(ctx), &labeled); diag.HasErrors() {
		return nil, diag
	}

	// Merge em
	return labeled.Labels, nil
}

// This returns a matching stage (if any) for the given context. The context
// is expected to have "labels" set in it.
//
// If ws is true, then the scope of the scopedStage will be compared to
// a label of value "waypoint/workspace" which is expected to always be
// present.
//
// This function can return (nil, nil) as a valid result. This means
// that no scopes matched (which is expected behavior to fallback to a
// default).
func scopeMatchStage(
	ctx *hcl.EvalContext, wsScopes []*scopedStage, labelScopes []*scopedStage,
) (*scopedStage, error) {
	// These are all scenarios where we can't possibly match any scope.
	if ctx == nil || ctx.Variables == nil {
		return nil, nil
	}

	// Get our labels. If we have none, we can never match.
	labels, ok := ctx.Variables["labels"]
	if !ok || labels.LengthInt() == 0 {
		return nil, nil
	}

	// If we're workspace matching, simplify this by looking up the
	// "waypoint/workspace" key.
	if len(wsScopes) > 0 {
		values := labels.AsValueMap()
		wsValue, ok := values["waypoint/workspace"]
		if !ok {
			// No workspace key we can't possibly match.
			return nil, nil
		}

		// Look for an exact match
		for _, s := range wsScopes {
			if s.Scope == wsValue.AsString() {
				return s, nil
			}
		}
	}

	// For label selectors, we want to sort by scope length so that
	// the longest label selectors match first. This is our rule for
	// tiebreaking.
	sort.Slice(labelScopes, func(i, j int) bool {
		x, y := labelScopes[i], labelScopes[j]
		return len(x.Scope) > len(y.Scope)
	})

	// Label selector matching.
	for _, s := range labelScopes {
		result, err := funcs.SelectorMatch(labels, cty.StringVal(s.Scope))
		if err != nil {
			return nil, err
		}

		if result.True() {
			return s, nil
		}
	}

	return nil, nil
}
