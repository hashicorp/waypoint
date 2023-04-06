// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clicontext

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestStorage returns a *Storage pointed at a temporary directory. This
// will cleanup automatically by using t.Cleanup.
func TestStorage(t testing.T) *Storage {
	td, err := ioutil.TempDir("", "waypoint-test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(td) })

	st, err := NewStorage(WithDir(td))
	require.NoError(t, err)

	return st
}
