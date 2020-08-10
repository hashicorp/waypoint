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

	// If we have an invite token, we have to exchange that and re-establish
	// the connection with the auth setup. If we have no token, we're done.
	if cfg.InviteToken == "" {
		ceb.logger.Warn("no auth token given, will use unauthenticated connection")
		return nil
	}

	// Exchange
	ceb.logger.Info("converting invite token to login token")
	resp, err := ceb.client.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
		Token: cfg.InviteToken,
	})
	if err != nil {
		return err
	}

	// We have our token, setup that usage
	grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(staticToken(resp.Token)))

	// Reconnect and return
	ceb.logger.Info("reconnecting to server with authentication")
	conn, err = grpc.DialContext(ctx, cfg.ServerAddr, grpcOpts...)
	if err != nil {
		return err
	}
	ceb.client = pb.NewWaypointClient(conn)

	return nil
}

// This is a weird type that only exists to satisify the interface required by
// grpc.WithPerRPCCredentials. That api is designed to incorporate things like OAuth
// but in our case, we really just want to send this static token through, but we still
// need to the dance.
type staticToken string

func (t staticToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(t),
	}, nil
}

func (t staticToken) RequireTransportSecurity() bool {
	return false
}
