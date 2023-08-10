// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ecrpull

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderConfig(t *testing.T) {
	t.Run("requires repository and tag", func(t *testing.T) {
		var b Builder
		cfg := &Config{}
		require.EqualError(t, b.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = Repository: cannot be blank; Tag: cannot be blank.")
	})

	t.Run("disallows unsupported architecture", func(t *testing.T) {
		var b Builder
		cfg := &Config{
			Repository:        "foo",
			Tag:               "latest",
			ForceArchitecture: "foobar",
		}

		require.EqualError(t, b.ConfigSet(cfg), "rpc error: code = InvalidArgument desc = ForceArchitecture: Unsupported force_architecture \"foobar\". Must be one of [\"x86_64\", \"arm64\"], or left blank.")
	})
}
