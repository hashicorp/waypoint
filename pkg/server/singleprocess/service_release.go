// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertRelease(
	ctx context.Context,
	req *pb.UpsertReleaseRequest,
) (*pb.UpsertReleaseResponse, error) {
	if err := ptypes.ValidateUpsertReleaseRequest(req); err != nil {
		return nil, err
	}

	result := req.Release

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				fmt.Errorf("uuid generation failed: %w", err),
				"failed to generate a uuid while upserting a release",
			)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state(ctx).ReleasePut(ctx, !insert, result); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to upsert release",
			"release_id",
			result.GetId(),
		)
	}

	return &pb.UpsertReleaseResponse{Release: result}, nil
}

// TODO: test
func (s *Service) ListReleases(
	ctx context.Context,
	req *pb.ListReleasesRequest,
) (*pb.ListReleasesResponse, error) {
	result, err := s.state(ctx).ReleaseList(ctx, req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing releases",
		)
	}

	for _, r := range result {
		if err := s.releasePreloadDetails(ctx, req.LoadDetails, r); err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error preloading release list details",
			)
		}
	}

	return &pb.ListReleasesResponse{Releases: result}, nil
}

// TODO: test
func (s *Service) GetLatestRelease(
	ctx context.Context,
	req *pb.GetLatestReleaseRequest,
) (*pb.Release, error) {
	if err := ptypes.ValidateGetLatestReleaseRequest(req); err != nil {
		return nil, err
	}

	r, err := s.state(ctx).ReleaseLatest(ctx, req.Application, req.Workspace)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting latest release",
		)
	}

	if err := s.releasePreloadDetails(ctx, req.LoadDetails, r); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error preloading release details",
		)
	}

	return r, nil
}

// GetRelease returns a Release based on ID
func (s *Service) GetRelease(
	ctx context.Context,
	req *pb.GetReleaseRequest,
) (*pb.Release, error) {
	if err := ptypes.ValidateGetReleaseRequest(req); err != nil {
		return nil, err
	}

	r, err := s.state(ctx).ReleaseGet(ctx, req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting release",
			"release",
			req.Ref.GetTarget(),
		)
	}

	if err := s.releasePreloadDetails(ctx, req.LoadDetails, r); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error preloading release details",
			"release",
			req.Ref.GetTarget(),
		)
	}

	return r, nil
}

func (s *Service) releasePreloadDetails(
	ctx context.Context,
	req pb.Release_LoadDetails,
	d *pb.Release,
) error {
	if req <= pb.Release_NONE {
		return nil
	}

	pd, err := s.state(ctx).DeploymentGet(ctx, &pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{
			Id: d.DeploymentId,
		},
	})
	if err != nil {
		return err
	}
	d.Preload.Deployment = pd

	if req > pb.Release_DEPLOYMENT {
		pa, err := s.state(ctx).ArtifactGet(ctx, &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: pd.ArtifactId,
			},
		})
		if err != nil {
			return err
		}
		d.Preload.Artifact = pa

		if req > pb.Release_ARTIFACT {
			build, err := s.state(ctx).BuildGet(ctx, &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: pa.BuildId,
				},
			})
			if err != nil {
				return err
			}

			d.Preload.Build = build
		}
	}

	return nil
}
