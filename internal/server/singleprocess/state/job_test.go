package state

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestJobAssign(t *testing.T) {
	t.Run("basic assignment with one", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		require.Equal(pb.Job_WAITING, job.State)

		// Should block if requesting another since none exist
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		job, err = s.JobAssignForRunner(ctx, &Runner{Id: "R_A"})
		require.Error(err)
		require.Equal(ctx.Err(), err)
	})

	t.Run("blocking on any", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		{
			job, err := s.JobAssignForRunner(context.Background(), &Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}

		// Get the next value in a goroutine
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var job *pb.Job
			var jerr error
			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				job, jerr = s.JobAssignForRunner(ctx, &Runner{Id: "R_A"})
			}()

			// We should be blocking
			select {
			case <-doneCh:
				t.Fatal("should wait")

			case <-time.After(50 * time.Millisecond):
			}

			// Insert another job
			require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
				Id: "B",
			})))

			// We should get a result
			select {
			case <-doneCh:

			case <-time.After(50 * time.Millisecond):
				t.Fatal("should have a result")
			}

			require.NoError(jerr)
			require.NotNil(job)
			require.Equal("B", job.Id)
		}
	})

	t.Run("basic assignment with two", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Create two builds slightly apart
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		// Assign it, we should get build A then B
		{
			job, err := s.JobAssignForRunner(ctx, &Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}
		{
			job, err := s.JobAssignForRunner(ctx, &Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
		}
	})

	t.Run("assignment by ID", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build by ID
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: "R_A",
					},
				},
			},
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		})))

		// Assign for R_B, which should get B since it won't match the earlier
		// assignment target.
		{
			job, err := s.JobAssignForRunner(context.Background(), &Runner{Id: "R_B"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
		}

		// Assign for R_A, which should get A since it matches the target.
		{
			job, err := s.JobAssignForRunner(context.Background(), &Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}
	})
}
