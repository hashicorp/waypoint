package handlertest

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["job"] = []testFunc{
		TestServiceJob,
		TestServiceJob_List,
		TestServiceQueueJob,
		TestServiceValidateJob,
		TestServiceGetJobStream_complete,
		TestServiceGetJobStream_bufferedData,
		TestServiceGetJobStream_completedBufferedData,
		TestServiceGetJobStream_expired,
		TestServiceQueueJob_odr_basic,
		TestServiceQueueJob_odr_default,
		TestServiceQueueJob_odr_target_id,
		TestServiceQueueJob_odr_target_labels,
	}
}

func TestServiceJob(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, impl := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	t.Run("create and get success", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Job should exist and be queued
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: resp.JobId})
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.State)
	})

	t.Run("fails to get non-existent job", func(t *testing.T) {
		require := require.New(t)

		// Job should not exist
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: "NotRealJobId"})
		require.Error(err)
		require.Nil(job)
	})

	_, err := impl.State(ctx).JobList(ctx, &pb.ListJobsRequest{})
	require.NoError(t, err)

}

func TestServiceJob_List(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	t.Run("create and list jobs", func(t *testing.T) {
		require := require.New(t)

		// No jobs
		jobList, err := client.ListJobs(ctx, &pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobList.Jobs, 0)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Create, should get an ID back
		resp, err = client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Create, should get an ID back
		resp, err = client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Three jobs
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobList.Jobs, 3)
	})

	t.Run("list jobs with filters", func(t *testing.T) {
		require := require.New(t)

		// 3 jobs from previous test
		jobList, err := client.ListJobs(ctx, &pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobList.Jobs, 3)

		// Three jobs filtered on workspace
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "w_test",
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 3)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, &pb.Job{
				Workspace: &pb.Ref_Workspace{
					Workspace: "prod",
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// One job filtered on workspace
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "prod",
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 1)
		require.Equal(jobList.Jobs[0].Workspace.Workspace, "prod")

		proj := serverptypes.TestProject(t, &pb.Project{Name: "new"})
		_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: proj,
		})
		require.NoError(err)

		resp, err = client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, &pb.Job{
				Application: &pb.Ref_Application{
					Application: "new",
					Project:     "new",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "dev",
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// One job filtered on dev workspace
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 1)
		require.Equal(jobList.Jobs[0].Workspace.Workspace, "dev")

		// No "new" Project apps in prod workspace.
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "prod",
			},
			Application: &pb.Ref_Application{
				Application: "new",
				Project:     "new",
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 0)

		// "new" project is in dev workspace
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
			Application: &pb.Ref_Application{
				Application: "new",
				Project:     "new",
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 1)
		require.Equal(jobList.Jobs[0].Workspace.Workspace, "dev")

		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			JobState: []pb.Job_State{pb.Job_QUEUED},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 5)

		// Default mocked jobs target Any runners
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{},
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 5)

		// Target a specific runner by id
		resp, err = client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, &pb.Job{
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "123",
						},
					},
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Filter all jobs on those who target the 123 runner
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: "123",
					},
				},
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 1)
		require.Equal(resp.JobId, jobList.Jobs[0].Id)

		// Queue a job who targets a runner with labels
		resp, err = client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, &pb.Job{
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Labels{
						Labels: &pb.Ref_RunnerLabels{
							Labels: map[string]string{"123": "yes", "456": "maybe", "789": "perhaps"},
						},
					},
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// A job with target runner labels will match
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{"123": "yes"},
					},
				},
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 1)
		require.Equal(resp.JobId, jobList.Jobs[0].Id)

		// A job with matching target runner label keys but different values
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{"123": "no"},
					},
				},
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 0)

		// A job with target runner labels but requested label key does not exist
		jobList, err = client.ListJobs(ctx, &pb.ListJobsRequest{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{"abc": "its easy as"},
					},
				},
			},
		})
		require.NoError(err)
		require.Len(jobList.Jobs, 0)
	})
}

