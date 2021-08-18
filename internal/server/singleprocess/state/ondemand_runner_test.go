package state

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestOnDemandRunnerConfig(t *testing.T) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		_, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{
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
		rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
			OciUrl: "h/w:s",
		})

		err := s.OnDemandRunnerConfigPut(rec)
		require.NoError(err)

		// Get exact
		{
			resp, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get case insensitive
		{
			resp, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{
				Id: strings.ToUpper(rec.Id),
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.OnDemandRunnerConfigList()
			require.NoError(err)
			require.Len(resp, 1)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{})

		err := s.OnDemandRunnerConfigPut(rec)
		require.NoError(err)

		// Read
		resp, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{
			Id: rec.Id,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.OnDemandRunnerConfigDelete(&pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.OnDemandRunnerConfigList()
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}
