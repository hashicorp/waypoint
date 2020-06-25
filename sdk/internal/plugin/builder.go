package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
	"github.com/hashicorp/waypoint/sdk/internal/plugincomponent"
	"github.com/hashicorp/waypoint/sdk/proto"
)

// BuilderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Builder component type.
type BuilderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Builder // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *BuilderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	base := &base{
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	}

	proto.RegisterBuilderServer(s, &builderServer{
		base: base,
		Impl: p.Impl,

		authenticatorServer: &authenticatorServer{
			base: base,
			Impl: p.Impl,
		},
	})
	return nil
}

func (p *BuilderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := &builderClient{
		client:  proto.NewBuilderClient(c),
		logger:  p.Logger,
		broker:  broker,
		mappers: p.Mappers,
	}

	authenticator := &authenticatorClient{
		Client:  client.client,
		Logger:  client.logger,
		Broker:  client.broker,
		Mappers: client.mappers,
	}
	if ok, err := authenticator.Implements(ctx); err != nil {
		return nil, err
	} else if ok {
		p.Logger.Info("builder plugin capable of auth")
	} else {
		authenticator = nil
	}

	result := &mix_Builder_Authenticator{
		ConfigurableNotify: client,
		Builder:            client,
		Authenticator:      authenticator,
	}

	return result, nil
}

// builderClient is an implementation of component.Builder that
// communicates over gRPC.
type builderClient struct {
	client  proto.BuilderClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *builderClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *builderClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *builderClient) BuildFunc() interface{} {
	// Get the build spec
	spec, err := c.client.BuildSpec(context.Background(), &proto.Empty{})
	if err != nil {
		return funcErr(err)
	}

	// We don't want to be a mapper
	spec.Result = nil

	return funcspec.Func(spec, c.build,
		argmapper.Logger(c.logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.broker,
			Mappers: c.mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *builderClient) build(
	ctx context.Context,
	args funcspec.Args,
) (component.Artifact, error) {
	// Call our function
	resp, err := c.client.Build(ctx, &proto.FuncSpec_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the
	return &plugincomponent.Artifact{Any: resp.Result}, nil
}

// builderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type builderServer struct {
	*base
	*authenticatorServer

	Impl component.Builder
}

func (s *builderServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *builderServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *builderServer) BuildSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.BuildFunc(),
		argmapper.Logger(s.Logger),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(s.internal()),
	)
}

func (s *builderServer) Build(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Build_Resp, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	encoded, err := callDynamicFuncAny2(s.Impl.BuildFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(ctx),
		argmapper.Typed(internal),
	)
	if err != nil {
		return nil, err
	}

	return &proto.Build_Resp{Result: encoded}, nil
}

var (
	_ plugin.Plugin                = (*BuilderPlugin)(nil)
	_ plugin.GRPCPlugin            = (*BuilderPlugin)(nil)
	_ proto.BuilderServer          = (*builderServer)(nil)
	_ component.Builder            = (*builderClient)(nil)
	_ component.Configurable       = (*builderClient)(nil)
	_ component.ConfigurableNotify = (*builderClient)(nil)
)
