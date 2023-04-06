// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

// UI_ListPipelineRuns returns pipeline runs for a given pipeline. While paginating is
// part of the request, this doesn't yet support pagination and will return
// everything every time.
func (s *Service) UI_ListPipelineRuns(
	ctx context.Context,
	req *pb.UI_ListPipelineRunsRequest,
) (*pb.UI_ListPipelineRunsResponse, error) {
	log := hclog.FromContext(ctx)

	if err := serverptypes.ValidateUIListPipelineRunsRequest(req); err != nil {
		return nil, err
	}

	// Create uninitialized array of pipeline run bundles
	var allPipelineRuns []*pb.UI_PipelineRunBundle

	// Get list of all pipeline runs
	pipelineRunListResponse, err := s.state(ctx).PipelineRunList(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error listing piplines",
		)
	}

	// Create bundles
	var pipelineRunBundle *pb.UI_PipelineRunBundle
	for _, pipelineRun := range pipelineRunListResponse {
		if len(pipelineRun.Jobs) != 0 {
			job, err := s.GetJob(ctx, &pb.GetJobRequest{
				JobId: pipelineRun.Jobs[0].Id,
			})
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to get first job for latest pipeline run",
				)
			}
			pipelineRunBundle = &pb.UI_PipelineRunBundle{
				PipelineRun:   pipelineRun,
				QueueTime:     job.QueueTime,
				Application:   job.Application,
				DataSourceRef: job.DataSourceRef,
			}
		}

		// Add pipeline bundle to uninitialized array
		allPipelineRuns = append(allPipelineRuns, pipelineRunBundle)

	}

	// Return the array
	return &pb.UI_ListPipelineRunsResponse{
		PipelineRunBundles: allPipelineRuns,
		Pagination:         &pb.PaginationResponse{},
	}, nil
}
