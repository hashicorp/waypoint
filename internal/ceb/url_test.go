// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ceb

import (
	"context"
	"testing"
	"time"

	hznpb "github.com/hashicorp/horizon/pkg/pb"
	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestCEB_url(t *testing.T) {
	ctx := context.Background()

	t.Run("with URL services disabled", func(t *testing.T) {
		require := require.New(t)

		// Start CEB
		data := TestCEB(t)

		// Should not have any
		resp, err := data.Horizon.ControlServer.ListServices(ctx, &hznpb.ListServicesRequest{
			Account: data.Horizon.Account,
		})
		require.NoError(err)
		require.Empty(resp.Services)
	})

	t.Run("with URL services enabled", func(t *testing.T) {
		require := require.New(t)

		// Start our server
		var hzn hzntest.DevSetup
		client := singleprocess.TestServer(t, singleprocess.TestWithURLService(t, &hzn))

		// Create our deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Start CEB
		testChenv(t, envDeploymentId, dep.Id)
		testChenv(t, "PORT", "1234")
		TestCEB(t, WithClient(client))

		// Should have services
		require.Eventually(func() bool {
			resp, err := hzn.ControlServer.ListServices(ctx, &hznpb.ListServicesRequest{
				Account: hzn.Account,
			})
			require.NoError(err)
			return len(resp.Services) > 0
		}, 5*time.Second, 20*time.Millisecond)
	})
}
