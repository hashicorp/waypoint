package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestProject(t *testing.T) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		_, err := s.ProjectGet(&pb.Ref_Project{
			Project: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "AbCdE",
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get case insensitive
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "abcDe",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.ProjectList()
			require.NoError(err)
			require.Len(resp, 1)
		}
	})

	t.Run("Put does not modify applications", func(t *testing.T) {
		require := require.New(t)

		const name = "AbCdE"
		ref := &pb.Ref_Project{Project: name}

		s := TestState(t)
		defer s.Close()

		// Set
		proj := serverptypes.TestProject(t, &pb.Project{Name: name})
		err := s.ProjectPut(proj)
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test",
			Project: ref,
		}))
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test2",
			Project: ref,
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.False(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}

		// Update the project
		proj.RemoteEnabled = true
		require.NoError(s.ProjectPut(proj))

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.True(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "AbCdE",
		}))
		require.NoError(err)

		// Read
		resp, err := s.ProjectGet(&pb.Ref_Project{
			Project: "AbCdE",
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.ProjectDelete(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.ProjectList()
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}
