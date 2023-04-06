// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Noop executes a noop operation. This is primarily for testing but is
// exported since it has its uses in verifying a runner is functioning
// properly.
//
// A noop operation will exercise the full logic of queueing a job,
// assigning it to a runner, dequeueing as a runner, executing, etc. It will
// use real remote runners if the client is configured to do so.
func (c *App) Noop(ctx context.Context) error {
	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Noop_{
		Noop: &pb.Job_Noop{},
	}

	// Execute it
	_, err := c.doJob(ctx, job)
	return err
}
