package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// TODO: test
func (s *service) UpsertOndemandRunnerConfig(
	ctx context.Context,
	req *pb.UpsertOndemandRunnerConfigRequest,
) (*pb.UpsertOndemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateUpsertOndemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	result := req.Config
	if err := s.state.OndemandRunnerConfigPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertOndemandRunnerConfigResponse{Config: result}, nil
}

// TODO: test
func (s *service) GetOndemandRunnerConfig(
	ctx context.Context,
	req *pb.GetOndemandRunnerConfigRequest,
) (*pb.GetOndemandRunnerConfigResponse, error) {
	result, err := s.state.OndemandRunnerConfigGet(req.Config)
	if err != nil {
		return nil, err
	}

	return &pb.GetOndemandRunnerConfigResponse{
		Config: result,
	}, nil
}

// TODO: test
func (s *service) ListOndemandRunnerConfigs(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOndemandRunnerConfigsResponse, error) {
	result, err := s.state.OndemandRunnerConfigList()
	if err != nil {
		return nil, err
	}

	return &pb.ListOndemandRunnerConfigsResponse{Configs: result}, nil
}