func TestServiceQueueJob(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	t.Run("create success", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Job should exist and be queued
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: resp.JobId})
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.State)
	})

	t.Run("expiration", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job:       serverptypes.TestJobNew(t, nil),
			ExpiresIn: "1ms",
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Job should exist and be queued
		require.Eventually(func() bool {
			job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: resp.JobId})
			require.NoError(err)
			return job.State == pb.Job_ERROR && job.CancelTime != nil
		}, 2*time.Second, 10*time.Millisecond)
	})
}

func TestServiceValidateJob(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.ValidateJobRequest

	t.Run("validate success, not assignable", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.ValidateJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.True(resp.Valid)
		require.False(resp.Assignable)
	})

	t.Run("invalid", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		job := serverptypes.TestJobNew(t, nil)
		job.Id = "HELLO"
		resp, err := client.ValidateJob(ctx, &Req{
			Job: job,
		})
		require.NoError(err)
		require.NotNil(resp)
		require.False(resp.Valid)
		require.False(resp.Assignable)
	})
}

func TestServiceGetJobStream_complete(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := runnerStream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		// The assigned runner ID has been set for the job
		require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)

		require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}

	// We should receive an initial state change
	{
		resp, err := stream.Recv()
		require.NoError(err)
		state, ok := resp.Event.(*pb.GetJobStreamResponse_State_)
		require.True(ok, "should be a state change")
		require.NotNil(state)

		require.Equal(pb.Job_UNKNOWN, state.State.Previous)
		require.NotNil(state.State.Job)
	}

	// We need to give the stream time to initialize the output readers
	// so our output below doesn't become buffered. This isn't really that
	// brittle and 100ms should be more than enough.
	time.Sleep(100 * time.Millisecond)

	// Send some output
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "hello",
							},
						},
					},
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "world",
							},
						},
					},
				},
			},
		},
	}))

	// Wait for output
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Terminal_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
		require.NotNil(event)
		require.False(event.Terminal.Buffered, 2)
		switch len(event.Terminal.Events) {
		case 2:
			// Expected case - events were batched.
			require.Equal("hello", event.Terminal.Events[0].Event.(*pb.GetJobStreamResponse_Terminal_Event_Line_).Line.Msg)
			require.Equal("world", event.Terminal.Events[1].Event.(*pb.GetJobStreamResponse_Terminal_Event_Line_).Line.Msg)
		case 1:
			// Not an error if they came in as two separate events, but they need to be in order.
			require.Equal("hello", event.Terminal.Events[0].Event.(*pb.GetJobStreamResponse_Terminal_Event_Line_).Line.Msg)

			// Read again, and should get another event
			resp2 := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Terminal_)(nil))
			event2 := resp2.Event.(*pb.GetJobStreamResponse_Terminal_)
			require.NotNil(event2)
			require.False(event.Terminal.Buffered, 2)
			require.Equal("world", event2.Terminal.Events[0].Event.(*pb.GetJobStreamResponse_Terminal_Event_Line_).Line.Msg)
		default:
			require.Fail("should have received one or two events, got: %d", len(event.Terminal.Events))
		}
	}

	// Send the download event. This realistically could happen after
	// output like above since data sources like Git will output download
	// progress to the UI first before the download is complete.
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Download{
			Download: &pb.GetJobStreamResponse_Download{
				DataSourceRef: &pb.Job_DataSource_Ref{
					Ref: &pb.Job_DataSource_Ref_Git{
						Git: &pb.Job_Git_Ref{
							Commit: "hello",
						},
					},
				},
			},
		},
	}))

	// Wait for download info
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Download_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Download_)
		require.NotNil(event)
		require.NotNil(event.Download.DataSourceRef)

		// We should also receive a job update
		jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Job)(nil))
	}

	// Send the config info event.
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_ConfigLoad_{
			ConfigLoad: &pb.RunnerJobStreamRequest_ConfigLoad{
				Config: &pb.Job_Config{
					Source: pb.Job_Config_SERVER,
				},
			},
		},
	}))

	// Wait for a job change event
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Job)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Job)
		require.NotNil(event)
		require.Equal(pb.Job_Config_SERVER, event.Job.Job.Config.Source)
	}

	// Send final variable values event.
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_VariableValuesSet_{
			VariableValuesSet: &pb.RunnerJobStreamRequest_VariableValuesSet{
				FinalValues: map[string]*pb.Variable_FinalValue{
					"test": {
						Value:  &pb.Variable_FinalValue_Str{Str: "hello"},
						Source: pb.Variable_FinalValue_CLI,
					},
				},
			},
		},
	}))

	// Wait for a job change event
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Job)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Job)
		require.NotNil(event)
		require.Equal(&pb.Variable_FinalValue_Str{Str: "hello"}, event.Job.Job.VariableFinalValues["test"].GetValue())
		require.Equal(pb.Variable_FinalValue_CLI, event.Job.Job.VariableFinalValues["test"].GetSource())
	}

	// Complete the job
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Wait for completion
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Complete_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.NotNil(event)
		require.Nil(event.Complete.Error)
	}
}

