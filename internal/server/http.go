package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/go-hclog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/hashicorp/waypoint/internal/server/httpapi"
	"github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
)

type httpServer struct {
	ln     net.Listener
	opts   *options
	log    hclog.Logger
	server *http.Server
}

// newHttpServer initializes a new http server.
// Uses grpc-web to wrap an existing grpc server.
func newHttpServer(grpcServer *grpc.Server, ln net.Listener, opts *options) *httpServer {
	log := opts.Logger.Named("http").With("ln", ln.Addr().String())

	// Wrap the grpc server so that it is grpc-web compatible
	grpcWrapped := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(string) bool { return true }),
		grpcweb.WithAllowNonRootResource(true),
	)

	// This is the http.Handler for the UI
	uifs := http.FileServer(&assetfs.AssetFS{
		Asset:     gen.Asset,
		AssetDir:  gen.AssetDir,
		AssetInfo: gen.AssetInfo,
		Prefix:    "ui/dist",
		Fallback:  "index.html",
	})

	// grpcAddr is the address that we can connect back to our own
	// gRPC server. This is used by the exec handler.
	grpcAddr := opts.GRPCListener.Addr().String()

	// Create grpc-gateway muxer
	grpcHandler := runtime.NewServeMux()

	grpcOpts := serverclient.BuildDialOptions()

	grpcOpts = append(grpcOpts,
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	)

	err := gen.RegisterWaypointHandlerFromEndpoint(opts.Context, grpcHandler, grpcAddr, grpcOpts)
	if err != nil {
		log.Error("Unable to register waypoint grpc gateway service")
	}

	// Create our full router
	r := mux.NewRouter()
	r.HandleFunc("/v1/exec", httpapi.HandleExec(grpcAddr, true))
	r.HandleFunc("/v1/trigger/{id:[a-zA-Z0-9]+}", httpapi.HandleTrigger(grpcAddr, true))
	r.PathPrefix("/grpc").Handler(grpcWrapped)
	r.PathPrefix("/v1").Handler(grpcHandler)
	r.PathPrefix("/").Handler(uifs)

	// Create our root handler which is just our router. We then wrap it
	// in various middlewares below.
	var rootHandler http.Handler = r

	// Wrap our handler to force TLS
	rootHandler = forceTLSHandler(rootHandler)

	// Wrap our handler to log
	rootHandler = httpLogHandler(rootHandler, log)

	// Create our http server
	return &httpServer{
		ln:   ln,
		opts: opts,
		log:  log,
		server: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       120 * time.Second,
			Handler:           rootHandler,
			BaseContext: func(net.Listener) context.Context {
				return opts.Context
			},
		},
	}
}

// start starts an http server
func (s *httpServer) start() error {
	// Serve traffic
	s.log.Info("starting HTTP server", "addr", s.ln.Addr().String())
	return s.server.Serve(s.ln)
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
