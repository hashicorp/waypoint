package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/datadir"
	"github.com/mitchellh/devflow/internal/mapper"
	"github.com/mitchellh/devflow/internal/pkg/status"
)

// App represents a single application and exposes all the operations
// that can be performed on an application.
//
// An App is only valid if it was returned by Project.App. The behavior of
// App if constructed in any other way is undefined and likely to result
// in crashes.
type App struct {
	Builder  component.Builder
	Registry component.Registry
	Platform component.Platform

	source        *component.Source
	logger        hclog.Logger
	dir           *datadir.App
	mappers       []*mapper.Func
	componentDirs map[interface{}]*datadir.Component
}

// newApp creates an App for the given project and configuration. This will
// initialize and configure all the components of this application. An error
// will be returned if this app fails to initialize: configuration is invalid,
// a component could not be found, etc.
func newApp(ctx context.Context, p *Project, cfg *config.App) (*App, error) {
	// Initialize
	app := &App{
		source:        &component.Source{App: cfg.Name, Path: "."},
		logger:        p.logger.Named("app").Named(cfg.Name),
		mappers:       p.mappers,
		componentDirs: make(map[interface{}]*datadir.Component),
	}

	// Setup our directory
	dir, err := p.dir.App(cfg.Name)
	if err != nil {
		return nil, err
	}
	app.dir = dir

	// Load all the components
	components := []struct {
		Target interface{}
		Type   component.Type
		Config *config.Component
	}{
		{&app.Builder, component.BuilderType, cfg.Build},
		{&app.Registry, component.RegistryType, cfg.Registry},
		{&app.Platform, component.PlatformType, cfg.Platform},
	}
	for _, c := range components {
		if c.Config == nil {
			// This component is not set, ignore.
			continue
		}

		err = app.initComponent(ctx, c.Type, c.Target, p.factories[c.Type], c.Config)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Build builds the artifact from source for this app.
// TODO(mitchellh): test
func (a *App) Build(ctx context.Context) (component.Artifact, error) {
	log := a.logger.Named("build")
	result, err := a.callDynamicFunc(ctx, log, a.Builder, a.Builder.BuildFunc())
	if err != nil {
		return nil, err
	}

	return result.(component.Artifact), nil
}

// Push pushes the given artifact to the registry.
// TODO(mitchellh): test
func (a *App) Push(ctx context.Context, artifact component.Artifact) (component.Artifact, error) {
	log := a.logger.Named("push")
	result, err := a.callDynamicFunc(ctx, log, a.Registry, a.Registry.PushFunc(), artifact)
	if err != nil {
		return nil, err
	}

	return result.(component.Artifact), nil
}

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, artifact component.Artifact) (component.Deployment, error) {
	log := a.logger.Named("platform")
	result, err := a.callDynamicFunc(ctx, log, a.Platform, a.Platform.DeployFunc(), artifact)
	if err != nil {
		return nil, err
	}

	return result.(component.Deployment), nil
}

// Exec using the deployer phase
// TODO(evanphx): test
func (a *App) Exec(ctx context.Context, updater status.Updater) error {
	log := a.logger.Named("platform")

	ep, ok := a.Platform.(component.ExecPlatform)
	if !ok {
		return fmt.Errorf("This platform does not support exec yet")
	}

	_, err := a.callDynamicFunc(ctx, log, a.Platform, ep.ExecFunc(), updater)
	if err != nil {
		return err
	}

	return nil
}

// Set config variables on the deployer phase
// TODO(evanphx): test
func (a *App) ConfigSet(ctx context.Context, key, val string) error {
	log := a.logger.Named("platform")

	ep, ok := a.Platform.(component.ConfigPlatform)
	if !ok {
		return fmt.Errorf("This platform does not support config yet")
	}

	cv := &component.ConfigVar{Name: key, Value: val}

	_, err := a.callDynamicFunc(ctx, log, a.Platform, ep.ConfigSetFunc(), cv)
	if err != nil {
		return err
	}

	return nil
}

// Get config variables on the deployer phase
// TODO(evanphx): test
func (a *App) ConfigGet(ctx context.Context, key string) (*component.ConfigVar, error) {
	log := a.logger.Named("platform")

	ep, ok := a.Platform.(component.ConfigPlatform)
	if !ok {
		return nil, fmt.Errorf("This platform does not support config yet")
	}

	cv := &component.ConfigVar{
		Name: key,
	}

	_, err := a.callDynamicFunc(ctx, log, a.Platform, ep.ConfigGetFunc(), cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

// Retrieve log viewer on the deployer phase
// TODO(evanphx): test
func (a *App) Logs(ctx context.Context) (component.LogViewer, error) {
	log := a.logger.Named("platform")

	ep, ok := a.Platform.(component.LogsPlatform)
	if !ok {
		return nil, fmt.Errorf("This platform does not support logs yet")
	}

	lv, err := a.callDynamicFunc(ctx, log, a.Platform, ep.LogsFunc())
	if err != nil {
		return nil, err
	}

	return lv.(component.LogViewer), nil
}

// callDynamicFunc calls a dynamic function which is a common pattern for
// our component interfaces. These are functions that are given to mapper,
// supplied with a series of arguments, dependency-injected, and then called.
//
// This always provides some common values for injection:
//
//   * *component.Source
//   * *datadir.Project
//
func (a *App) callDynamicFunc(
	ctx context.Context,
	log hclog.Logger,
	c interface{}, // component
	f interface{}, // function
	values ...interface{},
) (interface{}, error) {
	rawFunc, err := mapper.NewFunc(f)
	if err != nil {
		return nil, err
	}

	// Get the component directory
	cdir, ok := a.componentDirs[c]
	if !ok {
		return nil, fmt.Errorf("component dir not found for: %T", c)
	}

	// Make sure we have access to our context and logger and default args
	values = append(values,
		ctx,
		log,
		a.source,
		a.dir,
		cdir,
	)

	// Build the chain and call it
	chain, err := rawFunc.Chain(a.mappers, values...)
	if err != nil {
		return nil, err
	}
	log.Debug("function chain", "chain", chain.String())
	return chain.Call()
}

// initComponent initializes a component with the given factory and configuration
// and then sets it on the value pointed to by target.
func (a *App) initComponent(
	ctx context.Context,
	typ component.Type,
	target interface{},
	f *mapper.Factory,
	cfg *config.Component,
) error {
	// Before we do anything, the target should be a pointer. If so,
	// then we get the value of the pointer so we can set it later.
	targetV := reflect.ValueOf(target)
	if targetV.Kind() != reflect.Ptr {
		return fmt.Errorf("target value should be a pointer")
	}
	targetV = reflect.Indirect(targetV)

	// Get the factory function for this type
	fn := f.Func(cfg.Type)
	if fn == nil {
		return fmt.Errorf("unknown type: %q", cfg.Type)
	}

	// Create the data directory for this component
	cdir, err := a.dir.Component(strings.ToLower(typ.String()), cfg.Type)
	if err != nil {
		return err
	}

	// Call the factory to get our raw value (interface{} type)
	raw, err := fn.Call(ctx, a.source, a.logger, cdir)
	if err != nil {
		return err
	}

	// Store the component dir mapping
	a.componentDirs[raw] = cdir

	// We have our value so let's make sure it is the correct type.
	rawV := reflect.ValueOf(raw)
	if !rawV.Type().AssignableTo(targetV.Type()) {
		return fmt.Errorf("component %s not assigntable to type %s", rawV.Type(), targetV.Type())
	}

	// Configure the component. This will handle all the cases where no
	// config is given but required, vice versa, and everything in between.
	diag := component.Configure(raw, cfg.Body, nil)
	if diag.HasErrors() {
		return diag
	}

	// Assign our value now that we won't error anymore
	targetV.Set(rawV)

	return nil
}
