package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertPushedArtifact(
	ctx context.Context,
	req *pb.UpsertPushedArtifactRequest,
) (*pb.UpsertPushedArtifactResponse, error) {
	result := req.Artifact

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

	if err := s.state.ArtifactPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertPushedArtifactResponse{Artifact: result}, nil
}

// TODO: test
func (s *service) ListPushedArtifacts(
	ctx context.Context,
	req *pb.ListPushedArtifactsRequest,
) (*pb.ListPushedArtifactsResponse, error) {
	result, err := s.state.ArtifactList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListPushedArtifactsResponse{Artifacts: result}, nil
}

// TODO: test
func (s *service) GetLatestPushedArtifact(
	ctx context.Context,
	req *pb.GetLatestPushedArtifactRequest,
) (*pb.PushedArtifact, error) {
	return s.state.ArtifactLatest(req.Application, req.Workspace)
}
