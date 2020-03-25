package core

import (
	"context"
	"net"
	"path/filepath"

	"github.com/boltdb/bolt"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/internal/server/singleprocess"
)

// initLocalServer starts the local server and configures p.client to
// point to it. This also configures p.localClosers so that all the
// resources are properly cleaned up on Close.
//
// If this returns an error, all resources associated with this operation
// will be closed, but the project can retry.
func (p *Project) initLocalServer(ctx context.Context) error {
	log := p.logger.Named("server")

	// We use this pointer to accumulate things we need to clean up
	// in the case of an error. On success we nil this variable which
	// doesn't close anything.
	closers := &p.localClosers
	defer func() {
		if closers != nil {
			for _, c := range *closers {
				c.Close()
			}
			*closers = nil
		}
	}()

	path := filepath.Join(p.dir.DataDir(), "data.db")
	log.Debug("opening local mode DB", "path", path)

	// Open our database
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	p.localClosers = append(p.localClosers, db)

	// Create our server
	impl, err := singleprocess.New(db)
	if err != nil {
		return err
	}

	// We listen on a random locally bound port
	// TODO(mitchellh): we should use Unix domain sockets if supported
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		db.Close()
		return err
	}
	p.localClosers = append(p.localClosers, ln)

	// Create a new cancellation context so we can cancel in the case of an error
	ctx, cancel := context.WithCancel(ctx)

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
	)
	if err != nil {
		cancel()
		ln.Close()
		db.Close()
		return err
	}
	p.localClosers = append(p.localClosers, conn)

	// Init our client
	p.client = pb.NewDevflowClient(conn)

	// Success, nil our closers so we don't defer close them
	closers = nil

	return nil
}
