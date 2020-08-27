package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// workspaceDestroyerClient implements component.WorkspaceDestroyer for a service that
// has the destroy methods implemented.
type workspaceDestroyerClient struct {
	Client  workspaceDestroyerProtoClient
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
}

func (c *workspaceDestroyerClient) Implements(ctx context.Context) (bool, error) {
	if c == nil {
		return false, nil
	}

	resp, err := c.Client.IsWorkspaceDestroyer(ctx, &empty.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Implements, nil
}

func (c *workspaceDestroyerClient) DestroyWorkspaceFunc() interface{} {
	impl, err := c.Implements(context.Background())
	if err != nil {
		return funcErr(err)
	}
	if !impl {
		return nil
	}

	// Get the spec
	spec, err := c.Client.DestroyWorkspaceSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.destroy,
		argmapper.Logger(c.Logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.Broker,
			Mappers: c.Mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *workspaceDestroyerClient) destroy(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) error {
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	_, err := c.Client.DestroyWorkspace(ctx, &pb.FuncSpec_Args{Args: args})
	return err
}

// workspaceDestroyerServer implements the common WorkspaceDestroyer-related RPC calls.
// This should be embedded into the service implementation.
type workspaceDestroyerServer struct {
	*base
	Impl interface{}
}

func (s *workspaceDestroyerServer) IsWorkspaceDestroyer(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.ImplementsResp, error) {
	d, ok := s.Impl.(component.WorkspaceDestroyer)
	return &pb.ImplementsResp{
		Implements: ok && d.DestroyWorkspaceFunc() != nil,
	}, nil
}

func (s *workspaceDestroyerServer) DestroyWorkspaceSpec(
	ctx context.Context,
	args *empty.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.WorkspaceDestroyer).DestroyWorkspaceFunc(),
		//argmapper.WithNoOutput(), // we only expect an error value so ignore the rest
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *workspaceDestroyerServer) DestroyWorkspace(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*empty.Empty, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	_, err := callDynamicFunc2(s.Impl.(component.WorkspaceDestroyer).DestroyWorkspaceFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// workspaceDestroyerProtoClient is the interface we expect any gRPC service that
// supports destroy to implement.
type workspaceDestroyerProtoClient interface {
	IsWorkspaceDestroyer(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.ImplementsResp, error)
	DestroyWorkspaceSpec(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.FuncSpec, error)
	DestroyWorkspace(context.Context, *pb.FuncSpec_Args, ...grpc.CallOption) (*empty.Empty, error)
}

var (
	_ component.WorkspaceDestroyer = (*workspaceDestroyerClient)(nil)
)
