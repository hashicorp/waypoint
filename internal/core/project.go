package core

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/builtin"
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/mapper"
)

// Project represents a project with one or more applications.
type Project struct {
	logger    hclog.Logger
	apps      map[string]*App
	factories map[component.Type]*mapper.Factory
}

// NewProject creates a new Project with the given options.
func NewProject(ctx context.Context, os ...Option) (*Project, error) {
	// Defaults
	p := &Project{
		logger: hclog.L(),
		apps:   make(map[string]*App),
		factories: map[component.Type]*mapper.Factory{
			component.BuilderType:  builtin.Builders,
			component.RegistryType: builtin.Registries,
			component.PlatformType: builtin.Platforms,
		},
	}

	// Set our options
	var opts options
	for _, o := range os {
		o(p, &opts)
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
