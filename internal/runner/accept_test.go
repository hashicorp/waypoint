package runner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestRunnerAccept(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_cancel(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Set a blocker
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Cancel the context eventually. This isn't CI-sensitive cause
	// we'll block no matter what.
	time.AfterFunc(500*time.Millisecond, cancel)

	// Accept should complete with an error
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	require.Eventually(func() bool {
		job, err := client.GetJob(context.Background(), &pb.GetJobRequest{JobId: jobId})
		require.NoError(err)
		return job.State == pb.Job_ERROR
	}, 3*time.Second, 25*time.Millisecond)
}
