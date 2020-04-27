package ceb

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// dialServer connects to the server.
func (ceb *CEB) dialServer(ctx context.Context, cfg *config) error {
	// Build our options
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(5 * time.Second),
	}
	if cfg.ServerInsecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	// Connect to this server
	ceb.logger.Info("connecting to server", "addr", cfg.ServerAddr)
	conn, err := grpc.DialContext(ctx, cfg.ServerAddr, grpcOpts...)
	if err != nil {
		return err
	}
	ceb.logger.Trace("server connection successful")
	ceb.cleanup(func() { conn.Close() })

	// Init our client
	ceb.client = pb.NewWaypointClient(conn)
	return nil
}
