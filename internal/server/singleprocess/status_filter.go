package singleprocess

import (
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func statusFilterMatch(
	filters []*pb.StatusFilter,
	status *pb.Status,
) bool {
	if len(filters) == 0 {
		return true
	}

NEXT_FILTER:
	for _, group := range filters {
		for _, filter := range group.Filters {
			if !statusFilterMatchSingle(filter, status) {
				continue NEXT_FILTER
			}
		}

		// If any match we match (OR)
		return true
	}

	return false
}

func statusFilterMatchSingle(
	filter *pb.StatusFilter_Filter,
	status *pb.Status,
) bool {
	switch f := filter.Filter.(type) {
	case *pb.StatusFilter_Filter_State:
		return status.State == f.State

	default:
		// unknown filters never match
		return false
	}
}
