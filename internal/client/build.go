package client

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (c *App) Build(ctx context.Context, op *pb.Job_BuildOp) error {
	if op == nil {
		op = &pb.Job_BuildOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Build{
		Build: op,
	}

	// Execute it
	return c.doJob(ctx, job)
}
