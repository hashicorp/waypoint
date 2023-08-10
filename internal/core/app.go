// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/datadir"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// App represents a single application and exposes all the operations
// that can be performed on an application.
//
// An App is only valid if it was returned by Project.App. The behavior of
// App if constructed in any other way is undefined and likely to result
// in crashes.
type App struct {
	// UI is the UI that should be used for any output that is specific
	// to this app vs the project UI.
	UI terminal.UI

	project   *Project
	config    *config.App
	ref       *pb.Ref_Application
	workspace *pb.Ref_Workspace
	client    pb.WaypointClient
	source    *component.Source
	jobInfo   *component.JobInfo
	logger    hclog.Logger
	dir       *datadir.App
	mappers   []*argmapper.Func
	closers   []func() error
}

type appComponent struct {
	// Info is the protobuf metadata for this component.
	Info *pb.Component

	// Dir is the data directory for this component.
	Dir *datadir.Component

	// Labels are the set of labels that were set for this component.
	// This isn't merged yet with parent labels and app.mergeLabels must
	// be called.
	Labels map[string]string

	// Hooks are the hooks associated with this component keyed by their When value
	Hooks map[string][]*config.Hook
}

// newApp creates an App for the given project and configuration. This will
// initialize and configure all the components of this application. An error
// will be returned if this app fails to initialize: configuration is invalid,
// a component could not be found, etc.
func newApp(
	ctx context.Context,
	p *Project,
	cfg *config.App,
) (*App, error) {
	// Copy the job info structure, a shallow copy is fine.
	jobInfo := *p.jobInfo
	jobInfo.App = cfg.Name

	// Initialize
	app := &App{
		project: p,
		client:  p.client,
		source:  &component.Source{App: cfg.Name, Path: cfg.Path},
		jobInfo: &jobInfo,
		logger:  p.logger.Named("app").Named(cfg.Name),
		ref: &pb.Ref_Application{
			Application: cfg.Name,
			Project:     p.name,
		},
		workspace: p.WorkspaceRef(),
		config:    cfg,

		// very important below that we allocate a new slice since we modify
		mappers: append([]*argmapper.Func{}, p.mappers...),

		// set the UI, which for now is identical to project but in the
		// future should probably change as we do app-scoping, parallelization,
		// etc.
		UI: p.UI,
	}

	// Setup our directory
	dir, err := p.dir.App(cfg.Name)
	if err != nil {
		return nil, err
	}
	app.dir = dir

	// Initialize mappers if we have those
	if f, ok := p.factories[component.MapperType]; ok {
		err = app.initMappers(ctx, f)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Close is called to clean up any resources. This should be called
// whenever the app is done being used. This will be called by Project.Close.
func (a *App) Close() error {
	for _, c := range a.closers {
		c()
	}

	return nil
}

// Ref returns the reference to this application for us in API calls.
func (a *App) Ref() *pb.Ref_Application {
	return a.ref
}

// Components initializes and returns all the components that are defined
// for this app across all stages. The caller must call close on all the
// components to clean up resources properly.
func (a *App) Components(ctx context.Context) ([]*Component, error) {
	var results []*Component
	for _, cc := range componentCreatorMap {
		c, err := cc.CreateNoConfig(ctx, a)
		if status.Code(err) == codes.Unimplemented {
			c = nil
			err = nil
		}
		if err != nil {
			// Make sure we clean ourselves up in an error case.
			for _, r := range results {
				r.Close()
			}

			return nil, err
		}

		if c != nil {
			results = append(results, c)
		}
	}

	return results, nil
}

// mergeLabels merges the set of labels given. See project.mergeLabels.
// This is the app-specific version that adds the proper app-specific labels
// as necessary.
func (a *App) mergeLabels(ls ...map[string]string) map[string]string {
	ls = append([]map[string]string{a.config.Labels}, ls...)
	return a.project.mergeLabels(ls...)
}

// startPlugin starts a plugin with the given type and name. The returned
// value must be closed to clean up the plugin properly.
func (a *App) startPlugin(
	ctx context.Context,
	typ component.Type,
	f *factory.Factory,
	n string,
) (*plugin.Instance, error) {
	log := a.logger.Named(strings.ToLower(typ.String()))

	// Get the factory function for this type
	fn := f.Func(n)
	if fn == nil {
		return nil, fmt.Errorf("unknown type: %q", n)
	}

	// Call the factory to get our raw value (interface{} type)
	fnResult := fn.Call(argmapper.Typed(ctx, a.source, log))
	if err := fnResult.Err(); err != nil {
		return nil, err
	}
	log.Info("initialized component", "type", typ.String())
	raw := fnResult.Out(0)

	// If we have a plugin.Instance then we can extract other information
	// from this plugin. We accept pure factories too that don't return
	// this so we type-check here.
	pinst, ok := raw.(*plugin.Instance)
	if !ok {
		pinst = &plugin.Instance{
			Component: raw,
			Close:     func() {},
		}
	}

	return pinst, nil
}

// callDynamicFunc calls a dynamic function which is a common pattern for
// our component interfaces. These are functions that are given to mapper,
// supplied with a series of arguments, dependency-injected, and then called.
//
// This always provides some common values for injection:
//
//   - *component.Source
//   - *datadir.Project
//   - history.Client
func (a *App) callDynamicFunc(
	ctx context.Context,
	log hclog.Logger,
	result interface{}, // expected result type
	c *Component, // component
	f interface{}, // function
	args ...argmapper.Arg,
) (interface{}, error) {
	// We allow f to be a *mapper.Func because our plugin system creates
	// a func directly due to special argument types.
	// TODO: test
	rawFunc, ok := f.(*argmapper.Func)
	if !ok {
		var err error
		rawFunc, err = argmapper.NewFunc(f, argmapper.Logger(log))
		if err != nil {
			return nil, err
		}
	}

	// Be sure that the status is closed after every operation so we don't leak
	// weird output outside the normal execution.
	defer a.UI.Status().Close()

	// Make sure we have access to our context and logger and default args
	args = append(args,
		argmapper.ConverterFunc(a.mappers...),
		argmapper.Typed(
			ctx,
			log,
			a.source,
			a.jobInfo,
			a.dir,
			a.UI,
		),

		argmapper.Named("labels", &component.LabelSet{Labels: c.labels}),
	)

	// Build the chain and call it
	callResult := rawFunc.Call(args...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}
	raw := callResult.Out(0)

	// If we don't have an expected result type, then just return as-is.
	// Otherwise, we need to verify the result type matches properly.
	if result == nil {
		return raw, nil
	}

	// Verify
	interfaceType := reflect.TypeOf(result).Elem()
	if rawType := reflect.TypeOf(raw); !rawType.Implements(interfaceType) {
		return nil, status.Errorf(codes.FailedPrecondition,
			"operation expected result type %s, got %s",
			interfaceType.String(),
			rawType.String())
	}

	return raw, nil
}

// initMappers initializes plugins that are just mappers.
func (a *App) initMappers(
	ctx context.Context,
	f *factory.Factory,
) error {
	log := a.logger

	for _, name := range f.Registered() {
		plog := log.With("name", name)
		plog.Debug("loading mapper plugin")

		// Start the component
		pinst, err := a.startPlugin(ctx, component.MapperType, f, name)
		if err != nil {
			return err
		}

		// If we have no mappers in this plugin, then exit the plugin and
		// continue. This will keep us from having non-mapper plugins just
		// running in memory.
		if len(pinst.Mappers) == 0 {
			plog.Debug("no mappers advertised by plugin, closing")
			pinst.Close()
			continue
		}

		// We store the mappers
		plog.Debug("registered mappers", "len", len(pinst.Mappers))
		a.mappers = append(a.mappers, pinst.Mappers...)

		// Add this to our closer list
		a.closers = append(a.closers, func() error {
			pinst.Close()
			return nil
		})
	}

	return nil
}
