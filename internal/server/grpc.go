package server

import (
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type grpcServer struct {
	opts   *options
	log    hclog.Logger
	server *grpc.Server
}

// newGrpcServer initializes a new gRPC server
func newGrpcServer(opts *options) (*grpcServer, error) {
	log := opts.Logger.Named("grpc")

	// Get our server info immediately
	resp, err := opts.Service.GetVersionInfo(opts.Context, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var so []grpc.ServerOption
	so = append(so,
		grpc.ChainUnaryInterceptor(
			// Insert our logger and also log req/resp
			logUnaryInterceptor(log, false),

			// Protocol version negotiation
			versionUnaryInterceptor(resp.Info),
		),
		grpc.ChainStreamInterceptor(
			// Insert our logger and log
			logStreamInterceptor(log, false),

			// Protocol version negotiation
			versionStreamInterceptor(resp.Info),
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				// connections need to wait at least 20s before sending a
				// keepalive ping
				MinTime: 20 * time.Second,
				// allow runners to send keeplive pings even if there are no
				// active RCP streams.
				PermitWithoutStream: true,
			}),
	)

	if opts.AuthChecker != nil {
		so = append(so,
			grpc.ChainUnaryInterceptor(authUnaryInterceptor(opts.AuthChecker)),
			grpc.ChainStreamInterceptor(authStreamInterceptor(opts.AuthChecker)),
		)
	}

	s := grpc.NewServer(so...)

	return &grpcServer{
		opts:   opts,
		server: s,
		log:    log,
	}, nil
}

// start starts the grpc server
func (s *grpcServer) start() error {
	// Register the reflection service. This makes using tools like grpcurl
	// easier. It makes it slightly easier for malicious users to know about
	// the service but I think they'd figure out its a waypoint server
	// easy enough.
	reflection.Register(s.server)

	// Register our server
	pb.RegisterWaypointServer(s.server, s.opts.Service)
	// Serve traffic
	ln := s.opts.GRPCListener
	s.log.Info("starting gRPC server", "addr", ln.Addr().String())
	return s.server.Serve(ln)
}

// close stops the grpc server, gracefully if possible.
// Warning: before closing the GRPC server, the HTTP server must first be closed.
// Attempting to gracefully stop the GRPC server first will cause it to drain HTTP connections,
// which will panic.
func (s *grpcServer) close() {
	log := s.log
	// Graceful in a goroutine so we can timeout
	gracefulCh := make(chan struct{})
	go func() {
		defer close(gracefulCh)
		log.Debug("stopping")
		s.server.GracefulStop()
	}()

	select {
	case <-gracefulCh:
		log.Debug("exited gracefully")

	// After a timeout we just forcibly exit. Our gRPC endpoints should
	// be fairly quick and their operations are atomic so we just kill
	// the connections after a few seconds.
	case <-time.After(2 * time.Second):
		log.Debug("stopping forcefully after waiting unsuccessfully for graceful stop")
		s.server.Stop()
	}
}
