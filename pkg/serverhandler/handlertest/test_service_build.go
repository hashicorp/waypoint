package handlertest

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

func init() {
	tests["build"] = []testFunc{
		TestServiceBuild,
	}
}
func TestServiceBuild(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

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

	t.Run("create and delete", func(t *testing.T) {
		require := require.New(t)

		resp, err := client.UpsertBuild(ctx, &Req{
			Build: serverptypes.TestValidBuild(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Build
		require.NotEmpty(result.Id)

		_, err = client.DeleteBuild(ctx, &pb.DeleteBuildRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: resp.Build.Id,
				},
			},
		})
		require.Nil(err)
	})
}
