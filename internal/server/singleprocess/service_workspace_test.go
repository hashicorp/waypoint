package singleprocess

import (
	"context"
	"testing"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/stretchr/testify/require"
)

func TestUpsertWorkspace(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

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
	})
}
