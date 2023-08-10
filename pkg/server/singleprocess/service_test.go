// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"

	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func testServiceImpl(impl pb.WaypointServer) *Service {
	return impl.(*Service)
}

func testStateInmem(impl pb.WaypointServer) *boltdbstate.State {
	return testServiceImpl(impl).state(context.Background()).(*boltdbstate.State)
}
