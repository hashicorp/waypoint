// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package testvault

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestVault(t *testing.T) {
	if _, err := exec.LookPath("vault"); err != nil {
		t.Skip("vault not on path")
		return
	}

	client, closer := TestVault(t)
	defer closer()

	// Verify Vault works
	_, err := client.Sys().ListMounts()
	require.NoError(t, err)
}
