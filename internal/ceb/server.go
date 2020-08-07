package ceb

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// dialServer connects to the server.
func (ceb *CEB) dialServer(ctx context.Context, cfg *config) error {
	// Build our options
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(5 * time.Second),
	}
	if !cfg.ServerTls {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else {
		if cfg.ServerTlsSkipVerify {
			grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
			))
		}
	}

	// Connect to this server
	ceb.logger.Info("connecting to server",
		"addr", cfg.ServerAddr,
		"tls", cfg.ServerTls,
		"tls_skip_verify", cfg.ServerTlsSkipVerify,
	)
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
