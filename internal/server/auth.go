package server

import (
	"context"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// An interface implemented by something that wishes to authenticate the server
// actions.
type AuthChecker interface {
	// Called before each RPC to authenticate it. The implementation may
	// return a new context if they want to insert authentication information
	// into it (such as the current user). The implementation may return a nil
	// context and the existing context will be used.
	Authenticate(
		ctx context.Context,
		token, endpoint string,
		effects []string,
	) (context.Context, error)
}

var readonly = []string{"readonly"}

// Information about the effects of endpoints that are authenticated. If an endpoint
// is not listed, the DefaultEffect value is used.
var Effects = map[string][]string{
	"ListBuilds": readonly,
}

var DefaultEffects = []string{"mutable"}

// authUnaryInterceptor returns a gRPC unary interceptor that inspects the metadata
// attached to the context. A token is extracted from that metadata and the given
// AuthChecker is invoked to guard calling the target handler. Effectively
// it implements authentication in front of any unary call.
func authUnaryInterceptor(checker AuthChecker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// Allow reflection API to be unauthenticated
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection.v1alpha.ServerReflection/") {
			return handler(ctx, req)
		}

		name := filepath.Base(info.FullMethod)

		effects, ok := Effects[name]
		if !ok {
			effects = DefaultEffects
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		var token string
		if authHeader, ok := md["authorization"]; ok {
			token = authHeader[0]
		}

		newCtx, err := checker.Authenticate(ctx, token, name, effects)
		if err != nil {
			return nil, err
		}

		// If we were given a new context, use that for future requests
		if newCtx != nil {
			ctx = newCtx
		}

		return handler(ctx, req)
	}
}

// authStreamInterceptor returns a gRPC unary interceptor that inspects the metadata
// attached to the context. A token is extract from that metadata and the given
// AuthChecker is invoked to guard calling the target handler. Effectively
// it implements authentication in front of any stream call.
func authStreamInterceptor(checker AuthChecker) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		// Allow reflection API to be unauthenticated
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection.v1alpha.ServerReflection/") {
			return handler(srv, ss)
		}

		name := filepath.Base(info.FullMethod)

		effects, ok := Effects[name]
		if !ok {
			effects = DefaultEffects
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
		}

		token := authHeader[0]

		newCtx, err := checker.Authenticate(ss.Context(), token, name, effects)
		if err != nil {
			return err
		}

		if newCtx != nil {
			ss = &ssContextOverride{ServerStream: ss, Ctx: newCtx}
		}

		// Invoke the handler.
		return handler(srv, ss)
	}
}

// ssContextOverride implements grpc.ServerStream but only overrides the
// returned context.
type ssContextOverride struct {
	grpc.ServerStream
	Ctx context.Context
}

func (ss *ssContextOverride) Context() context.Context {
	return ss.Ctx
}

var _ grpc.ServerStream = (*ssContextOverride)(nil)
