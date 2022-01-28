package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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

		var cookie string
		{
			// Get
			cfg, err := s.ServerConfigGet()
			require.NoError(err)
			require.NotNil(cfg)
			require.NotNil(cfg.AdvertiseAddrs)
			require.NotEmpty(cfg.Cookie)

			cookie = cfg.Cookie
		}

		// Unset
		require.NoError(s.ServerConfigSet(nil))

		{
			// Get
			cfg, err := s.ServerConfigGet()
			require.NoError(err)
			require.NotNil(cfg)
			require.Nil(cfg.AdvertiseAddrs)
			require.NotEmpty(cfg.Cookie)
			require.Equal(cookie, cfg.Cookie)
		}
	})

	t.Run("cookie cannot be overwritten", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ServerConfigSet(&pb.ServerConfig{
			Cookie: "hello",
		}))

		{
			// Get
			cfg, err := s.ServerConfigGet()
			require.NoError(err)
			require.NotNil(cfg)
			require.NotEqual("hello", cfg.Cookie)
		}
	})
}
