package state

import (
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/waypoint/internal/serverstate"
	"github.com/hashicorp/waypoint/internal/serverstate/statetest"
)

func init() {
	// Seed our test randomness
	rand.Seed(time.Now().UnixNano())
}

func TestImpl(t *testing.T) {
	statetest.Test(t, func(t *testing.T) serverstate.Interface {
		return TestState(t)
	})
}
