// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["user"] = []testFunc{
		TestServiceUpdateUser,
	}
}

func TestServiceUpdateUser(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Get ourself
	getUser := func(t *testing.T) *pb.User {
		userResp, err := client.GetUser(ctx, &pb.GetUserRequest{})
		require.NoError(t, err)
		return userResp.User
	}

	// Simplify writing tests
	type Req = pb.UpdateUserRequest

	t.Run("update but do nothing", func(t *testing.T) {
		require := require.New(t)
		user := getUser(t)

		resp, err := client.UpdateUser(ctx, &Req{User: user})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(user.Id, resp.User.Id)
	})

	t.Run("update without a user set", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		_, err := client.UpdateUser(ctx, &Req{})
		require.Error(err)
	})

	t.Run("update with an invalid ID", func(t *testing.T) {
		require := require.New(t)

		user := getUser(t)
		user.Id = "NOPE"

		_, err := client.UpdateUser(ctx, &Req{User: user})
		require.Error(err)
	})

	t.Run("change username", func(t *testing.T) {
		require := require.New(t)

		user := getUser(t)
		user.Username = "tubes"

		resp, err := client.UpdateUser(ctx, &Req{User: user})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(user.Id, resp.User.Id)
		require.Equal("tubes", resp.User.Username)
	})
}
