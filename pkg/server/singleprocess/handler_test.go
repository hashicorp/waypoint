// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	handlertest.Test(t, func(t *testing.T) (pb.WaypointClient, handlertest.TestServerImpl) {
		impl := TestImpl(t)

		client := server.TestServer(t, impl)

		return client, &OSSTestServerImpl{service: impl.(*Service)}
	}, nil)
}
