// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeValidateOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	// The core loading currently handles all the validation for us since
	// we load all plugins and configurations.

	return &pb.Job_Result{
		Validate: &pb.Job_ValidateResult{},
	}, nil
}
