package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/ptypes"
)

// TODO: test
func (s *service) GetWorkspace(
	ctx context.Context,
	req *pb.GetWorkspaceRequest,
) (*pb.GetWorkspaceResponse, error) {
	if err := ptypes.ValidateGetWorkspaceRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state.WorkspaceGet(req.Workspace.Workspace)
	if err != nil {
		return nil, err
	}

	return &pb.GetWorkspaceResponse{Workspace: result}, nil
}

// TODO: test
func (s *service) ListWorkspaces(
	ctx context.Context,
	req *pb.ListWorkspacesRequest,
) (*pb.ListWorkspacesResponse, error) {
	var err error
	var result []*pb.Workspace

	switch v := req.Scope.(type) {
	case nil:
		// This is the same as Global for backwards compat reasons.
		result, err = s.state.WorkspaceList()

	case *pb.ListWorkspacesRequest_Global:
		result, err = s.state.WorkspaceList()

	case *pb.ListWorkspacesRequest_Project:
		result, err = s.state.WorkspaceListByProject(v.Project)

	case *pb.ListWorkspacesRequest_Application:
		result, err = s.state.WorkspaceListByApp(v.Application)

	default:
		return nil, status.Errorf(codes.FailedPrecondition,
			"unknown ListWorkspaces scope type: %T", req.Scope)
	}
	if err != nil {
		return nil, err
	}

	return &pb.ListWorkspacesResponse{Workspaces: result}, nil
}
