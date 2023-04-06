// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

		require.EqualError(t, p.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = Architecture: Unsupported function architecture \"foobar\". Must be one of [\"x86_64\", \"arm64\"], or left blank.")
	})

	t.Run("disallows invalid timeout", func(t *testing.T) {
		var p Platform
		{
			cfg := &Config{
				Timeout: 901,
			}
			require.EqualError(t, p.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = Timeout: Timeout must be less than or equal to 15 minutes.")
		}

		{
			cfg := &Config{
				Timeout: -1,
			}
			require.EqualError(t, p.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = Timeout: Timeout must not be negative.")
		}
	})

	t.Run("disallows invalid storagemb", func(t *testing.T) {
		var p Platform
		{
			cfg := &Config{
				StorageMB: 100,
			}
			require.EqualError(t, p.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = StorageMB: Storage must a value between 512 and 10240.")
		}

		{
			cfg := &Config{
				StorageMB: 20000,
			}
			require.EqualError(t, p.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = StorageMB: Storage must a value between 512 and 10240.")
		}
	})
}
