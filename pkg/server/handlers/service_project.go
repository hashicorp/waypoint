package handlers

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func UpsertProject(
	s Service,
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertProjectResponse, error) {
	if err := serverptypes.ValidateUpsertProjectRequest(req); err != nil {
		return nil, err
	}

	result := req.Project
	if err := s.State(ctx).ProjectPut(result); err != nil {
		return nil, err
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
			return nil, err
		}
	}

	return &pb.UpsertProjectResponse{Project: result}, nil
}

func GetProject(
	s Service,
	ctx context.Context,
	req *pb.GetProjectRequest,
) (*pb.GetProjectResponse, error) {
	if err := serverptypes.ValidateGetProjectRequest(req); err != nil {
		return nil, err
	}

	result, err := s.State(ctx).ProjectGet(req.Project)
	if err != nil {
		return nil, err
	}

	// Get all the workspaces that this project is part of
	workspaces, err := s.State(ctx).ProjectListWorkspaces(req.Project)
	if err != nil {
		return nil, err
	}

	return &pb.GetProjectResponse{
		Project:    result,
		Workspaces: workspaces,
	}, nil
}

func queueInitOps(s Service, ctx context.Context, project *pb.Project) error {
	workspaces, err := s.State(ctx).WorkspaceList()
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
