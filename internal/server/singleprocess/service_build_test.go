package singleprocess

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func TestServiceBuild(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	var id string
	{
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
		// List builds
		resp, err := client.ListBuilds(ctx, &empty.Empty{})
		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.Builds, 1)

		build := resp.Builds[0]
		require.Equal(id, build.Id)
	}
}
