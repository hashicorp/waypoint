package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertPipeline(
	ctx context.Context,
	req *pb.UpsertPipelineRequest,
) (*pb.UpsertPipelineResponse, error) {
	if err := serverptypes.ValidateUpsertPipelineRequest(req); err != nil {
		return nil, err
	}

	result := req.Pipeline
	if err := s.state(ctx).PipelinePut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertPipelineResponse{Pipeline: result}, nil
}
