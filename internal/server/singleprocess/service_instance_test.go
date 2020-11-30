package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestListInstances(t *testing.T) {
	ctx := context.Background()

	t.Run("deployment ID", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// List instances and it should be empty
		{
			resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
				Scope: &pb.ListInstancesRequest_DeploymentId{
					DeploymentId: dep.Id,
				},
				ConnectTimeout: "0s",
			})
			require.NoError(err)
			require.Len(resp.Instances, 0)
		}

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		// Wait for the first config so that we know we're registered
		_, err = stream.Recv()
		require.NoError(err)

		// List instances and it should exist
		{
			resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
				Scope: &pb.ListInstancesRequest_DeploymentId{
					DeploymentId: dep.Id,
				},
				ConnectTimeout: "0s",
			})
			require.NoError(err)
			require.Len(resp.Instances, 1)
		}
	})
}
