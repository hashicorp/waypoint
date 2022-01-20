package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/serverstate"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *service) UpsertPushedArtifact(
	ctx context.Context,
	req *pb.UpsertPushedArtifactRequest,
) (*pb.UpsertPushedArtifactResponse, error) {
	if err := serverptypes.ValidateUpsertPushedArtifactRequest(req); err != nil {
		return nil, err
	}

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

func (s *service) ListPushedArtifacts(
	ctx context.Context,
	req *pb.ListPushedArtifactsRequest,
) (*pb.ListPushedArtifactsResponse, error) {
	if err := serverptypes.ValidateListPushedArtifactsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state.ArtifactList(req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	if req.IncludeBuild {
		for _, a := range result {
			b, err := s.state.BuildGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: a.BuildId,
				},
			})

			if err != nil {
				return nil, err
			}

			a.Build = b
		}
	}

	return &pb.ListPushedArtifactsResponse{Artifacts: result}, nil
}

// TODO: test
func (s *service) GetLatestPushedArtifact(
	ctx context.Context,
	req *pb.GetLatestPushedArtifactRequest,
) (*pb.PushedArtifact, error) {
	if err := serverptypes.ValidateGetLatestPushedArtifactRequest(req); err != nil {
		return nil, err
	}

	return s.state.ArtifactLatest(req.Application, req.Workspace)
}

// GetPushedArtifact returns a PushedArtifact based on ID
func (s *service) GetPushedArtifact(
	ctx context.Context,
	req *pb.GetPushedArtifactRequest,
) (*pb.PushedArtifact, error) {
	if err := serverptypes.ValidateGetPushedArtifactRequest(req); err != nil {
		return nil, err
	}

	return s.state.ArtifactGet(req.Ref)
}
