package statetest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/go-memdb"

	"github.com/hashicorp/waypoint/pkg/serverstate"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["job"] = []testFunc{
		TestJobCreate_singleton,
		TestJobAssign,
		TestJobAck,
		TestJobComplete,
		TestJobTask_AckAndComplete,
		TestJobPipeline_AckAndComplete,
		TestJobIsAssignable,
		TestJobCancel,
		TestJobHeartbeat,
		TestJobHeartbeatOnRestart,
		TestJobUpdateRef,
		TestJobUpdateExpiry,
	}
}

func TestJobCreate_singleton(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("create with no existing", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		})))

		// Exactly one job should exist
		jobs, err := s.JobList(&pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobs, 1)
	})

	t.Run("create multiple with no existing", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:          "A",
			SingletonId: "1",
		}), serverptypes.TestJobNew(t, &pb.Job{
			Id:          "B",
			SingletonId: "2",
		})))

		// Exactly one job should exist
		jobs, err := s.JobList(&pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobs, 2)
	})

	t.Run("create with an existing queued", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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
		jobs, err := s.JobList(&pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be canceled
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_ERROR, job.State)

			oldQueueTime = job.QueueTime.AsTime()
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should be that of the old job, so that
			// we retain our position in the queue
			queueTime := job.QueueTime.AsTime()
			require.False(oldQueueTime.IsZero())
			require.True(oldQueueTime.Equal(queueTime))
		}
	})

	t.Run("create with an existing complete", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
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
		jobs, err := s.JobList(&pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be done
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_SUCCESS, job.State)

			oldQueueTime = job.QueueTime.AsTime()
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should NOT be equal
			queueTime := job.QueueTime.AsTime()
			require.NoError(err)
			require.False(queueTime.IsZero())
			require.False(oldQueueTime.Equal(queueTime))
		}
	})

	t.Run("create with an existing assigned", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
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
		jobs, err := s.JobList(&pb.ListJobsRequest{})
		require.NoError(err)
		require.Len(jobs, 2)

		// Job "A" should be done
		var oldQueueTime time.Time
		{
			job, err := s.JobById("A", nil)
			require.NoError(err)
			require.Equal(pb.Job_WAITING, job.State)

			oldQueueTime = job.QueueTime.AsTime()
		}

		// Job "B" should be queued
		{
			job, err := s.JobById("B", nil)
			require.NoError(err)
			require.Equal(pb.Job_QUEUED, job.State)

			// The queue time should NOT be equal
			queueTime := job.QueueTime.AsTime()
			require.NoError(err)
			require.False(queueTime.IsZero())
			require.False(oldQueueTime.Equal(queueTime))
		}
	})

	t.Run("create with a dependency cycle", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		err := s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "A",
			DependsOn: []string{"C"},
		}), serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		}), serverptypes.TestJobNew(t, &pb.Job{
			Id:        "C",
			DependsOn: []string{"B"},
		}))
		require.Error(err)
		require.Contains(err.Error(), "cycle")
	})

	t.Run("create in non-dependency order", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		err := s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "A",
			DependsOn: []string{"C"},
		}), serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		}), serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		}))
		require.NoError(err)
	})
}

