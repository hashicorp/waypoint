package singleprocess

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) SetServerConfig(
	ctx context.Context,
	req *pb.SetServerConfigRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateServerConfig(req.Config); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	if err := s.state(ctx).ServerConfigSet(ctx, req.Config); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to set server config", "platform", req.Config.Platform)
	}

	return &empty.Empty{}, nil
}

func (s *Service) GetServerConfig(
	ctx context.Context,
	req *empty.Empty,
) (*pb.GetServerConfigResponse, error) {
	cfg, err := s.state(ctx).ServerConfigGet(ctx)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get server config")
	}

	return &pb.GetServerConfigResponse{Config: cfg}, nil
}
