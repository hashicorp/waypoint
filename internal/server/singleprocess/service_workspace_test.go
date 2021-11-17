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
					Name: "default",
				}),
			})
			require.NoError(err)
			require.NotNil(resp)
		}
		// },
		// 	{
		// 		resp, err := client.UpsertAuthMethod(ctx, &Req{
		// 			AuthMethod: serverptypes.TestAuthMethod(t, &pb.AuthMethod{
		// 				Name: "B",
		// 			}),
		// 		})
		// 		require.NoError(err)
		// 		require.NotNil(resp)
		// 	}

		// 	// List
		// 	resp, err := client.ListOIDCAuthMethods(ctx, &empty.Empty{})
		// 	require.NoError(err)
		// 	require.NotNil(resp)
		// 	require.Len(resp.AuthMethods, 2)
		// 	for _, method := range resp.AuthMethods {
		// 		require.NotEmpty(method.Name)
		// 	}
	})
}
