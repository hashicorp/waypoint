package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlatformConfig(t *testing.T) {
	t.Run("empty is fine", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{},
		}

		require.NoError(t, p.ConfigSet(cfg))
	})

	t.Run("fine if only cert", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				CertificateId: "xyz",
			},
		}

		require.NoError(t, p.ConfigSet(cfg))
	})

	t.Run("fine if only listener", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				ListenerARN: "xyz",
			},
		}

		require.NoError(t, p.ConfigSet(cfg))
	})

	t.Run("errors if cert and listener are set", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				CertificateId: "xyz",
				ListenerARN:   "abc",
			},
		}

		require.Error(t, p.ConfigSet(cfg))
	})

	t.Run("errors if zone_id and fqdn and listener are set", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				ZoneId:      "xyz",
				FQDN:        "a.b",
				ListenerARN: "abc",
			},
		}

		require.Error(t, p.ConfigSet(cfg))
	})

	t.Run("errors if zone_id but not fqdn are set", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				ZoneId: "xyz",
			},
		}

		require.Error(t, p.ConfigSet(cfg))
	})

	t.Run("errors if fqdn but not zone_id are set", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				FQDN: "xyz",
			},
		}

		require.Error(t, p.ConfigSet(cfg))
	})

	t.Run("fine with just zone and fqdn", func(t *testing.T) {
		var p Platform

		cfg := &Config{
			ALB: &ALBConfig{
				ZoneId: "xyz",
				FQDN:   "a.b",
			},
		}

		require.NoError(t, p.ConfigSet(cfg))
	})
}
