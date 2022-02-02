package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServiceBuild(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.UpsertBuildRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertBuild(ctx, &Req{
			Build: serverptypes.TestValidBuild(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Build
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertBuild(ctx, &Req{
			Build: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Build
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertBuild(ctx, &Req{
			Build: serverptypes.TestValidBuild(t, &pb.Build{Id: "nope"}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}
