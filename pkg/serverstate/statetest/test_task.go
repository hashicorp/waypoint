package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["task"] = []testFunc{
		TestTask,
		TestTaskCancel,
	}
}

func TestTask(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		_, err := s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: "foo",
			},
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get by Task Id and Ref Job Id", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.TaskPut(&pb.Task{
			Id: "t_test",
		})
		require.Error(err) // no job id set
		err = nil

		// Set again
		err = s.TaskPut(&pb.Task{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get exact by job id
		{
			resp, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: "j_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Update
		err = s.TaskPut(&pb.Task{
			Id:       "t_test",
			TaskJob:  &pb.Ref_Job{Id: "j_test"},
			StartJob: &pb.Ref_Job{Id: "start_job"},
			StopJob:  &pb.Ref_Job{Id: "stop_job"},
		})
		require.NoError(err)

		// Get exact by id
		{
			resp, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.StartJob.Id, "start_job")
		}
	})

	t.Run("Deletion by Task Id and Job Id Ref", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.TaskPut(&pb.Task{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete it
		err = s.TaskDelete(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: "t_test",
			},
		})
		require.NoError(err)

		// It's gone
		{
			_, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: "t_test",
				},
			})
			require.Error(err)
		}
		err = nil

		// Set again
		err = s.TaskPut(&pb.Task{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get job id

		// Get exact by job id
		{
			resp, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: "j_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete it
		err = s.TaskDelete(&pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: "j_test",
			},
		})
		require.NoError(err)

		// It's gone
		{
			_, err := s.TaskGet(&pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: "j_test",
				},
			})
			require.Error(err)
		}
	})

	t.Run("Listing", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create more for listing
		err := s.TaskPut(&pb.Task{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)

		err = s.TaskPut(&pb.Task{
			Id:      "t_test_part2",
			TaskJob: &pb.Ref_Job{Id: "j2_test"},
		})
		require.NoError(err)

		err = s.TaskPut(&pb.Task{
			Id:      "t_test_part3",
			TaskJob: &pb.Ref_Job{Id: "j3_test"},
		})
		require.NoError(err)

		// List all
		{
			resp, err := s.TaskList(&pb.ListTaskRequest{})
			require.NoError(err)
			require.Len(resp, 3)
		}
	})

	t.Run("Getting task jobs by task ref", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "start_job",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "j_test",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "stop_job",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "watch_job",
		})))

		err := s.TaskPut(&pb.Task{
			Id:       "t_test",
			TaskJob:  &pb.Ref_Job{Id: "j_test"},
			StartJob: &pb.Ref_Job{Id: "start_job"},
			StopJob:  &pb.Ref_Job{Id: "stop_job"},
			WatchJob: &pb.Ref_Job{Id: "watch_job"},
		})
		require.NoError(err)

		task, err := s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: "j_test",
			},
		})
		require.NoError(err)
		require.NotNil(task)

		startJob, taskJob, stopJob, watchJob, err := s.JobsByTaskRef(task)
		require.NoError(err)
		require.Equal(startJob.Id, "start_job")
		require.Equal(taskJob.Id, "j_test")
		require.Equal(stopJob.Id, "stop_job")
		require.Equal(watchJob.Id, "watch_job")
	})
}

func TestTaskCancel(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("Canceling a Task by task id", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "start_job",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "j_test",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "watch_job",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "stop_job",
		})))

		err := s.TaskPut(&pb.Task{
			Id:       "t_test",
			TaskJob:  &pb.Ref_Job{Id: "j_test"},
			StartJob: &pb.Ref_Job{Id: "start_job"},
			WatchJob: &pb.Ref_Job{Id: "watch_job"},
			StopJob:  &pb.Ref_Job{Id: "stop_job"},
		})
		require.NoError(err)

		task, err := s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: "j_test",
			},
		})
		require.NoError(err)
		require.NotNil(task)

		require.NoError(s.TaskCancel(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: "t_test",
			},
		}))

		// Verify it is canceled
		job, err := s.JobById("start_job", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		job, err = s.JobById("j_test", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		job, err = s.JobById("watch_job", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		job, err = s.JobById("stop_job", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)
	})
}
