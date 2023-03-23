// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UI_ListProjects(
	ctx context.Context,
	req *pb.UI_ListProjectsRequest,
) (*pb.UI_ListProjectsResponse, error) {
	if err := serverptypes.ValidateUIListProjectsRequest(req); err != nil {
		return nil, err
	}

	count, err := s.state(ctx).ProjectCount(ctx)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to count projects",
		)
	}

	projectBundles, pagination, err := s.state(ctx).ProjectListBundles(ctx, req.Pagination)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list projects",
		)
	}

	return &pb.UI_ListProjectsResponse{
		ProjectBundles: projectBundles,
		Pagination:     pagination,
		TotalCount:     count,
	}, nil
}

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
