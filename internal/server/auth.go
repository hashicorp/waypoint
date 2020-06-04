package server

import (
	"context"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthChecker interface {
	Authenticate(ctx context.Context, token, endpoint string, effects []string) error
	DefaultToken() (string, error)
}

var readonly = []string{"readonly"}

var Effects = map[string][]string{
	"ListBuilds": readonly,
}

var defaultEffects = []string{"mutable"}

// logUnaryInterceptor returns a gRPC unary interceptor that inserts a hclog.Logger
// into the request context.
//
// Additionally, logUnaryInterceptor logs request and response metadata. If verbose
// is set to true, the request and response attributes are logged too.
func authUnaryInterceptor(checker AuthChecker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		name := filepath.Base(info.FullMethod)

		effects, ok := Effects[name]
		if !ok {
			effects = defaultEffects
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
		}

		token := authHeader[0]

		err := checker.Authenticate(ctx, token, name, effects)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// logUnaryInterceptor returns a gRPC unary interceptor that inserts a hclog.Logger
// into the request context.
//
// Additionally, logUnaryInterceptor logs request and response metadata. If verbose
// is set to true, the request and response attributes are logged too.
func authStreamInterceptor(checker AuthChecker) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		name := filepath.Base(info.FullMethod)

		effects, ok := Effects[name]
		if !ok {
			effects = defaultEffects
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

		err := checker.Authenticate(ss.Context(), token, name, effects)
		if err != nil {
			return err
		}

		// Invoke the handler.
		return handler(srv, ss)
	}
}
