// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	serverpkg "github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func TestRunnerAccept_happy(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

// Test runner accept and job execution with an adopted token.
func TestRunnerAccept_adopt(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Start runner
	serverImpl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, serverImpl)
	anonClient := serverpkg.TestServer(t, serverImpl, serverpkg.TestWithToken(""))
	runner := TestRunner(t,
		WithClient(anonClient),
		WithCookie(testCookie(t, client)),
	)
	defer runner.Close()

	// Start runner
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(ctx)
	}()

	// Wait for runner to show up
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// Adopt the runner
	_, err := client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
		RunnerId: runner.Id(),
		Adopt:    true,
	})
	require.NoError(err)

	// Runner should start
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("runner should start")

	case err := <-startErr:
		require.NoError(err)
	}

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_timeout(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client), WithAcceptTimeout(10*time.Millisecond))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// DO NOT QUEUE A JOB. Accept sholud timeout
	err := runner.Accept(ctx)
	require.Error(err)
	require.True(errors.Is(err, ErrTimeout))
}

// Test how accept behaves when the server is down to begin with.
func TestRunnerAccept_serverDown(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)

	// Setup our runner
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Shut it down
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: "A"})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Start accept
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Accept(ctx)
	}()

	// The runner should not error
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-errCh:
		t.Fatal("runner should not return")
	}

	// Restart
	restartCh <- struct{}{}

	// Accept should return
	select {
	case err := <-errCh:
		require.NoError(err)

	case <-time.After(5 * time.Second):
		t.Fatal("accept never returned")
	}

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

// Test how accept behaves when the server is down while waiting for
// assignment.
func TestRunnerAccept_serverDownAssign(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// New context that will be used for server restart.
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()

	// Setup the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(serverCtx),
		serverpkg.TestWithRestart(restartCh),
	)

	// Setup our runner
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Start accept
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Accept(ctx)
	}()

	// The runner should not error
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-errCh:
		t.Fatal("runner should not return")
	}

	// Shut it down
	serverCancel()
	serverCtx, serverCancel = context.WithCancel(context.Background())
	defer serverCancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// Queue a job
	var jobId string
	require.Eventually(func() bool {
		queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
			Job: serverptypes.TestJobNew(t, nil),
		})
		if err != nil {
			return false
		}

		jobId = queueResp.JobId
		return true
	}, 5*time.Second, 10*time.Millisecond)

	// Accept should return
	select {
	case err := <-errCh:
		require.NoError(err)

	case <-time.After(5 * time.Second):
		t.Fatal("accept never returned")
	}

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

// Server down during job execution.
func TestRunnerAccept_serverDownJobExec(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// New context that will be used for server restart.
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()

	// Setup the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(serverCtx),
		serverpkg.TestWithRestart(restartCh),
	)

	// Setup our runner
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Make sure our noop operation blocks
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Start accept
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Accept(ctx)
	}()

	// The runner should not error
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-errCh:
		t.Fatal("runner should not return")
	}

	// Wait for the job to be running, then we know it is acked.
	require.Eventually(func() bool {
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
		return err == nil && job.State == pb.Job_RUNNING
	}, 5*time.Second, 10*time.Millisecond)

	// Shut the server down
	serverCancel()
	serverCtx, serverCancel = context.WithCancel(context.Background())
	defer serverCancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// Let job complete
	close(noopCh)

	// Accept should return
	select {
	case err := <-errCh:
		require.NoError(err)

	case <-time.After(5 * time.Second):
		t.Fatal("accept never returned")
	}

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_closeCancelsAccept(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	done := make(chan struct{})

	go func() {
		if runner.Accept(ctx) != context.Canceled {
			close(done)
		}
	}()

	// To allow Accept to block on waiting for a job
	time.Sleep(time.Second)

	require.NoError(runner.Close())

	select {
	case <-done:
		// ok
	case <-ctx.Done():
		t.Error("close did not cancel accept")
	}
}

func TestRunnerAccept_closeHoldingJob(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	go func() {
		time.Sleep(time.Second)
		runner.Close()
		runner.logger.Info("runner closed")
	}()

	testRecvDelay = 2 * time.Second
	defer func() {
		testRecvDelay = 0
	}()

	// Accept should error, seeing the runner shutdown
	err = runner.Accept(ctx)
	require.Error(err)

	require.Contains(err.Error(), jobId)

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State, job.State.String())
}

func TestRunnerAccept_closeWaits(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	go func() {
		err := runner.Accept(ctx)
		runner.logger.Error("error in accept", "error", err)
	}()

	// To allow Accept to block on noopCh while executing the job
	time.Sleep(time.Second)

	time.AfterFunc(2*time.Second, func() {
		close(noopCh)
	})

	runner.runningCond.L.Lock()
	count := runner.runningJobs
	runner.runningCond.L.Unlock()

	require.Equal(1, count)

	ts := time.Now()
	require.NoError(runner.Close())
	dur := time.Since(ts)

	require.True(dur >= 2*time.Second, dur.String())

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State, job.State.String())
}

