// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/serverstate"
	"github.com/hashicorp/waypoint/pkg/serverstate/statetest"
)

func init() {
	// Seed our test randomness
	rand.Seed(time.Now().UnixNano())
}

func TestImpl(t *testing.T) {

	// Tests that are relevant, but are known to be failing.
	// It should be a priority to fix any test on this list.
	knownFailingStateTests := []string{
		"TestProjectPagination", // Failing b/c pagination not implemented in boltdb
		"TestJobListPagination", // Failing b/c pagination not implemented in boltdb
		"TestEvent",             //Failing b/c events aren't implemented in boltdb
	}

	// Tests for features that have not been implemented in OSS
	unimplementedTests := []string{
		"TestProjectTemplateFeatures",
	}

	statetest.Test(t, func(t *testing.T) serverstate.Interface {
		return TestState(t)
	}, func(t *testing.T, impl serverstate.Interface) serverstate.Interface {
		v, err := TestStateRestart(t, impl.(*State))
		require.NoError(t, err)
		return v
	}, append(knownFailingStateTests, unimplementedTests...))
}
