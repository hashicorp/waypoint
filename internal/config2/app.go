package config

import (
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/copystructure"
)

// App represents a single application.
type App struct {
	Name   string            `hcl:",label"`
	Path   string            `hcl:"path,optional"`
	Labels map[string]string `hcl:"labels,optional"`
	URL    *AppURL           `hcl:"url,block" default:"{}"`

	BuildRaw   *hclBuild `hcl:"build,block"`
	DeployRaw  *hclStage `hcl:"deploy,block"`
	ReleaseRaw *hclStage `hcl:"release,block"`

	ctx    *hcl.EvalContext
	config *Config
}

// AppURL configures the App-specific URL settings.
type AppURL struct {
	AutoHostname *bool `hcl:"auto_hostname,optional"`
}

type hclApp struct {
	Name string   `hcl:",label"`
	Path string   `hcl:"path,optional"`
	Body hcl.Body `hcl:",remain"`
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

	// Full decode
	app := App{
		Name: rawApp.Name,
		Path: appPath,
	}
	if diag := gohcl.DecodeBody(rawApp.Body, ctx, &app); diag.HasErrors() {
		return nil, diag
	}
	app.ctx = ctx
	app.config = c

	return &app, nil
}

// Build loads the Build section of the configuration.
func (c *App) Build(ctx *hcl.EvalContext) (*Build, error) {
	ctx = appendContext(c.ctx, ctx)

	var b Build
	if diag := gohcl.DecodeBody(c.BuildRaw.Body, ctx, &b); diag.HasErrors() {
		return nil, diag
	}

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
	if diag := gohcl.DecodeBody(c.BuildRaw.Registry.Body, ctx, &b); diag.HasErrors() {
		return nil, diag
	}

	return &b, nil
}

// Deploy loads the associated section of the configuration.
func (c *App) Deploy(ctx *hcl.EvalContext) (*Deploy, error) {
	ctx = appendContext(c.ctx, ctx)

	var b Deploy
	if diag := gohcl.DecodeBody(c.DeployRaw.Body, ctx, &b); diag.HasErrors() {
		return nil, diag
	}

	return &b, nil
}

// Release loads the associated section of the configuration.
func (c *App) Release(ctx *hcl.EvalContext) (*Release, error) {
	if c.ReleaseRaw == nil {
		return nil, nil
	}

	var b Release
	ctx = appendContext(c.ctx, ctx)
	if diag := gohcl.DecodeBody(c.ReleaseRaw.Body, ctx, &b); diag.HasErrors() {
		return nil, diag
	}

	return &b, nil
}
