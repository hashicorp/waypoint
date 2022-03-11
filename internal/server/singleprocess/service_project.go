package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/handlers"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *service) UpsertProject(
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertProjectResponse, error) {
	return handlers.UpsertProject(s, ctx, req)
}

func (s *service) GetProject(
	ctx context.Context,
	req *pb.GetProjectRequest,
) (*pb.GetProjectResponse, error) {
	return handlers.GetProject(s, ctx, req)
}

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

func (s *service) GetApplication(
	ctx context.Context,
	req *pb.GetApplicationRequest,
) (*pb.GetApplicationResponse, error) {
	if err := serverptypes.ValidateGetApplicationRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state.AppGet(req.Application)
	if err != nil {
		return nil, err
	}

	return &pb.GetApplicationResponse{
		Application: result,
	}, nil
}

func (s *service) UpsertApplication(
	ctx context.Context,
	req *pb.UpsertApplicationRequest,
) (*pb.UpsertApplicationResponse, error) {
	if err := serverptypes.ValidateUpsertApplicationRequest(req); err != nil {
		return nil, err
	}

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

func projectNeedsRemoteInit(project *pb.Project) bool {
	if project.DataSource == nil {
		return false
	}

	if project.DataSource.GetGit() == nil {
		return false
	}

	if project.DataSourcePoll != nil && project.DataSourcePoll.Enabled {
		return false
	}

	return true
}
