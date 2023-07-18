// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	serverpkg "github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestRunnerStart(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	client := singleprocess.TestServer(t)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithCookie(testCookie(t, client)),
	)
	require.NoError(err)
	defer runner.Close()

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))

	// Start it
	require.NoError(runner.Start(ctx))

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)

	// Close
	require.NoError(runner.Close())
	time.Sleep(100 * time.Millisecond)

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))
}

// Test that the runner reconnects after it successfully registered initially.
func TestRunnerStart_reconnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithCookie(testCookie(t, client)),
	)
	require.NoError(err)
	defer runner.Close()

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))

	// Start it
	require.NoError(runner.Start(ctx))

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)

	// Shut down the server
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// We should get deregistered
	require.Eventually(func() bool {
		_, err = impl.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err != nil && status.Code(err) == codes.NotFound
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// We should get re-registered
	require.Eventually(func() bool {
		_, err = impl.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}

// Test how the runner behaves on start if the server is down immediately.
func TestRunnerStart_serverDown(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)
	cookie := testCookie(t, client)

	// Shut it down
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: "A"})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithCookie(cookie),
	)
	require.NoError(err)
	defer runner.Close()

	// Start it
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Start(ctx)
	}()

	// Restart
	restartCh <- struct{}{}

	// Start should return
	select {
	case err := <-errCh:
		require.NoError(err)

	case <-time.After(5 * time.Second):
		t.Fatal("start never returned")
	}

	// We should get re-registered eventually
	require.Eventually(func() bool {
		_, err = impl.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}

// Test how the runner behaves on start if the server is down after initial syncing is complete.
func TestRunnerStart_serverDownPostRegistration(t *testing.T) {

	log := hclog.L()
	log.SetLevel(hclog.Trace)
	log = log.With("origin", "test")

	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)
	cookie := testCookie(t, client)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithCookie(cookie),
	)
	require.NoError(err)
	//defer runner.Close()

	// Move into a temporary directory
	td := testTempDir(t)
	testChdir(t, td)
	configpkg.TestConfigFile(t, configpkg.TestSource(t))
	// Start the runner
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Start(ctx)
	}()

	// Start should return
	select {
	case err := <-errCh:
		require.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("start never returned")
	}

	// Begin accepting jobs
	runnerCtx := context.Background()
	go runner.AcceptParallel(runnerCtx, 1)

	// Initialize our app
	projectName := serverptypes.TestJobNew(t, nil).Application.Project
	appName := serverptypes.TestJobNew(t, nil).Application.Application
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name: projectName,
			},
		})
		require.NoError(err)

		_, err = client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: projectName},
			Name:    appName,
		})
		require.NoError(err)
	}

	// Ensure jobs can run
	{
		// Queue the job
		queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
		require.NoError(err)
		require.NotNil(queueResp)
		require.NotEmpty(queueResp.JobId)

		// Wait for the job to execute
		require.Eventually(func() bool {
			jobResp, err := client.GetJob(ctx, &pb.GetJobRequest{
				JobId: queueResp.JobId,
			})
			require.NoError(err)
			require.Nil(jobResp.Error)
			return jobResp.State == pb.Job_SUCCESS
		}, 99999*time.Second, 5*time.Second)
	}

	log.Info("TEST: Done with the preamble - the runner/server works - now turning the server off...")

	// Shut down the server
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	//Wait to get an unavailable error from the runner so that we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: "A"})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	log.Info("TEST: The server is properly off now. Restarting it...")
	// Restart the server
	restartCh <- struct{}{}

	// We should get re-registered eventually
	require.Eventually(func() bool {
		_, err = impl.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	log.Info("TEST: The server is now back online. Queueing a new job...")

	// Queue the job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	log.Info("TEST: Waiting for the new job to execute...")

	// Wait for the job to execute
	require.Eventually(func() bool {
		jobResp, err := client.GetJob(ctx, &pb.GetJobRequest{
			JobId: queueResp.JobId,
		})
		require.NoError(err)
		require.Nil(jobResp.Error)
		return jobResp.State == pb.Job_SUCCESS
	}, 999999*time.Second, 5*time.Second)

	log.Info("TEST: Test complete! Shutting down...")
}

func TestRunnerStart_adoption(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	serverImpl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, serverImpl)

	// Client with no token
	anonClient := serverpkg.TestServer(t, serverImpl, serverpkg.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
		WithCookie(testCookie(t, client)),
	)
	require.NoError(err)
	defer runner.Close()

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))

	// Start
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(ctx)
	}()

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// The runner should not start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-startErr:
		t.Fatal("runner should not start")
	}

	// Adopt the runner
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
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

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)
	require.Equal(pb.Runner_ADOPTED, resp.AdoptionState)

	// Close
	require.NoError(runner.Close())
	time.Sleep(100 * time.Millisecond)
}

