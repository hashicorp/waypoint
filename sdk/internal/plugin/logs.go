package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/internal/funcspec"
	"github.com/mitchellh/devflow/sdk/proto"
)

// LogPlatformPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the LogPlatform component type.
type LogPlatformPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.LogPlatform // Impl is the concrete implementation
	Mappers []*mapper.Func        // Mappers
	Logger  hclog.Logger          // Logger
}

func (p *LogPlatformPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterLogPlatformServer(s, &logPlatformServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *LogPlatformPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &logPlatformClient{
		client: proto.NewLogPlatformClient(c),
		logger: p.Logger,
	}, nil
}

// logPlatformClient is an implementation of component.LogPlatform over gRPC.
type logPlatformClient struct {
	client proto.LogPlatformClient
	logger hclog.Logger
}

func (c *logPlatformClient) LogsFunc() interface{} {
	// Get the spec
	spec, err := c.client.LogsSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.push, funcspec.WithLogger(c.logger))
}

func (c *logPlatformClient) push(
	ctx context.Context,
	args funcspec.Args,
) (interface{}, error) {
	/*
		// Call our function
		resp, err := c.client.Deploy(ctx, &proto.Deploy_Args{Args: args})
		if err != nil {
			return nil, err
		}

		// We return the *any.Any directly.
		return resp.Result, nil
	*/
	return nil, nil
}

// logPlatformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type logPlatformServer struct {
	Impl    component.LogPlatform
	Mappers []*mapper.Func
	Logger  hclog.Logger
}

func (s *logPlatformServer) LogsSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.LogsFunc(),
		funcspec.WithMappers(s.Mappers),
		funcspec.WithLogger(s.Logger))
}

func (s *logPlatformServer) Logs(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Logs_Resp, error) {
	_, err := callDynamicFunc(ctx, s.Logger, args.Args, s.Impl.LogsFunc(), s.Mappers)
	if err != nil {
		return nil, err
	}

	return &proto.Logs_Resp{}, nil
}

var (
	_ plugin.Plugin           = (*LogPlatformPlugin)(nil)
	_ plugin.GRPCPlugin       = (*LogPlatformPlugin)(nil)
	_ proto.LogPlatformServer = (*logPlatformServer)(nil)
	_ component.LogPlatform   = (*logPlatformClient)(nil)
)
