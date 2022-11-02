package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["server_config"] = []testFunc{
		TestServerConfig,
	}
}

func TestServerConfig(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.ServerConfigSet(ctx, &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{},
		})
		if err != nil {
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.Unavailable, s.Code())
			return
		}

		var cookie string
		{
			// Get
			cfg, err := s.ServerConfigGet(ctx)
			require.NoError(err)
			require.NotNil(cfg)
			require.NotNil(cfg.AdvertiseAddrs)
			require.NotEmpty(cfg.Cookie)

			cookie = cfg.Cookie
		}

		// Unset
		require.NoError(s.ServerConfigSet(ctx, nil))

		{
			// Get
			cfg, err := s.ServerConfigGet(ctx)
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
		err := s.ServerConfigSet(ctx, &pb.ServerConfig{
			Cookie: "hello",
		})
		if err != nil {
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.Unavailable, s.Code())
			return
		}

		{
			// Get
			cfg, err := s.ServerConfigGet(ctx)
			require.NoError(err)
			require.NotNil(cfg)
			require.NotEqual("hello", cfg.Cookie)
		}
	})
}