// Tests that a client can connect to a job after the job has
// logs (but before the job has completed), and can get the buffered logs.
func TestServiceGetJobStream_bufferedData(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := runnerStream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Send some output
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "hello",
							},
						},
					},
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "world",
							},
						},
					},
				},
			},
		},
	}))

	// Get our job stream and verify we open
	var stream pb.Waypoint_GetJobStreamClient
	require.Eventually(func() bool {
		stream, err = client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
		if err != nil {
			t.Logf("retryable error connecting to job stream: %s", err)
			return false
		}

		// We use require below because no matter what this should always succeed.
		{
			resp, err := stream.Recv()
			require.NoError(err)
			open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
			require.True(ok, "should be an open")
			require.NotNil(open)

		}

		// Wait for buffered output.
		{
			resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Terminal_)(nil))
			event := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
			require.NotNil(event)

			if len(event.Terminal.Events) != 2 || !event.Terminal.Buffered {
				t.Logf("waiting for 2 buffered terminal events, got %d (buffered = %v)",
					len(event.Terminal.Events), event.Terminal.Buffered)
				return false
			}
		}

		return true
	}, 2*time.Second, 50*time.Millisecond)

	// Send a bit more output now that we're connected
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "hello",
							},
						},
					},
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "world",
							},
						},
					},
				},
			},
		},
	}))

	// Wait for output, verify it's unbuffered (live).
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Terminal_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
		require.NotNil(event)
		require.False(event.Terminal.Buffered)
		require.Len(event.Terminal.Events, 2)
	}

	// Complete the job
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Should be done
	_, err = runnerStream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// Wait for completion
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Complete_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.Nil(event.Complete.Error)
	}
}

// Tests that a client can connect to a completed job and still read
// its output.
// NOTE: this does not test that a client can read job logs for a
// long-completed job. Some server state implementation (namely this one)
// do not persist job logs, so streaming completed logs only works
// if the server hasn't restarted or pruned them from memory.
func TestServiceGetJobStream_completedBufferedData(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := runnerStream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Send some output
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "hello",
							},
						},
					},
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "world",
							},
						},
					},
				},
			},
		},
	}))

	// Complete the job
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Should be done
	_, err = runnerStream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

	}

	// Wait for output
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Terminal_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
		require.NotNil(event)
		require.True(event.Terminal.Buffered)
		require.Len(event.Terminal.Events, 2)

	}

	// Wait for completion
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Complete_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.Nil(event.Complete.Error)
	}
}

func TestServiceGetJobStream_expired(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil), ExpiresIn: "10ms"})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)

	// Wait for completion
	{
		resp := jobStreamRecv(t, stream, (*pb.GetJobStreamResponse_Complete_)(nil))
		event := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.NotNil(event)
		require.Equal(int32(codes.Canceled), event.Complete.Error.Code)
	}
}

// jobStreamRecv receives on the stream until an event of the given
// type is matched.
func jobStreamRecv(
	t *testing.T,
	stream pb.Waypoint_GetJobStreamClient,
	typ interface{},
) *pb.GetJobStreamResponse {
	match := reflect.TypeOf(typ)
	for {
		var resp *pb.GetJobStreamResponse
		var err error
		done := make(chan struct{})
		go func() {
			resp, err = stream.Recv()
			close(done)
		}()

		ticker := time.NewTicker(15 * time.Second)
		select {
		case <-ticker.C:
			t.Fatal("timeout receiving job stream event")
		case <-done:
			// request complete! We now have a resp or err
		}
		require.NoError(t, err)

		if reflect.TypeOf(resp.Event) == match {
			return resp
		}

		t.Logf("received event, but not the type we wanted: %T", resp.Event)
	}
}

