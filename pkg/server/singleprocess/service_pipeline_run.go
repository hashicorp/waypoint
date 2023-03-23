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

func (s *Service) ListPipelineRuns(
	ctx context.Context,
	req *pb.ListPipelineRunsRequest,
) (*pb.ListPipelineRunsResponse, error) {
	if err := serverptypes.ValidateListPipelineRunsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).PipelineRunList(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to push pipeline run",
		)
	}

	return &pb.ListPipelineRunsResponse{
		PipelineRuns: result,
	}, nil
}

func (s *Service) GetPipelineRun(
	ctx context.Context,
	req *pb.GetPipelineRunRequest,
) (*pb.GetPipelineRunResponse, error) {
	if err := serverptypes.ValidateGetPipelineRunRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).PipelineRunGet(ctx, req.Pipeline, req.Sequence)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get pipeline run",
		)
	}

	return &pb.GetPipelineRunResponse{
		PipelineRun: result,
	}, nil
}

func (s *Service) GetLatestPipelineRun(
	ctx context.Context,
	req *pb.GetPipelineRequest,
) (*pb.GetPipelineRunResponse, error) {
	if err := serverptypes.ValidateGetPipelineRequest(req); err != nil {
		return nil, err
	}

	pipeline, err := s.state(ctx).PipelineGet(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting pipeline",
		)
	}

	latestPipelineRun, err := s.state(ctx).PipelineRunGetLatest(ctx, pipeline.Id)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get latest pipeline run",
		)
	}

	return &pb.GetPipelineRunResponse{
		PipelineRun: latestPipelineRun,
	}, nil
}