func TestJobAssign(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("basic assignment with one", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Should get a peeked job
		job, err := s.JobPeekForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)

		// Assign it, we should get this build
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
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

		// Should not block if requested
		job, err = s.JobPeekForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.Nil(job)
	})

	t.Run("blocking on any", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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
			var job *serverstate.Job
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

			case <-time.After(3 * time.Second):
				t.Fatal("should have a result")
			}

			require.NoError(jerr)
			require.NotNil(job)
			require.Equal("B", job.Id)
		}
	})

	t.Run("blocking on matching app and workspace", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

		// Peeking still returns blocked jobs. We do this because otherwise
		// its possible for peek to return empty when there is an eligible job,
		// its just waiting.
		job, err := s.JobPeekForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)

		// Get the next value in a goroutine
		{
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var job *serverstate.Job
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

			case <-time.After(3 * time.Second):
				t.Fatal("should have a result")
			}

			require.NoError(jerr)
			require.NotNil(job)
			require.Equal("C", job.Id)
		}
	})

	t.Run("blocking on matching app and workspace (sequential)", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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
			var job *serverstate.Job
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

			case <-time.After(3 * time.Second):
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

		s := factory(t)
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

		s := factory(t)
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

		s := factory(t)
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

	t.Run("assignment by Labels", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build by labels
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"env": "test",
						},
					},
				},
			},
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"region": "testland-1",
						},
					},
				},
			},
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"useless": "label",
						},
					},
				},
			},
		})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{Id: "D"})))
		time.Sleep(1 * time.Millisecond)
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "E",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"region": "testland-2",
						},
					},
				},
			},
		})))

		// Assign for runner with completely matching labels
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{
				Labels: map[string]string{
					"env": "test",
				}})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}

		// Assign for runner with partially matching labels
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{
				Labels: map[string]string{
					"env":    "test",
					"region": "testland-1",
				}})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}

		// Runner with no labels. Should skip a job with labels and pick up a job with no labels.
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{})
			require.NoError(err)
			require.NotNil(job)
			require.NotEqual("C", job.Id)
			require.Equal("D", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}

		// Runner with no matching labels. Should skip a job with mismatched labels and pick up a later job with matching labels.
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{
				Labels: map[string]string{
					"region": "testland-2",
				}})
			require.NoError(err)
			require.NotNil(job)
			require.NotEqual("C", job.Id)
			require.Equal("E", job.Id)
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)
			require.NoError(s.JobComplete(job.Id, nil, nil))
		}
	})

	t.Run("any cannot be assigned to ByIdOnly runner", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

	t.Run("assignment with a dependent job", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))

		// Assign it, we should get job A
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

		// Peek should return something, even though it is blocked, to
		// show that we do have some job queued.
		job, err = s.JobPeekForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)

		// Ack it
		_, err = s.JobAck("A", true)
		require.NoError(err)

		// Complete it
		require.NoError(s.JobComplete("A", &pb.Job_Result{
			Build: &pb.Job_BuildResult{},
		}, nil))

		// Assign it, we should now get job B
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)
		require.Equal(pb.Job_WAITING, job.State)
	})

	t.Run("assignment with a multiple dependent jobs", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "C",
			DependsOn: []string{"A", "B"},
		})))

		{
			// Assign it, we should get job A
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
			require.Equal(pb.Job_WAITING, job.State)

			// Ack it
			_, err = s.JobAck("A", true)
			require.NoError(err)

			// Complete it
			require.NoError(s.JobComplete("A", &pb.Job_Result{
				Build: &pb.Job_BuildResult{},
			}, nil))
		}

		{
			// Assign it, we should get job B
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("B", job.Id)
			require.Equal(pb.Job_WAITING, job.State)

			// Should block if requesting another since none exist
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			badjob, err := s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
			require.Error(err)
			require.Nil(badjob)
			require.Equal(ctx.Err(), err)

			// Ack it
			_, err = s.JobAck(job.Id, true)
			require.NoError(err)

			// Complete it
			require.NoError(s.JobComplete(job.Id, &pb.Job_Result{
				Build: &pb.Job_BuildResult{},
			}, nil))
		}

		// Peek should return something, even though it is blocked, to
		// show that we do have some job queued.
		job, err := s.JobPeekForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("C", job.Id)

		// Assign it, we should now get job B
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("C", job.Id)
		require.Equal(pb.Job_WAITING, job.State)
	})

	t.Run("blocking on dependency", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))

		// Assign it, we should get this build
		{
			job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
			require.NoError(err)
			require.NotNil(job)
			require.Equal("A", job.Id)
		}

		// Get the next value in a goroutine
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var job *serverstate.Job
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

		// Ack it
		_, err := s.JobAck("A", true)
		require.NoError(err)

		// Complete it
		require.NoError(s.JobComplete("A", &pb.Job_Result{
			Build: &pb.Job_BuildResult{},
		}, nil))

		// We should get a result
		select {
		case <-doneCh:

		case <-time.After(3 * time.Second):
			t.Fatal("should have a result")
		}

		require.NoError(jerr)
		require.NotNil(job)
		require.Equal("B", job.Id)
	})

	t.Run("assignment on unadopted runner", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Create our runner
		r := &pb.Runner{Id: "R_A"}
		require.NoError(s.RunnerCreate(r))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), r)
		require.Error(err)
		require.Nil(job)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("blocking when runner is rejected", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create our runner and adopt it
		r := &pb.Runner{Id: "R_A"}
		require.NoError(s.RunnerCreate(r))
		require.NoError(s.RunnerAdopt(r.Id, false))

		// Get the job in a goroutine
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var job *serverstate.Job
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

		// Reject our runner
		require.NoError(s.RunnerReject(r.Id))

		// We should get a result
		select {
		case <-doneCh:

		case <-time.After(3 * time.Second):
			t.Fatal("should have a result")
		}

		require.Error(jerr)
		require.Nil(job)
	})
}