func TestServiceQueueJob_odr_basic(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Add a loud logger to our context
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "odr-test",
		Level: hclog.Trace,
	})
	ctx = hclog.WithContext(ctx, log)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create with no ODR should error
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: "fake",
			},
		}),
	})
	require.Error(err)
	require.Empty(queueResp)

	// Create an ODR profile
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	odr = cfgResp.Config
	log.Info("test odr profile", "id", odr.Id)

	// Update the project to include ondemand runner
	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	// Create, should get an ID back
	queueResp, err = client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: queueResp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)
	primaryJobId := job.Id

	// task should be PENDING
	taskResp, err := client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: queueResp.JobId,
		},
	}})
	require.NoError(err)
	task := taskResp.Task

	require.Equal(pb.Task_PENDING, task.JobState)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to start the job first.
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)

	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)
	startJobId := assignment.Assignment.Job.Id

	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}

	// Ack it and complete it
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Register our runner
	runnerId := st.Info.EnvironmentVariables["WAYPOINT_RUNNER_ID"]
	server.TestRunner(t, client, &pb.Runner{
		Id:       runnerId,
		ByIdOnly: true,
	})

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: runnerId,
			},
		},
	}))

	// task should be STARTING
	// sleep to ensure job stream request was sent, so that task state updates from
	// the JobAck. We do this a few times in this test to account for CI machine
	// slowness.
	time.Sleep(200 * time.Millisecond)
	taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	}})
	require.NoError(err)
	task = taskResp.Task
	require.Equal(pb.Task_STARTING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STARTED
	time.Sleep(200 * time.Millisecond)
	taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	}})
	require.NoError(err)
	task = taskResp.Task
	require.Equal(pb.Task_STARTED, task.JobState)

	// Wait for assignment and ack
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.Equal(runnerId, assignment.Assignment.Job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(runnerId, assignment.Assignment.Job.AssignedRunner.Id)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

		// task should be RUNNING
		time.Sleep(200 * time.Millisecond)
		taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: queueResp.JobId,
			},
		}})
		require.NoError(err)
		task = taskResp.Task
		require.Equal(pb.Task_RUNNING, task.JobState)
		require.Equal(primaryJobId, task.TaskJob.Id)
	}

	// Complete our run task job so that we can move on
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be COMPLETED
	time.Sleep(200 * time.Millisecond)
	taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	}})
	require.NoError(err)
	task = taskResp.Task
	require.Equal(pb.Task_COMPLETED, task.JobState)

	var watchId string
	{
		// Watch

		// Start a job request
		rs3, err := client.RunnerJobStream(ctx)
		require.NoError(err)
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Request_{
				Request: &pb.RunnerJobStreamRequest_Request{
					RunnerId: id,
				},
			},
		}))

		// We should get a task to stop the job
		resp, err = rs3.Recv()
		require.NoError(err)
		assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
		require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
		require.IsType(&pb.Job_WatchTask{}, assignment.Assignment.Job.Operation)
		watchId = assignment.Assignment.Job.Id

		watchTask := assignment.Assignment.Job.Operation.(*pb.Job_WatchTask).WatchTask
		require.Equal(startJobId, watchTask.StartJob.Id)

		// Ack it and complete it
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))

		// Complete our watch task job so that we can move on
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Complete_{
				Complete: &pb.RunnerJobStreamRequest_Complete{},
			},
		}))
	}

	// Stop the task

	// Start a job request
	rs3, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to stop the job
	resp, err = rs3.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StopTask{}, assignment.Assignment.Job.Operation)

	// Stop should depend on both the watch and source job
	require.Contains(assignment.Assignment.Job.DependsOn, watchId)
	require.Contains(assignment.Assignment.Job.DependsOn, primaryJobId)

	stopTask := assignment.Assignment.Job.Operation.(*pb.Job_StopTask).StopTask
	require.Equal(odr.PluginConfig, stopTask.Params.HclConfig)
	require.Equal(odr.PluginType, stopTask.Params.PluginType)

	// Ack it and complete it
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// task should be STOPPING
	time.Sleep(200 * time.Millisecond)
	taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	}})
	require.NoError(err)
	task = taskResp.Task
	require.Equal(pb.Task_STOPPING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STOPPED
	time.Sleep(200 * time.Millisecond)
	taskResp, err = client.GetTask(ctx, &pb.GetTaskRequest{Ref: &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	}})
	require.NoError(err)
	task = taskResp.Task
	require.Equal(pb.Task_STOPPED, task.JobState)
}

