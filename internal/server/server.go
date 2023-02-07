package server

import (
	"context"
	"net"
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Run initializes and starts the server. This will block until the server
// exits (by cancelling the associated context set with WithContext or due
// to an unrecoverable error).
func Run(opts ...Option) error {
	var cfg options
	for _, opt := range opts {
		opt(&cfg)
	}

	// Set defaults
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}
	if cfg.Logger == nil {
		cfg.Logger = hclog.L().Named("server")
	}
	log := cfg.Logger

	grpcServer, err := newGrpcServer(&cfg)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	gprcErrs := make(chan error, 1)
	wg.Add(1)
	go func() {
		err := grpcServer.start()
		gprcErrs <- err
		log.Debug("gRPC server has exited", "err", err)
		wg.Done()
	}()

	var httpServers []*httpServer
	httpErrs := make(chan error, len(cfg.HTTPListener))
	for _, ln := range cfg.HTTPListener {
		wg.Add(1)
		httpServer, err := newHttpServer(grpcServer.server, ln, &cfg)
		if err != nil {
			log.Error("failed setting up http server", "err", err)
			return err
		}
		go func() {
			err := httpServer.start()
			httpErrs <- err
			log.Debug("HTTP server has exited", "err", err)
			wg.Done()
		}()

		httpServers = append(httpServers, httpServer)
	}

	ctx, cancel := context.WithCancel(cfg.Context)
	defer cancel()

	defer wg.Wait() // Wait for server run processes to exit before returning
	select {
	case err := <-gprcErrs:
		// If the GRPC server errored, we can assume it's closed and shut down the http server if necessary
		for _, httpServer := range httpServers {
			httpServer.close()
		}
		log.Error("failed running the grpc server", "err", err)
		return err
	case err := <-httpErrs:
		// If the HTTP server errored, we can assume it's closed and shut down the grpc server
		grpcServer.close()
		log.Error("failed running the http server", "err", err)
		return err
	case <-cfg.Context.Done():
		// Received an external shutdown signal, and should close everything.
		// NOTE: must close HTTP server before GRPC server.
		for _, httpServer := range httpServers {
			httpServer.close()
		}
		grpcServer.close()
		return ctx.Err()
	}
}

// Option configures Run
type Option func(*options)

// options configure a server and are set by users only using the exported
// Option functions.
type options struct {
	// Context is the context to use for the server. When this is cancelled,
	// the server will be gracefully shutdown.
	Context context.Context

	// Logger is the logger to use. This will default to hclog.L() if not set.
	Logger hclog.Logger

	// Service is the backend service implementation to use for the server.
	Service pb.WaypointServer

	// GRPCListener will setup the gRPC server. If this is nil, then a
	// random loopback port will be chosen. The gRPC server must run since it
	// serves the HTTP endpoints as well.
	GRPCListener net.Listener

	// HTTPListener will setup the HTTP server. If this is nil, then
	// the HTTP-based API will be disabled.
	HTTPListener []net.Listener

	// AuthChecker, if set, activates authentication checking on the server.
	AuthChecker server.AuthChecker

	// BrowserUIEnabled determines if the browser UI should be mounted
	BrowserUIEnabled bool

	// TelemetryEnabled determines if the server should instrument itself to emit telemetry. Default false
	TelemetryEnabled bool
}

// WithContext sets the context for the server. When this context is cancelled,
// the server will be shut down.
func WithContext(ctx context.Context) Option {
	return func(opts *options) { opts.Context = ctx }
}

// WithLogger sets the logger.
func WithLogger(log hclog.Logger) Option {
	return func(opts *options) { opts.Logger = log }
}

// WithGRPC sets the GRPC listener. This listener must be closed manually
// by the caller. Prior to closing the listener, it is recommended that you
// cancel the context set with WithContext and wait for Run to return.
func WithGRPC(ln net.Listener) Option {
	return func(opts *options) { opts.GRPCListener = ln }
}

// WithHTTP sets the HTTP listener. This listener must be closed manually
// by the caller. Prior to closing the listener, it is recommended that you
// cancel the context set with WithContext and wait for Run to return.
//
// If this is called multiple times, multiple HTTP listeners are started.
func WithHTTP(ln net.Listener) Option {
	return func(opts *options) { opts.HTTPListener = append(opts.HTTPListener, ln) }
}

// WithImpl sets the service implementation to serve.
func WithImpl(impl pb.WaypointServer) Option {
	return func(opts *options) { opts.Service = impl }
}

// WithAuthentication configures the server to require authentication.
func WithAuthentication(ac server.AuthChecker) Option {
	return func(opts *options) { opts.AuthChecker = ac }
}

// WithBrowserUI configures the server to enable the browser UI.
func WithBrowserUI(enabled bool) Option {
	return func(opts *options) { opts.BrowserUIEnabled = enabled }
}

// WithTelemetry instructs the server to export OpenCensus traces for requests.
func WithTelemetry(enabled bool) Option {
	return func(opts *options) { opts.TelemetryEnabled = enabled }
}
