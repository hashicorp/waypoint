# Waypoint End-to-End Testing

## Requirements

For now, these tests assume that you already have Waypoint available on the path,
as well as Docker, K8s, and Nomad  installed and running. In the future, the
`run-tests.sh` script will set these up in a CI environment.

## How to run

The easiest way to run these tests is to use the shell runner with `./run-tests.sh`.

If you wish to run the script outside of this test directory, the following
environment variable must be set so that it knows where the `test-e2e` directory
is at:

* `WP_TESTE2E_DIR` (string) - Path to `test-e2e`.

For example, running these tests from the Waypoint root directory:

```
WP_TESTE2E_DIR="test-e2e" ./test-e2e/run-tests.sh
```

## How to write a new test

Create a new Go test file that has a relevant description for what the test will
be. For example, if you wanted to test and end-to-end scenario with Nomad, you
might call it `nomad_smoke_test.go`. The `run-tests.sh` script will pick up
and run this test file by defailt with `go test .`. But if you wish to run
it yourself without the shell runner, you can execute the test with
`go test <filename> util.go`. You need to include `util.go` to get the helper
functions for running the tests.

A simple test might be seeing that the Waypoint command line tool is available:

```go
package test

import (
	"strings"
	"testing"
)

func TestWaypointAvailable(t *testing.T) {
	wp := NewBinary(wpBinary, ".")
	stdout, stderr, err := wp.Run("version")
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
```

The `util.go` file inside this directory offers a simple way to execute any
binary.
