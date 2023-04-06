package boltdbstate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *State) CreateProjectTemplate(context.Context, *pb.ProjectTemplate) (*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}

func (s *State) UpdateProjectTemplate(context.Context, *pb.ProjectTemplate) (*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}

func (s *State) GetProjectTemplate(context.Context, *pb.Ref_ProjectTemplate) (*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
func (s *State) DeleteProjectTemplate(context.Context, *pb.Ref_ProjectTemplate) error {
	return status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
func (s *State) ListProjectTemplates(context.Context, *pb.ListProjectTemplatesRequest) ([]*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
