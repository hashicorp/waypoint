// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeConfigSyncOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	_, ok := job.Operation.(*pb.Job_ConfigSync)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	// Do the config sync
	if err := app.ConfigSync(ctx); err != nil {
		return nil, err
	}
	result := &pb.Job_Result{
		ConfigSync: &pb.Job_ConfigSyncResult{},
	}

	pipelineResult := make(map[string]*pb.Ref_Pipeline)
	for _, pipeline := range project.Pipelines() {
		if err := pipeline.ConfigSync(ctx); err != nil {
			return nil, err
		}

		pipelineResult[pipeline.Name()] = pipeline.Ref()
	}

	result.PipelineConfigSync = &pb.Job_PipelineConfigSyncResult{
		SyncedPipelines: pipelineResult,
	}

	return result, nil
}
