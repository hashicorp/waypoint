package server

import (
	"context"
	"net"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/pkg/protocolversion"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/tokenutil"
)

// TestServer starts a server and returns a gRPC client to that server.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T, impl pb.WaypointServer, opts ...TestOption) pb.WaypointClient {
	require := require.New(t)

	c := testConfig{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(&c)
	}

	// Listen on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)
	t.Cleanup(func() { ln.Close() })

	// We make run a function since we'll call it to restart too
	run := func(ctx context.Context) context.CancelFunc {
		ctx, cancel := context.WithCancel(ctx)

		server, err := newTestServer(ctx, impl)
		require.NoError(err)

		go func() {
			err := server.Serve(ln)
			require.NoError(err)
		}()

		go func() {
			<-ctx.Done()
			server.Stop()
		}()

		t.Cleanup(func() {
			cancel()
			server.Stop()
		})

		return cancel
	}

	// Create the server
	cancel := run(c.ctx)

	// If we have a restart channel, then listen to that for restarts.
	if c.restartCh != nil {
		doneCh := make(chan struct{})
		t.Cleanup(func() { close(doneCh) })

		go func() {
			for {
				select {
				case <-c.restartCh:
					// Cancel the old context
					cancel()

					// This can fail, but it probably won't. Can't think of
					// a cleaner way since gRPC force closes its listener.
					ln, err = net.Listen("tcp", ln.Addr().String())
					if err != nil {
						return
					}

					// Create a new one
					cancel = run(context.Background())

				case <-doneCh:
					return
				}
			}
		}()
	}

	// Get our version info we'll set on the client
	vsnInfo := TestVersionInfoResponse().Info

	// connect is a function since we need to connect multiple times:
	// once to bootstrap, then again with our auth information.
	connect := func(token string) (*grpc.ClientConn, error) {
		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(vsnInfo)),
			grpc.WithStreamInterceptor(protocolversion.StreamClientInterceptor(vsnInfo)),
			grpc.WithPerRPCCredentials(tokenutil.ContextToken(token)),
		}

		return grpc.DialContext(context.Background(), ln.Addr().String(), opts...)
	}

	// Connect, this should retry in the case Run is not going yet
	conn, err := connect("")
	require.NoError(err)
	client := pb.NewWaypointClient(conn)

	// Bootstrap
	tokenResp, err := client.BootstrapToken(context.Background(), &empty.Empty{})
	if status.Code(err) == codes.PermissionDenied {
		// Ignore bootstrap already complete errors
		err = nil
		tokenResp = &pb.NewTokenResponse{Token: ""}
	}
	conn.Close()
	require.NoError(err)

	// Reconnect with a token
	token := c.token
	if !c.tokenSet {
		token = tokenResp.Token
	}
	conn, err = connect(token)

	require.NoError(err)
	t.Cleanup(func() { conn.Close() })
	return pb.NewWaypointClient(conn)
}

func newTestServer(ctx context.Context, impl pb.WaypointServer) (*grpc.Server, error) {
	// Create and start a new GRPC server

	// Get our server info immediately
	resp, err := impl.GetVersionInfo(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var so []grpc.ServerOption
	so = append(so,
		grpc.ChainUnaryInterceptor(
			// Protocol version negotiation
			VersionUnaryInterceptor(resp.Info),
		),
		grpc.ChainStreamInterceptor(
			// Protocol version negotiation
			VersionStreamInterceptor(resp.Info),
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				// connections need to wait at least 20s before sending a
				// keepalive ping
				MinTime: 20 * time.Second,
				// allow runners to send keeplive pings even if there are no
				// active RCP streams.
				PermitWithoutStream: true,
			}),
	)

	if ac, ok := impl.(AuthChecker); ok {
		so = append(so,
			grpc.ChainUnaryInterceptor(AuthUnaryInterceptor(ac)),
			grpc.ChainStreamInterceptor(AuthStreamInterceptor(ac)),
		)
	}

	server := grpc.NewServer(so...)

	pb.RegisterWaypointServer(server, impl)
	return server, nil
}

// TestOption is used with TestServer to configure test behavior.
type TestOption func(*testConfig)

type testConfig struct {
	ctx       context.Context
	restartCh <-chan struct{}
	token     string
	tokenSet  bool
}

// TestWithContext specifies a context to use with the test server. When
// this is done then the server will exit.
func TestWithContext(ctx context.Context) TestOption {
	return func(c *testConfig) {
		c.ctx = ctx
	}
}

// TestWithRestart specifies a channel that will be sent to trigger
// a restart. The restart happens asynchronously. If you want to ensure the
// server is shutdown first, use TestWithContext, shut it down, wait for
// errors on the API, then restart.
func TestWithRestart(ch <-chan struct{}) TestOption {
	return func(c *testConfig) {
		c.restartCh = ch
	}
}

// TestWithToken specifies a specific token to use for auth.
func TestWithToken(token string) TestOption {
	return func(c *testConfig) {
		c.token = token
		c.tokenSet = true
	}
}

// TestVersionInfoResponse generates a valid version info response for testing
func TestVersionInfoResponse() *pb.GetVersionInfoResponse {
	return &pb.GetVersionInfoResponse{
		Info: &pb.VersionInfo{
			Api: &pb.VersionInfo_ProtocolVersion{
				Current: 10,
				Minimum: 1,
			},

			Entrypoint: &pb.VersionInfo_ProtocolVersion{
				Current: 10,
				Minimum: 1,
			},
		},
	}
}

// TestCookieContext returns a context with the cookie set.
func TestCookieContext(ctx context.Context, t testing.T, c pb.WaypointClient) context.Context {
	resp, err := c.GetServerConfig(ctx, &empty.Empty{})
	require.NoError(t, err)
	md := metadata.New(map[string]string{"wpcookie": resp.Config.Cookie})
	return metadata.NewOutgoingContext(ctx, md)
}

// TestRunner registers a runner and returns the ID and a function to
// deregister the runner. This uses t.Cleanup so that the runner will always
// be deregistered on test completion.
func TestRunner(t testing.T, client pb.WaypointClient, r *pb.Runner) (string, func()) {
	require := require.New(t)
	ctx := context.Background()

	// Get the runner
	if r == nil {
		r = &pb.Runner{}
	}
	id, err := Id()
	require.NoError(err)
	require.NoError(mergo.Merge(r, &pb.Runner{Id: id}))

	// Open the config stream
	stream, err := client.RunnerConfig(ctx)
	require.NoError(err)
	t.Cleanup(func() { stream.CloseSend() })

	// Register
	require.NoError(stream.Send(&pb.RunnerConfigRequest{
		Event: &pb.RunnerConfigRequest_Open_{
			Open: &pb.RunnerConfigRequest_Open{
				Runner: r,
			},
		},
	}))

	// Wait for first message to confirm we're registered
	_, err = stream.Recv()
	require.NoError(err)

	return id, func() { stream.CloseSend() }
}
