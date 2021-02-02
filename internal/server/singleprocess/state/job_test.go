package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestJobCreate_singleton(t *testing.T) {
	t.Run("create with no existing", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		})))

		// Exactly one job should exist
		jobs, err := s.JobList()
		require.NoError(err)
		require.Len(jobs, 1)
	})

	t.Run("create with an existing queued", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		})))

		// Create a different job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "B",
			SingletonId: "1",
		})))

		// Should have both jobs
		jobs, err := s.JobList()
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be canceled
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_ERROR, job.State)

			oldQueueTime, err = ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should be that of the old job, so that
			// we retain our position in the queue
			queueTime, err := ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
			require.False(oldQueueTime.IsZero())
			require.True(oldQueueTime.Equal(queueTime))
		}
	})

	t.Run("create with an existing complete", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		})))

		// Assign and complete A
		{
			job, err := s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}

		// Create a different job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "B",
			SingletonId: "1",
		})))

		// Should have both jobs
		jobs, err := s.JobList()
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be done
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_SUCCESS, job.State)

			oldQueueTime, err = ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should NOT be equal
			queueTime, err := ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
			require.False(queueTime.IsZero())
			require.False(oldQueueTime.Equal(queueTime))
		}
	})

	t.Run("create with an existing assigned", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := TestState(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		})))

		// Assign and complete A
		{
			job, err := s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)

			// Do NOT ack, do not complete, etc.
		}

		// Create a different job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "B",
			SingletonId: "1",
		})))

		// Should have both jobs
		jobs, err := s.JobList()
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be done
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_WAITING, job.State)

			oldQueueTime, err = ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should NOT be equal
			queueTime, err := ptypes.Timestamp(job.QueueTime)
			require.NoError(err)
			require.False(queueTime.IsZero())
			require.False(oldQueueTime.Equal(queueTime))
		}
	})
}

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
		require.Nil(job)
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

	t.Run("assignment by ID no candidates", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build by ID
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{
						Id: "R_B",
					},
				},
			},
		})))

		// Assign for R_A which should get nothing cause it doesn't match.
		// NOTE that using "R_A" here is very important. This fixes a bug
		// where our lower bound was picking up invalid IDs.
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			}()

			// We should be blocking
			select {
			case <-doneCh:
				t.Fatal("should wait")

			case <-time.After(500 * time.Millisecond):
			}
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
		require.Nil(job)
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
		require.NoError(s.JobCancel("A", false))

		// Verify it is canceled
		job, err := s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
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
		require.NoError(s.JobCancel("A", false))

		// Verify it is canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_WAITING, job.Job.State)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("assigned with force", func(t *testing.T) {
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
		require.NoError(s.JobCancel("A", true))

		// Verify it is canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("assigned with force clears assignedSet", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "A",
			Operation: &pb.Job_Deploy{},
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		require.Equal(pb.Job_WAITING, job.State)

		// Cancel it
		require.NoError(s.JobCancel("A", true))

		// Verify it is canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotEmpty(job.CancelTime)

		// Create a another job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			Operation: &pb.Job_Deploy{},
		})))

		ws := memdb.NewWatchSet()

		// Read it back to check the blocked status
		job2, err := s.JobById("B", ws)
		require.NoError(err)
		require.NotNil(job2)
		require.Equal("B", job2.Id)
		require.Equal(pb.Job_QUEUED, job2.State)
		require.False(job2.Blocked)
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
		require.NoError(s.JobCancel("A", false))

		// Verify it is not canceled
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_SUCCESS, job.Job.State)
		require.Empty(job.CancelTime)
	})
}

