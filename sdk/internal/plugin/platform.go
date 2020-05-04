package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal-shared/mapper"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
	"github.com/hashicorp/waypoint/sdk/internal/plugincomponent"
	"github.com/hashicorp/waypoint/sdk/proto"
)

// PlatformPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Platform component type.
type PlatformPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Platform // Impl is the concrete implementation
	Mappers []*mapper.Func     // Mappers
	Logger  hclog.Logger       // Logger
}

func (p *PlatformPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	base := &base{
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	}

	proto.RegisterPlatformServer(s, &platformServer{
		base: base,
		destroyerServer: &destroyerServer{
			base: base,
			Impl: p.Impl,
		},

		Impl: p.Impl,
	})

	return nil
}

func (p *PlatformPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	// Build our client to the platform service
	client := &platformClient{
		client:  proto.NewPlatformClient(c),
		logger:  p.Logger,
		broker:  broker,
		mappers: p.Mappers,
	}

	// Check if we also implement the LogPlatform
	var logPlatform component.LogPlatform
	resp, err := client.client.IsLogPlatform(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	if resp.Implements {
		raw, err := (&LogPlatformPlugin{
			Logger: p.Logger,
		}).GRPCClient(ctx, broker, c)
		if err != nil {
			return nil, err
		}

		p.Logger.Info("platform plugin capable of logs")
		logPlatform = raw.(component.LogPlatform)
	}

	// Compose destroyer
	destroyer := &destroyerClient{
		Client:  client.client,
		Logger:  client.logger,
		Broker:  client.broker,
		Mappers: client.mappers,
	}
	if ok, err := destroyer.Implements(ctx); err != nil {
		return nil, err
	} else if ok {
		p.Logger.Info("platform plugin capable of destroy")
	} else {
		destroyer = nil
	}

	// Figure out what we're returning
	var result interface{} = client
	switch {
	case logPlatform != nil && destroyer != nil:
		result = &mix_Platform_Log_Destroy{
			ConfigurableNotify: client,
			Platform:           client,
			LogPlatform:        logPlatform,
			Destroyer:          destroyer,
		}

	case logPlatform != nil:
		result = &mix_Platform_Log{
			ConfigurableNotify: client,
			Platform:           client,
			LogPlatform:        logPlatform,
		}

	case destroyer != nil:
		result = &mix_Platform_Destroy{
			ConfigurableNotify: client,
			Platform:           client,
			Destroyer:          destroyer,
		}
	}

	return result, nil
}

// platformClient is an implementation of component.Platform over gRPC.
type platformClient struct {
	client  proto.PlatformClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*mapper.Func
}

func (c *platformClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *platformClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *platformClient) DeployFunc() interface{} {
	// Get the spec
	spec, err := c.client.DeploySpec(context.Background(), &proto.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.deploy,
		funcspec.WithLogger(c.logger),
		funcspec.WithValues(&pluginargs.Internal{
			Broker:  c.broker,
			Mappers: c.mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *platformClient) deploy(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) (component.Deployment, error) {
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	resp, err := c.client.Deploy(ctx, &proto.Deploy_Args{Args: args})
	if err != nil {
		return nil, err
	}

	return &plugincomponent.Deployment{Any: resp.Result}, nil
}

// platformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type platformServer struct {
	*base
	*destroyerServer

	Impl    component.Platform
	Mappers []*mapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *platformServer) IsLogPlatform(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.ImplementsResp, error) {
	_, ok := s.Impl.(component.LogPlatform)
	return &proto.ImplementsResp{Implements: ok}, nil
}

func (s *platformServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *platformServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *platformServer) DeploySpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.DeployFunc(),
		funcspec.WithMappers(s.Mappers),
		funcspec.WithLogger(s.Logger),
		funcspec.WithValues(s.internal()),
	)
}

func (s *platformServer) Deploy(
	ctx context.Context,
	args *proto.Deploy_Args,
) (*proto.Deploy_Resp, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	encoded, err := callDynamicFuncAny(ctx, s.Logger, args.Args, s.Impl.DeployFunc(), s.Mappers, internal)
	if err != nil {
		return nil, err
	}

	return &proto.Deploy_Resp{Result: encoded}, nil
}

var (
	_ plugin.Plugin                = (*PlatformPlugin)(nil)
	_ plugin.GRPCPlugin            = (*PlatformPlugin)(nil)
	_ proto.PlatformServer         = (*platformServer)(nil)
	_ component.Platform           = (*platformClient)(nil)
	_ component.Configurable       = (*platformClient)(nil)
	_ component.ConfigurableNotify = (*platformClient)(nil)
)
