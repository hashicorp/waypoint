package client

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (c *Client) Build(ctx context.Context) error {
	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Build{
		Build: &pb.Job_BuildOp{},
	}

	// Execute it
	return c.doJob(ctx, job)
}
