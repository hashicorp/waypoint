package singleprocess

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertBuild(
	ctx context.Context,
	req *pb.UpsertBuildRequest,
) (*pb.UpsertBuildResponse, error) {
	if err := serverptypes.ValidateUpsertBuildRequest(req); err != nil {
		return nil, err
	}

	result := req.Build

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				fmt.Errorf("uuid generation failed: %w", err),
				"failed to generate a uuid while upserting a build",
			)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state(ctx).BuildPut(ctx, !insert, result); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to insert build for app", "app", req.Build.Application, "id", req.Build.Id)
	}

	return &pb.UpsertBuildResponse{Build: result}, nil
}

func (s *Service) ListBuilds(
	ctx context.Context,
	req *pb.ListBuildsRequest,
) (*pb.ListBuildsResponse, error) {
	if err := serverptypes.ValidateListBuildsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).BuildList(ctx, req.Application,
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithOrder(req.Order),
	)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to list builds for app", "app", req.Application.Application, "project", req.Application.Project)
	}

	return &pb.ListBuildsResponse{Builds: result}, nil
}

func (s *Service) GetLatestBuild(
	ctx context.Context,
	req *pb.GetLatestBuildRequest,
) (*pb.Build, error) {
	if err := serverptypes.ValidateGetLatestBuildRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).BuildLatest(ctx, req.Application, req.Workspace)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get latest build", "app", req.Application.Application, "project", req.Application.Project)
	}

	return result, nil
}

// GetBuild returns a Build based on ID
func (s *Service) GetBuild(
	ctx context.Context,
	req *pb.GetBuildRequest,
) (*pb.Build, error) {
	if err := serverptypes.ValidateGetBuildRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).BuildGet(ctx, req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get build", "id", req.Ref.Target)
	}

	return result, nil
}
