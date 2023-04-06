// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package testvault contains helpers for working with Vault in a test
// environment.
package testvault

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	testing "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/builtin/vault/freeport"
)

// TestVault starts a new dev-mode Vault cluster and returns a client that
// is configured to talk to the Vault cluster with root credentials.
//
// The returned func() should be called to shut down the Vault server and
// perform any additional cleanup.
func TestVault(t testing.T) (*api.Client, func()) {
	if _, err := exec.LookPath("vault"); err != nil {
		t.Skip("vault not on path")
		return nil, nil
	}

	tv := newTestVault(t)
	return tv.client, tv.Stop
}

// testVault is a test helper. It uses a fork/exec model to create a test Vault
// server instance in the background and can be initialized with policies, roles
// and backends mounted. The test Vault instances can be used to run a unit test
// and offers and easy API to tear itself down on test end. The only
// prerequisite is that the Vault binary is on the $PATH.
type testVault struct {
	cmd    *exec.Cmd
	t      testing.T
	waitCh chan error
	client *api.Client
}

// NewTestVault returns a new testVault instance that has been started
func newTestVault(t testing.T) *testVault {
	for i := 10; i >= 0; i-- {
		port := freeport.MustTake(1)[0]
		bind := fmt.Sprintf("-dev-listen-address=127.0.0.1:%d", port)
		http := fmt.Sprintf("http://127.0.0.1:%d", port)
		cmd := exec.Command("vault", "server", "-dev", bind, "-dev-root-token-id=test")

		// Build the config
		conf := api.DefaultConfig()
		conf.Address = http

		// Make the client and set the token to the root token
		client, err := api.NewClient(conf)
		if err != nil {
			t.Fatalf("failed to build Vault API client: %v", err)
		}
		client.SetToken("test")

		tv := &testVault{
			cmd:    cmd,
			t:      t,
			client: client,
		}

		if err := tv.cmd.Start(); err != nil {
			tv.t.Fatalf("failed to start vault: %v", err)
		}

		// Start the waiter
		tv.waitCh = make(chan error, 1)
		go func() {
			err := tv.cmd.Wait()
			tv.waitCh <- err
		}()

		// Ensure Vault started
		var startErr error
		select {
		case startErr = <-tv.waitCh:
		case <-time.After(2 * time.Second):
		}

		if startErr != nil && i == 0 {
			t.Fatalf("failed to start vault: %v", startErr)
		} else if startErr != nil {
			wait := time.Duration(rand.Int31n(2000)) * time.Millisecond
			time.Sleep(wait)
			continue
		}

		waitErr := tv.waitForAPI()
		if waitErr != nil && i == 0 {
			t.Fatalf("failed to start vault: %v", waitErr)
		} else if waitErr != nil {
			wait := time.Duration(rand.Int31n(2000)) * time.Millisecond
			time.Sleep(wait)
			continue
		}

		return tv
	}

	return nil

}

// Stop stops the test Vault server
func (tv *testVault) Stop() {
	if tv.cmd.Process == nil {
		return
	}

	if err := tv.cmd.Process.Kill(); err != nil {
		if strings.Contains(err.Error(), "process already finished") {
			// Process already killed
			return
		}

		tv.t.Errorf("err: %s", err)
	}
	if tv.waitCh != nil {
		select {
		case <-tv.waitCh:
			return
		case <-time.After(1 * time.Second):
			require.Fail(tv.t, "Timed out waiting for vault to terminate")
		}
	}
}

// waitForAPI waits for the Vault HTTP endpoint to start
// responding. This is an indication that the agent has started.
func (tv *testVault) waitForAPI() error {
	test := func() (bool, error) {
		inited, err := tv.client.Sys().InitStatus()
		if err != nil {
			return false, err
		}
		return inited, nil
	}

	var success bool
	var err error
	retries := 500
	for retries > 0 {
		success, err = test()
		if success {
			return nil
		}

		time.Sleep(10 * time.Millisecond)
		retries--
	}

	return err
}