// Test adoption works when the server is down on start
func TestRunnerStart_adoptionDownOnStart(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)
	cookie := testCookie(t, client)
	anonClient := serverpkg.TestServer(t, impl, serverpkg.TestWithToken(""))

	// Shut it down
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: "A"})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
		WithCookie(cookie),
	)
	require.NoError(err)
	defer runner.Close()

	// Start
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(ctx)
	}()

	// Restart
	restartCh <- struct{}{}

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// The runner should not start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-startErr:
		t.Fatal("runner should not start")
	}

	// Adopt the runner
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
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

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)
	require.Equal(pb.Runner_ADOPTED, resp.AdoptionState)

	// Close
	require.NoError(runner.Close())
	time.Sleep(100 * time.Millisecond)
}

// Test adoption works when the server goes down while blocked on adoption.
func TestRunnerStart_adoptionWaitDown(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, impl,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)
	cookie := testCookie(t, client)
	anonClient := serverpkg.TestServer(t, impl, serverpkg.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
		WithCookie(cookie),
	)
	require.NoError(err)
	defer runner.Close()

	// Start. We need a new context so our server shutdown doesn't cancel
	// the context we use for Start.
	startCtx, startCancel := context.WithCancel(context.Background())
	defer startCancel()
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(startCtx)
	}()

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// The runner should not start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-startErr:
		t.Fatal("runner should not start")
	}

	// Shut it down
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// The runner should not error on start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case err := <-startErr:
		t.Logf("err: %s %#[1]v", err)
		t.Fatal("runner should not error on start")
	}

	// Restart
	restartCh <- struct{}{}

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// Adopt the runner
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
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

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)
	require.Equal(pb.Runner_ADOPTED, resp.AdoptionState)

	// Close
	require.NoError(runner.Close())
	time.Sleep(100 * time.Millisecond)
}
func TestRunnerStart_rejection(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	serverImpl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, serverImpl)

	// Client with no token
	anonClient := serverpkg.TestServer(t, serverImpl, serverpkg.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
		WithCookie(testCookie(t, client)),
	)
	require.NoError(err)
	defer runner.Close()

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))

	// Start
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(ctx)
	}()

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// The runner should not start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-startErr:
		t.Fatal("runner should not start")
	}

	// Adopt the runner
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
		RunnerId: runner.Id(),
		Adopt:    false,
	})
	require.NoError(err)

	// Runner should start
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("runner should start")

	case err := <-startErr:
		require.Error(err)
	}
}

func TestRunnerStart_adoptionStateRestart(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Temp dir
	td, err := ioutil.TempDir("", "wprunner")
	require.NoError(err)
	defer os.RemoveAll(td)

	serverImpl := singleprocess.TestImpl(t)
	client := serverpkg.TestServer(t, serverImpl)

	// Client with no token
	anonClient := serverpkg.TestServer(t, serverImpl, serverpkg.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
		WithCookie(testCookie(t, client)),
		WithStateDir(td),
	)
	require.NoError(err)
	defer runner.Close()

	// The runner should not be registered
	_, err = client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))

	// Start
	startErr := make(chan error, 1)
	go func() {
		startErr <- runner.Start(ctx)
	}()

	// Wait for registration
	require.Eventually(func() bool {
		_, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	// The runner should not start
	select {
	case <-time.After(100 * time.Millisecond):
		// Good

	case <-startErr:
		t.Fatal("runner should not start")
	}

	// Adopt the runner
	_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
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

	// The runner should be registered
	resp, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: runner.Id()})
	require.NoError(err)
	require.Equal(runner.Id(), resp.Id)
	require.Equal(pb.Runner_ADOPTED, resp.AdoptionState)

	// Close
	require.NoError(runner.Close())
	time.Sleep(100 * time.Millisecond)

	// Restart
	runner, err = New(
		WithClient(anonClient),
		WithCookie(testCookie(t, client)),
		WithStateDir(td),
	)
	require.NoError(err)
	defer runner.Close()

	// Should start immediately
	require.NoError(runner.Start(ctx))
}

