package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func (s *service) UpsertOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.UpsertOnDemandRunnerConfigRequest,
) (*pb.UpsertOnDemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateUpsertOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	if req.Config.TargetRunner == nil {
		req.Config.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{},
		}
	} else {
		// Validate static runner exists on odr profile creation.
		if t, ok := req.Config.TargetRunner.Target.(*pb.Ref_Runner_Id); ok {
			_, err := s.GetRunner(ctx, &pb.GetRunnerRequest{
				RunnerId: t.Id.Id,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	result := req.Config
	if err := s.state.OnDemandRunnerConfigPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertOnDemandRunnerConfigResponse{Config: result}, nil
}

func (s *service) GetOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.GetOnDemandRunnerConfigRequest,
) (*pb.GetOnDemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateGetOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state.OnDemandRunnerConfigGet(req.Config)
	if err != nil {
		return nil, err
	}

	return &pb.GetOnDemandRunnerConfigResponse{
		Config: result,
	}, nil
}

func (s *service) ListOnDemandRunnerConfigs(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOnDemandRunnerConfigsResponse, error) {
	result, err := s.state.OnDemandRunnerConfigList()
	if err != nil {
		return nil, err
	}

	return &pb.ListOnDemandRunnerConfigsResponse{Configs: result}, nil
}
