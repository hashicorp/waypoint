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
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	t.Run("get and list", func(t *testing.T) {
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
		require.Equal(p.Pipeline.Id, run.PipelineRun.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id)
		require.Equal(len(run.PipelineRun.Jobs), len(resp.AllJobIds))
		require.Equal(resp.Sequence, run.PipelineRun.Sequence)
	})
}
