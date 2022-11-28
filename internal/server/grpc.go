package server

import (
	"time"

	"github.com/hashicorp/go-hclog"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/pkg/inlinekeepalive"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
			server.VersionUnaryInterceptor(resp.Info),

			// Nil protobuf "any" fields for gRPC-gateway since the JSON
			// encoding tries to decode the any.
			server.GWNullAnyUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			// Insert our logger and log
			logStreamInterceptor(log, false),

			// Send and receive keepalive messages along grpc streams.
			// Some loadbalancers (ALBs) don't respect http2 pings.
			// (https://stackoverflow.com/questions/66818645/http2-ping-frames-over-aws-alb-grpc-keepalive-ping)
			// This interceptor keeps low-traffic streams active and not timed out.
			// NOTE(izaak): long-term, we should ensure that all of our
			// streaming endpoints are robust to disconnect/resume.
			inlinekeepalive.KeepaliveServerStreamInterceptor(time.Duration(5)*time.Second),

			// Protocol version negotiation
			server.VersionStreamInterceptor(resp.Info),
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
			grpc.ChainUnaryInterceptor(server.AuthUnaryInterceptor(opts.AuthChecker)),
			grpc.ChainStreamInterceptor(server.AuthStreamInterceptor(opts.AuthChecker)),
		)
	}

	// This is the only place we wire telemetry into our grpc server.
	if opts.TelemetryEnabled {
		log.Debug("Enabling server ocgrpc stats handler")
		so = append(so, grpc.StatsHandler(&ocgrpc.ServerHandler{}))
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
	// the service but I think they'd figure out it's a waypoint server
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
