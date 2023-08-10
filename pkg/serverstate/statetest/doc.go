// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package statetest has a test suite for validating implementations of the
// serverstate.Interface interface. This must only be imported in "_test.go"
// files in other packages.
//
// IMPORTANT: This package MUST NOT be imported outside of test code. This
// package imports "testing" as part of the exported API, and importing this
// into any non-test-compiled files will introduce the global test flags
// into the embedding process. This is due to an issue with Go where it
// registers test flags globally in an Init function on import.
//
// These tests are not invoked directly, but are used via pkg/serverstate
// implementations, e.g. internal/server/boltdbstate/state.go
//
// To run these tests, run the tests for that package.
//
// Ex:
// $ go test -test.v ./internal/server/boltdbstate -count=1
//
// To run a specific test, use the -run flag with TestImpl. For example, to run
// the TestOnDemandRunnerConfig test defined in
// pkg/serverstate/statetest/test_runner_ondemand.go, use:
//
// $ go test -test.v ./internal/server/boltdbstate -count=1 -run=TestImpl/runner_ondemand/TestOnDemandRunnerConfig

package statetest
