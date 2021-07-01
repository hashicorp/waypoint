package config

import (
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/copystructure"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// App represents a single application.
type App struct {
	Name   string            `hcl:",label"`
	Path   string            `hcl:"path,optional"`
	Labels map[string]string `hcl:"labels,optional"`
	URL    *AppURL           `hcl:"url,block" default:"{}"`
	Config *genericConfig    `hcl:"config,block"`

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

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
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
	app.ctx = ctx
	app.config = c
	if app.Config != nil {
		app.Config.ctx = ctx
		app.Config.scopeFunc = func(cv *pb.ConfigVar) {
			cv.Scope = &pb.ConfigVar_Application{
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

	var b Build
	if diag := gohcl.DecodeBody(c.BuildRaw.Body, finalizeContext(ctx), &b); diag.HasErrors() {
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

	var b Registry
	ctx = appendContext(c.ctx, ctx)
	if diag := gohcl.DecodeBody(c.BuildRaw.Registry.Body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// Deploy loads the associated section of the configuration.
func (c *App) Deploy(ctx *hcl.EvalContext) (*Deploy, error) {
	ctx = appendContext(c.ctx, ctx)

	var b Deploy
	if diag := gohcl.DecodeBody(c.DeployRaw.Body, finalizeContext(ctx), &b); diag.HasErrors() {
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

	var b Release
	ctx = appendContext(c.ctx, ctx)
	if diag := gohcl.DecodeBody(c.ReleaseRaw.Body, finalizeContext(ctx), &b); diag.HasErrors() {
		return nil, diag
	}
	b.ctx = ctx

	return &b, nil
}

// BuildUse returns the plugin "use" value.
func (c *App) BuildUse() string {
	if c.BuildRaw == nil {
		return ""
	}

	return c.BuildRaw.Use.Type
}

// RegistryUse returns the plugin "use" value.
func (c *App) RegistryUse() string {
	if c.BuildRaw == nil || c.BuildRaw.Registry == nil {
		return ""
	}

	return c.BuildRaw.Registry.Use.Type
}

// DeployUse returns the plugin "use" value.
func (c *App) DeployUse() string {
	if c.DeployRaw == nil {
		return ""
	}

	return c.DeployRaw.Use.Type
}

// ReleaseUse returns the plugin "use" value.
func (c *App) ReleaseUse() string {
	if c.ReleaseRaw == nil {
		return ""
	}

	return c.ReleaseRaw.Use.Type
}
