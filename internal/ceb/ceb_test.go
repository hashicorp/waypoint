package ceb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"

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
	DeploymentId string
}
