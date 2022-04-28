package singleprocess

import (
	"math/rand"
	"testing"
	"time"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverhandler/handlertest"
)

func init() {
	// Seed our test randomness
	rand.Seed(time.Now().UnixNano())
}

func Test(t *testing.T) {
	handlertest.Test(t, func(t *testing.T) pb.WaypointServer {
		return TestImpl(t)
	}, func(t *testing.T, impl pb.WaypointServer) pb.WaypointServer {
		// TODO(izaak): figure this out.
		panic("TODO(izaak): figure out restarts")
	})
}
