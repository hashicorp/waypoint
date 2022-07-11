package handlertest

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

type (
	// Factory is the function type used to create a new serverhandler
	// implementation. To fail, this should fail the test.
	Factory func(*testing.T) (pb.WaypointClient, TestServerImpl)
)

// TestServerImpl is a wrapper around a server implementation that allows us
// to access otherwise-private fields and methods inside our tests. It should
// not be used outside of testing.
type TestServerImpl interface {

	// State returns the underlying serverstate implementation.
	State(ctx context.Context) serverstate.Interface
}

// Test runs a validation test suite for a pb.WaypointServer implementation.
// All server implementations should pass this suite with no errors to ensure
// the correct behavior of the server.
// skipTests are function names of tests in the serverstate package to skip
// (i.e. TestJobCreate_singleton). It cannot skip sub-tests (called by
// t.Run() inside a top-level test)
func Test(t *testing.T, f Factory, skipTests []string) {
	for name, funcs := range tests {
		t.Run(name, func(t *testing.T) {
			for _, tf := range funcs {
				name := runtime.FuncForPC(reflect.ValueOf(tf).Pointer()).Name()
				if idx := strings.LastIndexByte(name, '.'); idx >= 0 {
					name = name[idx+1:]
				}

				skip := false
				for _, skipTest := range skipTests {
					if name == skipTest {
						skip = true
						break
					}
				}
				if skip {
					t.Run(name, func(t *testing.T) {
						t.Skipf("Test %q is on the handler skip list - ignoring", name)
					})
					continue
				}

				t.Run(name, func(t *testing.T) {
					tf(t, f)
				})
			}
		})
	}
}

// TestGroup runs a specific group of validation tests for a state implementation.
func TestGroup(t *testing.T, name string, f Factory) {
	funcs, ok := tests[name]
	if !ok {
		panic(fmt.Sprintf("unknown test group: %s", name))
	}

	t.Run(name, func(t *testing.T) {
		for _, tf := range funcs {
			tf(t, f)
		}
	})
}

// tests is the list of tests to run.
var tests = map[string][]testFunc{}

// testFunc is the type of the function that a test that is run as part of
// Test implements. This is an internal only type.
type testFunc func(*testing.T, Factory)
