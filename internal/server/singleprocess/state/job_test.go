package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		require.Equal(pb.Job_WAITING, job.State)

		// We should not have an output buffer yet
		require.Nil(job.OutputBuffer)

		// Should block if requesting another since none exist
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		job, err = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
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
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
		})))

		// Assign it, we should get this build
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}

		// Get the next value in a goroutine
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var job *Job
			var jerr error
			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				job, jerr = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			}()

			// We should be blocking
			select {
			case <-doneCh:
				t.Fatal("should wait")

			case <-time.After(500 * time.Millisecond):
			}

			// Insert another job
			require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
				Id: "B",
				Workspace: &pb.Ref_Workspace{
					Workspace: "w2",
				},
			})))

			// We should get a result
			select {
			case <-doneCh:

			case <-time.After(500 * time.Millisecond):
				t.Fatal("should have a result")
			}

			require.NoError(jerr)
			require.NotNil(job)
			require.Equal("B", job.Id)
		}
	})

	t.Run("blocking on matching app and workspace", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create two builds for the same app/workspace
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
			Operation: &pb.Job_Deploy{
				Deploy: &pb.Job_DeployOp{},
			},
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
			Operation: &pb.Job_Deploy{
				Deploy: &pb.Job_DeployOp{},
			},
		})))

		// Assign it, we should get this build
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}

		// Get the next value in a goroutine
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var job *Job
			var jerr error
			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				job, jerr = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			}()

			// We should be blocking
			select {
			case <-doneCh:
				t.Fatal("should wait")

			case <-time.After(500 * time.Millisecond):
			}

			// Insert another job for a different workspace
			require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
				Id: "C",
				Workspace: &pb.Ref_Workspace{
					Workspace: "w2",
				},
				Operation: &pb.Job_Deploy{
					Deploy: &pb.Job_DeployOp{},
				},
			})))

			// We should get a result
			select {
			case <-doneCh:

			case <-time.After(500 * time.Millisecond):
				t.Fatal("should have a result")
			}

			require.NoError(jerr)
			require.NotNil(job)
			require.Equal("C", job.Id)
		}
	})

	t.Run("blocking on matching app and workspace (sequential)", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create two builds for the same app/workspace
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
			Operation: &pb.Job_Deploy{
				Deploy: &pb.Job_DeployOp{},
			},
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
			Operation: &pb.Job_Deploy{
				Deploy: &pb.Job_DeployOp{},
			},
		})))

		// Assign it, we should get this build
		job1, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job1)
		require.Equal("A", job1.Id)

		// Get the next value in a goroutine
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var job *Job
			var jerr error
			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				job, jerr = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			}()

			// We should be blocking
			select {
			case <-doneCh:
				t.Fatal("should wait")

			case <-time.After(500 * time.Millisecond):
			}

			// Complete the job
			_, err = s.JobAck(job1.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job1.Id, nil, nil))

			// We should get a result
			select {
			case <-doneCh:

			case <-time.After(500 * time.Millisecond):
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
			job, err := s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}
		{
			job, err := s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
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
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_B"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}

		// Assign for R_A, which should get A since it matches the target.
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}
	})

	t.Run("any cannot be assigned to ByIdOnly runner", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		r := &pb.Runner{Id: "R_A", ByIdOnly: true}

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Should block because none direct assign
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		job, err := s.JobAssignForRunner(ctx, r)
		require.Error(err)
		require.Equal(ctx.Err(), err)

		// Create a target
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: "R_A",
					},
				},
			},
		})))

		// Assign it, we should get this build
		job, err = s.JobAssignForRunner(context.Background(), r)
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)
	})
}

func TestJobAck(t *testing.T) {
	t.Run("ack", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// We should have an output buffer
		require.NotNil(job.OutputBuffer)
	})

	t.Run("ack negative", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, false)
		require.NoError(err)

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.State)

		// We should not have an output buffer
		require.Nil(job.OutputBuffer)
	})

	t.Run("timeout before ack should requeue", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := jobWaitingTimeout
		defer func() { jobWaitingTimeout = old }()
		jobWaitingTimeout = 5 * time.Millisecond

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Sleep too long
		time.Sleep(100 * time.Millisecond)

		// Verify it is queued
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.Job.State)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.Error(err)
	})
}

func TestJobComplete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it
		require.NoError(s.JobComplete(job.Id, &pb.Job_Result{
			Build: &pb.Job_BuildResult{},
		}, nil))

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_SUCCESS, job.State)
		require.Nil(job.Error)
		require.NotNil(job.Result)
		require.NotNil(job.Result.Build)
	})

	t.Run("error", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it
		require.NoError(s.JobComplete(job.Id, nil, fmt.Errorf("bad")))

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)

		st := status.FromProto(job.Error)
		require.Equal(codes.Unknown, st.Code())
		require.Contains(st.Message(), "bad")
	})
}

func TestJobIsAssignable(t *testing.T) {
	t.Run("no runners", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Create a build
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		}))
		require.NoError(err)
		require.False(result)
	})

	t.Run("any target, runners exist", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Register a runner
		require.NoError(s.RunnerCreate(serverptypes.TestRunner(t, nil)))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		}))
		require.NoError(err)
		require.True(result)
	})

	t.Run("any target, runners ByIdOnly", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Register a runner
		require.NoError(s.RunnerCreate(serverptypes.TestRunner(t, &pb.Runner{
			ByIdOnly: true,
		})))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		}))
		require.NoError(err)
		require.False(result)
	})

	t.Run("ID target, no match", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Register a runner
		require.NoError(s.RunnerCreate(serverptypes.TestRunner(t, nil)))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: "R_A",
					},
				},
			},
		}))
		require.NoError(err)
		require.False(result)
	})

	t.Run("ID target, match", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Register a runner
		runner := serverptypes.TestRunner(t, nil)
		require.NoError(s.RunnerCreate(runner))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: runner.Id,
					},
				},
			},
		}))
		require.NoError(err)
		require.True(result)
	})
}

func TestJobCancel(t *testing.T) {
	t.Run("queued", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Cancel it
		require.NoError(s.JobCancel("A"))

		// Verify it is canceled
		job, err := s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_SUCCESS, job.Job.State)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("assigned", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		require.Equal(pb.Job_WAITING, job.State)

		// Cancel it
		require.NoError(s.JobCancel("A"))

		// Verify it is canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_WAITING, job.Job.State)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("completed", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		require.Equal(pb.Job_WAITING, job.State)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it
		require.NoError(s.JobComplete(job.Id, nil, nil))

		// Cancel it
		require.NoError(s.JobCancel("A"))

		// Verify it is not canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_SUCCESS, job.Job.State)
		require.Empty(job.CancelTime)
	})
}
