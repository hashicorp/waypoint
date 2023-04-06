// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func TestServicePipeline_Basic(t *testing.T) {
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
					Id: "doesnotexist",
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
					Id: "test",
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

func TestServicePipeline_Run(t *testing.T) {
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
					Id: pipeResp.Pipeline.Id,
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

		pRef := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		}

		// Pipeline Runs should exist
		runs, err := client.ListPipelineRuns(ctx, &pb.ListPipelineRunsRequest{
			Pipeline: pRef,
		})
		require.NoError(err)
		require.NotEmpty(runs)
		require.Len(runs.PipelineRuns, 1)

		// Get pipeline run
		run, err := client.GetPipelineRun(ctx, &pb.GetPipelineRunRequest{
			Pipeline: pRef,
			Sequence: 1,
		})
		require.NoError(err)
		require.Equal(pipeline.Id, run.PipelineRun.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
		require.Equal(len(run.PipelineRun.Jobs), len(resp.AllJobIds))
		require.Equal(resp.Sequence, run.PipelineRun.Sequence)
	})

	t.Run("runs a pipeline with workspace scoped steps by request", func(t *testing.T) {
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
			Workspace: &pb.Ref_Workspace{
				Workspace: "staging",
			},
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
			Workspace: &pb.Ref_Workspace{
				Workspace: "default",
			},
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
					Id: pipeResp.Pipeline.Id,
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

		// check all workspaces equal what we expect
		for stepName, stepSrc := range pipeline.Steps {
			// find the job that matches this step
			var jobId string
			for id, jobStep := range resp.JobMap {
				if stepName == jobStep.Step {
					jobId = id
					break
				}
			}
			stepJob, err := client.GetJob(ctx, &pb.GetJobRequest{
				JobId: jobId,
			})
			require.NoError(err)
			require.NotEmpty(stepJob)
			require.Equal(stepJob.Id, jobId)

			// the default jobs we're using for tests come with a default
			// workspace "w_test", see TestJobNew usage
			expectedWorkspaceVal := "w_test"
			if stepSrc.Workspace != nil {
				expectedWorkspaceVal = stepSrc.Workspace.Workspace
			}

			require.Equal(stepJob.Workspace.Workspace, expectedWorkspaceVal)
		}

		pRef := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		}

		// Pipeline Runs should exist
		runs, err := client.ListPipelineRuns(ctx, &pb.ListPipelineRunsRequest{
			Pipeline: pRef,
		})
		require.NoError(err)
		require.NotEmpty(runs)
		require.Len(runs.PipelineRuns, 1)

		// Get pipeline run
		run, err := client.GetPipelineRun(ctx, &pb.GetPipelineRunRequest{
			Pipeline: pRef,
			Sequence: 1,
		})
		require.NoError(err)
		require.Equal(pipeline.Id, run.PipelineRun.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
		require.Equal(len(run.PipelineRun.Jobs), len(resp.AllJobIds))
		require.Equal(resp.Sequence, run.PipelineRun.Sequence)
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
							Id: "embed",
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
					Id: pipeResp.Pipeline.Id,
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
		require.Len(resp.JobMap, 6)
		require.Len(resp.AllJobIds, 6)
	})

	t.Run("runs a pipeline with embedded pipeline and workspace", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a job to use in TestJobNew that specifies the default
		// workspace. "default" is used for operations that do not have a
		// workspace specified, and we need to check for this sentinel value
		// when overriding embedded pipeline steps workspace value.
		defaultJob := &pb.Job{
			Workspace: &pb.Ref_Workspace{
				Workspace: "default",
			},
		}
		// Initialize our app
		TestApp(t, client, serverptypes.TestJobNew(t, defaultJob).Application)

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
			Workspace: &pb.Ref_Workspace{
				Workspace: "embedded-workspace",
			},
			Kind: &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Id{
							Id: "embed",
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
				"otherws": {
					Name: "otherws",
					Workspace: &pb.Ref_Workspace{
						Workspace: "otherws",
					},
					DependsOn: []string{"second"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"third": {
					Name:      "third",
					DependsOn: []string{"second"},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
			},
		}

		// Create, should get an ID back
		embeddedResp, err := client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: embedPipeline,
		})
		require.NoError(err)

		// Build our job template
		jobTemplate := serverptypes.TestJobNew(t, defaultJob)
		resp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: pipeResp.Pipeline.Id,
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
		require.Len(resp.JobMap, 8)
		require.Len(resp.AllJobIds, 8)

		// find our embedded jobs and verify they used the correct workspace
		for id, jobStep := range resp.JobMap {
			if jobStep.PipelineId == embeddedResp.Pipeline.Id {
				// all the embedded steps should inherit the parent steps
				// workspace, except the step named "otherws", which has itself
				// a specified workspace
				expectedWs := "embedded-workspace"
				if jobStep.Step == "otherws" {
					expectedWs = "otherws"
				}
				embeddedJob, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: id})
				require.NoError(err)
				require.Equal(expectedWs, embeddedJob.Workspace.Workspace)
			}
		}
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
									Id: "embed",
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
									Id: "test",
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
					Id: pipeResp.Pipeline.Id,
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
							Id: "embed",
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
					Id: pipeResp.Pipeline.Id,
				},
			},
			JobTemplate: jobTemplate,
		})

		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.AllJobIds, 6)
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
									Id: "uzumaki",
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
									Id: "embed",
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
									Id: "test",
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
									Id: "uzumaki",
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
									Id: "test",
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
					Id: "test",
				},
			},
			JobTemplate: jobTemplate,
		})

		// You shouldn't be able to run a cyclic pipeline
		require.Error(err)
		require.Nil(resp)
	})

	t.Run("returns job ids in expected order for a single pipeline", func(t *testing.T) {
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
					Id: pipeResp.Pipeline.Id,
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
		var allJobs []*pb.Job
		for _, id := range resp.AllJobIds {
			require.Contains(resp.JobMap, id)
			names = append(names, resp.JobMap[id].Step)

			j, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: id})
			require.NoError(err)
			allJobs = append(allJobs, j)
		}
		require.Equal([]string{"root", "B", "C", "D", "E", "F", "G"}, names)

		// Loop through all Job Ids returned from RunPipeline and verify that the
		// order lines up with the expected order of the pipeline for the test.
		for i, job := range allJobs {
			require.Equal(job.Pipeline.PipelineId, pipeline.Id)
			require.Equal(job.Pipeline.PipelineName, pipeline.Name)

			switch i {
			case 0:
				require.Equal(job.Pipeline.Step, "root")
			case 1:
				require.Equal(job.Pipeline.Step, "B")
			case 2:
				require.Equal(job.Pipeline.Step, "C")
			case 3:
				require.Equal(job.Pipeline.Step, "D")
			case 4:
				require.Equal(job.Pipeline.Step, "E")
			case 5:
				require.Equal(job.Pipeline.Step, "F")
			case 6:
				require.Equal(job.Pipeline.Step, "G")
			}
		}

	})

	t.Run("returns job ids in expected order for an embedded pipeline", func(t *testing.T) {
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

		// Having two embedded pipeline step references is what is causing the bug
		// at the moment, i.e. https://github.com/hashicorp/waypoint/issues/3869
		pipeline.Steps["Embed"] = &pb.Pipeline_Step{
			Name:      "Embed",
			DependsOn: []string{"C"},
			Kind: &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Id{
							Id: "embed",
						},
					},
				},
			},
		}
		pipeline.Steps["AnotherEmbed"] = &pb.Pipeline_Step{
			Name:      "AnotherEmbed",
			DependsOn: []string{"Embed"},
			Kind: &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Id{
							Id: "twoembed",
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

		// Create our pipeline
		embedTwoPipeline := &pb.Pipeline{
			Id:   "twoembed",
			Name: "twoembed",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"one": {
					Name: "one",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image: "hashicorp/waypoint",
						},
					},
				},
				"two": {
					Name:      "two",
					DependsOn: []string{"one"},
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
			Pipeline: embedTwoPipeline,
		})
		require.NoError(err)

		// Build our job template
		jobTemplate := serverptypes.TestJobNew(t, nil)
		resp, err := client.RunPipeline(ctx, &pb.RunPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: pipeResp.Pipeline.Id,
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
		require.Len(resp.JobMap, 9)
		require.Len(resp.AllJobIds, 9)

		var allJobs []*pb.Job
		for _, id := range resp.AllJobIds {
			require.Contains(resp.JobMap, id)

			j, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: id})
			require.NoError(err)
			allJobs = append(allJobs, j)
		}

		idx := func(p1, s1 string) int {
			for i, job := range allJobs {
				if job.Pipeline.PipelineName == p1 && job.Pipeline.Step == s1 {
					return i
				}
			}

			return -1
		}

		before := func(p1, s1, p2, s2 string) {
			a := idx(p1, s1)
			b := idx(p2, s2)

			require.True(a < b, "expected %s.%s before %s.%s", p1, s1, p2, s2)
		}

		// Ensure that flattened list of jobs matches expected order for embedded pipelines
		before("test", "root", "test", "B")
		before("test", "B", "test", "C")
		before("test", "C", "test", "Embed")
		before("embed", "first", "test", "Embed")
		before("embed", "second", "test", "Embed")
		before("test", "Embed", "test", "AnotherEmbed")
		before("twoembed", "one", "test", "AnotherEmbed")
		before("twoembed", "two", "test", "AnotherEmbed")
	})

}

func TestServicePipeline_List(t *testing.T) {
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
