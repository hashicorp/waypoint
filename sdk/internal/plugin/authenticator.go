package plugin

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
	"github.com/hashicorp/waypoint/sdk/internal/testproto"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// authenticatorProtoClient is the interface implemented by all gRPC services that
// have the authenticator RPC methods.
type authenticatorProtoClient interface {
	IsAuthenticator(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.ImplementsResp, error)
	Auth(context.Context, *pb.FuncSpec_Args, ...grpc.CallOption) (*pb.Auth_AuthResponse, error)
	ValidateAuth(context.Context, *pb.FuncSpec_Args, ...grpc.CallOption) (*empty.Empty, error)
	AuthSpec(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.FuncSpec, error)
	ValidateAuthSpec(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.FuncSpec, error)
}

// authenticatorClient implements component.Authenticator for a service that
// has the authenticator methods implemented
type authenticatorClient struct {
	Client  authenticatorProtoClient
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
}

func (c *authenticatorClient) Implements(ctx context.Context) (bool, error) {
	resp, err := c.Client.IsAuthenticator(ctx, &empty.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Implements, nil
}

func (c *authenticatorClient) AuthFunc() interface{} {
	// Get the spec
	spec, err := c.Client.AuthSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	// We don't want to be a mapper
	spec.Result = nil

	return funcspec.Func(spec, c.auth,
		argmapper.Logger(c.Logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.Broker,
			Mappers: c.Mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *authenticatorClient) ValidateAuthFunc() interface{} {
	// Get the spec
	spec, err := c.Client.ValidateAuthSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	// We don't want to be a mapper
	spec.Result = nil

	return funcspec.Func(spec, c.validateAuth,
		argmapper.Logger(c.Logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.Broker,
			Mappers: c.Mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *authenticatorClient) auth(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) (*component.AuthResult, error) {
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	resp, err := c.Client.Auth(ctx, &pb.FuncSpec_Args{Args: args})
	if err != nil {
		return nil, err
	}

	return &component.AuthResult{
		Authenticated: resp.Authenticated,
	}, nil
}

func (c *authenticatorClient) validateAuth(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) error {
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	_, err := c.Client.ValidateAuth(ctx, &pb.FuncSpec_Args{Args: args})
	return err
}

// authenticatorServer implements the common Authenticator-related RPC calls.
type authenticatorServer struct {
	*base
	Impl interface{}
}

func (s *authenticatorServer) IsAuthenticator(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.ImplementsResp, error) {
	_, ok := s.Impl.(component.Authenticator)
	return &pb.ImplementsResp{Implements: ok}, nil
}

func (s *authenticatorServer) AuthSpec(
	ctx context.Context,
	args *empty.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.Authenticator).AuthFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),

		// We expect a auth result.
		argmapper.FilterOutput(argmapper.FilterOr(
			argmapper.FilterType(reflect.TypeOf((*component.AuthResult)(nil))),

			// We expect this for tests.
			argmapper.FilterType(reflect.TypeOf((*testproto.Data)(nil))),
		)),
	)
}

func (s *authenticatorServer) Auth(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*pb.Auth_AuthResponse, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	raw, err := callDynamicFunc2(s.Impl.(component.Authenticator).AuthFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	result, ok := raw.(*component.AuthResult)
	if !ok {
		return &pb.Auth_AuthResponse{
			Authenticated: false,
		}, nil
	}

	return &pb.Auth_AuthResponse{
		Authenticated: result.Authenticated,
	}, nil
}

func (s *authenticatorServer) ValidateAuthSpec(
	ctx context.Context,
	args *empty.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.Authenticator).ValidateAuthFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *authenticatorServer) ValidateAuth(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*empty.Empty, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	_, err := callDynamicFunc2(s.Impl.(component.Authenticator).ValidateAuthFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

var (
	_ component.Authenticator = (*authenticatorClient)(nil)
)
