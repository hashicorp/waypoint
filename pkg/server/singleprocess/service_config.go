// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) SetConfig(
	ctx context.Context,
	req *pb.ConfigSetRequest,
) (*pb.ConfigSetResponse, error) {
	if err := s.state(ctx).ConfigSet(ctx, req.Variables...); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to set config")
	}

	return &pb.ConfigSetResponse{}, nil
}

func (s *Service) GetConfig(
	ctx context.Context,
	req *pb.ConfigGetRequest,
) (*pb.ConfigGetResponse, error) {
	if err := ptypes.ValidateGetConfigRequest(req); err != nil {
		return nil, err
	}

	vars, err := s.state(ctx).ConfigGet(ctx, req)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get config")
	}

	return &pb.ConfigGetResponse{Variables: vars}, nil
}

func (s *Service) SetConfigSource(
	ctx context.Context,
	req *pb.SetConfigSourceRequest,
) (*empty.Empty, error) {
	if err := ptypes.ValidateSetConfigSourceRequest(req); err != nil {
		return nil, err
	}

	if err := s.state(ctx).ConfigSourceSet(ctx, req.ConfigSource); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to set config source",
		)
	}

	return &empty.Empty{}, nil
}

func (s *Service) GetConfigSource(
	ctx context.Context,
	req *pb.GetConfigSourceRequest,
) (*pb.GetConfigSourceResponse, error) {
	if err := ptypes.ValidateGetConfigSourceRequest(req); err != nil {
		return nil, err
	}

	vars, err := s.state(ctx).ConfigSourceGet(ctx, req)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get config source",
		)
	}

	return &pb.GetConfigSourceResponse{ConfigSources: vars}, nil
}
