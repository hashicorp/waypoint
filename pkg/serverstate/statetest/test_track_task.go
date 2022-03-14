package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["tracktask"] = []testFunc{
		TestTrackTask,
	}
}

func TestTrackTask(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		_, err := s.TrackTaskGet(&pb.Ref_TrackTask{
			Ref: &pb.Ref_TrackTask_Id{
				Id: "foo",
			},
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get by TrackTask Id", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.TrackTaskPut(&pb.TrackTask{
			Id: "t_test",
		})
		require.Error(err) // no job id set
		err = nil

		// Set again
		err = s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get exact by job id
		{
			resp, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_JobId{
					JobId: "j_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Update
		err = s.TrackTaskPut(&pb.TrackTask{
			Id:       "t_test",
			TaskJob:  &pb.Ref_Job{Id: "j_test"},
			StartJob: &pb.Ref_Job{Id: "start_job"},
			StopJob:  &pb.Ref_Job{Id: "stop_job"},
		})
		require.NoError(err)

		// Get exact by id
		{
			resp, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.StartJob.Id, "start_job")
		}
	})

	t.Run("Deletion by TrackTask Id and Job Id Ref", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_Id{
					Id: "t_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete it
		err = s.TrackTaskDelete(&pb.Ref_TrackTask{
			Ref: &pb.Ref_TrackTask_Id{
				Id: "t_test",
			},
		})
		require.NoError(err)

		// It's gone
		{
			_, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_Id{
					Id: "t_test",
				},
			})
			require.Error(err)
		}
		err = nil

		// Set again
		err = s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)
		// Get job id

		// Get exact by job id
		{
			resp, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_JobId{
					JobId: "j_test",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete it
		err = s.TrackTaskDelete(&pb.Ref_TrackTask{
			Ref: &pb.Ref_TrackTask_JobId{
				JobId: "j_test",
			},
		})
		require.NoError(err)

		// It's gone
		{
			_, err := s.TrackTaskGet(&pb.Ref_TrackTask{
				Ref: &pb.Ref_TrackTask_JobId{
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
		err := s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test",
			TaskJob: &pb.Ref_Job{Id: "j_test"},
		})
		require.NoError(err)

		err = s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test_part2",
			TaskJob: &pb.Ref_Job{Id: "j2_test"},
		})
		require.NoError(err)

		err = s.TrackTaskPut(&pb.TrackTask{
			Id:      "t_test_part3",
			TaskJob: &pb.Ref_Job{Id: "j3_test"},
		})
		require.NoError(err)

		// List all
		{
			resp, err := s.TrackTaskList()
			require.NoError(err)
			require.Len(resp, 3)
		}
	})
}
