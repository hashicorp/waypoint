package statetest

import (
	"crypto/rand"
	"testing"

	ulidpkg "github.com/oklog/ulid"

	"github.com/hashicorp/waypoint/internal/serverstate"
)

// Factory is the function type used to create a new serverstate
// implementation. To fail, this should fail the test.
type Factory func(t *testing.T) serverstate.Interface

// Test runs a validation test suite for a state implementation. All
// state implementations should pass this suite with no errors to ensure
// the correct behavior of the state when Waypoint uses it.
func Test(t *testing.T, f Factory) {
	for name, funcs := range tests {
		t.Run(name, func(t *testing.T) {
			for _, tf := range funcs {
				tf(t, f)
			}
		})
	}
}

// tests is the list of tests to run.
var tests = map[string][]testFunc{}

// testFunc is the type of the function that a test that is run as part of
// Test implements. This is an internal only type.
type testFunc func(*testing.T, Factory)

// ulid returns a unique ULID.
func ulid() (string, error) {
	id, err := ulidpkg.New(ulidpkg.Now(), rand.Reader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
