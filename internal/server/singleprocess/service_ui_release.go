package singleprocess

import (
	"context"

	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TODO: test
func (s *service) UI_ListReleases(
	cxt context.Context,
	req *pb.UI_ListReleasesRequest,
) (*pb.UI_ListReleasesResponse, error) {
	releaseList, err := s.state.ReleaseList(req.Application,
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	// TODO: make this more efficient. We should be able to just
	// grab the relevant status report in one go, not have to iterate
	// over all of them.
	statusReports, err := s.state.StatusReportList(
		req.Application,
		state.ListWithOrder(&pb.OperationOrder{
			// We only ever care about the latest status report for each operation,
			// so if we sort this way we can stop when we hit the first match.
			Order: pb.OperationOrder_COMPLETE_TIME,
			Desc:  true,
		}),
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	var releaseBundles []*pb.UI_ReleaseBundle

	for _, release := range releaseList {
		var matchingStatusReport *pb.StatusReport
	MATCH_STATUS_REPORT_LOOP:
		for _, statusReport := range statusReports {
			switch target := statusReport.TargetId.(type) {
			case *pb.StatusReport_ReleaseId:
				if target.ReleaseId == release.Id {
					matchingStatusReport = statusReport
					break MATCH_STATUS_REPORT_LOOP
				}
			}
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
