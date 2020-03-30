package mapper

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestFunc_basic(t *testing.T) {
	require := require.New(t)

	addTwo := func(a int) int { return a + 2 }
	f, err := NewFunc(addTwo)
	require.NoError(err)
	result, err := f.Call(1)
	require.NoError(err)
	require.Equal(result, 3)
}

func TestFunc_error(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		require := require.New(t)

		addTwo := func(a int) (int, error) { return a + 2, nil }
		f, err := NewFunc(addTwo)
		require.NoError(err)
		result, err := f.Call(1)
		require.NoError(err)
		require.Equal(result, 3)
	})

	t.Run("nil error", func(t *testing.T) {
		require := require.New(t)

		addTwo := func(a int) (int, error) { return a + 2, errors.New("error!") }
		f, err := NewFunc(addTwo)
		require.NoError(err)
		result, err := f.Call(1)
		require.Error(err)
		require.Equal(result, 3)
	})
}

func TestFunc_hclog(t *testing.T) {
	require := require.New(t)

	factory := func(log hclog.Logger) int { return 42 }
	f, err := NewFunc(factory)
	require.NoError(err)
	result, err := f.Call(hclog.L())
	require.NoError(err)
	require.Equal(result, 42)
}

func TestFunc_values(t *testing.T) {
	require := require.New(t)

	f := TestFunc(t, func(a int) int { return a + 2 },
		WithValues(int(12)))
	result, err := f.Call()
	require.NoError(err)
	require.Equal(14, result)
}

func mustFunc(t *testing.T, f interface{}) *Func {
	t.Helper()

	result, err := NewFunc(f)
	require.NoError(t, err)
	return result
}
