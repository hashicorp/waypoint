// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package serverstate exports the verification harness for
// implementing a new waypoint server protobuf (i.e. pb.WaypointServer)
// implementations.
//
// Additional server implementations is not an officially supported extension
// of Waypoint, so you do this at your own peril. It will not be supported.
// However, this package is exported because we maintain our own private server
// implementations for other use cases of Waypoint.
package serverhandler
