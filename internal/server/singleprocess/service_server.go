package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *service) SetServerConfig(
	ctx context.Context,
	req *pb.SetServerConfigRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateServerConfig(req.Config); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	if err := s.state.ServerConfigSet(req.Config); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) GetServerConfig(
	ctx context.Context,
	req *empty.Empty,
) (*pb.GetServerConfigResponse, error) {
	cfg, err := s.state.ServerConfigGet()
	if err != nil {
		return nil, err
	}

	return &pb.GetServerConfigResponse{Config: cfg}, nil
}
