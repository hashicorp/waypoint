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

// ReleaseManagerPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the ReleaseManager component type.
type ReleaseManagerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.ReleaseManager // Impl is the concrete implementation
	Mappers []*argmapper.Func        // Mappers
	Logger  hclog.Logger             // Logger
}

func (p *ReleaseManagerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterReleaseManagerServer(s, &releaseManagerServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *ReleaseManagerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &releaseManagerClient{
		client: proto.NewReleaseManagerClient(c),
		logger: p.Logger,
	}, nil
}

// releaseManagerClient is an implementation of component.ReleaseManager that
// communicates over gRPC.
type releaseManagerClient struct {
	client proto.ReleaseManagerClient
	logger hclog.Logger
}

func (c *releaseManagerClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *releaseManagerClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *releaseManagerClient) ReleaseFunc() interface{} {
	// Get the build spec
	spec, err := c.client.ReleaseSpec(context.Background(), &proto.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.build, argmapper.Logger(c.logger))
}

func (c *releaseManagerClient) build(
	ctx context.Context,
	args funcspec.Args,
) (component.Release, error) {
	// Call our function
	resp, err := c.client.Release(ctx, &proto.Release_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the
	return &plugincomponent.Release{
		Any:     resp.Result,
		Release: resp.Release,
	}, nil
}

// releaseManagerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type releaseManagerServer struct {
	Impl    component.ReleaseManager
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (s *releaseManagerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *releaseManagerServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *releaseManagerServer) ReleaseSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcspec.Spec(s.Impl.ReleaseFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger))
}

func (s *releaseManagerServer) Release(
	ctx context.Context,
	args *proto.Release_Args,
) (*proto.Release_Resp, error) {
	raw, err := callDynamicFunc2(s.Impl.ReleaseFunc(), args.Args,
		argmapper.Typed(ctx),
		argmapper.ConverterFunc(s.Mappers...),
	)
	if err != nil {
		return nil, err
	}
	encoded, err := component.ProtoAny(raw)
	if err != nil {
		return nil, err
	}

	release := raw.(component.Release)
	return &proto.Release_Resp{
		Result: encoded,
		Release: &proto.Release{
			Url: release.URL(),
		},
	}, nil
}

var (
	_ plugin.Plugin                = (*ReleaseManagerPlugin)(nil)
	_ plugin.GRPCPlugin            = (*ReleaseManagerPlugin)(nil)
	_ proto.ReleaseManagerServer   = (*releaseManagerServer)(nil)
	_ component.ReleaseManager     = (*releaseManagerClient)(nil)
	_ component.Configurable       = (*releaseManagerClient)(nil)
	_ component.ConfigurableNotify = (*releaseManagerClient)(nil)
)
