package singleprocess

import (
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func testServiceImpl(impl pb.WaypointServer) *service {
	return impl.(*service)
}

func testStateInmem(impl pb.WaypointServer) *state.State {
	return testServiceImpl(impl).state.(*state.State)
}
