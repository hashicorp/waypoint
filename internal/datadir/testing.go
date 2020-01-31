package datadir

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestDir returns a Dir for testing.
func TestDir(t testing.T) (Dir, func()) {
	t.Helper()

	td, err := ioutil.TempDir("", "datadir-test")
	require.NoError(t, err)

	dir, err := newRootDir(td)
	require.NoError(t, err)

	return dir, func() { os.RemoveAll(td) }
}
