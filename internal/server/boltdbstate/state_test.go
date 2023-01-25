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
	}

	statetest.Test(t, func(t *testing.T) serverstate.Interface {
		return TestState(t)
	}, func(t *testing.T, impl serverstate.Interface) serverstate.Interface {
		v, err := TestStateRestart(t, impl.(*State))
		require.NoError(t, err)
		return v
	}, knownFailingStateTests)
}
