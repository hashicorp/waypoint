// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestApp creates the app in the DB.
func TestApp(t *testing.T, client pb.WaypointClient, ref *pb.Ref_Application) {
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name: ref.Project,
			},
		})
		require.NoError(t, err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(t, err)
	}
}
