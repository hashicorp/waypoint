package singleprocess

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func TestServiceBuild(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	var id string
	{
		require := require.New(t)

		// Start a build
		resp, err := client.CreateBuild(ctx, &pb.CreateBuildRequest{
			Component: &pb.Component{
				Type: pb.Component_BUILDER,
				Name: "packer",
			},
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.Id)
		id = resp.Id
	}

	{
		require := require.New(t)

		// List builds
		resp, err := client.ListBuilds(ctx, &empty.Empty{})
		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.Builds, 1)

		build := resp.Builds[0]
		require.Equal(id, build.Id)
		require.Equal(pb.Status_RUNNING, build.Status.State)
	}

	{
		require := require.New(t)

		// Complete the build
		_, err := client.CompleteBuild(ctx, &pb.CompleteBuildRequest{
			Id: id,
			Result: &pb.CompleteBuildRequest_Error{
				Error: status.Newf(codes.DataLoss, "oh no!").Proto(),
			},
		})
		require.NoError(err)

		// Get the build to verify the state
		resp, err := client.ListBuilds(ctx, &empty.Empty{})
		require.NoError(err)
		require.Len(resp.Builds, 1)
		build := resp.Builds[0]
		require.Equal(id, build.Id)
		require.Equal(pb.Status_ERROR, build.Status.State)
		require.Equal(codes.DataLoss, status.FromProto(build.Status.Error).Code())
	}

	t.Run("complete a non-existent build", func(t *testing.T) {
		require := require.New(t)

		resp, err := client.CompleteBuild(ctx, &pb.CompleteBuildRequest{
			Id: "nope",
		})
		require.Error(err)
		require.Nil(resp)

		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}
