package copy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidLink(t *testing.T) {
	t.Run("normal relative", func(t *testing.T) {
		r := require.New(t)

		target, err := evalSymlink("/a/b", "c")
		r.NoError(err)

		r.Equal("/a/c", target)
	})

	t.Run("normal upward relative", func(t *testing.T) {
		r := require.New(t)

		target, err := evalSymlink("/a/b", "../a/c")
		r.NoError(err)

		r.Equal("/a/c", target)
	})

	t.Run("normal relative upward and downward", func(t *testing.T) {
		r := require.New(t)

		target, err := evalSymlink("/a/b", "d/../c")
		r.NoError(err)

		r.Equal("/a/c", target)
	})

	t.Run("normal absolute", func(t *testing.T) {
		r := require.New(t)

		target, err := evalSymlink("/a/b", "/a/c")
		r.NoError(err)

		r.Equal("/a/c", target)
	})

	t.Run("invalid relative", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", "../../c")
		r.Error(err)
	})

	t.Run("invalid relative root dir", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", ".")
		r.Error(err)
	})

	t.Run("invalid absolute", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", "/c")
		r.Error(err)
	})

	t.Run("invalid absolute root dir", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", "/a")
		r.Error(err)
	})

	t.Run("invalid relative upward and downward", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", "b/../../c")
		r.Error(err)
	})

	t.Run("invalid absolute with upward", func(t *testing.T) {
		r := require.New(t)

		_, err := evalSymlink("/a/b", "/a/../c")
		r.Error(err)
	})

}
