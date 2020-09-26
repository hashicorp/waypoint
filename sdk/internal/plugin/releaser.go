package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/docs"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
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
	base := &base{
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	}

	proto.RegisterReleaseManagerServer(s, &releaseManagerServer{
		base: base,
		Impl: p.Impl,

		authenticatorServer: &authenticatorServer{
			base: base,
			Impl: p.Impl,
		},

		destroyerServer: &destroyerServer{
			base: base,
			Impl: p.Impl,
		},

		workspaceDestroyerServer: &workspaceDestroyerServer{
			base: base,
			Impl: p.Impl,
		},
	})
	return nil
}

func (p *ReleaseManagerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := &releaseManagerClient{
		client:  proto.NewReleaseManagerClient(c),
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
		p.Logger.Info("release plugin capable of auth")
	} else {
		authenticator = nil
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
		p.Logger.Info("release plugin capable of destroy")
	} else {
		destroyer = nil
	}

	// Compose workspace destroyer
	wsDestroyer := &workspaceDestroyerClient{
		Client:  client.client,
		Logger:  client.logger,
		Broker:  client.broker,
		Mappers: client.mappers,
	}
	if ok, err := wsDestroyer.Implements(ctx); err != nil {
		return nil, err
	} else if ok {
		p.Logger.Info("platform plugin capable of destroy")
	} else {
		wsDestroyer = nil
	}

	result := &mix_ReleaseManager_Authenticator{
		ConfigurableNotify: client,
		ReleaseManager:     client,
		Authenticator:      authenticator,
		Destroyer:          destroyer,
		WorkspaceDestroyer: wsDestroyer,
		Documented:         client,
	}

	return result, nil
}

// releaseManagerClient is an implementation of component.ReleaseManager that
// communicates over gRPC.
type releaseManagerClient struct {
	client  proto.ReleaseManagerClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *releaseManagerClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *releaseManagerClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *releaseManagerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *releaseManagerClient) ReleaseFunc() interface{} {
	if c == nil || c.client == nil {
		return nil
	}

	// Get the build spec
	spec, err := c.client.ReleaseSpec(context.Background(), &proto.Empty{})
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

func (c *releaseManagerClient) build(
	ctx context.Context,
	args funcspec.Args,
) (component.Release, error) {
	// Call our function
	resp, err := c.client.Release(ctx, &proto.FuncSpec_Args{Args: args})
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
	*base
	*authenticatorServer
	*destroyerServer
	*workspaceDestroyerServer

	Impl component.ReleaseManager
}

func (s *releaseManagerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *releaseManagerServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
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
	if s.Impl == nil {
		return nil, status.Errorf(codes.Unimplemented, "plugin does not implement: release manager")
	}

	return funcspec.Spec(s.Impl.ReleaseFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *releaseManagerServer) Release(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Release_Resp, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	raw, err := callDynamicFunc2(s.Impl.ReleaseFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(ctx),
		argmapper.Typed(internal),
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
	_ component.Documented         = (*releaseManagerClient)(nil)
	_ component.ConfigurableNotify = (*releaseManagerClient)(nil)
)
