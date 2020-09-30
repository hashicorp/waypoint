package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/protocol"
)

// versionUnaryInterceptor returns a gRPC unary interceptor that inserts a hclog.Logger
// into the request context.
//
// Additionally, logUnaryInterceptor logs request and response metadata. If verbose
// is set to true, the request and response attributes are logged too.
func versionUnaryInterceptor(serverInfo *pb.VersionInfo) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := versionContext(ctx, protocol.Api, serverInfo)
		if err != nil {
			return nil, err
		}

		ctx, err = versionContext(ctx, protocol.Entrypoint, serverInfo)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// versionUnaryInterceptor returns a gRPC unary interceptor that inserts a hclog.Logger
// into the request context.
//
// Additionally, versionUnaryInterceptor logs request and response metadata. If verbose
// is set to true, the request and response attributes are logged too.
func versionStreamInterceptor(serverInfo *pb.VersionInfo) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := ss.Context()

		ctx, err := versionContext(ctx, protocol.Api, serverInfo)
		if err != nil {
			return err
		}

		ctx, err = versionContext(ctx, protocol.Entrypoint, serverInfo)
		if err != nil {
			return err
		}

		// Invoke the handler.
		return handler(srv, &versionStream{
			ServerStream: ss,
			context:      ctx,
		})
	}
}

// versionContext
func versionContext(
	ctx context.Context,
	typ protocol.Type,
	info *pb.VersionInfo,
) (context.Context, error) {
	var header string
	var server *pb.VersionInfo_ProtocolVersion
	switch typ {
	case protocol.Api:
		header = protocol.HeaderClientApiProtocol
		server = info.Api

	case protocol.Entrypoint:
		header = protocol.HeaderClientEntrypointProtocol
		server = info.Entrypoint

	default:
		return nil, status.Errorf(codes.Internal, "invalid protocol type")
	}

	// Get our metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	// Get the client version information
	vs := md[header]
	if len(vs) != 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"required header %s is not set", header)
	}
	min, current, err := protocol.ParseHeader(vs[0])
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"header %q: %s", header, err)
	}

	// Negotiate the version to use
	version, err := protocol.Negotiate(&pb.VersionInfo_ProtocolVersion{
		Current: current,
		Minimum: min,
	}, server)

	// Invoke the handler.
	return protocol.WithContext(ctx, typ, version), nil
}

type versionStream struct {
	grpc.ServerStream
	context context.Context
}

func (s *versionStream) Context() context.Context {
	return s.context
}
