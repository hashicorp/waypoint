// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// This test ensures that empty hashes return the expected value. This
// is important because we set a base key to ensure empty maps don't hash
// to 0.
func TestRunnerLabelHash_empty(t *testing.T) {
	{
		h, err := RunnerLabelHash(nil)
		require.NoError(t, err)
		require.Equal(t, h, uint64(0x85d03bbbdf8bbf66))
	}

	{
		h, err := RunnerLabelHash(map[string]string{})
		require.NoError(t, err)
		require.Equal(t, h, uint64(0x85d03bbbdf8bbf66))
	}
}
