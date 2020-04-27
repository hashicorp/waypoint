package singleprocess

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func testServiceImpl(impl pb.WaypointServer) *service {
	return impl.(*service)
}