func TestJobHeartbeat(t *testing.T) {
	t.Run("times out after ack", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set a short timeout
		old := jobHeartbeatTimeout
		defer func() { jobHeartbeatTimeout = old }()
		jobHeartbeatTimeout = 5 * time.Millisecond

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

		// Should time out
		require.Eventually(func() bool {
			// Verify it is canceled
			job, err = s.JobById("A", nil)
			require.NoError(err)
			return job.Job.State == pb.Job_ERROR
		}, 1*time.Second, 10*time.Millisecond)
	})

	t.Run("doesn't time out if heartbeating", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := jobHeartbeatTimeout
		defer func() { jobHeartbeatTimeout = old }()
		jobHeartbeatTimeout = 250 * time.Millisecond

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

		// Start heartbeating
		ctx, cancel := context.WithCancel(context.Background())
		doneCh := make(chan struct{})
		defer func() {
			cancel()
			<-doneCh
		}()
		go func() {
			defer close(doneCh)

			tick := time.NewTicker(20 * time.Millisecond)
			defer tick.Stop()

			for {
				select {
				case <-tick.C:
					s.JobHeartbeat(job.Id)

				case <-ctx.Done():
					return
				}
			}
		}()

		// Sleep for a bit
		time.Sleep(1 * time.Second)

		// Verify it is running
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// Stop it
		require.NoError(s.JobComplete(job.Id, nil, nil))
	})

	t.Run("times out if heartbeating stops", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := jobHeartbeatTimeout
		defer func() { jobHeartbeatTimeout = old }()
		jobHeartbeatTimeout = 250 * time.Millisecond

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

		// Start heartbeating
		ctx, cancel := context.WithCancel(context.Background())
		doneCh := make(chan struct{})
		defer func() {
			cancel()
			<-doneCh
		}()
		go func() {
			defer close(doneCh)

			tick := time.NewTicker(20 * time.Millisecond)
			defer tick.Stop()

			for {
				select {
				case <-tick.C:
					s.JobHeartbeat(job.Id)

				case <-ctx.Done():
					return
				}
			}
		}()

		// Sleep for a bit
		time.Sleep(1 * time.Second)

		// Verify it is running
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// Stop heartbeating
		cancel()

		// Should time out
		require.Eventually(func() bool {
			// Verify it is canceled
			job, err = s.JobById("A", nil)
			require.NoError(err)
			return job.Job.State == pb.Job_ERROR
		}, 1*time.Second, 10*time.Millisecond)
	})

	t.Run("times out if running state loaded on restart", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := jobHeartbeatTimeout
		defer func() { jobHeartbeatTimeout = old }()
		jobHeartbeatTimeout = 250 * time.Millisecond

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

		// Start heartbeating
		ctx, cancel := context.WithCancel(context.Background())
		doneCh := make(chan struct{})
		defer func() {
			cancel()
			<-doneCh
		}()
		go func(s *State) {
			defer close(doneCh)

			tick := time.NewTicker(20 * time.Millisecond)
			defer tick.Stop()

			for {
				select {
				case <-tick.C:
					s.JobHeartbeat(job.Id)

				case <-ctx.Done():
					return
				}
			}
		}(s)

		// Reinit the state as if we crashed
		s = TestStateReinit(t, s)
		defer s.Close()

		// Verify it exists
		job, err = s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// Should time out
		require.Eventually(func() bool {
			// Verify it is canceled
			job, err = s.JobById("A", nil)
			require.NoError(err)
			return job.Job.State == pb.Job_ERROR
		}, 2*time.Second, 10*time.Millisecond)
	})
}

func TestJobUpdateRef(t *testing.T) {
	t.Run("running", func(t *testing.T) {
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

		// Create a watchset on the job
		ws := memdb.NewWatchSet()

		// Verify it was changed
		job, err = s.JobById(job.Id, ws)
		require.NoError(err)
		require.Nil(job.DataSourceRef)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Update the ref
		require.NoError(s.JobUpdateRef(job.Id, &pb.Job_DataSource_Ref{
			Ref: &pb.Job_DataSource_Ref_Git{
				Git: &pb.Job_Git_Ref{
					Commit: "hello",
				},
			},
		}))

		// Should be triggered. This is a very important test because
		// we need to ensure that the watchers can detect ref changes.
		require.False(ws.Watch(time.After(100 * time.Millisecond)))

		// Verify it was changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.NotNil(job.DataSourceRef)

		ref := job.DataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git).Git
		require.Equal(ref.Commit, "hello")
	})
}
