package state

import (
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestState returns an initialized State for testing.
func TestState(t testing.T) *State {
	result, err := New()
	require.NoError(t, err)
	return result
}