func TestServiceQueueJob_odr_customTask(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, impl := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create an ODR profile
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	require.NoError(err)
	odr = cfgResp.Config

	// Update the project to include ondemand runner
	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	// Create, should get an ID back
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
			OndemandRunnerTask: &pb.Job_TaskOverride{
				LaunchInfo: &pb.TaskLaunchInfo{
					OciUrl: "special",
				},
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := impl.State(ctx).JobById(ctx, queueResp.JobId, nil)
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)
	primaryJobId := job.Id

	// task should be PENDING
	task, err := impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: queueResp.JobId,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_PENDING, task.JobState)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to start the job first.
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)

	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)
	require.Equal("special", st.Info.OciUrl)
	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}
	startJobId := assignment.Assignment.Job.Id

	// Ack it and complete it
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Register our runner
	runnerId := st.Info.EnvironmentVariables["WAYPOINT_RUNNER_ID"]
	server.TestRunner(t, client, &pb.Runner{
		Id:       runnerId,
		ByIdOnly: true,
	})

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: runnerId,
			},
		},
	}))

	// task should be STARTING
	// sleep to ensure job stream request was sent, so that task state updates from
	// the JobAck. We do this a few times in this test to account for CI machine
	// slowness.
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STARTING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STARTED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STARTED, task.JobState)

	// Wait for assignment and ack
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.Equal(runnerId, assignment.Assignment.Job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(runnerId, assignment.Assignment.Job.AssignedRunner.Id)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

		// task should be RUNNING
		time.Sleep(200 * time.Millisecond)
		task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: queueResp.JobId,
			},
		})
		require.NoError(err)
		require.Equal(pb.Task_RUNNING, task.JobState)
	}

	// Complete our run task job so that we can move on
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be COMPLETED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_COMPLETED, task.JobState)

	var watchId string
	{
		// Watch

		// Start a job request
		rs3, err := client.RunnerJobStream(ctx)
		require.NoError(err)
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Request_{
				Request: &pb.RunnerJobStreamRequest_Request{
					RunnerId: id,
				},
			},
		}))

		// We should get a task to stop the job
		resp, err = rs3.Recv()
		require.NoError(err)
		assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
		require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
		require.IsType(&pb.Job_WatchTask{}, assignment.Assignment.Job.Operation)
		watchId = assignment.Assignment.Job.Id

		watchTask := assignment.Assignment.Job.Operation.(*pb.Job_WatchTask).WatchTask
		require.Equal(startJobId, watchTask.StartJob.Id)

		// Ack it and complete it
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))

		// Complete our watch task job so that we can move on
		require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Complete_{
				Complete: &pb.RunnerJobStreamRequest_Complete{},
			},
		}))
	}

	// Stop the task

	// Start a job request
	rs3, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to stop the job
	resp, err = rs3.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StopTask{}, assignment.Assignment.Job.Operation)

	// Stop should depend on both the watch and source job
	require.Contains(assignment.Assignment.Job.DependsOn, watchId)
	require.Contains(assignment.Assignment.Job.DependsOn, primaryJobId)

	stopTask := assignment.Assignment.Job.Operation.(*pb.Job_StopTask).StopTask
	require.Equal(odr.PluginConfig, stopTask.Params.HclConfig)
	require.Equal(odr.PluginType, stopTask.Params.PluginType)

	// Ack it and complete it
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// task should be STOPPING
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STOPPING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STOPPED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STOPPED, task.JobState)
}

