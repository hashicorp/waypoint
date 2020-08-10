package server

import (
	"time"

	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// grpcInit initializes the gRPC server and adds it to the run group.
func grpcInit(group *run.Group, opts *options) error {
	log := opts.Logger.Named("grpc")

	var so []grpc.ServerOption

	if opts.AuthChecker != nil {
		so = append(so,
			grpc.ChainUnaryInterceptor(authUnaryInterceptor(opts.AuthChecker)),
			grpc.ChainStreamInterceptor(authStreamInterceptor(opts.AuthChecker)),
		)
	}

	so = append(so,
		grpc.ChainUnaryInterceptor(
			// Insert our logger and also log req/resp
			logUnaryInterceptor(log, false),
		),
		grpc.ChainStreamInterceptor(
			// Insert our logger and log
			logStreamInterceptor(log, false),
		),
	)

	s := grpc.NewServer(so...)
	opts.grpcServer = s

	// Register the reflection service. This makes using tools like grpcurl
	// easier. It makes it slightly easier for malicious users to know about
	// the service but I think they'd figure out its a waypoint server
	// easy enough.
	reflection.Register(s)

	// Register our server
	pb.RegisterWaypointServer(s, opts.Service)

	// Add our gRPC server to the run group
	group.Add(func() error {
		// Serve traffic
		ln := opts.GRPCListener
		log.Info("starting gRPC server", "addr", ln.Addr().String())
		return s.Serve(ln)
	}, func(err error) {
		// Graceful in a goroutine so we can timeout
		gracefulCh := make(chan struct{})
		go func() {
			defer close(gracefulCh)
			log.Info("shutting down gRPC server")
			s.GracefulStop()
		}()

		select {
		case <-gracefulCh:

		// After a timeout we just forcibly exit. Our gRPC endpoints should
		// be fairly quick and their operations are atomic so we just kill
		// the connections after a few seconds.
		case <-time.After(2 * time.Second):
			s.Stop()
		}
	})

	return nil
}
