package ceb

import (
	"context"
	"os"
	"testing"
	"time"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := singleprocess.TestServer(t)
	ceb := testRun(t, ctx, &testRunOpts{Client: client})

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
}

var (
	testExec      = os.Args[0]
	envHelperMode = "TEST_HELPER_MODE"
)

func TestMain(m *testing.M) {
	switch os.Getenv(envHelperMode) {
	case "":
		// Normal test mode
		os.Exit(m.Run())

	case "sleep":
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
