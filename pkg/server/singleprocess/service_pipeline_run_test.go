package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServicePipelineRun(t *testing.T) {
	t.Run("get and list", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Initialize our app
		TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

		// Create pipeline
		p, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipeline(t, nil),
		})
		require.NoError(err)
		require.NotNil(p)
		require.NotEmpty(p.Pipeline.Id)

		pRef := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{
					Id: p.Pipeline.Id,
				},
			},
		}

		// Run Pipeline once
		jobTemplate := serverptypes.TestJobNew(t, nil)
		resp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline:    pRef,
			JobTemplate: jobTemplate,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Get pipeline run
		run, err := client.GetPipelineRun(ctx, &pb.GetPipelineRunRequest{
			Pipeline: pRef,
			Sequence: 1,
		})
		require.NoError(err)
		require.Equal(p.Pipeline.Id, run.PipelineRun.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
		require.Equal(len(run.PipelineRun.Jobs), len(resp.AllJobIds))
		require.Equal(resp.Sequence, run.PipelineRun.Sequence)

		// Run Pipeline again
		resp, err = client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline:    pRef,
			JobTemplate: jobTemplate,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Two pipeline runs should exist
		runs, err := client.ListPipelineRuns(ctx, &pb.ListPipelineRunsRequest{
			Pipeline: pRef,
		})
		require.NoError(err)
		require.NotEmpty(runs)
		require.Len(runs.PipelineRuns, 2)

		// Run Pipeline again
		resp, err = client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline:    pRef,
			JobTemplate: jobTemplate,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Three pipeline runs should exist
		runs, err = client.ListPipelineRuns(ctx, &pb.ListPipelineRunsRequest{
			Pipeline: pRef,
		})
		require.NoError(err)
		require.NotEmpty(runs)
		require.Len(runs.PipelineRuns, 3)
		require.Equal(uint64(3), runs.PipelineRuns[len(runs.PipelineRuns)-1].Sequence)
	})
}
