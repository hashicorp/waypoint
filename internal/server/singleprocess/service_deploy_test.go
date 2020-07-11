package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceDeployment(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.UpsertDeploymentRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertDeployment(ctx, &Req{
			Deployment: serverptypes.TestValidDeployment(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Deployment
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertDeployment(ctx, &Req{
			Deployment: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Deployment
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertDeployment(ctx, &Req{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
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

func TestServiceDeployment_GetDeployment(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Best way to mock for now is to make a request
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.GetDeploymentRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		deployment, err := client.GetDeployment(ctx, &Req{
			DeploymentId: resp.Deployment.Id,
		})
		require.NoError(err)
		require.NotNil(deployment)
		require.NotEmpty(deployment.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetDeployment(ctx, &Req{
			DeploymentId: "nope",
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}