func TestServiceQueueJob_odr_customTaskSkipOp(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, impl := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create an ODR profile
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	require.NoError(err)
	odr = cfgResp.Config

	// Update the project to include ondemand runner
	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	// Create, should get an ID back
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
			OndemandRunnerTask: &pb.Job_TaskOverride{
				LaunchInfo: &pb.TaskLaunchInfo{
					OciUrl: "special",
				},

				SkipOperation: true,
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := impl.State(ctx).JobById(ctx, queueResp.JobId, nil)
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)

	// task should be PENDING
	task, err := impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: queueResp.JobId,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_PENDING, task.JobState)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to start the job first.
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)

	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)
	require.Equal("special", st.Info.OciUrl)
	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}

	// Ack it and complete it
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// task should be STARTING
	// sleep to ensure job stream request was sent, so that task state updates from
	// the JobAck. We do this a few times in this test to account for CI machine
	// slowness.
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STARTING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STARTED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STARTED, task.JobState)

	// Wait for assignment and ack
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_WatchTask{}, assignment.Assignment.Job.Operation)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

		// task should be RUNNING
		time.Sleep(200 * time.Millisecond)
		task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: queueResp.JobId,
			},
		})
		require.NoError(err)
		require.Equal(pb.Task_RUNNING, task.JobState)
	}

	// Complete our run task job so that we can move on
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be COMPLETED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_COMPLETED, task.JobState)

	// Stop the task

	// Start a job request
	rs3, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to stop the job
	resp, err = rs3.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StopTask{}, assignment.Assignment.Job.Operation)

	stopTask := assignment.Assignment.Job.Operation.(*pb.Job_StopTask).StopTask
	require.Equal(odr.PluginConfig, stopTask.Params.HclConfig)
	require.Equal(odr.PluginType, stopTask.Params.PluginType)

	// Ack it and complete it
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// task should be STOPPING
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STOPPING, task.JobState)

	// Complete our launch task job so that we can move on
	require.NoError(rs3.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// task should be STOPPED
	time.Sleep(200 * time.Millisecond)
	task, err = impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: job.Id,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_STOPPED, task.JobState)
}

func TestServiceQueueJob_odr_default(t *testing.T, factory Factory) {
	require := require.New(t)

	ctx := context.Background()

	log := hclog.New(&hclog.LoggerOptions{
		Name:  "odr-test",
		Level: hclog.Trace,
	})

	ctx = hclog.WithContext(ctx, log)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
		Default: true,
	})

	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})

	odr = cfgResp.Config

	log.Info("test odr", "id", odr.Id)

	// Update the project to include ondemand runner
	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})

	// Create, should get an ID back
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			// Note, not setting OnDemandRunnerConfig here. This is the difference between
			// this test and the previous one.
		}),
	})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: queueResp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)

	// Register our runner
	id, _ := server.TestRunner(t, client, nil)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and ack
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)

	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)
	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)

	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}

	runnerId := st.Info.EnvironmentVariables["WAYPOINT_RUNNER_ID"]

	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Register our runner
	server.TestRunner(t, client, &pb.Runner{
		Id:       runnerId,
		ByIdOnly: true,
	})

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: runnerId,
			},
		},
	}))

	// Wait for assignment and ack
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.Equal(runnerId, assignment.Assignment.Job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}
}

func TestServiceQueueJob_odr_target_id(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	log := hclog.New(&hclog.LoggerOptions{
		Name:  "odr-test",
		Level: hclog.Trace,
	})

	ctx = hclog.WithContext(ctx, log)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create an ODR profile with target runner ID
	runnerId := "test_r"
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		Name:         "test",
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: runnerId,
				},
			},
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	odr = cfgResp.Config
	log.Info("test odr profile", "name", odr.Name)

	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	// Create and queue job
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: queueResp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)

	// Register our static runner
	server.TestRunner(t, client, &pb.Runner{Id: runnerId})

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: runnerId,
			},
		},
	}))

	// We should get a task to start the job first.
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(runnerId, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)

	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)

	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}

	// Ack and complete it
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Register on-demand runner
	odrId := st.Info.EnvironmentVariables["WAYPOINT_RUNNER_ID"]
	server.TestRunner(t, client, &pb.Runner{Id: odrId})

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: odrId,
			},
		},
	}))

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Wait for assignment and ack.
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.Equal(odrId, assignment.Assignment.Job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(odrId, assignment.Assignment.Job.AssignedRunner.Id)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}
}

