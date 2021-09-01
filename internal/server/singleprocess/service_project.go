package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// TODO: test
func (s *service) UpsertProject(
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertProjectResponse, error) {
	if err := serverptypes.ValidateUpsertProjectRequest(req); err != nil {
		return nil, err
	}

	result := req.Project
	if err := s.state.ProjectPut(result); err != nil {
		return nil, err
	}

	if result.DataSource != nil && result.DataSource.GetGit() != nil {
		// The project is connected to a data source so let’s
		// try to queue some remote init operations to ensure the
		// application list is populated and ready to work with.

		// TODO: error handling
		queueInitOps(s, ctx, result)
	}

	return &pb.UpsertProjectResponse{Project: result}, nil
}

// TODO: test
func (s *service) GetProject(
	ctx context.Context,
	req *pb.GetProjectRequest,
) (*pb.GetProjectResponse, error) {
	result, err := s.state.ProjectGet(req.Project)
	if err != nil {
		return nil, err
	}

	// Get all the workspaces that this project is part of
	workspaces, err := s.state.ProjectListWorkspaces(req.Project)
	if err != nil {
		return nil, err
	}

	return &pb.GetProjectResponse{
		Project:    result,
		Workspaces: workspaces,
	}, nil
}

// TODO: test
func (s *service) ListProjects(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListProjectsResponse, error) {
	result, err := s.state.ProjectList()
	if err != nil {
		return nil, err
	}

	return &pb.ListProjectsResponse{Projects: result}, nil
}

// TODO: test
func (s *service) UpsertApplication(
	ctx context.Context,
	req *pb.UpsertApplicationRequest,
) (*pb.UpsertApplicationResponse, error) {
	// Get the project
	praw, err := s.state.ProjectGet(req.Project)
	if err != nil {
		return nil, err
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

	app, err = s.state.AppPut(app)
	if err != nil {
		return nil, err
	}

	return &pb.UpsertApplicationResponse{Application: app}, nil
}

func queueInitOps(s *service, ctx context.Context, project *pb.Project) error {
	workspaces, err := s.state.WorkspaceList()

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
