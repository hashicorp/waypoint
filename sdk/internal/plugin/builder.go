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
	"github.com/hashicorp/waypoint/sdk/internal/plugincomponent"
	"github.com/hashicorp/waypoint/sdk/proto"
)

// BuilderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Builder component type.
type BuilderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Builder // Impl is the concrete implementation
	Mappers []*mapper.Func    // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *BuilderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBuilderServer(s, &builderServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *BuilderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &builderClient{
		client: proto.NewBuilderClient(c),
		logger: p.Logger,
	}, nil
}

// builderClient is an implementation of component.Builder that
// communicates over gRPC.
type builderClient struct {
	client proto.BuilderClient
	logger hclog.Logger
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

	return funcspec.Func(spec, c.build, funcspec.WithLogger(c.logger))
}

func (c *builderClient) build(
	ctx context.Context,
	args funcspec.Args,
) (component.Artifact, error) {
	// Call our function
	resp, err := c.client.Build(ctx, &proto.Build_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the
	return &plugincomponent.Artifact{Any: resp.Result}, nil
}

// builderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type builderServer struct {
	Impl    component.Builder
	Mappers []*mapper.Func
	Logger  hclog.Logger
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
		funcspec.WithMappers(s.Mappers),
		funcspec.WithLogger(s.Logger))
}

func (s *builderServer) Build(
	ctx context.Context,
	args *proto.Build_Args,
) (*proto.Build_Resp, error) {
	encoded, err := callDynamicFuncAny(ctx, s.Logger, args.Args, s.Impl.BuildFunc(), s.Mappers)
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
