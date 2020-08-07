package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertRelease(
	ctx context.Context,
	req *pb.UpsertReleaseRequest,
) (*pb.UpsertReleaseResponse, error) {
	result := req.Release

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

	if err := s.state.ReleasePut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertReleaseResponse{Release: result}, nil
}

// TODO: test
func (s *service) ListReleases(
	ctx context.Context,
	req *pb.ListReleasesRequest,
) (*pb.ListReleasesResponse, error) {
	result, err := s.state.ReleaseList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
		state.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListReleasesResponse{Releases: result}, nil
}

// TODO: test
func (s *service) GetLatestRelease(
	ctx context.Context,
	req *pb.GetLatestReleaseRequest,
) (*pb.Release, error) {
	return s.state.ReleaseLatest(req.Application, req.Workspace)
}

// GetRelease returns a Release based on ID
func (s *service) GetRelease(
	ctx context.Context,
	req *pb.GetReleaseRequest,
) (*pb.Release, error) {
	return s.state.ReleaseGet(req.Ref)
}
