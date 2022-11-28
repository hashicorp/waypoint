package server

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hashicorp/waypoint/pkg/tokenutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthChecker - An interface implemented by something that wishes to authenticate the server
// actions.
type AuthChecker interface {
	// Authenticate is called before each RPC to authenticate it. The implementation may
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

// Effects is information about the effects of endpoints that are authenticated. If an endpoint
// is not listed, the DefaultEffect value is used.
var Effects = map[string][]string{
	"ListBuilds": readonly,
}

var DefaultEffects = []string{"mutable"}

// AuthUnaryInterceptor returns a gRPC unary interceptor that inspects the metadata
// attached to the context. A token is extracted from that metadata and the given
// AuthChecker is invoked to guard calling the target handler. Effectively
// it implements authentication in front of any unary call.
func AuthUnaryInterceptor(checker AuthChecker) grpc.UnaryServerInterceptor {
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
		// First, let's see if the token is in our dedicated key.
		if tokenHeader, ok := md[tokenutil.MetadataKey]; ok {
			token = tokenHeader[0]
		} else {
			// Otherwise, see if it's been stuffed into Authorization.
			if authHeader, ok := md["authorization"]; ok {
				token = authHeader[0]
			}
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

// AuthStreamInterceptor returns a gRPC unary interceptor that inspects the metadata
// attached to the context. A token is extract from that metadata and the given
// AuthChecker is invoked to guard calling the target handler. Effectively
// it implements authentication in front of any stream call.
func AuthStreamInterceptor(checker AuthChecker) grpc.StreamServerInterceptor {
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

		var token string

		// First, let's see if the token is in our dedicated key.
		if tokenHeader, ok := md[tokenutil.MetadataKey]; ok {
			token = tokenHeader[0]
		} else {
			// Otherwise, see if it's been stuffed into Authorization.
			if authHeader, ok := md["authorization"]; ok {
				token = authHeader[0]
			}
		}

		if token == "" {
			return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
		}

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
