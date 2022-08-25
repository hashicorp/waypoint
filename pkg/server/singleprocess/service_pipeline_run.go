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
