// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeDestroyOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Destroy)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	switch t := op.Destroy.Target.(type) {
	case *pb.Job_DestroyOp_Deployment:
		err = app.DestroyDeploy(ctx, t.Deployment)

	case *pb.Job_DestroyOp_Workspace:
		err = app.Destroy(ctx)

	default:
		err = fmt.Errorf("unknown destruction target: %T", op.Destroy.Target)
	}
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{}, nil
}
