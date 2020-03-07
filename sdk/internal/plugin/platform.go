package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
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
	proto.RegisterPlatformServer(s, &platformServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *PlatformPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &platformClient{
		client: proto.NewPlatformClient(c),
		logger: p.Logger,
	}, nil
}

// platformClient is an implementation of component.Platform over gRPC.
type platformClient struct {
	client proto.PlatformClient
	logger hclog.Logger
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

	return specToFunc(c.logger, spec, c.push)
}

func (c *platformClient) push(
	ctx context.Context,
	args dynamicArgs,
) (interface{}, error) {
	// Call our function
	resp, err := c.client.Deploy(ctx, &proto.Deploy_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the *any.Any directly.
	return resp.Result, nil
}

// platformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type platformServer struct {
	Impl    component.Platform
	Mappers []*mapper.Func
	Logger  hclog.Logger
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
	return funcToSpec(s.Logger, s.Impl.DeployFunc(), s.Mappers)
}

func (s *platformServer) Deploy(
	ctx context.Context,
	args *proto.Deploy_Args,
) (*proto.Deploy_Resp, error) {
	encoded, err := callDynamicFunc(ctx, s.Logger, args.Args, s.Impl.DeployFunc(), s.Mappers)
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
