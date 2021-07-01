package server

import (
	"context"
	"net"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

//go:generate sh -c "protoc -I../../vendor/proto/api-common-protos -I ../.. ../../internal/server/proto/server.proto --go_out=plugins=grpc:../.. --go-json_out=../.."
//go:generate mv ./proto/server.pb.json.go ./gen
//go:generate mockery -all -case underscore -dir ./gen -output ./gen/mocks

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
		cfg.Logger = hclog.L()
	}

	grpcServer, err := newGrpcServer(&cfg)
	if err != nil {
		return err
	}

	errch := make(chan error)
	go func() {
		if err := grpcServer.start(); err != nil {
			errch <- err
		}
	}()

	httpServer := newHttpServer(grpcServer.server, &cfg)
	go func() {
		if err := httpServer.start(); err != nil {
			errch <- err
		}
	}()

	ctx, cancel := context.WithCancel(cfg.Context)
	defer cancel()

	select {
	case err := <-errch:
		return err
	case <-cfg.Context.Done():
		// Must shut down the http server first, as the grpc server can't drain http connections
		httpServer.close()
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
	HTTPListener net.Listener

	// AuthChecker, if set, activates authentication checking on the server.
	AuthChecker AuthChecker

	// BrowserUIEnabled determines if the browser UI should be mounted
	BrowserUIEnabled bool
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
func WithHTTP(ln net.Listener) Option {
	return func(opts *options) { opts.HTTPListener = ln }
}

// WithImpl sets the service implementation to serve.
func WithImpl(impl pb.WaypointServer) Option {
	return func(opts *options) { opts.Service = impl }
}

// WithAuthentication configures the server to require authentication.
func WithAuthentication(ac AuthChecker) Option {
	return func(opts *options) { opts.AuthChecker = ac }
}

// WithBrowserUI configures the server to enable the browser UI.
func WithBrowserUI(enabled bool) Option {
	return func(opts *options) { opts.BrowserUIEnabled = enabled }
}
