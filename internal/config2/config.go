package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	*hclConfig

	ctx      *hcl.EvalContext
	pathData map[string]string
}

type hclConfig struct {
	Project string            `hcl:"project,attr"`
	Runner  *Runner           `hcl:"runner,block" default:"{}"`
	Labels  map[string]string `hcl:"labels,optional"`
	Plugin  []*Plugin         `hcl:"plugin,block"`
	Apps    []*hclApp         `hcl:"app,block"`
}

// Runner is the configuration for supporting runners in this project.
type Runner struct {
	// Enabled is whether or not runners are enabled. If this is false
	// then the "-remote" flag will not work.
	Enabled bool `hcl:"enabled,optional"`

	// DataSource is the default data source when a remote job is queued.
	DataSource *DataSource `hcl:"data_source,block"`
}

// DataSource configures the data source for the runner.
type DataSource struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

// Load loads the configuration file from the given path.
//
// Configuration loading in Waypoint is lazy. This will load just the amount
// necessary to return the initial Config structure. Additional functions on
// Config must be called to load additional parts of the Config.
//
// This also means that the config may be invalid. To validate the config
// call the Validate method.
func Load(path string, pwd string) (*Config, error) {
	// We require an absolute path for the path so we can set the path vars
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	// If we have no pwd, then create a temporary directory
	if pwd == "" {
		td, err := ioutil.TempDir("", "waypoint-config")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(td)
		pwd = td
	}

	// Setup our initial variable set
	pathData := map[string]string{
		"pwd":     pwd,
		"project": filepath.Dir(path),
	}

	// Build our context
	ctx := EvalContext(nil, pwd).NewChild()
	addPathValue(ctx, pathData)

	// Decode
	var cfg hclConfig
	if err := hclsimple.DecodeFile(path, ctx, &cfg); err != nil {
		return nil, err
	}

	return &Config{
		hclConfig: &cfg,
		ctx:       ctx,
		pathData:  pathData,
	}, nil
}

// HCLContext returns the eval context for this configuration.
func (c *Config) HCLContext() *hcl.EvalContext {
	return c.ctx.NewChild()
}
