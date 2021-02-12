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
	UseType    func(*App) (string, error)
	ConfigFunc func(*App, *hcl.EvalContext) (interface{}, error)
}

// componentCreatorMap contains all the components that can be initialized
// for an app.
var componentCreatorMap = map[component.Type]*componentCreator{
	component.BuilderType: {
		Type: component.BuilderType,
		UseType: func(a *App) (string, error) {
			return a.config.BuildUse(), nil
		},
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Build(ctx)
		},
	},

	component.RegistryType: {
		Type: component.RegistryType,
		UseType: func(a *App) (string, error) {
			return a.config.RegistryUse(), nil
		},
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Registry(ctx)
		},
	},

	component.PlatformType: {
		Type: component.PlatformType,
		UseType: func(a *App) (string, error) {
			return a.config.DeployUse(), nil
		},
		ConfigFunc: func(a *App, ctx *hcl.EvalContext) (interface{}, error) {
			return a.config.Deploy(ctx)
		},
	},

	component.ReleaseManagerType: {
		Type: component.ReleaseManagerType,
		UseType: func(a *App) (string, error) {
			return a.config.ReleaseUse(), nil
		},
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
	return cc.create(ctx, app, true, hclCtx)
}

// CreateNoConfig creates a component with no config loaded. This means
// hooks also won't be available.
func (cc *componentCreator) CreateNoConfig(
	ctx context.Context,
	app *App,
) (*Component, error) {
	return cc.create(ctx, app, false, nil)
}

func (cc *componentCreator) create(
	ctx context.Context,
	app *App,
	loadConfig bool,
	hclCtx *hcl.EvalContext,
) (*Component, error) {
	useType, err := cc.UseType(app)
	if err != nil {
		return nil, err
	}
	if useType == "" {
		return nil, status.Errorf(codes.Unimplemented,
			"no plugin type declared for type: %s", cc.Type.String())
	}

	// Start the plugin
	pinst, err := app.startPlugin(
		ctx,
		cc.Type,
		app.project.factories[cc.Type],
		useType,
	)
	if err != nil {
		return nil, err
	}

	result := &Component{
		Value: pinst.Component,
		Info: &pb.Component{
			Type: pb.Component_Type(cc.Type),
			Name: useType,
		},

		mappers: pinst.Mappers,
		plugin:  pinst,
	}

	if loadConfig {
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

		// If we have a config, configure
		// Configure the component. This will handle all the cases where no
		// config is given but required, vice versa, and everything in between.
		diag := opCfg.Configure(pinst.Component, hclCtx)
		if diag.HasErrors() {
			pinst.Close()
			return nil, diag
		}

		// Setup hooks
		hooks := map[string][]*config.Hook{}
		for _, h := range opCfg.Hooks {
			hooks[h.When] = append(hooks[h.When], h)
		}

		result.hooks = hooks
		result.labels = opCfg.Labels
	}

	return result, nil
}