func TestJobAck(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("ack", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

		s := factory(t)
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
		old := serverstate.JobWaitingTimeout
		defer func() { serverstate.JobWaitingTimeout = old }()
		serverstate.JobWaitingTimeout = time.Second

		s := factory(t)
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
		time.Sleep(3 * time.Second)

		// Verify it is queued
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.Job.State)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.Error(err)
	})
}

func TestJobComplete(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("success", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

		s := factory(t)
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

	t.Run("error cascades to dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it with error
		require.NoError(s.JobComplete(job.Id, nil, fmt.Errorf("bad")))

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)

		st := status.FromProto(job.Error)
		require.Equal(codes.Unknown, st.Code())
		require.Contains(st.Message(), "bad")

		// Should block if requesting another since none exist
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		job, err = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
		require.Error(err)
		require.Nil(job)
		require.Equal(ctx.Err(), err)

		// Dependent should be error
		job, err = s.JobById("B", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)
	})

	t.Run("error cascades to grandchildren+ dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "C",
			DependsOn: []string{"B"},
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it with error
		require.NoError(s.JobComplete(job.Id, nil, fmt.Errorf("bad")))

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)

		st := status.FromProto(job.Error)
		require.Equal(codes.Unknown, st.Code())
		require.Contains(st.Message(), "bad")

		// Should block if requesting another since none exist
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		job, err = s.JobAssignForRunner(ctx, &pb.Runner{Id: "R_A"})
		require.Error(err)
		require.Nil(job)
		require.Equal(ctx.Err(), err)

		// Dependent should be error
		job, err = s.JobById("B", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)

		// Grandchild should be error
		job, err = s.JobById("C", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)
	})

	t.Run("error does not cascade to allowed fail dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:                    "B",
			DependsOn:             []string{"A"},
			DependsOnAllowFailure: []string{"A"},
		})))

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Complete it with error
		require.NoError(s.JobComplete(job.Id, nil, fmt.Errorf("bad")))

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.State)
		require.NotNil(job.Error)

		st := status.FromProto(job.Error)
		require.Equal(codes.Unknown, st.Code())
		require.Contains(st.Message(), "bad")

		// Should block if requesting another since none exist
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)
	})

	t.Run("partially failed dependendents but allowed failure", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create jobs
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:                    "C",
			DependsOn:             []string{"A", "B"},
			DependsOnAllowFailure: []string{"B"},
		})))

		// Get A, success
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("A", job.Id)
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)
		require.NoError(s.JobComplete(job.Id, nil, nil))

		// Get B, failure
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("B", job.Id)
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)
		require.NoError(s.JobComplete(job.Id, nil, fmt.Errorf("bad")))

		// Should get C, even though partial failure.
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("C", job.Id)
	})
}

func TestJobTask_AckAndComplete(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("ack and complete on-demand runner jobs", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		task_id := "t_test"

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "start_job",
			Task: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: task_id,
				},
			},
		})))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "task_job",
			Task: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: task_id,
				},
			},
		})))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "stop_job",
			Task: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: task_id,
				},
			},
		})))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "watch_job",
			Task: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: task_id,
				},
			},
		})))

		// Create a pending task. note that `service_job` does this when it wraps
		// a requested job with an on-demand runner job triple.
		err := s.TaskPut(&pb.Task{
			Id:       task_id,
			TaskJob:  &pb.Ref_Job{Id: "task_job"},
			StartJob: &pb.Ref_Job{Id: "start_job"},
			StopJob:  &pb.Ref_Job{Id: "stop_job"},
			WatchJob: &pb.Ref_Job{Id: "watch_job"},
		})
		require.NoError(err)

		// Assign it, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("start_job", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		task, err := s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_STARTING, task.JobState)

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// We should have an output buffer
		require.NotNil(job.OutputBuffer)

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

		task, err = s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_STARTED, task.JobState)

		// Assign it, we should get this build
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_B"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("task_job", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		task, err = s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_RUNNING, task.JobState)

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// We should have an output buffer
		require.NotNil(job.OutputBuffer)

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

		task, err = s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_COMPLETED, task.JobState)

		// Assign it, we should get this build
		job, err = s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_C"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal("stop_job", job.Id)

		// Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		task, err = s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_STOPPING, task.JobState)

		// Verify it is changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		// We should have an output buffer
		require.NotNil(job.OutputBuffer)

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

		task, err = s.TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task_id,
			},
		})
		require.NoError(err)
		require.NotNil(task)
		require.Equal(pb.Task_STOPPED, task.JobState)
	})
}

