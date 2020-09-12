package ceb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{Client: client})

	// We should get registered
	require.Eventually(func() bool {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 2*time.Second, 10*time.Millisecond)

	// Shut down the server
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// We should get deregistered
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 0
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// We should get re-registered
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 5*time.Second, 10*time.Millisecond)
}

// Test how the CEB behaves when the server is down on startup.
func TestRun_serverDown(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Create a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)
	dep := resp.Deployment

	// But actually shut it down
	cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: "doesn't matter",
			},
		})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client:       client,
		DeploymentId: dep.Id,
		Helper:       "write-file",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// We should get re-registered eventually
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 5*time.Second, 10*time.Millisecond)
}

func TestRun_serverDownRequired(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Create a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)
	dep := resp.Deployment

	// But actually shut it down
	cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: "doesn't matter",
			},
		})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	testRun(t, context.Background(), &testRunOpts{
		Client:       client,
		DeploymentId: dep.Id,
		Helper:       "write-file",
		HelperEnv: map[string]string{
			envServerRequired: "1",
			"HELPER_PATH":     path,
		},
	})

	// The child should NOT start up
	time.Sleep(1 * time.Second)
	_, err = ioutil.ReadFile(path)
	require.Error(err)

	// Restart
	restartCh <- struct{}{}

	// The child should start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}

var (
	testExec      = os.Args[0]
	envHelperMode = "TEST_HELPER_MODE"
)

func TestMain(m *testing.M) {
	switch os.Getenv(envHelperMode) {
	case "":
		// Log
		hclog.L().SetLevel(hclog.Trace)

		// Normal test mode
		os.Exit(m.Run())

	case "sleep":
		time.Sleep(10 * time.Minute)

	case "logs-stdout":
		for {
			time.Sleep(250 * time.Millisecond)
			fmt.Println(time.Now().String())
		}

	case "write-file":
		path := os.Getenv("HELPER_PATH")
		if path == "" {
			panic("bad")
		}

		ioutil.WriteFile(path, []byte("hello"), 0600)
		time.Sleep(10 * time.Minute)

	default:
		panic("invalid helperfunc")
	}
}

func testRun(t *testing.T, ctx context.Context, opts *testRunOpts) *CEB {
	if opts == nil {
		opts = &testRunOpts{}
	}

	if opts.Client == nil {
		opts.Client = singleprocess.TestServer(t)
	}

	if opts.Helper == "" {
		opts.Helper = "sleep"
	}

	if opts.DeploymentId == "" {
		resp, err := opts.Client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, nil),
		})
		require.NoError(t, err)
		opts.DeploymentId = resp.Deployment.Id
	}

	// Setup our deployment
	require.NoError(t, os.Setenv(envDeploymentId, opts.DeploymentId))
	t.Cleanup(func() { os.Setenv(envDeploymentId, "") })

	// Setup our helper
	require.NoError(t, os.Setenv(envHelperMode, opts.Helper))
	t.Cleanup(func() { os.Setenv(envHelperMode, "") })

	// Setup helper env vars
	for k, v := range opts.HelperEnv {
		kcopy := k
		require.NoError(t, os.Setenv(k, v))
		t.Cleanup(func() { os.Setenv(kcopy, "") })
	}

	// This is so we can wait to get it set
	cebCh := make(chan *CEB, 1)

	// Run it
	go Run(ctx,
		WithExec([]string{testExec}),
		WithClient(opts.Client),
		WithEnvDefaults(),
		withCEBValue(cebCh),
	)

	return <-cebCh
}

type testRunOpts struct {
	Client       pb.WaypointClient
	Helper       string
	HelperEnv    map[string]string
	DeploymentId string
}
