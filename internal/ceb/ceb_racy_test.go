// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// These tests fail the race detector, and should eventually be fixed.
//go:build !race

// Package ceb contains the core logic for the custom entrypoint binary ("ceb").
//
// The CEB does not work on Windows.
package ceb

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test how the CEB behaves when the server is down on startup.
func TestRun_serverDownNoConnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a listener that will refuse connections
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)
	ln.Close()

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	testRun(t, ctx, &testRunOpts{
		ClientDisable: true,
		DeploymentId:  "ABCD1234",
		Helper:        "write-file",
		HelperEnv: map[string]string{
			envServerAddr: ln.Addr().String(),
			"HELPER_PATH": path,
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 10*time.Second, 10*time.Millisecond)
}
