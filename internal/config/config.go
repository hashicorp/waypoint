// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/hashicorp/waypoint/internal/config/variables"
	"github.com/hashicorp/waypoint/internal/pkg/defaults"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type Config struct {
	hclConfig

	ctx      *hcl.EvalContext
	path     string
	pathData map[string]string

	InputVariables map[string]*variables.Variable
}

type hclConfig struct {
	Project   string                   `hcl:"project,optional"`
	Runner    *Runner                  `hcl:"runner,block" default:"{}"`
	Labels    map[string]string        `hcl:"labels,optional"`
	Variables []*variables.HclVariable `hcl:"variable,block"`
	Plugin    []*Plugin                `hcl:"plugin,block"`
	Config    *genericConfig           `hcl:"config,block"`
	Apps      []*hclApp                `hcl:"app,block"`
	Pipelines []*hclPipeline           `hcl:"pipeline,block"`
	Body      hcl.Body                 `hcl:",body"`
}

// Runner is the configuration for supporting runners in this project.
type Runner struct {
	// Profile is the name of the on-demand runner configuration.
	Profile string `hcl:"profile,optional"`

	// Note (XX): The other properties in this struct are only used on init,
	// and don't really make sense being here.
	// DataSource, for example, should be set at the project level, not with the runner...

	// Enabled is whether runners are enabled. If this is false
	// then the "-remote" flag will not work.
	Enabled bool `hcl:"enabled,optional"`

	// DataSource is the default data source when a remote job is queued.
	DataSource *DataSource `hcl:"data_source,block"`

	// Poll are the settings related to polling.
	Poll *Poll `hcl:"poll,block"`
}

// DataSource configures the data source for the runner.
type DataSource struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

// Poll configures the polling settings for a project.
type Poll struct {
	Enabled  bool   `hcl:"enabled,optional"`
	Interval string `hcl:"interval,optional"`
}

// LoadOptions should be set for the Load function.
type LoadOptions struct {
	// Pwd is the current working directory. This is used to setup the
	// `path.pwd` variable and also makes things such as the Git rules work.
	Pwd string

	// Workspace is the workspace that we are executing in. This is used to
	// setup `workspace.name` variables.
	Workspace string
}

// Load loads the configuration file from the given path.
//
// Configuration loading in Waypoint is lazy. This will load just the amount
// necessary to return the initial Config structure. Additional functions on
// Config must be called to load additional parts of the Config.
//
// This also means that the config may be invalid. To validate the config
// call the Validate method.
func Load(path string, opts *LoadOptions) (*Config, error) {
	if opts == nil {
		opts = &LoadOptions{}
	}

	// Unpack these cause they're used a lot
	pwd := opts.Pwd

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
	addWorkspaceValue(ctx, opts.Workspace)

	// Decode
	var cfg hclConfig
	if err := hclsimple.DecodeFile(path, finalizeContext(ctx), &cfg); err != nil {
		return nil, err
	}

	// Decode variable blocks
	schema, _ := gohcl.ImpliedBodySchema(&hclConfig{})
	content, diags := cfg.Body.Content(schema)
	if diags.HasErrors() {
		return nil, diags
	}
	vs, diags := variables.DecodeVariableBlocks(ctx, content)
	if diags.HasErrors() {
		return nil, diags
	}

	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}

	// Set some values
	if cfg.Config != nil {
		cfg.Config.ctx = ctx
		cfg.Config.scopeFunc = func(cv *pb.ConfigVar) {
			cv.Target.AppScope = &pb.ConfigVar_Target_Project{
				Project: &pb.Ref_Project{Project: cfg.Project},
			}
		}
	}

	return &Config{
		hclConfig:      cfg,
		ctx:            ctx,
		path:           filepath.Dir(path),
		pathData:       pathData,
		InputVariables: vs,
	}, nil
}

// HCLContext returns the eval context for this configuration.
func (c *Config) HCLContext() *hcl.EvalContext {
	return c.ctx.NewChild()
}

// ConfigPath returns the path to the directory that contains the config file (waypoint.hcl)
func (c *Config) ConfigPath() string {
	return c.path
}

// ConfigRunner returns the runner stanza on a project
func (c *Config) ConfigRunner() *Runner {
	return c.Runner
}

// ConfigAppRunner returns the runner stanza on an application
func (c *Config) ConfigAppRunner(name string) *Runner {
	for _, app := range c.hclConfig.Apps {
		if app.Name == name {
			return app.Runner
		}
	}
	return nil
}
