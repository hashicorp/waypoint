// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServiceDeployment_URLService(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)), TestWithURLService(t, nil))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.UpsertDeploymentRequest

	t.Run("default hostname", func(t *testing.T) {
		require := require.New(t)

		deploy := serverptypes.TestValidDeployment(t, nil)

		// Create, should get an ID back
		resp, err := client.UpsertDeployment(ctx, &Req{
			Deployment: deploy,
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Deployment
		require.NotEmpty(result.Id)

		// Should have the hostname
		{
			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{})
			require.NoError(err)
			require.Len(resp.Hostnames, 1)
		}

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

		// Should have the hostname
		{
			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{})
			require.NoError(err)
			require.Len(resp.Hostnames, 1)
		}
	})
}
