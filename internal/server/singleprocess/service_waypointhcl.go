package singleprocess

import (
	"context"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) WaypointHclFmt(
	ctx context.Context,
	req *pb.WaypointHclFmtRequest,
) (*pb.WaypointHclFmtResponse, error) {
	result, err := configpkg.Format(req.WaypointHcl, "<input>")
	if err != nil {
		return nil, err
	}

	return &pb.WaypointHclFmtResponse{WaypointHcl: result}, nil
}
