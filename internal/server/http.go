package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

type httpServer struct {
	opts   *options
	log    hclog.Logger
	server *http.Server
}

// newHttpServer initializes a new http server.
// Uses grpc-web to wrap an existing grpc server.
func newHttpServer(grpcServer *grpc.Server, opts *options) *httpServer {
	log := opts.Logger.Named("http")
	if opts.HTTPListener == nil {
		log.Info("HTTP listener not specified, HTTP API is disabled")
		return nil
	}

	// Wrap the grpc server so that it is grpc-web compatible
	grpcWrapped := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(string) bool { return true }),
		grpcweb.WithAllowNonRootResource(true),
	)

	uifs := http.FileServer(&assetfs.AssetFS{
		Asset:     gen.Asset,
		AssetDir:  gen.AssetDir,
		AssetInfo: gen.AssetInfo,
		Prefix:    "ui/dist",
		Fallback:  "index.html",
	})

	// If the path has a grpc prefix we assume it's a GRPC gateway request,
	// otherwise fall back to serving the UI from the filesystem
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/grpc") {
			grpcWrapped.ServeHTTP(w, r)
		} else if opts.BrowserUIEnabled {
			uifs.ServeHTTP(w, r)
		}
	})

	// Create our http server
	return &httpServer{
		opts: opts,
		log:  log,
		server: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       120 * time.Second,
			Handler:           httpLogHandler(rootHandler, log),
			BaseContext: func(net.Listener) context.Context {
				return opts.Context
			},
		},
	}
}

// start starts an http server
func (s *httpServer) start() error {
	// Serve traffic
	ln := s.opts.HTTPListener
	s.log.Info("starting HTTP server", "addr", ln.Addr().String())
	return s.server.Serve(ln)
}

// close stops the grpc server, gracefully if possible. Should be called exactly once.
// Warning: before closing the GRPC server, this HTTP server must first be closed.
// Attempting to gracefully stop the GRPC server first will cause it to drain HTTP connections,
// which will panic.
func (s *httpServer) close() {
	log := s.log
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Graceful in a goroutine so we can timeout
	gracefulCh := make(chan struct{})
	go func() {
		defer close(gracefulCh)
		log.Debug("stopping")
		if err := s.server.Shutdown(ctx); err != nil {
			log.Error("failed graceful shutdown: %s", err)
		}
	}()

	select {
	case <-gracefulCh:
		log.Debug("exited gracefully")

	// After a timeout we just forcibly exit. Our HTTP endpoints should
	// be fairly quick and their operations are atomic so we just kill
	// the connections after a few seconds.
	case <-time.After(2 * time.Second):
		log.Debug("stopping forcefully after waiting unsuccessfully for graceful stop")
		cancelFunc()
	}
}