func TestServiceQueueJob_odr_target_labels(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	log := hclog.New(&hclog.LoggerOptions{
		Name:  "odr-test",
		Level: hclog.Trace,
	})

	ctx = hclog.WithContext(ctx, log)

	// Create our server
	client, _ := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create an ODR profile with target runner ID
	labels := map[string]string{
		"env": "test",
	}
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		Name:         "test",
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Labels{
				Labels: &pb.Ref_RunnerLabels{
					Labels: labels,
				},
			},
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	odr = cfgResp.Config
	log.Info("test odr profile", "name", odr.Name)

	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	// Create and queue job
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Job should exist and be queued
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: queueResp.JobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)

	// Register our static runner
	id, _ := server.TestRunner(t, client, &pb.Runner{Labels: labels})

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// We should get a task to start the job first.
	resp, err := runnerStream.Recv()
	require.NoError(err)
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(id, assignment.Assignment.Job.AssignedRunner.Id)
	require.NotEqual(queueResp.JobId, assignment.Assignment.Job.Id)
	require.IsType(&pb.Job_StartTask{}, assignment.Assignment.Job.Operation)

	st := assignment.Assignment.Job.Operation.(*pb.Job_StartTask).StartTask
	require.Equal(odr.PluginConfig, st.Params.HclConfig)
	require.Equal(odr.PluginType, st.Params.PluginType)

	for k, v := range odr.EnvironmentVariables {
		require.Equal(v, st.Info.EnvironmentVariables[k])
	}

	// Ack and complete it
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Register on-demand runner
	odrId := st.Info.EnvironmentVariables["WAYPOINT_RUNNER_ID"]
	server.TestRunner(t, client, &pb.Runner{Id: odrId})

	// Start a job request
	rs2, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: odrId,
			},
		},
	}))

	// Complete our launch task job so that we can move on
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Wait for assignment and ack.
	resp, err = rs2.Recv()
	require.NoError(err)
	assignment, ok = resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(ok, "should be an assignment")
	require.NotNil(assignment)
	require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)
	require.Equal(odrId, assignment.Assignment.Job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(odrId, assignment.Assignment.Job.AssignedRunner.Id)

	require.NoError(rs2.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}
}

// Test that a job with dependencies has an ODR started that depends on
// the same things.
func TestServiceQueueJob_odr_depends(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, impl := factory(t)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "app",
			Project:     "proj",
		},
	}).Application)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	// Create an ODR profile
	odr := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		PluginType:   "magic-carpet",
		PluginConfig: []byte("foo = 1"),
		EnvironmentVariables: map[string]string{
			"CARPET_DRIVER": "apu",
		},
	})
	cfgResp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: odr,
	})
	odr = cfgResp.Config

	// Update the project to include ondemand runner
	proj := serverptypes.TestProject(t, &pb.Project{Name: "proj"})
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	require.NoError(err)

	queueA, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
		}),
	})
	require.NoError(err)
	queueB, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
		}),
	})
	require.NoError(err)

	// Create, should get an ID back
	queueResp, err := client.QueueJob(ctx, &Req{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "proj",
			},
			DependsOn:             []string{queueA.JobId, queueB.JobId},
			DependsOnAllowFailure: []string{queueB.JobId},
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{
				Name: odr.Name,
			},
		}),
	})
	require.NoError(err)
	require.NotEmpty(queueResp)

	// Get the task to get the start job
	task, err := impl.State(ctx).TaskGet(ctx, &pb.Ref_Task{
		Ref: &pb.Ref_Task_JobId{
			JobId: queueResp.JobId,
		},
	})
	require.NoError(err)
	require.Equal(pb.Task_PENDING, task.JobState)

	// Get the start job
	job, err := impl.State(ctx).JobById(ctx, task.StartJob.Id, nil)
	require.NoError(err)
	require.IsType(&pb.Job_StartTask{}, job.Operation)
	require.Equal([]string{queueA.JobId, queueB.JobId}, job.DependsOn)
	require.Equal([]string{queueB.JobId}, job.DependsOnAllowFailure)
}
