package singleprocess

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverhandler/handlertest"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func init() {
	// Seed our test randomness
	rand.Seed(time.Now().UnixNano())
}

type OSSTestServerImpl struct {
	service *Service
}

func (o *OSSTestServerImpl) State(ctx context.Context) serverstate.Interface {
	return o.service.state(ctx)
}

// TestHandlers run the service handler tests that depend exclusively on the protobuf
// interfaces.
func TestHandlers(t *testing.T) {
	// Tests that are relevant, but are known to be failing.
	// It should be a priority to fix any test on this list.
	knownFailingStateTests := []string{
		//Failing b/c UI_ListEvents in the service layer calls the state layer EventListBundles, which is not implemented in boltdb
		"TestEvent",
	}

	handlertest.Test(t, func(t *testing.T) (pb.WaypointClient, handlertest.TestServerImpl) {
		impl := TestImpl(t)

		client := server.TestServer(t, impl)

		return client, &OSSTestServerImpl{service: impl.(*Service)}
	}, knownFailingStateTests)
}
