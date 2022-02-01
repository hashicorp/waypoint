package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	tests["server_urltoken"] = []testFunc{
		TestServerURLToken,
	}
}

func TestServerURLToken(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("set and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ServerURLTokenSet("foo"))

		str, err := s.ServerURLTokenGet()
		require.NoError(err)
		require.Equal("foo", str)

	})
}
