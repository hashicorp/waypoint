package boltdbstate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *State) ProjectTemplatePut(context.Context, *pb.ProjectTemplate) error {
	return status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}

func (s *State) ProjectTemplateGet(context.Context, *pb.Ref_ProjectTemplate) (*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
func (s *State) ProjectTemplateDelete(context.Context, *pb.Ref_ProjectTemplate) error {
	return status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
func (s *State) ProjectTemplateList(context.Context, *pb.ListProjectTemplatesRequest) ([]*pb.ProjectTemplate, error) {
	return nil, status.Errorf(codes.Unimplemented, "ProjectTemplate Unimplemented")
}
