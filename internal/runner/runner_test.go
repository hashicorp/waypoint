package runner

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestRunnerStart(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	client := singleprocess.TestServer(t)

	// Initialize our runner
	runner, err := New(
		WithClient(client),
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

func TestRunnerStart_adoption(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	serverImpl := singleprocess.TestImpl(t)
	client := server.TestServer(t, serverImpl)

	// Client with no token
	anonClient := server.TestServer(t, serverImpl, server.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
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

func TestRunnerStart_rejection(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	serverImpl := singleprocess.TestImpl(t)
	client := server.TestServer(t, serverImpl)

	// Client with no token
	anonClient := server.TestServer(t, serverImpl, server.TestWithToken(""))

	// Initialize our runner
	runner, err := New(
		WithClient(anonClient),
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
