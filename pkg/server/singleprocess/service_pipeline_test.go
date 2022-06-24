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

	t.Run("get an missing pipeline by id returns nothing", func(t *testing.T) {
		require := require.New(t)

		pResp, err := client.GetPipeline(ctx, &pb.GetPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: &pb.Ref_PipelineId{
						Id: "doesnotexist",
					},
				},
			},
		})
		require.Error(err)
		require.Nil(pResp)
	})

	t.Run("get an existing pipeline by id", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipeline(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Pipeline
		require.NotEmpty(result.Id)

		pResp, err := client.GetPipeline(ctx, &pb.GetPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: &pb.Ref_PipelineId{
						Id: "test",
					},
				},
			},
		})
		require.NoError(err)
		require.NotNil(pResp)
		require.Equal(pResp.RootStep, "root")

		pipeline := pResp.Pipeline
		require.Equal(pipeline.Name, "test")
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

	// Create our pipeline
	pipeline := serverptypes.TestPipeline(t, nil)
	pipeline.Steps["B"] = &pb.Pipeline_Step{
		Name:      "B",
		DependsOn: []string{"root"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}
	pipeline.Steps["C"] = &pb.Pipeline_Step{
		Name:      "C",
		DependsOn: []string{"B"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}
	pipeline.Steps["D"] = &pb.Pipeline_Step{
		Name:      "D",
		DependsOn: []string{"C"},
		Kind: &pb.Pipeline_Step_Build_{
			Build: &pb.Pipeline_Step_Build{},
		},
	}
	pipeline.Steps["E"] = &pb.Pipeline_Step{
		Name:      "E",
		DependsOn: []string{"D"},
		Kind: &pb.Pipeline_Step_Deploy_{
			Deploy: &pb.Pipeline_Step_Deploy{},
		},
	}
	pipeline.Steps["F"] = &pb.Pipeline_Step{
		Name:      "F",
		DependsOn: []string{"E"},
		Kind: &pb.Pipeline_Step_Release_{
			Release: &pb.Pipeline_Step_Release{},
		},
	}

	// Create, should get an ID back
	pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: pipeline,
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

	// Job should exist
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: resp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)

	// We should have all the job IDs
	require.Len(resp.AllJobIds, 6)
	var names []string
	for _, id := range resp.AllJobIds {
		require.Contains(resp.JobMap, id)
		names = append(names, resp.JobMap[id].Step)
	}
	require.Equal([]string{"root", "B", "C", "D", "E", "F"}, names)
}

func TestServiceListPipelines(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create our pipeline
	pipeline := serverptypes.TestPipeline(t, nil)
	pipeline.Steps["B"] = &pb.Pipeline_Step{
		Name:      "B",
		DependsOn: []string{"root"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}
	pipeline.Steps["C"] = &pb.Pipeline_Step{
		Name:      "C",
		DependsOn: []string{"B"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}

	// Create, should get an ID back
	pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: pipeline,
	})
	require.NoError(err)
	require.NotNil(pipeResp)
	require.Equal(pipeResp.Pipeline.Name, "test")
	require.Equal(pipeResp.Pipeline.Id, "test")

	pipelinesResp, err := client.ListPipelines(ctx, &pb.ListPipelinesRequest{
		Project: &pb.Ref_Project{
			Project: "project",
		},
	})
	require.NoError(err)
	require.Len(pipelinesResp.Pipelines, 1)
	require.Equal(pipelinesResp.Pipelines[0].Name, "test")

	// Create some more, list some more.

	// Create our pipeline
	pipeline = serverptypes.TestPipeline(t, &pb.Pipeline{
		Name: "another",
		Id:   "another",
	})
	pipeline.Steps["B"] = &pb.Pipeline_Step{
		Name:      "B",
		DependsOn: []string{"root"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}
	pipeline.Steps["C"] = &pb.Pipeline_Step{
		Name:      "C",
		DependsOn: []string{"B"},
		Kind: &pb.Pipeline_Step_Exec_{
			Exec: &pb.Pipeline_Step_Exec{
				Image: "hashicorp/waypoint",
			},
		},
	}

	// Create, should get an ID back
	pipeResp, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
		Pipeline: pipeline,
	})
	require.NoError(err)
	require.NotNil(pipeResp)
	require.Equal(pipeResp.Pipeline.Name, "another")
	require.Equal(pipeResp.Pipeline.Id, "another")

	pipelinesResp, err = client.ListPipelines(ctx, &pb.ListPipelinesRequest{
		Project: &pb.Ref_Project{
			Project: "project",
		},
	})
	require.NoError(err)
	require.Len(pipelinesResp.Pipelines, 2)

	// Order dependent tests, there might be a better way to test the pipelines
	// we get back match the ones we inserted.
	require.Equal(pipelinesResp.Pipelines[0].Name, "another")
	require.Equal(pipelinesResp.Pipelines[1].Name, "test")
}
