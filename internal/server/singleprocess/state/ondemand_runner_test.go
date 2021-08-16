package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestOndemandRunner(t *testing.T) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		_, err := s.OndemandRunnerGet(&pb.Ref_OndemandRunner{
			Id: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.OndemandRunnerPut(serverptypes.TestOndemandRunner(t, &pb.OndemandRunner{
			Id:     "foo",
			OciUrl: "h/w:s",
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.OndemandRunnerGet(&pb.Ref_OndemandRunner{
				Id: "foo",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get case insensitive
		{
			resp, err := s.OndemandRunnerGet(&pb.Ref_OndemandRunner{
				Id: "Foo",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.OndemandRunnerList()
			require.NoError(err)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.OndemandRunnerPut(serverptypes.TestOndemandRunner(t, &pb.OndemandRunner{
			Id: "AbCdE",
		}))
		require.NoError(err)

		// Read
		resp, err := s.OndemandRunnerGet(&pb.Ref_OndemandRunner{
			Id: "AbCdE",
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.OndemandRunnerDelete(&pb.Ref_OndemandRunner{
				Id: "AbCdE",
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.OndemandRunnerGet(&pb.Ref_OndemandRunner{
				Id: "AbCdE",
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.OndemandRunnerList()
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}