func TestJobPipeline_AckAndComplete(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("ack and complete on-demand runner jobs", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		jobRef := &pb.Ref_Job{Id: "root_job"}

		// Write project
		ref := &pb.Ref_Project{Project: "project"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		p := serverptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)
		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: p.Id}}

		// Create a new pipeline run
		pr := &pb.PipelineRun{Pipeline: pipeline}
		r := serverptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Create a job
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: jobRef.Id,
			Pipeline: &pb.Ref_PipelineStep{
				PipelineId:  p.Id,
				RunSequence: 1,
			},
		})))

		r.Jobs = append(r.Jobs, jobRef)
		err = s.PipelineRunPut(r)
		require.NoError(err)
		require.Equal(uint64(1), r.Sequence)
		require.Equal(pb.PipelineRun_PENDING, r.State)

		// Assign the job, we should get this build
		job, err := s.JobAssignForRunner(context.Background(), &pb.Runner{Id: "R_A"})
		require.NoError(err)
		require.NotNil(job)
		require.Equal(jobRef.Id, job.Id)

		//Ack it
		_, err = s.JobAck(job.Id, true)
		require.NoError(err)

		// Verify job and pipeline states changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_RUNNING, job.Job.State)

		run, err := s.PipelineRunGetById(r.Id)
		require.NoError(err)
		require.NotNil(run)
		require.Equal(pb.PipelineRun_RUNNING, run.State)
		require.Equal(job.Pipeline.RunSequence, run.Sequence)

		// We should have an output buffer
		require.NotNil(job.OutputBuffer)

		// Complete the job
		require.NoError(s.JobComplete(job.Id, &pb.Job_Result{
			Build: &pb.Job_BuildResult{},
		}, nil))

		// Verify job and pipeline status changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.Equal(pb.Job_SUCCESS, job.State)
		require.Nil(job.Error)
		require.NotNil(job.Result)

		run, err = s.PipelineRunGetById(pr.Id)
		require.NoError(err)
		require.NotNil(run)
		require.Equal(pb.PipelineRun_SUCCESS, run.State)
		require.Equal(job.Pipeline.RunSequence, run.Sequence)
	})
}

func TestJobIsAssignable(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("no runners", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
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

		s := factory(t)
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

		s := factory(t)
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

		s := factory(t)
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

		s := factory(t)
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

	t.Run("Labels target, all labels", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
		defer s.Close()

		// Register a runner with labels
		runner := serverptypes.TestRunner(t, &pb.Runner{
			Labels: map[string]string{
				"env":    "test",
				"region": "testland-1",
			},
		})
		require.NoError(s.RunnerCreate(runner))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: runner.Labels,
					},
				},
			},
		}))
		require.NoError(err)
		require.True(result)
	})

	t.Run("Labels target, partial", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
		defer s.Close()

		// Register a runner with labels
		runner := serverptypes.TestRunner(t, &pb.Runner{
			Labels: map[string]string{
				"env":    "test",
				"region": "testland-1",
			},
		})
		require.NoError(s.RunnerCreate(runner))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"env": "test",
						},
					},
				},
			},
		}))
		require.NoError(err)
		require.True(result)
	})

	t.Run("Labels target, no match", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()

		s := factory(t)
		defer s.Close()

		// Register a runner with labels
		runner := serverptypes.TestRunner(t, &pb.Runner{
			Labels: map[string]string{
				"env":    "test",
				"region": "testland-1",
			},
		})
		require.NoError(s.RunnerCreate(runner))

		// Should be assignable
		result, err := s.JobIsAssignable(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: map[string]string{
							"region": "outer-space",
						},
					},
				},
			},
		}))
		require.NoError(err)
		require.False(result)
	})
}

