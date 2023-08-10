// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.UpsertOnDemandRunnerConfigRequest,
) (*pb.UpsertOnDemandRunnerConfigResponse, error) {
	log := hclog.FromContext(ctx)
	if err := serverptypes.ValidateUpsertOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	if req.Config.TargetRunner == nil {
		req.Config.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{},
		}
	}

	result, err := s.state(ctx).OnDemandRunnerConfigPut(ctx, req.Config)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed setting on-demand runner config", "id", req.Config.Id, "name", req.Config.Name)
	}

	return &pb.UpsertOnDemandRunnerConfigResponse{Config: result}, nil
}

func (s *Service) GetOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.GetOnDemandRunnerConfigRequest,
) (*pb.GetOnDemandRunnerConfigResponse, error) {
	log := hclog.FromContext(ctx)
	if err := serverptypes.ValidateGetOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).OnDemandRunnerConfigGet(ctx, req.Config)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to get on-demand runner config", "id", req.Config.Id, "name", req.Config.Name)
	}

	return &pb.GetOnDemandRunnerConfigResponse{
		Config: result,
	}, nil
}

func (s *Service) GetDefaultOnDemandRunnerConfig(
	ctx context.Context,
	req *empty.Empty,
) (*pb.GetOnDemandRunnerConfigResponse, error) {
	log := hclog.FromContext(ctx)

	results, err := s.state(ctx).OnDemandRunnerConfigDefault(ctx)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to get default on-demand runner config")
	}

	var result *pb.OnDemandRunnerConfig
	if len(results) > 0 {
		// NOTE(briancain): we only ever care about the *First* default runner profile,
		// because we've set ODR profiles to only ever allow for ONE default. The fact
		// that the state version returns a slice is an artifact of when it was first
		// implemented and you could have multiple defaults.
		odr := results[0]

		result, err = s.state(ctx).OnDemandRunnerConfigGet(ctx, odr)
		if err != nil {
			return nil, hcerr.Externalize(log, err, "failed to get on-demand runner config", "id", odr.Id, "name", odr.Name)
		}
	}

	return &pb.GetOnDemandRunnerConfigResponse{
		Config: result,
	}, nil
}

func (s *Service) ListOnDemandRunnerConfigs(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOnDemandRunnerConfigsResponse, error) {
	log := hclog.FromContext(ctx)
	result, err := s.state(ctx).OnDemandRunnerConfigList(ctx)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to list on-demand runner configs")
	}

	return &pb.ListOnDemandRunnerConfigsResponse{Configs: result}, nil
}

func (s *Service) DeleteOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.DeleteOnDemandRunnerConfigRequest,
) (*pb.DeleteOnDemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateDeleteOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	// Check that runner config exists
	resp, err := s.GetOnDemandRunnerConfig(ctx, &pb.GetOnDemandRunnerConfigRequest{Config: req.Config})
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to generate get on-demand runner config when trying to delete it",
		)
	}

	// Delete the runner config
	err = s.state(ctx).OnDemandRunnerConfigDelete(ctx, req.Config)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to delete on-demand runner config",
		)
	}
	result := resp.Config

	return &pb.DeleteOnDemandRunnerConfigResponse{Config: result}, nil
}
