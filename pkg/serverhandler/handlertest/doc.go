// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package handlertest has a test suite for validating implementations of the
// pb.WaypointService interface. This must only be imported in "_test.go"
// files in other packages.
//
// IMPORTANT: This package MUST NOT be imported outside of test code. This
// package imports "testing" as part of the exported API, and importing this
// into any non-test-compiled files will introduce the global test flags
// into the embedding process. This is due to an issue with Go where it
// registers test flags globally in an Init function on import.
//
// These tests are not invoked directly, but are used via pb.WaypointServer
// implementations, e.g. pkg/server/singleprocess
//
// To run these tests, run the tests for that package.
//
// Ex:
// $ go test -test.v ./pkg/server/singleprocess -count=1
//
// To run a specific test, use the -run flag with TestHandlers. For example, to run
// the TestWorkspace_Upsert test defined in
// pkg/serverhandler/handlertest/test_service_workspace.go, use:
//
// $ go test -test.v ./pkg/server/singleprocess -count=1 -run=TestHandlers/workspace/TestWorkspace_Upsert

package handlertest
