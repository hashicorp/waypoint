package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"time"

	"github.com/hashicorp/waypoint/internal/protocolversion"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"

	bolt "go.etcd.io/bbolt"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

type LocalServer struct {
	ctx     context.Context
	closers []io.Closer
	log     hclog.Logger
}

func NewLocalServer(ctx context.Context, log hclog.Logger) *LocalServer {
	return &LocalServer{
		ctx: ctx,
		log: log.Named("server"),
	}
}

// Start starts the local server and returns a connection to it.
//
// If this returns an error, all resources associated with this operation
// will be closed, but the project can retry.
func (s *LocalServer) Start() (*grpc.ClientConn, error) {
	log := s.log

	// We use this pointer to accumulate things we need to clean up
	// in the case of an error. On success, we nil this variable which
	// doesn't close anything.
	var closers []io.Closer
	defer func() {
		for _, c := range closers {
			c.Close()
		}
	}()

	// TODO(mitchellh): path to this
	path := filepath.Join("data.db")
	log.Debug("opening local mode DB", "path", path)

	// Open our database
	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed opening boltdb path %s. Is another server already running against this db?: %w", path, err)
	}
	closers = append(closers, db)

	// Create our server
	impl, err := singleprocess.New(
		singleprocess.WithDB(db),
		singleprocess.WithLogger(log.Named("singleprocess")),
		singleprocess.WithSuperuser(), // local mode has no auth
	)
	if err != nil {
		return nil, err
	}

	// We listen on a random locally bound port
	// TODO(mitchellh): we should use Unix domain sockets if supported
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, err
	}
	closers = append(closers, ln)

	// Create a new cancellation context so we can cancel in the case of an error
	ctx, cancel := context.WithCancel(s.ctx)

	// Run the server
	log.Info("starting built-in server for local operations", "addr", ln.Addr().String())
	go server.Run(server.WithContext(ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithImpl(impl),
	)

	// Connect to the local server
	conn, err := grpc.DialContext(ctx, ln.Addr().String(),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(protocolversion.Current())),
		grpc.WithStreamInterceptor(protocolversion.StreamClientInterceptor(protocolversion.Current())),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	// Setup our server config. The configuration is specifically set so
	// that there is no advertise address which will disable the CEB
	// completely.
	client := pb.NewWaypointClient(conn)
	_, err = client.SetServerConfig(ctx, &pb.SetServerConfigRequest{
		Config: &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
				{
					Addr: "",
				},
			},
		},
	})
	if err != nil {
		cancel()
		return nil, err
	}

	// Success, persist the closers
	s.closers = closers
	closers = nil
	_ = cancel // pacify vet lostcancel

	return conn, nil
}

func (s *LocalServer) Close() error {
	for _, c := range s.closers {
		c.Close()
	}

	// TODO(izaak): do we need to also close the cancel context on the local server? We weren't before...
	return nil
}
