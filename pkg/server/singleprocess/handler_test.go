package singleprocess

import (
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverhandler/handlertest"
)

func init() {
	// Seed our test randomness
	rand.Seed(time.Now().UnixNano())
}

// TestHandlers run the service handler tests that depend exclusively on the protobuf
// interfaces.
func TestHandlers(t *testing.T) {
	handlertest.Test(t, func(t *testing.T) (pb.WaypointServer, pb.WaypointClient) {
		impl := TestImpl(t)
		client := server.TestServer(t, impl)
		return impl, client
	}, nil)
}
