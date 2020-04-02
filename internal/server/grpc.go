package server

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oklog/run"
	"google.golang.org/grpc"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// grpcInit initializes the gRPC server and adds it to the run group.
func grpcInit(group *run.Group, opts *options) error {
	log := opts.Logger.Named("grpc")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// Insert our logger and also log req/resp
			logInterceptor(log, false),
		)),
	)

	// Register our server
	pb.RegisterDevflowServer(s, opts.Service)

	// Create a cancellation context we'll use to stop our gRPC server
	ctx, cancel := context.WithCancel(opts.Context)

	// Add our gRPC server to the run group
	group.Add(func() error {
		// Start a goroutine that waits for cancellation
		go func() {
			<-ctx.Done()
			log.Info("shutting down gRPC server")
			s.GracefulStop()
		}()

		// Serve traffic
		ln := opts.GRPCListener
		log.Info("starting gRPC server", "addr", ln.Addr().String())
		return s.Serve(ln)
	}, func(err error) {
		cancel()
	})

	return nil
}
