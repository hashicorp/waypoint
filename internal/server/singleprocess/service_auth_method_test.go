package singleprocess

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestListOIDCAuthMethods(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.UpsertAuthMethodRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create
		{
			resp, err := client.UpsertAuthMethod(ctx, &Req{
				AuthMethod: serverptypes.TestAuthMethod(t, &pb.AuthMethod{
					Name: "A",
				}),
			})
			require.NoError(err)
			require.NotNil(resp)
		}
		{
			resp, err := client.UpsertAuthMethod(ctx, &Req{
				AuthMethod: serverptypes.TestAuthMethod(t, &pb.AuthMethod{
					Name: "B",
				}),
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		resp, err := client.ListOIDCAuthMethods(ctx, &empty.Empty{})
		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.AuthMethods, 2)
		for _, method := range resp.AuthMethods {
			require.NotEmpty(method.Name)
		}
	})
}
