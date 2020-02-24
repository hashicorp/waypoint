package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/plugin"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	"github.com/mitchellh/devflow/sdk/protomappers"
)

// Project represents a project with one or more applications.
type Project struct {
	logger    hclog.Logger
	apps      map[string]*App
	factories map[component.Type]*mapper.Factory
	dir       *datadir.Project
	mappers   []*mapper.Func
}

// NewProject creates a new Project with the given options.
func NewProject(ctx context.Context, os ...Option) (*Project, error) {
	// Defaults
	p := &Project{
		logger: hclog.L(),
		apps:   make(map[string]*App),
		factories: map[component.Type]*mapper.Factory{
			component.BuilderType:  plugin.Builders,
			component.RegistryType: plugin.Registries,
			component.PlatformType: plugin.Platforms,
		},
	}

	// Set our options
	var opts options
	for _, o := range os {
		o(p, &opts)
	}

	// Defaults
	if len(p.mappers) == 0 {
		p.mappers = protomappers.AllFuncs
	}

	// Validation
	if p.dir == nil {
		return nil, fmt.Errorf("WithDataDir must be specified")
	}

	// Initialize all the applications and load all their components.
	for _, appConfig := range opts.Config.Apps {
		app, err := newApp(ctx, p, appConfig)
		if err != nil {
			return nil, err
		}

		p.apps[appConfig.Name] = app
	}

	return p, nil
}

// App initializes and returns the app with the given name.
func (p *Project) App(name string) (*App, error) {
	return p.apps[name], nil
}

// options is the configuration to construct a new Project. Some
// configuration is set directly on the Project. This is only used for
// intermediate values that need to be processed further before initializing
// the project.
type options struct {
	Config *config.Config
}

// Option is used to set options for NewProject.
type Option func(*Project, *options)

// WithConfig uses the given project configuration for initializing the
// Project. This configuration must be validated already prior to using this
// option.
func WithConfig(c *config.Config) Option {
	return func(p *Project, opts *options) { opts.Config = c }
}

// WithDataDir sets the datadir that will be used for this project.
func WithDataDir(dir *datadir.Project) Option {
	return func(p *Project, opts *options) { p.dir = dir }
}

// WithLogger sets the logger to use with the project. If this option
// is not provided, a default logger will be used (`hclog.L()`).
func WithLogger(log hclog.Logger) Option {
	return func(p *Project, opts *options) { p.logger = log }
}

// WithFactory sets a factory for a component type. If this isn't set for
// any component type, then the builtin mapper will be used.
func WithFactory(t component.Type, f *mapper.Factory) Option {
	return func(p *Project, opts *options) { p.factories[t] = f }
}

// WithMappers adds the mappers to the list of mappers.
func WithMappers(m ...*mapper.Func) Option {
	return func(p *Project, opts *options) { p.mappers = append(p.mappers, m...) }
}