func TestRunnerAccept_cancelContext(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Set a blocker
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Cancel the context eventually. This isn't CI-sensitive cause
	// we'll block no matter what.
	time.AfterFunc(500*time.Millisecond, cancel)

	// Accept should complete with no error
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	require.Eventually(func() bool {
		job, err := client.GetJob(context.Background(), &pb.GetJobRequest{JobId: jobId})
		require.NoError(err)
		return job.State == pb.Job_ERROR
	}, 3*time.Second, 25*time.Millisecond)
}

func TestRunnerAccept_cancelJob(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Set a blocker
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Cancel the context eventually. This isn't CI-sensitive cause
	// we'll block no matter what.
	time.AfterFunc(500*time.Millisecond, func() {
		_, err := client.CancelJob(ctx, &pb.CancelJobRequest{
			JobId: jobId,
		})
		require.NoError(err)
	})

	// Accept should complete with an error
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	require.Eventually(func() bool {
		job, err := client.GetJob(context.Background(), &pb.GetJobRequest{JobId: jobId})
		require.NoError(err)
		return job.State == pb.Job_ERROR
	}, 3*time.Second, 25*time.Millisecond)
}

func TestRunnerAccept_gitData(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	require := require.New(t)
	ctx := context.Background()

	// Get a repo path
	path := testGitFixture(t, "git-noop")

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: path,
					},
				},
			},
		}),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
	require.NotNil(job.DataSourceRef)
}

func TestRunnerAccept_noConfig(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_ERROR, job.State)
	require.NotNil(job.Error)

	st := status.FromProto(job.Error)
	require.Equal(codes.FailedPrecondition, st.Code())
}

func TestRunnerAccept_noConfig_serverHcl(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Initialize our app
	ref := serverptypes.TestJobNew(t, nil).Application
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name:        ref.Project,
				WaypointHcl: []byte(configpkg.TestSource(t)),
			},
		})
		require.NoError(err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(err)
	}

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_noConfig_serverHclJson(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Initialize our app
	ref := serverptypes.TestJobNew(t, nil).Application
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name:              ref.Project,
				WaypointHcl:       []byte(configpkg.TestSourceJSON(t)),
				WaypointHclFormat: pb.Hcl_JSON,
			},
		})
		require.NoError(err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(err)
	}

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_jobHcl(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Create the project/app
	ref := serverptypes.TestJobNew(t, nil).Application
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name: ref.Project,
			},
		})
		require.NoError(err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(err)
	}

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			WaypointHcl: &pb.Hcl{
				Contents: []byte(configpkg.TestSourceJSON(t)),
				Format:   pb.Hcl_JSON,
			},
		}),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
	require.Equal(pb.Job_Config_JOB, job.Config.Source)
}

func TestRunnerAcceptParallel(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start(ctx))

	// Block our noop jobs so we can inspect their state
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue jobs
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Workspace: &pb.Ref_Workspace{Workspace: "w1"},
		}),
	})
	require.NoError(err)
	jobId_1 := queueResp.JobId

	queueResp, err = client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Workspace: &pb.Ref_Workspace{Workspace: "w2"},
		}),
	})
	require.NoError(err)
	jobId_2 := queueResp.JobId

	// Accept should complete
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		runner.AcceptParallel(ctx, 2)
	}()

	// Both jobs should be running at once eventually
	require.Eventually(func() bool {
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId_1})
		require.NoError(err)
		if job.State != pb.Job_RUNNING {
			return false
		}

		job, err = client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId_2})
		require.NoError(err)
		return job.State == pb.Job_RUNNING
	}, 3*time.Second, 10*time.Millisecond)

	// Jobs should complete
	close(noopCh)
	require.Eventually(func() bool {
		job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId_1})
		require.NoError(err)
		if job.State != pb.Job_SUCCESS {
			return false
		}

		job, err = client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId_2})
		require.NoError(err)
		return job.State == pb.Job_SUCCESS
	}, 3*time.Second, 10*time.Millisecond)

	// Loop should exit
	cancel()
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("accept should exit")

	default:
	}
}

// testGitFixture MUST be called before TestRunner since TestRunner
// changes our working directory.
func testGitFixture(t *testing.T, n string) string {
	t.Helper()

	// We need to get our working directory since the TestRunner call
	// changes it.
	wd, err := os.Getwd()
	require.NoError(t, err)
	wd, err = filepath.Abs(wd)
	require.NoError(t, err)
	path := filepath.Join(wd, "testdata", n)

	// Look for a DOTgit
	original := filepath.Join(path, "DOTgit")
	_, err = os.Stat(original)
	require.NoError(t, err)

	// Rename it
	newPath := filepath.Join(path, ".git")
	require.NoError(t, os.Rename(original, newPath))
	t.Cleanup(func() { os.Rename(newPath, original) })

	return path
}
