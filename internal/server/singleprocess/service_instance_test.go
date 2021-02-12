package singleprocess

import (
	"context"
	"testing"
	"time"

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
			})
			require.NoError(err)
			require.Len(resp.Instances, 1)
		}
	})

	t.Run("deployment ID with wait timeout defined", func(t *testing.T) {
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
		resultCh := make(chan []*pb.Instance, 1)
		errCh := make(chan error, 1)
		go func() {
			resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
				Scope: &pb.ListInstancesRequest_DeploymentId{
					DeploymentId: dep.Id,
				},
				WaitTimeout: "3s",
			})

			if err != nil {
				errCh <- err
			} else {
				resultCh <- resp.Instances

			}
		}()

		select {
		case <-resultCh:
			t.Fatal("Should not have got value from listing instances")
		case <-errCh:
			t.Fatal("Should not have got error from listing instances")
		case <-time.After(250 * time.Millisecond):
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
		select {
		case result := <-resultCh:
			require.Len(result, 1)
		case err := <-errCh:
			require.NoError(err)
		case <-time.After(1 * time.Second):
			t.Fatal("We should have got value from listing instances")
		}
	})
}
