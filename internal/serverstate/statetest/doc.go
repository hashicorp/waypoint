// Package statetest has a test suite for validating implementations of the
// serverstate.Interface interface. This must only be imported in "_test.go"
// files in other packages.
//
// IMPORTANT: This package MUST NOT be imported outside of test code. This
// package imports "testing" as part of the exported API, and importing this
// into any non-test-compiled files will introduce the global test flags
// into the embedding process. This is due to an issue with Go where it
// registers test flags globally in an Init function on import.
package statetest
