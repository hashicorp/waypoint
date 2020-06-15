package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) SetConfig(
	ctx context.Context,
	req *pb.ConfigSetRequest,
) (*pb.ConfigSetResponse, error) {
	if err := s.state.ConfigSet(req.Variables...); err != nil {
		return nil, err
	}

	return &pb.ConfigSetResponse{}, nil
}

func (s *service) GetConfig(
	ctx context.Context,
	req *pb.ConfigGetRequest,
) (*pb.ConfigGetResponse, error) {
	vars, err := s.state.ConfigGet(req)
	if err != nil {
		return nil, err
	}

	return &pb.ConfigGetResponse{Variables: vars}, nil
}
