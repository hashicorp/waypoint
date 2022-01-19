package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal_nomore/server/gen"
)

func init() {
	tests["server_config"] = []testFunc{
		TestServerConfig,
	}
}

func TestServerConfig(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ServerConfigSet(&pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{},
		}))

		{
			// Get
			cfg, err := s.ServerConfigGet()
			require.NoError(err)
			require.NotNil(cfg)
			require.NotNil(cfg.AdvertiseAddrs)
		}

		// Unset
		require.NoError(s.ServerConfigSet(nil))

		{
			// Get
			cfg, err := s.ServerConfigGet()
			require.NoError(err)
			require.NotNil(cfg)
			require.Nil(cfg.AdvertiseAddrs)
		}
	})
}
