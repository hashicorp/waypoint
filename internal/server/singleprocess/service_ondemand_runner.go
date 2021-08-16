package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// TODO: test
func (s *service) UpsertOndemandRunner(
	ctx context.Context,
	req *pb.UpsertOndemandRunnerRequest,
) (*pb.UpsertOndemandRunnerResponse, error) {
	if err := serverptypes.ValidateUpsertOndemandRunnerRequest(req); err != nil {
		return nil, err
	}

	result := req.OndemandRunner
	if err := s.state.OndemandRunnerPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertOndemandRunnerResponse{OndemandRunner: result}, nil
}

// TODO: test
func (s *service) GetOndemandRunner(
	ctx context.Context,
	req *pb.GetOndemandRunnerRequest,
) (*pb.GetOndemandRunnerResponse, error) {
	result, err := s.state.OndemandRunnerGet(req.OndemandRunner)
	if err != nil {
		return nil, err
	}

	return &pb.GetOndemandRunnerResponse{
		OndemandRunner: result,
	}, nil
}

// TODO: test
func (s *service) ListOndemandRunners(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOndemandRunnersResponse, error) {
	result, err := s.state.OndemandRunnerList()
	if err != nil {
		return nil, err
	}

	return &pb.ListOndemandRunnersResponse{OndemandRunners: result}, nil
}
