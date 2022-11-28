package boltdbstate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestJobAck(t *testing.T) {
	ctx := context.Background()
	t.Run("A job nack unsets the job's assigned runner", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		job := serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})
		require.NoError(s.JobCreate(ctx, job))

		ctx := context.Background()
		runner := &pb.Runner{
			Id: "test_runner_id",
		}

		assignedJob, err := s.JobAssignForRunner(ctx, runner)
		require.NoError(err)
		require.Equal(assignedJob.AssignedRunner.Id, runner.Id)

		nackedJob, err := s.JobAck(ctx, job.Id, false)
		require.NoError(err)
		require.Nil(nackedJob.Job.AssignedRunner)
	})
}

func TestJobCreateMulti(t *testing.T) {
	ctx := context.Background()

	t.Run("creates one job", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		jobList := make([]*pb.Job, 0, 1)
		jobList = append(jobList, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		}))

		err := s.JobCreate(ctx, jobList...)
		require.NoError(err)

		require.Equal(1, s.indexedJobs)
	})

	t.Run("creates the same number of jobs that were requested", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		jobList := make([]*pb.Job, 0, 3)
		jobList = append(jobList, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		}))
		jobList = append(jobList, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		}))
		jobList = append(jobList, serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		}))

		err := s.JobCreate(ctx, jobList...)
		require.NoError(err)

		require.Equal(3, s.indexedJobs)
	})
}

func TestJobAssignForRunner(t *testing.T) {
	ctx := context.Background()

	t.Run("job assignment sets the job's assigned runner id", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		ctx := context.Background()
		runner := &pb.Runner{
			Id: "test_runner_id",
		}

		job, err := s.JobAssignForRunner(ctx, runner)
		require.NoError(err)
		require.Equal(job.AssignedRunner.Id, runner.Id)
	})
}

func TestJobsPrune(t *testing.T) {
	ctx := context.Background()
	t.Run("removes only completed jobs", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Cancel it
		require.NoError(s.JobCancel(ctx, "A", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		// Leave B running

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 0)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(1, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById(ctx, "A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("does nothing there are fewer than the maximum", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Cancel it
		require.NoError(s.JobCancel(ctx, "A", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		// Leave B running

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		require.Equal(2, s.indexedJobs)
		cnt, err := s.jobsPruneOld(memTxn, 10)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(0, cnt)
		require.Equal(2, s.indexedJobs)

		val, err := s.JobById(ctx, "A", nil)
		require.NoError(err)
		require.NotNil(val)

		val, err = s.JobById(ctx, "B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("stops when the maximum are pruned", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel(ctx, "A", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel(ctx, "B", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 1)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(1, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById(ctx, "A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("prunes according to the queue time", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel(ctx, "A", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel(ctx, "B", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		})))

		require.NoError(s.JobCancel(ctx, "C", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 1)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(2, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById(ctx, "A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "B", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "C", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("can prune all jobs", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel(ctx, "A", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel(ctx, "B", false))

		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		})))

		require.NoError(s.JobCancel(ctx, "C", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 0)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(3, cnt)
		require.Equal(0, s.indexedJobs)

		val, err := s.JobById(ctx, "A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "B", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById(ctx, "C", nil)
		require.NoError(err)
		require.Nil(val)
	})
}

func TestJobsProjectScopedRequest(t *testing.T) {
	ctx := context.Background()
	t.Run("returns error if no project ref found", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		const name = "proj"
		ref := &pb.Ref_Project{Project: name}

		jobTemplate := &pb.Job{
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Operation: &pb.Job_Init{},
		}

		{
			resp, err := s.JobProjectScopedRequest(ctx, ref, jobTemplate)
			require.Error(err)
			require.Nil(resp)
		}
	})

	t.Run("returns a list of queued job request messages for all apps in project", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		const name = "proj"
		ref := &pb.Ref_Project{Project: name}

		proj := serverptypes.TestProject(t, &pb.Project{Name: name})
		err := s.ProjectPut(ctx, proj)
		require.NoError(err)
		_, err = s.AppPut(ctx, serverptypes.TestApplication(t, &pb.Application{
			Name:    "test",
			Project: ref,
		}))
		require.NoError(err)
		_, err = s.AppPut(ctx, serverptypes.TestApplication(t, &pb.Application{
			Name:    "test2",
			Project: ref,
		}))
		require.NoError(err)

		jobTemplate := &pb.Job{
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Operation: &pb.Job_Init{},
		}

		{
			resp, err := s.JobProjectScopedRequest(ctx, ref, jobTemplate)
			require.NoError(err)
			require.Len(resp, 2)
		}
	})
}
