package runner

import (
	"context"
	"testing"

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

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept())

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}
