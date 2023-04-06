// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package grpcmetadata contains functions for reading and writing waypoint specific
// metadata to contexts, which is transmitted by RPC calls.
package grpcmetadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// The metadata key that stores the runner id associated with a client. This is
// used by the CLI to advertise it's local client for when the server needs to
// spawn jobs back on that client in response to an RPC.
const grpcMetadataRunnerId = "waypoint-runner-id"

// AddRunner adds gRPC metadata to an outgoing context to indicate that RPCs sent
// with the returned context having the given runner (specified by id) attached to
// the sending client, allow the server to target jobs back to the client that
// performed an RPC call.
func AddRunner(ctx context.Context, id string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, grpcMetadataRunnerId, id)
}

// RunnerId returns the runner id attached to the incoming context as grpc Metadata.
// This would be set by the client to indicate there is a runner attached
// directly to it.
func RunnerId(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	val := md.Get(grpcMetadataRunnerId)
	if len(val) == 0 {
		return "", false
	}

	return val[0], true
}

// OutgoingRunnerId returns the runner id attached to the context as grpc Metadata.
// This is primarily used in tests, to validate that a context was set correctly.
func OutgoingRunnerId(ctx context.Context) (string, bool) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", false
	}

	val := md.Get(grpcMetadataRunnerId)
	if len(val) == 0 {
		return "", false
	}

	return val[0], true
}
