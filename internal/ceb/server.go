package ceb

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/protocolversion"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// client returns the Waypoint client or blocks until it is set or the
// ceb is exiting. Once this returns, users should ALWAYS check if an exit
// condition was triggered to avoid nil panics.
func (ceb *CEB) waitClient() pb.WaypointClient {
	ceb.clientMu.Lock()
	defer ceb.clientMu.Unlock()

	for ceb.client == nil {
		ceb.clientCond.Wait()
	}

	return ceb.client
}

// initClient initializes the client connection to the server. This will
// attempt to synchronously connect once, and then reattempt connection in
// the background.
//
// Users of the client should use the waitClient() function to wait
// for the client to be set.
func (ceb *CEB) initClient(ctx context.Context, log hclog.Logger, cfg *config, retry bool) error {
	if ceb.client != nil {
		return nil
	}

	if cfg.ServerAddr == "" {
		log.Info("no waypoint server configured, disabled entrypoint")
		return nil
	}

RETRY_INIT:
	err := ceb.dialServer(ctx, cfg, retry)
	if status.Code(err) == codes.Unavailable {
		// If we require a server connection, then just retry.
		if cfg.ServerRequired {
			log.Warn("server unavailable but ceb configured to require it, retrying synchronously")
			retry = true
			goto RETRY_INIT
		}

		// If we don't require a server connection, then we start a
		// goroutine to retry and eventually connect (hopefully).
		log.Warn("server unavailable, will retry in the background")
		go ceb.initClient(ctx, log, cfg, true)

		// We also mark that we can begin executing the child command.
		// We usually don't do this because we wait for initial config, but
		// if we fail to connect to the client, we can just start it.
		ceb.markChildCmdReady()

		return nil
	}

	return err
}

// dialServer connects to the server.
func (ceb *CEB) dialServer(ctx context.Context, cfg *config, isRetry bool) error {
	// Build our options
	grpcOpts := []grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
		grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(protocolversion.Current())),
		grpc.WithStreamInterceptor(protocolversion.StreamClientInterceptor(protocolversion.Current())),
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
	ceb.logger.Debug("connecting to server",
		"addr", cfg.ServerAddr,
		"tls", cfg.ServerTls,
		"tls_skip_verify", cfg.ServerTlsSkipVerify,
	)
	conn, err := grpc.DialContext(ctx, cfg.ServerAddr, grpcOpts...)
	if err != nil {
		ceb.logger.Warn("failed to connect to server", "err", err)
		return err
	}

	// When we commit to keeping conn, we should set it to nil.
	defer func() {
		if conn != nil {
			conn.Close()
			ceb.client = nil
		}
	}()

	// Verify we're connected. We do this loop so that we can
	// fail fast in the non-isRetry case.
	ceb.logger.Debug("waiting on server connection state to become ready")
	for {
		s := conn.GetState()
		ceb.logger.Trace("connection state", "state", s.String())

		// If we're ready then we're done!
		if s == connectivity.Ready {
			ceb.logger.Debug("connection is ready")
			break
		}

		// If we have a transient error and we're not retrying, then we're done.
		if s == connectivity.TransientFailure && !isRetry {
			ceb.logger.Warn("failed to connect to the server, temporary network error")
			conn.Close()
			return status.Errorf(codes.Unavailable, "server is unavailable")
		}

		if !conn.WaitForStateChange(ctx, s) {
			return ctx.Err()
		}
	}

	// Init our client
	client := pb.NewWaypointClient(conn)

	// If we have an invite token, we have to exchange that and re-establish
	// the connection with the auth setup. If we have no token, we're done.
	if cfg.InviteToken != "" {
		// Exchange
		ceb.logger.Debug("converting invite token to login token")
		resp, err := client.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
			Token: cfg.InviteToken,
		}, grpc.WaitForReady(isRetry))
		if err != nil {
			return err
		}

		// We have our token, setup that usage
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(staticToken(resp.Token)))

		// Reconnect and return
		conn.Close()
		ceb.logger.Debug("reconnecting to server with authentication")
		conn, err = grpc.DialContext(ctx, cfg.ServerAddr, grpcOpts...)
		if err != nil {
			return err
		}
		client = pb.NewWaypointClient(conn)
	}

	// Negotiate API version
	ceb.logger.Trace("requesting version info from server")
	vsnResp, err := client.GetVersionInfo(ctx, &empty.Empty{}, grpc.WaitForReady(isRetry))
	if err != nil {
		return err
	}

	ceb.logger.Info("server version info",
		"version", vsnResp.Info.Version,
		"api_min", vsnResp.Info.Api.Minimum,
		"api_current", vsnResp.Info.Api.Current,
		"entrypoint_min", vsnResp.Info.Entrypoint.Minimum,
		"entrypoint_current", vsnResp.Info.Entrypoint.Current,
	)

	vsn, err := protocolversion.Negotiate(protocolversion.Current().Entrypoint, vsnResp.Info.Entrypoint)
	if err != nil {
		return err
	}
	ceb.logger.Debug("negotiated entrypoint protocol version", "version", vsn)

	// Commit to using this client
	ceb.clientMu.Lock()
	defer ceb.clientMu.Unlock()
	ceb.client = client
	ceb.clientCond.Broadcast()
	connCopy := conn
	conn = nil
	ceb.cleanup(func() { connCopy.Close() })

	return nil
}

// This is a weird type that only exists to satisify the interface required by
// grpc.WithPerRPCCredentials. That api is designed to incorporate things like OAuth
// but in our case, we really just want to send this static token through, but we still
// need to do the dance.
type staticToken string

func (t staticToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(t),
	}, nil
}

func (t staticToken) RequireTransportSecurity() bool {
	return false
}
