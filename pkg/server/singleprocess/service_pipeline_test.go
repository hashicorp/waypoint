package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	t.Run("get an existing pipeline that has cycles returns error", func(t *testing.T) {
		require := require.New(t)

		// Create, should return an error about a cycle
		_, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: serverptypes.TestPipelineCycle(t, nil),
		})
		require.Error(err)
		require.Equal(codes.InvalidArgument, status.Code(err))
	})

}

func TestServiceRunPipeline(t *testing.T) {
	t.Run("runs a pipeline by request", func(t *testing.T) {
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
		pipeline.Steps["G"] = &pb.Pipeline_Step{
			Name:      "G",
			DependsOn: []string{"F"},
			Kind: &pb.Pipeline_Step_Up_{
				Up: &pb.Pipeline_Step_Up{},
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
		require.Len(resp.AllJobIds, 7)
		var names []string
		for _, id := range resp.AllJobIds {
			require.Contains(resp.JobMap, id)
			names = append(names, resp.JobMap[id].Step)
		}
		require.Equal([]string{"root", "B", "C", "D", "E", "F", "G"}, names)
	})

	t.Run("runs a pipeline with embedded pipelines by request", func(t *testing.T) {
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
		pipeline.Steps["Embed"] = &pb.Pipeline_Step{
			Name:      "Embed",
			DependsOn: []string{"C"},
			Kind: &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Id{
							Id: &pb.Ref_PipelineId{
								Id: "embed",
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		require.NoError(err)

		// Create another pipeline that references the first one
		// Create our pipeline
		embedPipeline := &pb.Pipeline{
			Id:   "embed",
			Name: "embed",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"first": {
					Name: "first",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"second": {
					Name:      "second",
					DependsOn: []string{"first"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: embedPipeline,
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
		require.Len(resp.JobMap, 5)
		require.Len(resp.AllJobIds, 5)
	})

	t.Run("returns an error if theres a cycle detected", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Initialize our app
		TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

		// Create our pipeline
		pipeline := &pb.Pipeline{
			Id:   "test",
			Name: "test",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"A": {
					Name: "A",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"B": {
					Name:      "B",
					DependsOn: []string{"A"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"Embed": {
					Name:      "Embed",
					DependsOn: []string{"B"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "embed",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		require.NoError(err)

		// Create another pipeline that references the first one
		// Create our pipeline
		embedPipeline := &pb.Pipeline{
			Id:   "embed",
			Name: "embed",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"first": {
					Name: "first",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"second": {
					Name:      "second",
					DependsOn: []string{"first"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"embedpt2": {
					Name:      "embedpt2",
					DependsOn: []string{"second"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "test",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: embedPipeline,
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

		// You shouldn't be able to run a cyclic pipeline
		require.Error(err)
		require.Nil(resp)
	})

	// This is likely an internal error if it does error. When we go to construct
	// a full graph including any embedded pipeline graphs, we rename the step nodes
	// to have their own unique ids and keep track of which node id is associated to
	// which pipeline and step name.
	t.Run("does not error when multiple pipelines have the same step names", func(t *testing.T) {
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
		pipeline.Steps["Embed"] = &pb.Pipeline_Step{
			Name:      "Embed",
			DependsOn: []string{"C"},
			Kind: &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Id{
							Id: &pb.Ref_PipelineId{
								Id: "embed",
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		pipeResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		require.NoError(err)

		// Create another pipeline that references the first one
		// Create our pipeline
		embedPipeline := &pb.Pipeline{
			Id:   "embed",
			Name: "embed",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"B": {
					Name: "B",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"C": {
					Name:      "C",
					DependsOn: []string{"B"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: embedPipeline,
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
		require.Len(resp.AllJobIds, 5)
	})

	// This is a nasty cycle graph
	t.Run("returns an error if theres a cycle detected on a deeply cyclic graph", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Initialize our app
		TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

		// Create our pipeline
		pipeline := &pb.Pipeline{
			Id:   "test",
			Name: "test",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"A": {
					Name: "A",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"Uzumaki": {
					Name:      "Uzumaki",
					DependsOn: []string{"Embed"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "uzumaki",
									},
								},
							},
						},
					},
				},
				"B": {
					Name:      "B",
					DependsOn: []string{"A"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"Embed": {
					Name:      "Embed",
					DependsOn: []string{"B"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "embed",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		require.NoError(err)

		// Create another pipeline that references the first one
		// Create our pipeline
		embedPipeline := &pb.Pipeline{
			Id:   "embed",
			Name: "embed",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"first": {
					Name: "first",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"second": {
					Name:      "second",
					DependsOn: []string{"first"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"embedpt2": {
					Name:      "embedpt2",
					DependsOn: []string{"second"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "test",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: embedPipeline,
		})
		require.NoError(err)

		// Create another pipeline that references the first one
		// Create our pipeline
		uzumakiPipeline := &pb.Pipeline{
			Id:   "uzumaki",
			Name: "uzumaki",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"spiral": {
					Name: "spiral",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"self": {
					Name:      "self",
					DependsOn: []string{"spiral"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "uzumaki",
									},
								},
							},
						},
					},
				},
				"vortex": {
					Name:      "vortex",
					DependsOn: []string{"spiral"},
					Kind: &pb.Pipeline_Step_Pipeline_{
						Pipeline: &pb.Pipeline_Step_Pipeline{
							Ref: &pb.Ref_Pipeline{
								Ref: &pb.Ref_Pipeline_Id{
									Id: &pb.Ref_PipelineId{
										Id: "test",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create, should get an ID back
		_, err = client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: uzumakiPipeline,
		})
		require.NoError(err)

		// Build our job template
		jobTemplate := serverptypes.TestJobNew(t, nil)
		resp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: &pb.Ref_PipelineId{
						Id: "test",
					},
				},
			},
			JobTemplate: jobTemplate,
		})

		// You shouldn't be able to run a cyclic pipeline
		require.Error(err)
		require.Nil(resp)
	})

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
