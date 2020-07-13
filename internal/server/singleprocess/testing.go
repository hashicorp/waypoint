package singleprocess

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	hznhub "github.com/hashicorp/horizon/pkg/hub"
	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"
	wphzn "github.com/hashicorp/waypoint-hzn/pkg/server"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// TestServer starts a singleprocess server and returns the connected client.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T, opts ...Option) pb.WaypointClient {
	impl, err := New(append(
		[]Option{WithDB(testDB(t))},
		opts...,
	)...)
	require.NoError(t, err)
	return server.TestServer(t, impl)
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
	if out != nil {
		*out = *setup
	}

	// Make our test registration API
	wphzndata := wphzn.TestServer(t)

	return func(s *service, cfg *config) error {
		if cfg.serverConfig == nil {
			cfg.serverConfig = &configpkg.ServerConfig{}
		}

		cfg.serverConfig.URL = &configpkg.URL{
			Enabled:        true,
			APIAddress:     wphzndata.Addr,
			APIInsecure:    true,
			ControlAddress: fmt.Sprintf("dev://%s", setup.HubAddr),
			Token:          setup.AgentToken,
		}

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
	id, err := server.Id()
	require.NoError(mergo.Merge(r, &pb.Runner{Id: id}))

	// Open the config stream
	stream, err := client.RunnerConfig(ctx)
	require.NoError(err)
	t.Cleanup(func() { stream.CloseSend() })

	// Register
	require.NoError(err)
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
