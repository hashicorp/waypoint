package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServicePipeline(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	type Req = pb.UpsertPipelineRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipeline(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Pipeline
		require.NotEmpty(result.Id)

		// Let's write some data
		testName := "TestyTest"
		result.Name = testName
		resp, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Pipeline
		require.Equal(result.Name, testName)
	})

}

func TestServiceRunPipeline(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create, should get an ID back
	pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: serverptypes.TestPipeline(t, nil),
	})
	require.NoError(err)

	// Build our job template
	jobTemplate := serverptypes.TestJobNew(t, nil)
	resp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{
					Id: pipeResp.Pipeline.Id,
				},
			},
		},
		JobTemplate: jobTemplate,
	})
	require.NoError(err)
	require.NotNil(resp)

	// Job should exit
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: resp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)
}
