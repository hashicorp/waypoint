// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	tests["HMAC"] = []testFunc{TestHMAC}
}

func TestHMAC(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Get returns nil if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		result, err := s.HMACKeyGet(ctx, "foo")
		require.NoError(err)
		require.Nil(result)
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		key, err := s.HMACKeyCreateIfNotExist(ctx, "foo", 32)
		require.NoError(err)
		require.NotNil(key)
		require.NotEmpty(key.Key)

		// Get exact
		{
			resp, err := s.HMACKeyGet(ctx, "foo")
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Key, key.Key)
		}

		// Get case insensitive
		{
			resp, err := s.HMACKeyGet(ctx, "fOo")
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Key, key.Key)
		}

		{
			// Set should return identical key
			key2, err := s.HMACKeyCreateIfNotExist(ctx, "foo", 32)
			require.NoError(err)
			require.NotNil(key2)
			require.Equal(key2.Key, key.Key)
		}
	})
}
