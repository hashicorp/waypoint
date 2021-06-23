package env

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetEnvBool(t *testing.T) {
	envVarTestKey := "WAYPOINT_GET_ENV_BOOL_TEST"
	require := require.New(t)

	t.Run("Unset env var returns default", func (t *testing.T) {
		require.True(GetEnvBool(envVarTestKey, true))
		require.False(GetEnvBool(envVarTestKey, false))
	})

	t.Run("Empty env var returns default", func (t *testing.T) {
		os.Setenv(envVarTestKey, "")
		require.True(GetEnvBool(envVarTestKey, true))
		require.False(GetEnvBool(envVarTestKey, false))
	})

	t.Run("Non-truthy env var returns default", func (t *testing.T) {
		os.Setenv(envVarTestKey, "unparseable")
		require.True(GetEnvBool(envVarTestKey, true))
		require.False(GetEnvBool(envVarTestKey, false))
	})

	t.Run("true/false env vars return non-default", func (t *testing.T) {
		os.Setenv(envVarTestKey, "true")
		require.True(GetEnvBool(envVarTestKey, false))

		os.Setenv(envVarTestKey, "false")
		require.False(GetEnvBool(envVarTestKey, true))
	})

	t.Run("1 evaluates as true", func (t *testing.T) {
		os.Setenv(envVarTestKey, "1")
		require.True(GetEnvBool(envVarTestKey, false))
		require.True(GetEnvBool(envVarTestKey, true))
	})
}
