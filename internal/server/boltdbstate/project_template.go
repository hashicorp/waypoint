package boltdbstate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *State) ProjectTemplatePut(context.Context, *pb.ProjectTemplate) error {
	return status.Error(codes.Unimplemented, "project template put unimplemented")
}

func (s *State) ProjectTemplateGet(context.Context, *pb.Ref_ProjectTemplate) (*pb.ProjectTemplate, error) {
	return nil, status.Error(codes.Unimplemented, "project template get unimplemented")
}

func (s *State) ProjectTemplateList(context.Context, *pb.PaginationRequest) ([]*pb.ProjectTemplate, *pb.PaginationResponse, error) {
	return nil, nil, status.Error(codes.Unimplemented, "project template list unimplemented")
}

func (s *State) ProjectTemplateDelete(ctx context.Context, template *pb.Ref_ProjectTemplate) error {
	return status.Error(codes.Unimplemented, "project template delete unimplemented")
}
