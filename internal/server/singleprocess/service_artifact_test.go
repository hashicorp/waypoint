package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func TestServiceArtifact(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.UpsertPushedArtifactRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPushedArtifact(ctx, &Req{
			Artifact: &pb.PushedArtifact{},
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Artifact
		require.NotEmpty(result.Id)
		require.Nil(result.Status)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertPushedArtifact(ctx, &Req{
			Artifact: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Artifact
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPushedArtifact(ctx, &Req{
			Artifact: &pb.PushedArtifact{Id: "nope"},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}
