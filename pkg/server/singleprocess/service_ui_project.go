package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UI_GetProject(
	ctx context.Context,
	req *pb.UI_GetProjectRequest,
) (*pb.UI_GetProjectResponse, error) {
	if err := serverptypes.ValidateUIGetProjectRequest(req); err != nil {
		return nil, err
	}

	project, err := s.state(ctx).ProjectGet(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting project",
		)
	}

	latestInitJob, err := s.state(ctx).JobLatestInit(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting latest init job for project",
		)
	}

	return &pb.UI_GetProjectResponse{
		Project:       project,
		LatestInitJob: latestInitJob,
	}, nil
}