func TestJobCancel(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("queued", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

	t.Run("queued with dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
		})))

		// Cancel it
		require.NoError(s.JobCancel("A", false))

		// Verify it is canceled
		job, err := s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		// Verify dependent is canceled
		job, err = s.JobById("B", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("queued with pipeline and dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a new pipeline run
		p := serverptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)
		pr := &pb.PipelineRun{
			Pipeline: &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: p.Id}},
		}
		r := serverptypes.TestPipelineRun(t, pr)

		// Create jobs
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:       "A",
			Pipeline: &pb.Ref_PipelineStep{PipelineId: p.Id, RunSequence: 1},
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			Pipeline:  &pb.Ref_PipelineStep{PipelineId: p.Id, RunSequence: 1},
			DependsOn: []string{"A"},
		})))

		// Update pipeline run with jobs
		r.Jobs = []*pb.Ref_Job{{Id: "A"}, {Id: "B"}}
		err = s.PipelineRunPut(r)
		require.NoError(err)
		require.Equal(pb.PipelineRun_PENDING, r.State)

		// Cancel parent
		require.NoError(s.JobCancel("A", false))

		// Verify it is cancelled
		job, err := s.JobById("A", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		// Verify dependent is cancelled
		job, err = s.JobById("B", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		// Verify pipeline run is cancelled
		run, err := s.PipelineRunGetById(r.Id)
		require.NoError(err)
		require.NotNil(run)
		require.NotEmpty(run.Jobs)
		require.Equal(uint64(1), run.Sequence)
		require.Equal(pb.PipelineRun_CANCELLED, run.State)
	})

	t.Run("assigned", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

		s := factory(t)
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

		s := factory(t)
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

	t.Run("assigned with force cancels dependents", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id:        "B",
			DependsOn: []string{"A"},
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
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)

		// Verify dependent is canceled. Even if it was assigned and forced
		// we should have canceled all dependents.
		job, err = s.JobById("B", nil)
		require.NoError(err)
		require.Equal(pb.Job_ERROR, job.Job.State)
		require.NotNil(job.Job.Error)
		require.NotEmpty(job.CancelTime)
	})

	t.Run("completed", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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

func TestJobHeartbeat(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("times out after ack", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set a short timeout
		old := serverstate.JobHeartbeatTimeout
		defer func() { serverstate.JobHeartbeatTimeout = old }()
		serverstate.JobHeartbeatTimeout = time.Second

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
		}, 4*time.Second, time.Second)
	})

	t.Run("doesn't time out if heartbeating", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := serverstate.JobHeartbeatTimeout
		defer func() { serverstate.JobHeartbeatTimeout = old }()
		serverstate.JobHeartbeatTimeout = time.Second

		s := factory(t)
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

			tick := time.NewTicker(500 * time.Millisecond)
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
		old := serverstate.JobHeartbeatTimeout
		defer func() { serverstate.JobHeartbeatTimeout = old }()
		serverstate.JobHeartbeatTimeout = time.Second

		s := factory(t)
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

			tick := time.NewTicker(500 * time.Millisecond)
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
		}, 4*time.Second, time.Second)
	})
}

func TestJobHeartbeatOnRestart(t *testing.T, factory Factory, rf RestartFactory) {

	t.Run("times out if running state loaded on restart", func(t *testing.T) {
		require := require.New(t)

		// Set a short timeout
		old := serverstate.JobHeartbeatTimeout
		defer func() { serverstate.JobHeartbeatTimeout = old }()
		serverstate.JobHeartbeatTimeout = time.Second

		s := factory(t)
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
		go func(s serverstate.Interface) {
			defer close(doneCh)

			tick := time.NewTicker(500 * time.Millisecond)
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
		s = rf(t, s)
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
		}, 4*time.Second, time.Second)
	})
}

func TestJobUpdateRef(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("running", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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
		require.False(ws.Watch(time.After(3 * time.Second)))

		// Verify it was changed
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.NotNil(job.DataSourceRef)

		ref := job.DataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git).Git
		require.Equal(ref.Commit, "hello")
	})
}

func TestJobUpdateExpiry(t *testing.T, factory Factory, rf RestartFactory) {
	t.Run("new expire time", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
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
		require.Nil(job.ExpireTime)

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

		// Update the expire time
		dur, err := time.ParseDuration("60s")
		newExpireTime := timestamppb.New(time.Now().Add(dur))
		require.NoError(s.JobUpdateExpiry(job.Id, newExpireTime))

		// Should be triggered. This is a very important test because
		// we need to ensure that the watchers can detect ref changes.
		require.False(ws.Watch(time.After(3 * time.Second)))

		// Verify it was set
		job, err = s.JobById(job.Id, nil)
		require.NoError(err)
		require.NotNil(job.ExpireTime)
	})
}