func TestRunnerStart_config(t *testing.T) {
	t.Run("set and unset", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()
		client := singleprocess.TestServer(t)

		cfgVar := &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Global{
					Global: &pb.Ref_Global{},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "I_AM_A_TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "1234567890"},
		}

		// Initialize our runner
		runner := TestRunner(t, WithClient(client))
		defer runner.Close()
		require.NoError(runner.Start(ctx))

		// Verify it is not set
		require.Empty(os.Getenv(cfgVar.Name))

		// Set some config
		_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be set
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == cfgVar.Value.(*pb.ConfigVar_Static).Static
		}, 2000*time.Millisecond, 50*time.Millisecond)

		// Unset
		cfgVar.Value = &pb.ConfigVar_Static{Static: ""}
		_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be unset
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == ""
		}, 2000*time.Millisecond, 50*time.Millisecond)
	})

	t.Run("unset with original env", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()
		client := singleprocess.TestServer(t)

		cfgVar := &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Global{
					Global: &pb.Ref_Global{},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "I_AM_A_TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "1234567890"},
		}

		// Set a value
		require.NoError(os.Setenv(cfgVar.Name, "ORIGINAL"))
		defer os.Unsetenv(cfgVar.Name)

		// Initialize our runner
		runner := TestRunner(t, WithClient(client))
		defer runner.Close()
		require.NoError(runner.Start(ctx))

		// Set some config
		_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be set
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == cfgVar.Value.(*pb.ConfigVar_Static).Static
		}, 2000*time.Millisecond, 50*time.Millisecond)

		// Unset
		cfgVar.Value = &pb.ConfigVar_Static{Static: ""}
		_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be unset back to original value
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == "ORIGINAL"
		}, 2000*time.Millisecond, 50*time.Millisecond)
	})

	t.Run("files", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()
		client := singleprocess.TestServer(t)

		// Create a temp dir with a filepath in it
		td, err := ioutil.TempDir("", "waypoint")
		require.NoError(err)
		defer os.RemoveAll(td)
		name := filepath.Join(td, "config")

		cfgVar := &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Global{
					Global: &pb.Ref_Global{},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:       name,
			NameIsPath: true,
			Value:      &pb.ConfigVar_Static{Static: "1234567890"},
		}

		// Initialize our runner
		runner := TestRunner(t, WithClient(client))
		defer runner.Close()
		require.NoError(runner.Start(ctx))

		// Verify it does not exist
		_, err = os.Stat(name)
		require.Error(err)
		require.True(os.IsNotExist(err))

		// Set some config
		_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be set
		require.Eventually(func() bool {
			data, err := ioutil.ReadFile(name)
			return err == nil && cfgVar.Value.(*pb.ConfigVar_Static).Static == string(data)
		}, 2000*time.Millisecond, 50*time.Millisecond)
	})
}

func TestRunner_stateId(t *testing.T) {
	require := require.New(t)
	client := singleprocess.TestServer(t)

	// Temp dir
	td, err := ioutil.TempDir("", "wprunner")
	require.NoError(err)
	defer os.RemoveAll(td)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithStateDir(td),
	)
	require.NoError(err)
	defer runner.Close()

	// Should have some ID
	id := runner.Id()
	require.NotEmpty(id)

	// Init again, should have same ID
	runner2, err := New(
		WithClient(client),
		WithStateDir(td),
	)
	require.NoError(err)
	defer runner2.Close()

	// Should have some ID
	require.Equal(id, runner2.Id())
}

func testCookie(t *testing.T, c pb.WaypointClient) string {
	resp, err := c.GetServerConfig(context.Background(), &empty.Empty{})
	require.NoError(t, err)
	return resp.Config.Cookie
}
