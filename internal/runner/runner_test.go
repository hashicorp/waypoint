package runner

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	require.NoError(runner.Start())

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

func TestRunnerStart_config(t *testing.T) {
	t.Run("set and unset", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()
		client := singleprocess.TestServer(t)

		cfgVar := &pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
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
		require.NoError(runner.Start())

		// Verify it is not set
		require.Empty(os.Getenv(cfgVar.Name))

		// Set some config
		_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be set
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == cfgVar.Value.(*pb.ConfigVar_Static).Static
		}, 1000*time.Millisecond, 50*time.Millisecond)

		// Unset
		cfgVar.Value = &pb.ConfigVar_Static{Static: ""}
		_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be unset
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == ""
		}, 1000*time.Millisecond, 50*time.Millisecond)
	})

	t.Run("unset with original env", func(t *testing.T) {
		require := require.New(t)
		ctx := context.Background()
		client := singleprocess.TestServer(t)

		cfgVar := &pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
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
		require.NoError(runner.Start())

		// Set some config
		_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be set
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == cfgVar.Value.(*pb.ConfigVar_Static).Static
		}, 1000*time.Millisecond, 50*time.Millisecond)

		// Unset
		cfgVar.Value = &pb.ConfigVar_Static{Static: ""}
		_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{Variables: []*pb.ConfigVar{cfgVar}})
		require.NoError(err)

		// Should be unset back to original value
		require.Eventually(func() bool {
			return os.Getenv(cfgVar.Name) == "ORIGINAL"
		}, 1000*time.Millisecond, 50*time.Millisecond)
	})
}
