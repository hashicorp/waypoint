package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertBuild(
	ctx context.Context,
	req *pb.UpsertBuildRequest,
) (*pb.UpsertBuildResponse, error) {
	result := req.Build

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state.BuildPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertBuildResponse{Build: result}, nil
}

func (s *service) ListBuilds(
	ctx context.Context,
	req *pb.ListBuildsRequest,
) (*pb.ListBuildsResponse, error) {
	result, err := s.state.BuildList(req.Application,
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListBuildsResponse{Builds: result}, nil
}

func (s *service) GetLatestBuild(
	ctx context.Context,
	req *pb.GetLatestBuildRequest,
) (*pb.Build, error) {
	return s.state.BuildLatest(req.Application, req.Workspace)
}
