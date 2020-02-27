package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
)

// RegistryPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Registry component type.
type RegistryPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Registry // Impl is the concrete implementation
	Mappers []*mapper.Func     // Mappers
}

func (p *RegistryPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRegistryServer(s, &registryServer{Impl: p.Impl, Mappers: p.Mappers})
	return nil
}

func (p *RegistryPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &registryClient{client: proto.NewRegistryClient(c)}, nil
}

// registryClient is an implementation of component.Registry over gRPC.
type registryClient struct {
	client proto.RegistryClient
}

func (c *registryClient) PushFunc() interface{} {
	// Get the spec
	spec, err := c.client.PushSpec(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}

	return specToFunc(spec, c.push)
}

func (c *registryClient) push(
	ctx context.Context,
	args dynamicArgs,
) (interface{}, error) {
	// Call our function
	resp, err := c.client.Push(ctx, &proto.Push_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the *any.Any directly.
	return resp.Result, nil
}

// registryServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type registryServer struct {
	Impl    component.Registry
	Mappers []*mapper.Func
}

func (s *registryServer) PushSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcToSpec(s.Impl.PushFunc(), s.Mappers)
}

func (s *registryServer) Push(
	ctx context.Context,
	args *proto.Push_Args,
) (*proto.Push_Resp, error) {
	encoded, err := callDynamicFunc(ctx, args.Args, s.Impl.PushFunc(), s.Mappers)
	if err != nil {
		return nil, err
	}

	return &proto.Push_Resp{Result: encoded}, nil
}

var (
	_ plugin.Plugin        = (*RegistryPlugin)(nil)
	_ plugin.GRPCPlugin    = (*RegistryPlugin)(nil)
	_ proto.RegistryServer = (*registryServer)(nil)
	_ component.Registry   = (*registryClient)(nil)
)
