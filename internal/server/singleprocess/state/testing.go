package state

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestState returns an initialized State for testing.
func TestState(t testing.T) *State {
	result, err := New(hclog.L(), testDB(t))
	require.NoError(t, err)
	return result
}

func testDB(t testing.T) *bolt.DB {
	t.Helper()

	// Temporary directory for the database
	td, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(td) })

	// Create the DB
	db, err := bolt.Open(filepath.Join(td, "test.db"), 0600, nil)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	return db
}
