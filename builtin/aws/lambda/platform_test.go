package lambda

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlatformConfig(t *testing.T) {
	t.Run("empty is fine", func(t *testing.T) {
		var p Platform
		cfg := &Config{}
		require.NoError(t, p.ConfigSet(cfg))
	})

	t.Run("disallows unsupported architecture", func(t *testing.T) {
		var p Platform
		cfg := &Config{
			Architecture: "foobar",
		}
		require.Error(t, p.ConfigSet(cfg))
	})

	t.Run("disallows invalid timeout", func(t *testing.T) {
		var p Platform
		{
			cfg := &Config{
				Timeout: 901,
			}
			require.Error(t, p.ConfigSet(cfg))
		}

		{
			cfg := &Config{
				Timeout: -1,
			}
			require.Error(t, p.ConfigSet(cfg))
		}
	})
}
