package server

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

// logInterceptor returns a gRPC unary interceptor that inserts a hclog.Logger
// into the request context.
//
// Additionally, logInterceptor logs request and response metadata. If verbose
// is set to true, the request and response attributes are logged too.
func logInterceptor(logger hclog.Logger, verbose bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Log the request.
		{
			var reqLogArgs []interface{}
			// Log the request's attributes only if verbose is set to true.
			if verbose {
				reqLogArgs = append(reqLogArgs, "request", req)
			}
			logger.Info(info.FullMethod+" request", reqLogArgs...)
		}

		// Invoke the handler.
		ctx = hclog.WithContext(ctx, logger)
		resp, err := handler(ctx, req)

		// Log the response.
		{
			respLogArgs := []interface{}{
				"error", err,
				"duration", time.Since(start).String(),
			}
			// Log the response's attributes only if verbose is set to true.
			if verbose {
				respLogArgs = append(respLogArgs, "response", resp)
			}
			logger.Info(info.FullMethod+" response", respLogArgs...)
		}

		return resp, err
	}
}
