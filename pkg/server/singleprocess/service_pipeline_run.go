package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) ListPipelineRuns(
	ctx context.Context,
	req *pb.ListPipelineRunsRequest,
) (*pb.ListPipelineRunsResponse, error) {
	if err := serverptypes.ValidateListPipelineRunsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).PipelineRunList(req.Pipeline)
	if err != nil {
		return nil, err
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

	result, err := s.state(ctx).PipelineRunGet(req.Pipeline, req.Sequence)
	if err != nil {
		return nil, err
	}

	return &pb.GetPipelineRunResponse{
		PipelineRun: result,
	}, nil
}
