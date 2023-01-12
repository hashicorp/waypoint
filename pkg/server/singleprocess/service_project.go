package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"

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

	if projectNeedsRemoteInit(result) {
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
