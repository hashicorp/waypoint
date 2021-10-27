package singleprocess

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func testServiceImpl(impl pb.WaypointServer) *service {
	return impl.(*service)
}

func testStateInmem(impl pb.WaypointServer) *state.State {
	return testServiceImpl(impl).state.(*state.State)
}
