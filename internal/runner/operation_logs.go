// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeLogsOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	r.logger.Debug("Logs operation starting")
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Logs)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	err = app.Logs(ctx,
		op.Logs.InstanceId,
		op.Logs.Deployment,
		op.Logs.StartTime.AsTime(),
		int(op.Logs.Limit),
	)

	return nil, err
}
