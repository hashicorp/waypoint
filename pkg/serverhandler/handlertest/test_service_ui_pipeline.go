package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["ui_pipeline"] = []testFunc{
		TestServiceUI_ListPipelines,
	}
}

func TestServiceUI_ListPipelines(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create project
	var project *pb.Project
	{
		resp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: serverptypes.TestProject(t, &pb.Project{}),
		})
		require.NoError(err)
		project = resp.Project
	}

	// Create some pipelines in the project
	for _, name := range []string{"alpha", "beta"} {
		_, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipeline(t, &pb.Pipeline{
				Id:   name,
				Name: name,
				Owner: &pb.Pipeline_Project{
					Project: &pb.Ref_Project{
						Project: project.Name,
					},
				},
			}),
		})
		require.NoError(err)
	}

	t.Run("with no runs", func(t *testing.T) {
		// Call the method
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: project.Name,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.Len(resp.Pipelines, 2)
		require.Nil(resp.Pipelines[0].LastRun)
		// TODO: require.EqualValues(2, resp.TotalCount)
	})

	t.Run("with some runs", func(t *testing.T) {
		// TODO: Add some runs somehow

		// Call the method
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: project.Name,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Pipelines[0].LastRun)
		require.NotNil(resp.Pipelines[1].LastRun)
	})
}
