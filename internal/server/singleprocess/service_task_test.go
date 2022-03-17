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

func TestServiceTask(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	type Req = pb.UpsertTaskRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: serverptypes.TestValidTask(t, &pb.Task{TaskJob: &pb.Ref_Job{Id: "run_job"}}),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Task
		require.NotEmpty(result.Id)

		// Let's write some data
		result.StartJob = &pb.Ref_Job{Id: "start_job"}
		result.StopJob = &pb.Ref_Job{Id: "stop_job"}
		resp, err = client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Task
		require.Equal(result.StartJob.Id, "start_job")
	})

	t.Run("update non-existent creates a new task", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTask(ctx, &Req{
			Task: serverptypes.TestValidTask(t, &pb.Task{
				Id:      "newone",
				TaskJob: &pb.Ref_Job{Id: "newone"},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(resp.Task.Id, "newone")
	})
}

func TestServiceTask_GetTask(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
	}).Application)

	// Create, should get an ID back
	jobResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	startJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	runJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	stopJobId := jobResp.JobId

	resp, err := client.UpsertTask(ctx, &pb.UpsertTaskRequest{
		Task: serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	})
	require.NoError(t, err)
	taskId := resp.Task.Id

	// Create, should get an ID back
	t.Run("get existing by task id", func(t *testing.T) {
		require := require.New(t)

		type JobReq = pb.QueueJobRequest

		// Get, should return a task
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: resp.Task.Id,
				},
			},
		})
		require.NoError(err)
		require.NotNil(resp.Task)
		require.NotEmpty(resp.Task.Id)
		require.Equal(taskId, resp.Task.Id)
	})

	t.Run("get existing by run job id", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a task
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: runJobId,
				},
			},
		})
		require.NoError(err)
		require.NotNil(resp.Task)
		require.NotEmpty(resp.Task.Id)
		require.Equal(taskId, resp.Task.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "nope",
				},
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("get non-existing by job id", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: "nope",
				},
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceTask_ListTaskSimple(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
	}).Application)

	// Create, should get an ID back
	jobResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	startJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	runJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	stopJobId := jobResp.JobId

	_, err = client.UpsertTask(ctx, &pb.UpsertTaskRequest{
		Task: serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	})
	_, err = client.UpsertTask(ctx, &pb.UpsertTaskRequest{
		Task: serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	})
	_, err = client.UpsertTask(ctx, &pb.UpsertTaskRequest{
		Task: serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	})

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTask(ctx, &pb.ListTaskRequest{})
		require.NoError(err)
		require.Equal(len(respList.Tasks), 3)
	})
}

func TestServiceTask_DeleteTask(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Initialize our app
	TestApp(t, client, serverptypes.TestJobNew(t, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
	}).Application)

	// Create, should get an ID back
	jobResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	startJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	runJobId := jobResp.JobId

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	stopJobId := jobResp.JobId

	resp, err := client.UpsertTask(ctx, &pb.UpsertTaskRequest{
		Task: serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	})
	taskId := resp.Task.Id

	t.Run("get existing then delete", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a task
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: resp.Task.Id,
				},
			},
		})
		require.NoError(err)
		require.NotNil(resp.Task)
		require.NotEmpty(resp.Task.Id)
		require.Equal(taskId, resp.Task.Id)

		_, err = client.DeleteTask(ctx, &pb.DeleteTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: taskId,
				},
			},
		})
		require.NoError(err)

		// get, should fail
		resp, err = client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: taskId,
				},
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("delete non-existing", func(t *testing.T) {
		require := require.New(t)

		resp, err := client.DeleteTask(ctx, &pb.DeleteTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "nope",
				},
			},
		})
		require.NoError(err)
		require.NotNil(resp)
	})
}
