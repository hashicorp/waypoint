package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/ptypes"
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
	if err := ptypes.ValidateGetConfigRequest(req); err != nil {
		return nil, err
	}

	vars, err := s.state.ConfigGet(req)
	if err != nil {
		return nil, err
	}

	return &pb.ConfigGetResponse{Variables: vars}, nil
}

func (s *service) SetConfigSource(
	ctx context.Context,
	req *pb.SetConfigSourceRequest,
) (*empty.Empty, error) {
	if err := ptypes.ValidateSetConfigSourceRequest(req); err != nil {
		return nil, err
	}

	if err := s.state.ConfigSourceSet(req.ConfigSource); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) GetConfigSource(
	ctx context.Context,
	req *pb.GetConfigSourceRequest,
) (*pb.GetConfigSourceResponse, error) {
	if err := ptypes.ValidateGetConfigSourceRequest(req); err != nil {
		return nil, err
	}

	vars, err := s.state.ConfigSourceGet(req)
	if err != nil {
		return nil, err
	}

	return &pb.GetConfigSourceResponse{ConfigSources: vars}, nil
}
