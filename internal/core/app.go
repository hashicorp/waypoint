package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/datadir"
	"github.com/mitchellh/devflow/internal/mapper"
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

	source  *component.Source
	logger  hclog.Logger
	dir     *datadir.App
	mappers []*mapper.Func
}

// newApp creates an App for the given project and configuration. This will
// initialize and configure all the components of this application. An error
// will be returned if this app fails to initialize: configuration is invalid,
// a component could not be found, etc.
func newApp(ctx context.Context, p *Project, cfg *config.App) (*App, error) {
	// Initialize
	app := &App{
		source:  &component.Source{App: cfg.Name, Path: "."},
		logger:  p.logger.Named("app").Named(cfg.Name),
		mappers: p.mappers,
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
	}
	for _, c := range components {
		if c.Config == nil {
			// This component is not set, ignore.
			continue
		}

		err := app.initComponent(ctx, c.Target, p.factories[c.Type], c.Config)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Build builds the artifact from source for this app.
func (a *App) Build(ctx context.Context) (component.Artifact, error) {
	log := a.logger.Named("build")

	buildFunc, err := mapper.NewFunc(a.Builder.BuildFunc())
	if err != nil {
		return nil, err
	}

	chain, err := buildFunc.Chain(a.mappers,
		ctx,
		log,
		a.source,
		a.dir,
	)
	if err != nil {
		return nil, err
	}
	log.Debug("function chain", "chain", chain.String())

	buildArtifact, err := chain.Call()
	if err != nil {
		return nil, err
	}

	return buildArtifact.(component.Artifact), nil
}

func (a *App) Push(component.Artifact) (component.Artifact, error) { return nil, nil }

func (a *App) Deploy(component.Artifact) (component.Deployment, error) { return nil, nil }

// initComponent initializes a component with the given factory and configuration
// and then sets it on the value pointed to by target.
func (a *App) initComponent(
	ctx context.Context,
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

	// Call the factory to get our raw value (interface{} type)
	raw, err := fn.Call(ctx, a.source, a.logger)
	if err != nil {
		return err
	}

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
