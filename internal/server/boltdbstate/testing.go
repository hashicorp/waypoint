// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

// TestState returns an initialized State for testing.
func TestState(t testing.T) *State {
	result, err := New(hclog.L(), testDB(t))
	require.NoError(t, err)
	return result
}

// TestStateRestart closes the given state and restarts it against the
// same DB file.
func TestStateRestart(t testing.T, s *State) (*State, error) {
	path := s.db.Path()
	require.NoError(t, s.Close())

	// Open the new DB
	db, err := bolt.Open(path, 0600, nil)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Init new state
	return New(hclog.L(), db)
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
