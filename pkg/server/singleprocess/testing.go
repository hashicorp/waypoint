// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	hznhub "github.com/hashicorp/horizon/pkg/hub"
	hznpb "github.com/hashicorp/horizon/pkg/pb"
	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"
	hzntoken "github.com/hashicorp/horizon/pkg/token"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	wphzn "github.com/hashicorp/waypoint-hzn/pkg/server"
	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	"github.com/hashicorp/waypoint/pkg/serverstate"

	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

// TestServer starts a singleprocess server and returns the connected client.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T, opts ...Option) pb.WaypointClient {
	return server.TestServer(t, TestImpl(t, opts...))
}

// TestImpl returns the waypoint server implementation. This can be used
// with server.TestServer. It is easier to just use TestServer directly.
func TestImpl(t testing.T, opts ...Option) pb.WaypointServer {
	impl, err := New(append(
		[]Option{WithDB(testDB(t))},
		opts...,
	)...)
	require.NoError(t, err)
	if c, ok := impl.(io.Closer); ok {
		t.Cleanup(func() { c.Close() })
	}
	return impl
}

// TestWithURLService is an Option for testing only that creates an
// in-memory URL service server. This requires access to an external
// postgres server.
//
// If out is non-nil, it will be written to with the DevSetup info.
func TestWithURLService(t testing.T, out *hzntest.DevSetup) Option {
	// Create the test server. On test end we close the channel which quits
	// the Horizon test server.
	setupCh := make(chan *hzntest.DevSetup, 1)
	closeCh := make(chan struct{})
	t.Cleanup(func() { close(closeCh) })
	go hzntest.Dev(t, func(setup *hzntest.DevSetup) {
		hubclient, err := hznhub.NewHub(hclog.L(), setup.ControlClient, setup.HubToken)
		require.NoError(t, err)
		go hubclient.Run(context.Background(), setup.ClientListener)

		setupCh <- setup
		<-closeCh
	})
	setup := <-setupCh

	// Make our test registration API
	wphzndata := wphzn.TestServer(t,
		wphzn.WithNamespace("/"),
		wphzn.WithHznControl(setup.MgmtClient),
	)

	// Get our account token.
	wpaccountResp, err := wphzndata.Client.RegisterGuestAccount(
		context.Background(),
		&wphznpb.RegisterGuestAccountRequest{
			ServerId: "A",
		},
	)
	require.NoError(t, err)

	// We need to get the account since that is what we need to query with
	tokenInfo, err := setup.MgmtClient.GetTokenPublicKey(context.Background(), &hznpb.Noop{})
	require.NoError(t, err)
	token, err := hzntoken.CheckTokenED25519(wpaccountResp.Token, tokenInfo.PublicKey)
	require.NoError(t, err)
	setup.Account = token.Account()

	// Copy our setup config
	if out != nil {
		*out = *setup
	}

	return func(s *Service, cfg *config) error {
		if cfg.serverConfig == nil {
			cfg.serverConfig = &serverconfig.Config{}
		}

		cfg.serverConfig.URL = &serverconfig.URL{
			Enabled:              true,
			APIAddress:           wphzndata.Addr,
			APIInsecure:          true,
			APIToken:             wpaccountResp.Token,
			ControlAddress:       fmt.Sprintf("dev://%s", setup.HubAddr),
			AutomaticAppHostname: true,
		}

		return nil
	}
}

// TestWithURLServiceGuestAccount sets the API token to empty to force
// getting a guest account with the URL service. This can ONLY be set if
// TestWithURLService is set before this.
func TestWithURLServiceGuestAccount(t testing.T) Option {
	return func(s *Service, cfg *config) error {
		// Set the API token to empty which will force a guest account registration.
		cfg.serverConfig.URL.APIToken = ""

		return nil
	}
}

func TestEntrypoint(t testing.T, client pb.WaypointClient) (string, string, func()) {
	instanceId, err := server.Id()
	require.NoError(t, err)

	ctx := context.Background()

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},
		}),
	})
	require.NoError(t, err)

	dep := resp.Deployment

	// Create the config
	stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		InstanceId:   instanceId,
		DeploymentId: dep.Id,
	})
	require.NoError(t, err)

	// Wait for the first config so that we know we're registered
	_, err = stream.Recv()
	require.NoError(t, err)

	return instanceId, dep.Id, func() {
		stream.CloseSend()
	}
}

func TestEntrypointPlugin(t testing.T, client pb.WaypointClient) (string, string, func()) {
	instanceId, err := server.Id()
	require.NoError(t, err)

	ctx := context.Background()

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},

			HasExecPlugin: true,
		}),
	})
	require.NoError(t, err)

	dep := resp.Deployment

	// Create the config
	stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		InstanceId:   instanceId,
		DeploymentId: dep.Id,
	})
	require.NoError(t, err)

	// Wait for the first config so that we know we're registered
	_, err = stream.Recv()
	require.NoError(t, err)

	return instanceId, dep.Id, func() {
		stream.CloseSend()
	}
}

// TestRunnerAdopted registers a runner using the adoption process and
// returns a client that uses the token for this specific runner. This
// uses t.Cleanup so the runner will be cleaned up on completion.
func TestRunnerAdopted(
	t testing.T,
	impl pb.WaypointServer,
	client pb.WaypointClient,
	r *pb.Runner,
) (string, pb.WaypointClient) {
	require := require.New(t)
	ctx := context.Background()

	// Get the runner
	if r == nil {
		r = &pb.Runner{}
	}
	id, err := server.Id()
	require.NoError(err)
	require.NoError(mergo.Merge(r, &pb.Runner{Id: id}))

	// Get our cookied context
	ctx = server.TestCookieContext(ctx, t, client)

	// Reconnect with no token
	anonClient := server.TestServer(t, impl, server.TestWithToken(""))

	// Start getting the resp
	var resp *pb.RunnerTokenResponse
	var respErr error
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		resp, respErr = anonClient.RunnerToken(ctx, &pb.RunnerTokenRequest{
			Runner: r,
		})
	}()

	// Should block
	select {
	case <-time.After(50 * time.Millisecond):
	case <-doneCh:
		t.Fatal("should block")
	}

	// Adopt it
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
		RunnerId: id,
		Adopt:    true,
	})
	require.NoError(err)

	// Should be done
	select {
	case <-doneCh:
	case <-time.After(1 * time.Second):
		t.Fatal("should return")
	}

	// Verify token resp
	require.NoError(respErr)
	require.NotNil(resp)
	require.NotEmpty(resp.Token)

	// Reconnect with the token
	authedClient := server.TestServer(t, impl, server.TestWithToken(resp.Token))

	// Open the config stream
	stream, err := authedClient.RunnerConfig(ctx)
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

	return id, authedClient
}

// TestApp creates the app in the DB.
func TestApp(t testing.T, client pb.WaypointClient, ref *pb.Ref_Application) {
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name: ref.Project,
			},
		})
		require.NoError(t, err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(t, err)
	}
}

func testDB(t testing.T) *bolt.DB {
	t.Helper()

	// Temporary directory for the database
	td, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(td) })

	// Create the DB
	db, err := bolt.Open(filepath.Join(td, "test.db"), 0600, nil)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	return db
}

// WithDB return a server option with a boltdb state provider
func WithDB(db *bolt.DB) Option {
	return func(s *Service, cfg *config) error {
		// Initialize our state
		state, err := boltdbstate.New(hclog.Default(), db)
		if err != nil {
			return err
		}
		cfg.stateProvider = func(_ context.Context) serverstate.Interface {
			return state
		}
		return nil
	}
}
