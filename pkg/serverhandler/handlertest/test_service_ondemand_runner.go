package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["runner_ondemand"] = []testFunc{
		TestServiceOnDemandRunnerConfig,
		TestServiceOnDemandRunnerConfig_GetOnDemandRunnerConfig,
		TestServiceOnDemandRunnerConfig_ListOnDemandRunnerConfigs,
	}
}

func TestServiceOnDemandRunnerConfig(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertOnDemandRunnerConfigRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertOnDemandRunnerConfig(ctx, &Req{
			Config: serverptypes.TestOnDemandRunnerConfig(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Config
		result.PluginType = "blargh"
		require.NotEmpty(result.Id)

		// Let's write some data
		resp, err = client.UpsertOnDemandRunnerConfig(ctx, &Req{
			Config: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Config

		require.Equal("blargh", result.PluginType)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertOnDemandRunnerConfig(ctx, &Req{
			Config: serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
				Id: "nope",
			}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("create with target runner labels", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertOnDemandRunnerConfig(ctx, &Req{
			Config: serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Labels{
						Labels: &pb.Ref_RunnerLabels{
							Labels: map[string]string{
								"test": "test",
							},
						},
					},
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Config
		require.NotEmpty(result.Id)

		resp, err = client.UpsertOnDemandRunnerConfig(ctx, &Req{
			Config: result,
		})
		require.NoError(err)
		require.NotNil(resp)
	})
}

func TestServiceOnDemandRunnerConfig_GetOnDemandRunnerConfig(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Best way to mock for now is to make a request
	resp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: serverptypes.TestOnDemandRunnerConfig(t, nil),
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.GetOnDemandRunnerConfigRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		resp, err := client.GetOnDemandRunnerConfig(ctx, &Req{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Id: resp.Config.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.Config)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetOnDemandRunnerConfig(ctx, &Req{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Id: "nope",
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceOnDemandRunnerConfig_ListOnDemandRunnerConfigs(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	dep := serverptypes.TestOnDemandRunnerConfig(t, nil)

	// Best way to mock for now is to make a request
	resp, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: dep,
	})

	require.NoError(t, err)

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a ondemand runner config
		deployments, err := client.ListOnDemandRunnerConfigs(ctx, &emptypb.Empty{})
		require.NoError(err)
		require.NotEmpty(deployments)
		require.Equal(deployments.Configs[0].Id, resp.Config.Id)
	})
}
