package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// authenticatorClient is the interface implemented by all gRPC services that
// have the authenticator RPC methods.
type authenticatorProtoClient interface {
	Auth(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error)
	AuthValidate(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error)
}

// authenticatorClient implements component.Authenticator for a service that
// has the authenticator methods implemented
type authenticatorClient struct {
	Client  authenticatorProtoClient
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
}

// authenticatorServer implements the common Authenticator-related RPC calls.
type authenticatorServer struct {
	*base
	Impl interface{}
}

func (s *authenticatorServer) AuthSpec(
	ctx context.Context,
	args *pb.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.Authenticator).AuthFunc(),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *authenticatorServer) Auth(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*empty.Empty, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	_, err := callDynamicFunc2(s.Impl.(component.Authenticator).AuthFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
