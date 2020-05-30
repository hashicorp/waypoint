package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-argmapper"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/plugincomponent"
	"github.com/hashicorp/waypoint/sdk/proto"
)

// RegistryPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Registry component type.
type RegistryPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Registry // Impl is the concrete implementation
	Mappers []*argmapper.Func  // Mappers
	Logger  hclog.Logger       // Logger
}

func (p *RegistryPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRegistryServer(s, &registryServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *RegistryPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &registryClient{
		client: proto.NewRegistryClient(c),
		logger: p.Logger,
	}, nil
}

// registryClient is an implementation of component.Registry over gRPC.
type registryClient struct {
	client proto.RegistryClient
	logger hclog.Logger
}

func (c *registryClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *registryClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *registryClient) PushFunc() interface{} {
	// Get the spec
	spec, err := c.client.PushSpec(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}

	// We don't want to be a mapper
	spec.Result = nil

	return funcspec.Func(spec, c.push, argmapper.Logger(c.logger))
}

func (c *registryClient) push(
	ctx context.Context,
	args funcspec.Args,
) (component.Artifact, error) {
	// Call our function
	resp, err := c.client.Push(ctx, &proto.Push_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the *any.Any directly.
	return &plugincomponent.Artifact{Any: resp.Result}, nil
}

// registryServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type registryServer struct {
	Impl    component.Registry
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (s *registryServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *registryServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *registryServer) PushSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.PushFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger))
}

func (s *registryServer) Push(
	ctx context.Context,
	args *proto.Push_Args,
) (*proto.Push_Resp, error) {
	encoded, err := callDynamicFuncAny2(s.Impl.PushFunc(), args.Args,
		argmapper.Typed(ctx),
		argmapper.ConverterFunc(s.Mappers...),
	)
	if err != nil {
		return nil, err
	}

	return &proto.Push_Resp{Result: encoded}, nil
}

var (
	_ plugin.Plugin                = (*RegistryPlugin)(nil)
	_ plugin.GRPCPlugin            = (*RegistryPlugin)(nil)
	_ proto.RegistryServer         = (*registryServer)(nil)
	_ component.Registry           = (*registryClient)(nil)
	_ component.Configurable       = (*registryClient)(nil)
	_ component.ConfigurableNotify = (*registryClient)(nil)
)
