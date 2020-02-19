package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/mapper"
	"github.com/mitchellh/devflow/sdk/plugin/proto"
)

// BuilderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Builder component type.
type BuilderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Builder // Impl is the concrete implementation
}

func (p *BuilderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBuilderServer(s, &builderServer{Impl: p.Impl})
	return nil
}

func (p *BuilderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &builderClient{client: proto.NewBuilderClient(c)}, nil
}

// builderClient is an implementation of component.Builder that
// communicates over gRPC.
type builderClient struct {
	client proto.BuilderClient
}

func (c *builderClient) BuildFunc() interface{} {
	// Get the build spec
	spec, err := c.client.BuildSpec(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}

	return specToFunc(spec, c.build)
}

func (c *builderClient) build(
	ctx context.Context,
	src *proto.Args_Source,
	args dynamicArgs,
) interface{} {
	return nil
}

// builderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type builderServer struct {
	Impl component.Builder
}

func (s *builderServer) BuildSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return &proto.FuncSpec{}, nil
}

func (s *builderServer) Build(
	ctx context.Context,
	args *proto.Build_Args,
) (*proto.Build_Resp, error) {
	return nil, nil
}

// specToFunc takes a FuncSpec and returns a mapper.Func that can be used.
func specToFunc(s *proto.FuncSpec, cb interface{}) *mapper.Func {
	// Build the function
	f, err := mapper.NewFunc(cb, mapper.WithType(dynamicArgsType, makeDynamicArgsMapperType(s)))
	if err != nil {
		panic(err)
	}

	return f
}

var (
	_ plugin.Plugin       = (*BuilderPlugin)(nil)
	_ plugin.GRPCPlugin   = (*BuilderPlugin)(nil)
	_ proto.BuilderServer = (*builderServer)(nil)
	_ component.Builder   = (*builderClient)(nil)
)
