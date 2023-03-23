// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverconfig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientEnv(t *testing.T) {
	listToMap := func(t *testing.T, vs []string) map[string]string {
		result := map[string]string{}
		for _, v := range vs {
			idx := strings.Index(v, "=")
			key := v[:idx]
			value := v[idx+1:]
			result[key] = value
		}

		return result
	}

	t.Run("no require auth, token set", func(t *testing.T) {
		require := require.New(t)

		env := listToMap(t, (&Client{
			Address:       "foo",
			Tls:           true,
			TlsSkipVerify: true,
			AuthToken:     "bar",
		}).Env())

		require.Equal(env["WAYPOINT_SERVER_ADDR"], "foo")
		require.Equal(env["WAYPOINT_SERVER_TLS"], "true")
		require.Equal(env["WAYPOINT_SERVER_TLS_SKIP_VERIFY"], "true")
		require.Empty(env["WAYPOINT_SERVER_TOKEN"])
	})

	t.Run("require auth, token set", func(t *testing.T) {
		require := require.New(t)

		env := listToMap(t, (&Client{
			Address:       "foo",
			Tls:           true,
			TlsSkipVerify: true,
			AuthToken:     "bar",
			RequireAuth:   true,
		}).Env())

		require.Equal(env["WAYPOINT_SERVER_ADDR"], "foo")
		require.Equal(env["WAYPOINT_SERVER_TLS"], "true")
		require.Equal(env["WAYPOINT_SERVER_TLS_SKIP_VERIFY"], "true")
		require.Equal(env["WAYPOINT_SERVER_TOKEN"], "bar")
	})
}
