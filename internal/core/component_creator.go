package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/hcl/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type Component struct {
	Value interface{}
	Info  *pb.Component

	// These fields can be accessed internally
	hooks   map[string][]*config.Hook
	labels  map[string]string
	mappers []*argmapper.Func

	// These are private, please do not access them ever except as an
	// internal Component implementation detail.
	closed bool
	plugin *plugin.Instance
}

// Close cleans up any resources associated with the Component. Close should
// always be called when the component is done being used.
func (c *Component) Close() error {
	if c == nil {
		return nil
	}

	// If we closed already do nothing.
	if c.closed {
		return nil
	}

	c.closed = true
	if c.plugin != nil {
		c.plugin.Close()
	}

	return nil
}

// componentCreator represents the configuration to initialize the component
// for a given application.
type componentCreator struct {
	Type       component.Type
	ConfigFunc func(*App, *hcl.EvalContext) (interface{}, error)
}

// componentCreatorMap contains all the components that can be initialized
// for an app.
var componentCreatorMap = map[component.Type]*componentCreator{
	component.BuilderType: {
		Type: component.BuilderType,
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Build(ctx)
		},
	},

	component.RegistryType: {
		Type: component.RegistryType,
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Registry(ctx)
		},
	},

	component.PlatformType: {
		Type: component.PlatformType,
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Deploy(ctx)
		},
	},

	component.ReleaseManagerType: {
		Type: component.ReleaseManagerType,
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Release(ctx)
		},
	},
}

// Create creates the component of the given type.
func (cc *componentCreator) Create(
	ctx context.Context,
	app *App,
	hclCtx *hcl.EvalContext,
) (*Component, error) {
	cfg, err := cc.ConfigFunc(app, hclCtx)
	if err != nil {
		return nil, err
	}

	// If we have no configuration or the use is nil or type is empty then
	// we return an error. We have to use the reflect trick here because we
	// may get a non-nil interface but nil value.
	if cfg == nil || !reflect.Indirect(reflect.ValueOf(cfg)).IsValid() {
		return nil, status.Errorf(codes.Unimplemented,
			"component type %s is not configured", cc.Type)
	}

	// This should represent an operartion otherwise we have nothing to do.
	opCfger, ok := cfg.(interface {
		Operation() *config.Operation
	})
	if !ok {
		panic(fmt.Sprintf("config %T should turn into operation", cfg))
	}
	opCfg := opCfger.Operation()

	// Start the plugin
	pinst, err := app.startPlugin(
		ctx,
		cc.Type,
		app.project.factories[cc.Type],
		opCfg.Use.Type,
	)
	if err != nil {
		return nil, err
	}

	// If we have a config, configure
	// Configure the component. This will handle all the cases where no
	// config is given but required, vice versa, and everything in between.
	diag := component.Configure(pinst.Component, opCfg.Use.Body, hclCtx)
	if diag.HasErrors() {
		pinst.Close()
		return nil, diag
	}

	// Setup hooks
	hooks := map[string][]*config.Hook{}
	for _, h := range opCfg.Hooks {
		hooks[h.When] = append(hooks[h.When], h)
	}

	return &Component{
		Value: pinst.Component,
		Info: &pb.Component{
			Type: pb.Component_Type(cc.Type),
			Name: opCfg.Use.Type,
		},

		hooks:   hooks,
		labels:  opCfg.Labels,
		mappers: pinst.Mappers,
		plugin:  pinst,
	}, nil
}
