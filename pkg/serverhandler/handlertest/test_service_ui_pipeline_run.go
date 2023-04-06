// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["ui_pipeline_run"] = []testFunc{
		TestServiceUI_ListPipelineRuns,
	}
}

func TestServiceUI_ListPipelineRuns(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create project with application
	jobTemplate := serverptypes.TestJobNew(t, nil)
	appRef := jobTemplate.Application
	dataSourceRef := jobTemplate.DataSourceRef
	TestApp(t, client, appRef)

	// Create a pipeline in the project
	_, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: serverptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "alpha",
			Name: "alpha",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: appRef.Project,
				},
			},
		}),
	})
	require.NoError(err)

	t.Run("list no runs", func(t *testing.T) {
		resp, err := client.UI_ListPipelineRuns(ctx, &pb.UI_ListPipelineRunsRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: "alpha",
				},
			},
		})
		require.NoError(err)
		require.Len(resp.PipelineRunBundles, 0)
	})

	t.Run("list with runs", func(t *testing.T) {
		// Create runs
		for i := 0; i < 3; i++ {
			client.RunPipeline(ctx, &pb.RunPipelineRequest{
				Pipeline: &pb.Ref_Pipeline{
					Ref: &pb.Ref_Pipeline_Id{
						Id: "alpha",
					},
				},
				JobTemplate: jobTemplate,
			})
		}
		resp, err := client.UI_ListPipelineRuns(ctx, &pb.UI_ListPipelineRunsRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: "alpha",
				},
			},
		})
		require.NoError(err)
		for i := 1; i < len(resp.PipelineRunBundles); i++ {
			require.NotNil(resp.PipelineRunBundles[i].QueueTime)
			require.Equal(appRef.Application, resp.PipelineRunBundles[i].Application.Application)
			require.Truef(proto.Equal(dataSourceRef, resp.PipelineRunBundles[i].DataSourceRef), "expected %#v to equal %#v", dataSourceRef, resp.PipelineRunBundles[i].DataSourceRef)
		}
		require.Len(resp.PipelineRunBundles, 3)
	})

	t.Run("with page size request", func(t *testing.T) {
		t.Skip("TODO: implement pagination for UI_ListPipelineRuns")
		// create 8 pipeline runs in addition to the 3 already created above
		for i := 1; i < 9; i++ {
			client.RunPipeline(ctx, &pb.RunPipelineRequest{
				Pipeline: &pb.Ref_Pipeline{
					Ref: &pb.Ref_Pipeline_Id{
						Id: "alpha",
					},
				},
				JobTemplate: jobTemplate,
			})
			require.NoError(err)
		}
		resp, err := client.UI_ListPipelineRuns(ctx, &pb.UI_ListPipelineRunsRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: "alpha",
				},
			},
			Pagination: &pb.PaginationRequest{
				PageSize: 10,
			},
		})
		require.NoError(err)
		require.Len(resp.PipelineRunBundles, 10)
		require.NotNil(resp.Pagination.NextPageToken)
	})
}
