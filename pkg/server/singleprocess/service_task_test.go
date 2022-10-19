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

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	watchJobId := jobResp.JobId

	taskId := "test_task"
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			Id:       taskId,
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
		}),
	)
	require.NoError(t, err)

	// Create, should get an ID back
	t.Run("get existing by task id", func(t *testing.T) {
		require := require.New(t)

		type JobReq = pb.QueueJobRequest

		// Get, should return a task
		resp, err := client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: taskId,
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

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	watchJobId := jobResp.JobId

	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
		}),
	)
	require.NoError(t, err)
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
		}),
	)
	require.NoError(t, err)
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
		}),
	)
	require.NoError(t, err)

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTask(ctx, &pb.ListTaskRequest{})
		require.NoError(err)
		require.Equal(len(respList.Tasks), 3)

		for _, t := range respList.Tasks {
			run := false
			if t.TaskJob.Id == runJobId {
				run = true
			}
			require.True(run)
		}
	})
}

func TestServiceTask_ListTaskFilters(t *testing.T) {
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

	jobResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, jobResp)
	require.NotEmpty(t, jobResp.JobId)
	watchJobId := jobResp.JobId

	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
			JobState: pb.Task_STOPPED,
		}),
	)
	require.NoError(t, err)
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
			JobState: pb.Task_RUNNING,
		}),
	)
	require.NoError(t, err)
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
			WatchJob: &pb.Ref_Job{Id: watchJobId},
		}),
	)
	require.NoError(t, err)

	t.Run("list filter on job state", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTask(ctx, &pb.ListTaskRequest{
			TaskState: []pb.Task_State{pb.Task_STOPPED},
		})
		require.NoError(err)
		require.Equal(len(respList.Tasks), 1)

		respList, err = client.ListTask(ctx, &pb.ListTaskRequest{
			TaskState: []pb.Task_State{pb.Task_STOPPED, pb.Task_RUNNING},
		})
		require.NoError(err)
		require.Equal(len(respList.Tasks), 2)

	})
}

func TestServiceTask_CancelTask(t *testing.T) {
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

	taskId := "test_task"
	err = testServiceImpl(impl).state(ctx).TaskPut(ctx,
		serverptypes.TestValidTask(t, &pb.Task{
			Id:       taskId,
			TaskJob:  &pb.Ref_Job{Id: runJobId},
			StartJob: &pb.Ref_Job{Id: startJobId},
			StopJob:  &pb.Ref_Job{Id: stopJobId},
		}),
	)
	require.NoError(t, err)

	t.Run("cancel existing by task id", func(t *testing.T) {
		require := require.New(t)

		type JobReq = pb.QueueJobRequest

		_, err = client.CancelTask(ctx, &pb.CancelTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: taskId,
				},
			},
		})
		require.NoError(err)

		job, err := testServiceImpl(impl).state(ctx).JobById(startJobId, nil)
		require.NoError(err)
		require.True(job.State == pb.Job_ERROR && job.CancelTime != nil)

		job, err = testServiceImpl(impl).state(ctx).JobById(runJobId, nil)
		require.NoError(err)
		require.True(job.State == pb.Job_ERROR && job.CancelTime != nil)

		job, err = testServiceImpl(impl).state(ctx).JobById(stopJobId, nil)
		require.NoError(err)
		require.True(job.State == pb.Job_ERROR && job.CancelTime != nil)
	})
}
