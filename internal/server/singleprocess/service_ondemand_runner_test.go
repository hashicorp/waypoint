package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServiceOnDemandRunnerConfig(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

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
}

func TestServiceOnDemandRunnerConfig_GetOnDemandRunnerConfig(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

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

func TestServiceOnDemandRunnerConfig_ListOnDemandRunnerConfigs(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

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
