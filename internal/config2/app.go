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

	Build   *Build   `hcl:"build,block"`
	Deploy  *Deploy  `hcl:"deploy,block"`
	Release *Release `hcl:"release,block"`

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

// App returns the configured app named n. If the app doesn't exist, this
// will return (nil, nil).
func (c *Config) App(n string, ctx *hcl.EvalContext) (*App, error) {
	ctx = c.ctx

	// Find the app by progressively decoding
	var rawApp *hclApp
	for _, app := range c.raw.Apps {
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
