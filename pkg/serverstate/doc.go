// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package serverstate exports the interface and verification harness for
// implementing a new state storage backend for the server.
//
// Additional state storage backends is not an officially supported extension
// of Waypoint, so you do this at your own peril. It will not be supported.
// However, this package is exported because we maintain our own private state
// backends for other use cases of Waypoint.
package serverstate
