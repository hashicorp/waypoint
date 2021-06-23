package env

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetEnvBool(t *testing.T) {
	envVarTestKey := "WAYPOINT_GET_ENV_BOOL_TEST"
	require := require.New(t)

	t.Run("Unset env var returns default", func(t *testing.T) {
		b, err := GetEnvBool(envVarTestKey, true)
		require.NoError(err)
		require.True(b)

		b, err = GetEnvBool(envVarTestKey, false)
		require.NoError(err)
		require.False(b)
	})

	t.Run("Empty env var returns default", func(t *testing.T) {
		os.Setenv(envVarTestKey, "")
		b, err := GetEnvBool(envVarTestKey, true)
		require.NoError(err)
		require.True(b)

		b, err = GetEnvBool(envVarTestKey, false)
		require.NoError(err)
		require.False(b)
	})

	t.Run("Non-truthy env var returns an error", func(t *testing.T) {
		os.Setenv(envVarTestKey, "unparseable")
		_, err := GetEnvBool(envVarTestKey, true)
		require.Error(err)
	})

	t.Run("true/false env vars return non-default", func(t *testing.T) {
		os.Setenv(envVarTestKey, "true")
		b, err := GetEnvBool(envVarTestKey, false)
		require.NoError(err)
		require.True(b)

		os.Setenv(envVarTestKey, "false")
		b, err = GetEnvBool(envVarTestKey, true)
		require.NoError(err)
		require.False(b)
	})

	t.Run("boolean parsing is generous with capitalization", func(t *testing.T) {
		os.Setenv(envVarTestKey, "tRuE")
		b, err := GetEnvBool(envVarTestKey, false)
		require.NoError(err)
		require.True(b)
	})

	t.Run("1 evaluates as true", func(t *testing.T) {
		os.Setenv(envVarTestKey, "1")
		b, err := GetEnvBool(envVarTestKey, false)
		require.NoError(err)
		require.True(b)
	})
}
