package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) UI_ListReleases(
	ctx context.Context,
	req *pb.UI_ListReleasesRequest,
) (*pb.UI_ListReleasesResponse, error) {
	var result []*pb.UI_ReleaseBundle

	return &pb.UI_ListReleasesResponse{Releases: result}, nil
}

func (s *service) UI_GetRelease(
	ctx context.Context,
	req *pb.UI_GetReleaseRequest,
) (*pb.UI_ReleaseBundle, error) {
	return nil, nil
}
