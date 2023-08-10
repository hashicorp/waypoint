// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UI_ListReleases(
	ctx context.Context,
	req *pb.UI_ListReleasesRequest,
) (*pb.UI_ListReleasesResponse, error) {
	releaseList, err := s.state(ctx).ReleaseList(ctx, req.Application,
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing releases",
		)
	}

	// TODO: make this more efficient. We should be able to just grab the relevant status report in one go, not have to
	// iterate over all of them.
	// NOTE(brian): We need to redo how GetLatestStatusReport is implemented. Right now it just calls its inherited func
	// from app operation to get the latest item in the database. For us to target a Release or release status report
	// we'll have to not use that abstraction and instead write our own query for grabbing a status report if a target
	// is requested.
	statusReports, err := s.state(ctx).StatusReportList(
		ctx,
		req.Application,
		// NOTE(izaak): the only implemented order for list is pb.OperationOrder_COMPLETE_TIME, which doesn't apply here.
		serverstate.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing status reports",
			"application",
			req.GetApplication(),
			"workspace",
			req.GetWorkspace(),
		)
	}

	var releaseBundles []*pb.UI_ReleaseBundle

	for _, release := range releaseList {
		var matchingStatusReport *pb.StatusReport
		for _, statusReport := range statusReports {
			switch target := statusReport.TargetId.(type) {
			case *pb.StatusReport_ReleaseId:
				if target.ReleaseId == release.Id {
					// We need to find the _latest_ status report that matches. Another opportunity for efficiency by improving the statue query
					if matchingStatusReport == nil || statusReport.GeneratedTime.GetSeconds() > matchingStatusReport.GeneratedTime.Seconds {
						matchingStatusReport = statusReport
					}
				}
			}
		}

		// Always pre-populate release details for bundles
		if err := s.releasePreloadDetails(ctx, pb.Release_BUILD, release); err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error reloading data for releases",
			)
		}

		bundle := pb.UI_ReleaseBundle{
			Release:            release,
			LatestStatusReport: matchingStatusReport,
		}
		releaseBundles = append(releaseBundles, &bundle)
	}
	return &pb.UI_ListReleasesResponse{
		Releases: releaseBundles,
	}, nil
}
