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

	// Create project with application
	jobTemplate := serverptypes.TestJobNew(t, nil)
	appRef := jobTemplate.Application
	TestApp(t, client, appRef)

	// Create some pipelines in the project
	for _, name := range []string{"alpha", "beta"} {
		_, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipeline(t, &pb.Pipeline{
				Id:   name,
				Name: name,
				Owner: &pb.Pipeline_Project{
					Project: &pb.Ref_Project{
						Project: appRef.Project,
					},
				},
			}),
		})
		require.NoError(err)
	}

	t.Run("list", func(t *testing.T) {
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.Len(resp.Pipelines, 2)
	})

	t.Run("with no runs", func(t *testing.T) {
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.Nil(resp.Pipelines[0].LastRun)
		require.EqualValues(0, resp.Pipelines[0].TotalRuns)
	})

	t.Run("with some runs", func(t *testing.T) {
		// Add some runs
		client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: "alpha",
				},
			},
			JobTemplate: jobTemplate,
		})
		client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: "beta",
				},
			},
			JobTemplate: jobTemplate,
		})

		// Call the method
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Pipelines[0].LastRun)
		require.NotNil(resp.Pipelines[1].LastRun)
		require.EqualValues(1, resp.Pipelines[0].TotalRuns)
	})
}
