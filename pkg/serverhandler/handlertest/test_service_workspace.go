// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["workspace"] = []testFunc{
		TestWorkspace_Upsert,
	}
}

func TestWorkspace_Upsert(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertWorkspaceRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create
		{
			resp, err := client.UpsertWorkspace(ctx, &Req{
				Workspace: serverptypes.TestWorkspace(t, &pb.Workspace{
					Name: "staging",
				}),
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Create another
		{
			resp, err := client.UpsertWorkspace(ctx, &Req{
				Workspace: serverptypes.TestWorkspace(t, &pb.Workspace{
					Name: "dev",
				}),
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := client.ListWorkspaces(ctx, &pb.ListWorkspacesRequest{})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Workspaces, 2)
			for _, workspace := range resp.Workspaces {
				require.NotEmpty(workspace.Name)
			}
		}

		// Get dev
		{
			resp, err := client.GetWorkspace(ctx, &pb.GetWorkspaceRequest{
				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Workspace.Name, "dev")
		}

		// Fail with bad Workspace name
		{
			resp, err := client.UpsertWorkspace(ctx, &Req{
				Workspace: serverptypes.TestWorkspace(t, &pb.Workspace{
					Name: "a bad name",
				}),
			})
			require.Error(err)
			require.Nil(resp)
		}
	})
}
