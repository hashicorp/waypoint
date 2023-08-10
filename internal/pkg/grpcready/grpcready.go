// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package grpcready

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
)

// Conn returns nil if the connection is in the ready state. If wait is
// set, this will wait for the connection to become ready or until the context
// is cancelled.
//
// This is useful over `WithBlock` when connecting because it allows a
// fail-fast first attempt to connect. This makes it easily to synchronously
// attempt connection once before falling back to a retry loop in the background.
func Conn(
	ctx context.Context,
	log hclog.Logger,
	conn *grpc.ClientConn,
	wait bool,
) error {
	for {
		s := conn.GetState()
		log.Trace("connection state", "state", s.String())

		// If we're ready then we're done!
		if s == connectivity.Ready {
			log.Debug("connection is ready")
			return nil
		}

		// If we have a transient error and we're not retrying, then we're done.
		if s == connectivity.TransientFailure && !wait {
			log.Warn("failed to connect to the server, temporary network error")
			conn.Close()
			return status.Errorf(codes.Unavailable, "server is unavailable")
		}

		if !conn.WaitForStateChange(ctx, s) {
			return ctx.Err()
		}
	}
}
