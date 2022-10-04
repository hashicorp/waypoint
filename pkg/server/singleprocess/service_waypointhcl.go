package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	configpkg "github.com/hashicorp/waypoint/pkg/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
)

func (s *Service) WaypointHclFmt(
	ctx context.Context,
	req *pb.WaypointHclFmtRequest,
) (*pb.WaypointHclFmtResponse, error) {
	result, err := configpkg.Format(req.WaypointHcl, "<input>")
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to format waypoint hcl",
		)
	}

	return &pb.WaypointHclFmtResponse{WaypointHcl: result}, nil
}
