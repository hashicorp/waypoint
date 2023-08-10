// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertProject(
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertProjectResponse, error) {
	if err := serverptypes.ValidateUpsertProjectRequest(req); err != nil {
		return nil, err
	}

	result := req.Project
	if err := s.state(ctx).ProjectPut(ctx, result); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to upsert project",
			"project_name",
			result.GetName(),
		)
	}

	if hasUsableWaypointHCL(result) {
		proj, err := s.serverSideProjectInit(ctx, result)
		if err != nil {
			return nil, err // already externalized
		}
		result = proj
	} else if projectNeedsRemoteInit(result) {
		// The project is connected to a data source but doesn’t use
		// automatic polling, so let’s queue some remote init operations
		// to ensure the application list is populated.

		// TODO(jgwhite): only queue init ops if the relevant fields have *changed*

		err := queueInitOps(s, ctx, result)

		if err != nil {
			// An error here indicates a failure to enqueue an
			// InitOp, not a failure during the operation itself,
			// which happen out-of-band.
			return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed queueing init job")
		}
	}

	return &pb.UpsertProjectResponse{Project: result}, nil
}

func (s *Service) GetProject(
	ctx context.Context,
	req *pb.GetProjectRequest,
) (*pb.GetProjectResponse, error) {
	if err := serverptypes.ValidateGetProjectRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).ProjectGet(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get project",
			"project_name",
			result.GetName(),
		)
	}

	// Get all the workspaces that this project is part of
	workspaces, err := s.state(ctx).ProjectListWorkspaces(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list workspaces for project",
			"project_name",
			result.GetName(),
		)
	}

	return &pb.GetProjectResponse{
		Project:    result,
		Workspaces: workspaces,
	}, nil
}

func (s *Service) ListProjects(
	ctx context.Context,
	req *pb.ListProjectsRequest,
) (*pb.ListProjectsResponse, error) {
	if err := serverptypes.ValidateListProjectsRequest(req); err != nil {
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

	result, pagination, err := s.state(ctx).ProjectList(ctx, req.Pagination)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list projects",
		)
	}

	return &pb.ListProjectsResponse{Projects: result, Pagination: pagination, TotalCount: count}, nil
}

func (s *Service) DestroyProject(
	ctx context.Context,
	req *pb.DestroyProjectRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDestroyProjectRequest(req); err != nil {
		return nil, err
	}

	err := s.state(ctx).ProjectDelete(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to delete project",
			"project",
			req.Project.Project,
		)
	}

	return &empty.Empty{}, nil
}

func (s *Service) GetApplication(
	ctx context.Context,
	req *pb.GetApplicationRequest,
) (*pb.GetApplicationResponse, error) {
	if err := serverptypes.ValidateGetApplicationRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).AppGet(ctx, req.Application)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get application",
			"application",
			req.Application.Application,
		)
	}

	return &pb.GetApplicationResponse{
		Application: result,
	}, nil
}

func (s *Service) UpsertApplication(
	ctx context.Context,
	req *pb.UpsertApplicationRequest,
) (*pb.UpsertApplicationResponse, error) {
	if err := serverptypes.ValidateUpsertApplicationRequest(req); err != nil {
		return nil, err
	}

	// Get the project
	praw, err := s.state(ctx).ProjectGet(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get project in application upsert",
		)
	}

	var app *pb.Application

	// If the project has the application already then we're done.
	p := serverptypes.Project{Project: praw}
	if idx := p.App(req.Name); idx >= 0 {
		app = p.Applications[idx]
	} else {
		app = &pb.Application{
			Project: req.Project,
			Name:    req.Name,
		}
	}

	app.FileChangeSignal = req.FileChangeSignal

	app, err = s.state(ctx).AppPut(ctx, app)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get upsert application",
			"application_name",
			app.GetName(),
		)
	}

	return &pb.UpsertApplicationResponse{Application: app}, nil
}

func queueInitOps(s *Service, ctx context.Context, project *pb.Project) error {
	workspaces, err := s.state(ctx).WorkspaceList(ctx)
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		workspaces = append(workspaces, &pb.Workspace{Name: "default"})
	}

	for _, workspace := range workspaces {
		_, err := s.QueueJob(ctx, &pb.QueueJobRequest{
			Job: &pb.Job{
				Application: &pb.Ref_Application{
					Project: project.Name,
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: workspace.Name,
				},
				Operation: &pb.Job_Init{},
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{},
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func projectNeedsRemoteInit(project *pb.Project) bool {
	if project.State != pb.Project_ACTIVE {
		return false
	}

	if project.DataSource == nil {
		return false
	}

	if project.DataSource.GetGit() == nil && project.DataSource.GetRemote() == nil {
		return false
	}

	if project.DataSourcePoll != nil && project.DataSourcePoll.Enabled {
		return false
	}

	return true
}

// hasUsableWaypointHCL verifies that a project has waypoint.hcl contents
// of type HCL
func hasUsableWaypointHCL(project *pb.Project) bool {
	return len(project.WaypointHcl) > 0 && project.WaypointHclFormat == pb.Hcl_HCL
}

// serverSideProjectInit initializes a project that directly contains a waypoint.hcl.
// "init" currently consists of getting the list of apps on a project, and upserting each one.
// If the waypoint.hcl is not on the project but is in VCS, you must enqueue an init job
// rather than attempting a serverside init.
// Returns externalized errors
func (s *Service) serverSideProjectInit(ctx context.Context, project *pb.Project) (*pb.Project, error) {
	if !hasUsableWaypointHCL(project) {
		return nil, fmt.Errorf("cannot init a project without a stored waypoint.hcl with hcl type contents serverside")
	}
	file, _ := hclsyntax.ParseConfig(project.WaypointHcl, "<waypoint-hcl>", hcl.Pos{})
	content, _ := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "app", LabelNames: []string{"name"}},
		},
	})
	projRef := &pb.Ref_Project{Project: project.Name}

	for _, b := range content.Blocks.ByType()["app"] {
		name := b.Labels[0]
		_, err := s.state(ctx).AppPut(ctx, &pb.Application{
			Project: projRef,
			Name:    name,
		})
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"failed to register app %q while creating project %q",
				name, project.GetName(),
			)
		}
	}

	// Reload the project to populate the newly-added apps
	result, err := s.state(ctx).ProjectGet(ctx, projRef)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to reload project %q",
			project.GetName(),
		)
	}

	return result, nil
}
