// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func init() {
	tests["user"] = []testFunc{
		TestUser,
	}
}

func TestUser(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	id, err := ulid()
	require.NoError(t, err)

	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		_, err := s.UserGet(ctx, &pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: "foo"},
			},
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Empty", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		{
			empty, err := s.UserEmpty(ctx)
			require.NoError(err)
			require.True(empty)
		}

		// Set
		err = s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		{
			empty, err := s.UserEmpty(ctx)
			require.NoError(err)
			require.False(empty)
		}
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err = s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.UserGet(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get by username
		{
			resp, err := s.UserGet(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Username{
					Username: &pb.Ref_UserUsername{
						Username: "foo",
					},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get by username, case insensitive
		{
			resp, err := s.UserGet(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Username{
					Username: &pb.Ref_UserUsername{
						Username: "Foo",
					},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.UserList(ctx)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// We need two users
		require.NoError(s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Username: "bar",
		})))

		// Set
		err := s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Read
		resp, err := s.UserGet(ctx, &pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: id},
			},
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.UserDelete(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.UserGet(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.UserList(ctx)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete last user fails", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Read
		resp, err := s.UserGet(ctx, &pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: id},
			},
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.UserDelete(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.Error(err)
			require.Equal(codes.FailedPrecondition, status.Code(err))
		}
	})

	t.Run("Delete default user fails", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       serverstate.DefaultUserId,
			Username: "foo",
		}))
		require.NoError(err)

		// Delete
		{
			err := s.UserDelete(ctx, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: serverstate.DefaultUserId},
				},
			})
			require.Error(err)
			require.Equal(codes.FailedPrecondition, status.Code(err))
		}
	})

	t.Run("User lookup by OIDC", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.UserPut(ctx, serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
			Links: []*pb.User_Link{
				{
					Method: &pb.User_Link_Oidc{
						Oidc: &pb.User_Link_OIDC{
							Iss: "A",
							Sub: "B",
						},
					},
				},
			},
		})))

		// Read
		{
			resp, err := s.UserGetOIDC(ctx, "A", "B")
			require.NoError(err)
			require.NotNil(resp)
		}

		// Not matching issuer
		{
			resp, err := s.UserGetOIDC(ctx, "B", "B")
			require.Error(err)
			require.Nil(resp)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// Not matching sub
		{
			resp, err := s.UserGetOIDC(ctx, "A", "C")
			require.Error(err)
			require.Nil(resp)
			require.Equal(codes.NotFound, status.Code(err))
		}
	})
}
