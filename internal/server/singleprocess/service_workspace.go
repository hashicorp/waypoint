package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TODO: test
func (s *service) GetWorkspace(
	ctx context.Context,
	req *pb.GetWorkspaceRequest,
) (*pb.GetWorkspaceResponse, error) {
	result, err := s.state.WorkspaceGet(req.Workspace.Workspace)
	if err != nil {
		return nil, err
	}

	return &pb.GetWorkspaceResponse{Workspace: result}, nil
}

// TODO: test
func (s *service) ListWorkspaces(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListWorkspacesResponse, error) {
	result, err := s.state.WorkspaceList()
	if err != nil {
		return nil, err
	}

	return &pb.ListWorkspacesResponse{Workspaces: result}, nil
}
