package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestUser(t *testing.T) {
	id, err := ulid()
	require.NoError(t, err)

	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		_, err := s.UserGet(&pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: "foo"},
			},
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err = s.UserPut(serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.UserGet(&pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get by username
		{
			resp, err := s.UserGet(&pb.Ref_User{
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
			resp, err := s.UserGet(&pb.Ref_User{
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
			resp, err := s.UserList()
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// We need two users
		require.NoError(s.UserPut(serverptypes.TestUser(t, &pb.User{
			Username: "bar",
		})))

		// Set
		err := s.UserPut(serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Read
		resp, err := s.UserGet(&pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: id},
			},
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.UserDelete(&pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.UserGet(&pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.UserList()
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete last user fails", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.UserPut(serverptypes.TestUser(t, &pb.User{
			Id:       id,
			Username: "foo",
		}))
		require.NoError(err)

		// Read
		resp, err := s.UserGet(&pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: id},
			},
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.UserDelete(&pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: id},
				},
			})
			require.Error(err)
			require.Equal(codes.FailedPrecondition, status.Code(err))
		}
	})
}
