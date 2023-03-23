// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"strings"
	"testing"
)

func TestWaypointInstall(t *testing.T) {
	t.Logf("Testing waypoint is available...")
	wp := NewBinary(t, wpBinary, ".")
	stdout, stderr, err := wp.RunRaw("version")
	if err != nil {
		t.Errorf("unexpected error getting version: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output getting version: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint v") {
		t.Errorf("No version output detected:\n%s", stdout)
	}
}
