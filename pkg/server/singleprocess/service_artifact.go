package singleprocess

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertPushedArtifact(
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
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				fmt.Errorf("uuid generation failed: %s", err),
				"failed to generate a uuid while upserting a pushed artifact",
			)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state(ctx).ArtifactPut(!insert, result); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to upsert artifact")
	}

	return &pb.UpsertPushedArtifactResponse{Artifact: result}, nil
}

func (s *Service) ListPushedArtifacts(
	ctx context.Context,
	req *pb.ListPushedArtifactsRequest,
) (*pb.ListPushedArtifactsResponse, error) {
	if err := serverptypes.ValidateListPushedArtifactsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).ArtifactList(req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list artifacts",
		)
	}

	if req.IncludeBuild {
		for _, a := range result {
			b, err := s.state(ctx).BuildGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: a.BuildId,
				},
			})

			if err != nil {
				return nil, hcerr.Externalize(
					hclog.FromContext(ctx),
					err,
					"failed to get artifact",
					"build_id",
					a.BuildId,
				)
			}

			a.Build = b
		}
	}

	return &pb.ListPushedArtifactsResponse{Artifacts: result}, nil
}

// TODO: test
func (s *Service) GetLatestPushedArtifact(
	ctx context.Context,
	req *pb.GetLatestPushedArtifactRequest,
) (*pb.PushedArtifact, error) {
	if err := serverptypes.ValidateGetLatestPushedArtifactRequest(req); err != nil {
		return nil, err
	}

	a, err := s.state(ctx).ArtifactLatest(req.Application, req.Workspace)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get latest artifact",
			"application",
			req.Application.Application,
		)
	}
	return a, nil
}

// GetPushedArtifact returns a PushedArtifact based on ID
func (s *Service) GetPushedArtifact(
	ctx context.Context,
	req *pb.GetPushedArtifactRequest,
) (*pb.PushedArtifact, error) {
	if err := serverptypes.ValidateGetPushedArtifactRequest(req); err != nil {
		return nil, err
	}

	a, err := s.state(ctx).ArtifactGet(req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get pushed artifact",
			"reference",
			req.Ref,
		)
	}
	return a, nil
}
