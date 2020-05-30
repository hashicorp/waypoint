package plugin

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-argmapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/proto"
)

// LogPlatformPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the LogPlatform component type.
type LogPlatformPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.LogPlatform // Impl is the concrete implementation
	Mappers []*argmapper.Func     // Mappers
	Logger  hclog.Logger          // Logger
}

func (p *LogPlatformPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterLogPlatformServer(s, &logPlatformServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
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
		broker: broker,
	}, nil
}

// logPlatformClient is an implementation of component.LogPlatform over gRPC.
type logPlatformClient struct {
	client proto.LogPlatformClient
	logger hclog.Logger
	broker *plugin.GRPCBroker
}

func (c *logPlatformClient) LogsFunc() interface{} {
	// Get the spec
	spec, err := c.client.LogsSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	// We don't want to be a mapper
	spec.Result = nil

	return funcspec.Func(spec, c.logs, argmapper.Logger(c.logger))
}

func (c *logPlatformClient) logs(
	ctx context.Context,
	args funcspec.Args,
) (component.LogViewer, error) {
	// Call our function
	resp, err := c.client.Logs(ctx, &proto.FuncSpec_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// Get the stream ID and connect to it
	conn, err := c.broker.Dial(resp.StreamId)
	if err != nil {
		return nil, err
	}

	return &logViewerClient{
		Client: proto.NewLogViewerClient(conn),
		Logger: c.logger.Named("logviewer"),
	}, nil
}

// logPlatformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type logPlatformServer struct {
	Impl    component.LogPlatform
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *logPlatformServer) LogsSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.LogsFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),

		// We expect a component.LogViewer output type and not a proto.Message
		argmapper.FilterOutput(argmapper.FilterType(
			reflect.TypeOf((*component.LogViewer)(nil)).Elem()),
		),
	)
}

func (s *logPlatformServer) Logs(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Logs_Resp, error) {
	result, err := callDynamicFunc2(s.Impl.LogsFunc(), args.Args,
		argmapper.Typed(ctx),
		argmapper.ConverterFunc(s.Mappers...))
	if err != nil {
		return nil, err
	}

	lv, ok := result.(component.LogViewer)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"plugin Logs function should've returned a component.LogViewer, got %T",
			result)
	}

	// Get the ID for the server we're going to start to run our viewer
	id := s.Broker.NextId()

	// Start our server
	go s.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
		server := plugin.DefaultGRPCServer(opts)
		proto.RegisterLogViewerServer(server, &logViewerServer{
			Impl:   lv,
			Logger: s.Logger.Named("logviewer"),
		})
		return server
	})

	return &proto.Logs_Resp{StreamId: id}, nil
}

var (
	_ plugin.Plugin           = (*LogPlatformPlugin)(nil)
	_ plugin.GRPCPlugin       = (*LogPlatformPlugin)(nil)
	_ proto.LogPlatformServer = (*logPlatformServer)(nil)
	_ component.LogPlatform   = (*logPlatformClient)(nil)
)
