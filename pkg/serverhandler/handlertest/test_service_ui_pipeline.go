package handlertest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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
	dataSourceRef := jobTemplate.DataSourceRef
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
		})
		require.NoError(err)
		require.Len(resp.Pipelines, 2)
	})

	t.Run("with no runs", func(t *testing.T) {
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
		})
		require.NoError(err)
		require.Nil(resp.Pipelines[0].LastRun)
		require.EqualValues(0, resp.Pipelines[0].TotalRuns)
	})

	t.Run("with some runs", func(t *testing.T) {
		// Add some runs
		for i := 0; i < 3; i++ {
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
		}

		// Call the method
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
		})
		require.NoError(err)

		require.Len(resp.Pipelines, 2)

		for _, p := range resp.Pipelines {
			require.NotNil(p.LastRun)
			require.NotNil(p.LastRun.QueueTime)
			require.Equal(appRef.Application, p.LastRun.Application.Application)
			require.Truef(proto.Equal(dataSourceRef, p.LastRun.DataSourceRef), "expected %#v to equal %#v", dataSourceRef, p.LastRun.DataSourceRef)
			require.EqualValues(3, p.TotalRuns)
		}
	})

	t.Run("with page size request", func(t *testing.T) {
		t.Skip("TODO: implement pagination for UI_ListPipelines")
		// create 9 pipelines in addition to the 2 already created above
		for i := 1; i < 10; i++ {
			name := fmt.Sprintf("pipeline-%d", i)
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
		resp, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
			Project: &pb.Ref_Project{
				Project: appRef.Project,
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.Len(resp.Pipelines, 10)
		require.NotNil(resp.Pagination.NextPageToken)
	})
}
