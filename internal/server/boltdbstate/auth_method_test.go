// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestAuthMethod(t *testing.T) {
	ctx := context.Background()
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		_, err := s.AuthMethodGet(ctx, &pb.Ref_AuthMethod{Name: "foo"})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.AuthMethodPut(ctx, serverptypes.TestAuthMethod(t, &pb.AuthMethod{
			Name: "foo",
		}))
		require.NoError(err)

		// Get by name
		{
			resp, err := s.AuthMethodGet(ctx, &pb.Ref_AuthMethod{Name: "foo"})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get by name, case insensitive
		{
			resp, err := s.AuthMethodGet(ctx, &pb.Ref_AuthMethod{Name: "Foo"})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.AuthMethodList(ctx)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// We need two methods
		require.NoError(s.AuthMethodPut(ctx, serverptypes.TestAuthMethod(t, &pb.AuthMethod{
			Name: "bar",
		})))

		// Set
		err := s.AuthMethodPut(ctx, serverptypes.TestAuthMethod(t, &pb.AuthMethod{
			Name: "baz",
		}))
		require.NoError(err)

		// Read
		resp, err := s.AuthMethodGet(ctx, &pb.Ref_AuthMethod{Name: "bar"})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.AuthMethodDelete(ctx, &pb.Ref_AuthMethod{Name: "bar"})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.AuthMethodGet(ctx, &pb.Ref_AuthMethod{Name: "bar"})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.AuthMethodList(ctx)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})
}
