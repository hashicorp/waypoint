package core

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/internal-shared/protomappers"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Project represents a project with one or more applications.
//
// The Close function should be called when finished with the project
// to properly clean up any open resources.
type Project struct {
	logger    hclog.Logger
	apps      map[string]*App
	factories map[component.Type]*factory.Factory
	dir       *datadir.Project
	mappers   []*argmapper.Func
	client    pb.WaypointClient
	dconfig   component.DeploymentConfig

	// name is the name of the project
	name string

	// labels is the list of labels that are assigned to this project.
	labels map[string]string

	// workspace is the workspace that this project will work in.
	workspace string

	// This lock only needs to be held currently to protect localClosers.
	lock sync.Mutex

	// The below are resources we need to close when Close is called, if non-nil
	localClosers []io.Closer

	// UI is the terminal UI to use for messages related to the project
	// as a whole. These messages will show up unprefixed for example compared
	// to the app-specific UI.
	UI terminal.UI

	// overrideLabels are the labels specified via the CLI to override
	// all other conflicting keys.
	overrideLabels map[string]string
}

// NewProject creates a new Project with the given options.
func NewProject(ctx context.Context, os ...Option) (*Project, error) {
	// Defaults
	p := &Project{
		UI: &terminal.BasicUI{},

		logger:    hclog.L(),
		workspace: "default",
		apps:      make(map[string]*App),
		factories: map[component.Type]*factory.Factory{
			component.BuilderType:        plugin.Builders,
			component.RegistryType:       plugin.Registries,
			component.PlatformType:       plugin.Platforms,
			component.ReleaseManagerType: plugin.Releasers,
		},
	}

	// Set our options
	var opts options
	for _, o := range os {
		o(p, &opts)
	}

	// Defaults
	if len(p.mappers) == 0 {
		var err error
		p.mappers, err = argmapper.NewFuncList(protomappers.All,
			argmapper.Logger(p.logger),
		)
		if err != nil {
			return nil, err
		}
	}

	// Validation
	if p.dir == nil {
		return nil, fmt.Errorf("WithDataDir must be specified")
	}
	if err := opts.Config.Validate(); err != nil {
		return nil, err
	}
	if errs := config.ValidateLabels(p.overrideLabels); len(errs) > 0 {
		return nil, multierror.Append(nil, errs...)
	}

	// Init our server connection. This may be in-process if we're in
	// local mode.
	if p.client == nil {
		panic("p.client should never be nil")
	}

	// Set our labels
	p.labels = opts.Config.Labels

	if opts.Config.URL != nil {
		p.dconfig.UrlControlAddr = opts.Config.URL.ControlAddr
		p.dconfig.UrlToken = opts.Config.URL.Token
	}

	// Initialize all the applications and load all their components.
	for _, appConfig := range opts.Config.Apps {
		app, err := newApp(ctx, p, appConfig)
		if err != nil {
			return nil, err
		}

		p.apps[appConfig.Name] = app
	}

	p.logger.Info("project initialized", "workspace", p.workspace)
	return p, nil
}

// App initializes and returns the app with the given name.
func (p *Project) App(name string) (*App, error) {
	return p.apps[name], nil
}

// Client returns the API client for the backend server.
func (p *Project) Client() pb.WaypointClient {
	return p.client
}

// Ref returns the project ref for API calls.
func (p *Project) Ref() *pb.Ref_Project {
	return &pb.Ref_Project{Project: p.name}
}

// WorkspaceRef returns the project ref for API calls.
func (p *Project) WorkspaceRef() *pb.Ref_Workspace {
	return &pb.Ref_Workspace{
		Workspace: p.workspace,
	}
}

// Close is called to clean up resources allocated by the project.
// This should be called and blocked on to gracefully stop the project.
func (p *Project) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.logger.Debug("closing project")

	// Stop all our apps
	for name, app := range p.apps {
		p.logger.Trace("closing app", "app", name)
		if err := app.Close(); err != nil {
			p.logger.Warn("error closing app", "err", err)
		}
	}

	// If we're running in local mode, close our local resources we started
	for _, c := range p.localClosers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	p.localClosers = nil

	return nil
}

// mergeLabels merges the set of labels given. This will set the project
// labels as a base automatically and then merge ls in order.
func (p *Project) mergeLabels(ls ...map[string]string) map[string]string {
	result := map[string]string{}

	// Set our builtin labels
	result["waypoint/workspace"] = p.workspace

	// Set our project labels
	for k, v := range p.labels {
		result[k] = v
	}

	// Set any labels given
	for _, lm := range ls {
		for k, v := range lm {
			result[k] = v
		}
	}

	// Set any overrides
	for k, v := range p.overrideLabels {
		result[k] = v
	}

	return result
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

// WithClient sets the API client to use.
func WithClient(client pb.WaypointClient) Option {
	return func(p *Project, opts *options) {
		p.client = client
	}
}

// WithConfig uses the given project configuration for initializing the
// Project. This configuration must be validated already prior to using this
// option.
func WithConfig(c *config.Config) Option {
	return func(p *Project, opts *options) {
		opts.Config = c
		p.name = c.Project
	}
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
func WithFactory(t component.Type, f *factory.Factory) Option {
	return func(p *Project, opts *options) { p.factories[t] = f }
}

// WithComponents sets the factories for components.
func WithComponents(fs map[component.Type]*factory.Factory) Option {
	return func(p *Project, opts *options) { p.factories = fs }
}

// WithMappers adds the mappers to the list of mappers.
func WithMappers(m ...*argmapper.Func) Option {
	return func(p *Project, opts *options) { p.mappers = append(p.mappers, m...) }
}

// WithLabels sets the labels that will override any other labels set.
func WithLabels(m map[string]string) Option {
	return func(p *Project, opts *options) { p.overrideLabels = m }
}

// WithWorkspace sets the workspace we'll be working in.
func WithWorkspace(ws string) Option {
	return func(p *Project, opts *options) {
		if ws != "" {
			p.workspace = ws
		}
	}
}

// WithUI sets the UI to use. If this isn't set, a BasicUI is used.
func WithUI(ui terminal.UI) Option {
	return func(p *Project, opts *options) { p.UI = ui }
}
