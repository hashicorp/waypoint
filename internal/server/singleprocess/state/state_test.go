package state

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
	statetest.Test(t, func(t *testing.T) serverstate.Interface {
		return TestState(t)
	}, func(t *testing.T, impl serverstate.Interface) serverstate.Interface {
		v, err := TestStateRestart(t, impl.(*State))
		require.NoError(t, err)
		return v
	})
}
