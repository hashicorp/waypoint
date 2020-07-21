package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TODO: test
func (s *service) UpsertProject(
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertProjectResponse, error) {
	result := req.Project
	if err := s.state.ProjectPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertProjectResponse{Project: result}, nil
}

// TODO: test
func (s *service) GetProject(
	ctx context.Context,
	req *pb.GetProjectRequest,
) (*pb.GetProjectResponse, error) {
	result, err := s.state.ProjectGet(req.Project)
	if err != nil {
		return nil, err
	}

	return &pb.GetProjectResponse{Project: result}, nil
}
