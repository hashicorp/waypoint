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
		TestServiceUI_GetPipelineRun,
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

func TestServiceUI_GetPipelineRun(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create project with application
	jobTemplate := serverptypes.TestJobNew(t, nil)
	appRef := jobTemplate.Application
	TestApp(t, client, appRef)

	// Create a pipeline
	pipelineResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: &pb.Pipeline{
			Id:   "test",
			Name: "test",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: appRef.Project,
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"hello": {
					DependsOn: []string{},
					Name:      "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"Hello"},
						},
					},
				},
				"bye": {
					DependsOn: []string{"hello"},
					Name:      "bye",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"Bye"},
						},
					},
				},
			},
		},
	})
	require.NoError(err)
	pipeline := pipelineResp.Pipeline

	// Create a run
	runResp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		},
		JobTemplate: jobTemplate,
	})
	require.NoError(err)
	seq := runResp.Sequence

	// Call UI_GetPipelineRun
	resp, err := client.UI_GetPipelineRun(ctx, &pb.UI_GetPipelineRunRequest{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		},
		Sequence: seq,
	})
	require.NoError(err)
	require.NotNil(resp.PipelineRun)
	require.Equal(seq, resp.PipelineRun.Sequence)

	hello := resp.RootTreeNode
	require.NotNil(hello)
	require.Equal("hello", hello.Step.Name)
	require.Equal(pb.UI_PipelineRunTreeNode_QUEUED, hello.State)
	require.Equal(runResp.AllJobIds[0], hello.Job.Id)
	require.Equal(pb.UI_PipelineRunTreeNode_Children_SERIAL, hello.Children.Mode)
	require.Len(hello.Children.Nodes, 1)

	bye := resp.RootTreeNode.Children.Nodes[0]
	require.NotNil(bye)
	require.Equal("bye", bye.Step.Name)
	require.Equal(pb.UI_PipelineRunTreeNode_QUEUED, bye.State)
	require.Equal(runResp.AllJobIds[1], bye.Job.Id)
	require.Equal(pb.UI_PipelineRunTreeNode_Children_SERIAL, bye.Children.Mode)
	require.Len(bye.Children.Nodes, 0)
}
