// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["auth_method"] = []testFunc{
		TestListOIDCAuthMethods,
	}
}

func TestListOIDCAuthMethods(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

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
