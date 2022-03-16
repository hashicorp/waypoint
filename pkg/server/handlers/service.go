package handlers

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// Service is an implementation of pb.WaypointServer capable of
// using shared service handlers in this package.
type Service interface {
	pb.WaypointServer

	// State returns something capable of persisting
	// waypoint's state (i.e. a boltdb-backed state struct)
	State(ctx context.Context) serverstate.Interface

	// SuperUser - if true, forces all API actions to behave as if a superuser
	// made them. This is usually turned on for local mode only.
	SuperUser() bool

	// DecodeId takes a string that contains an ID (likely created
	// with EncodeId), and returns only the waypoint-relevant ID.
	DecodeId(encodedId string) (id string, err error)

	// EncodeId takes a waypoint ID (user id, runner id, etc.),
	// uses the provided context to encode additional metadata
	// (if present), and returns an ID that can be decoded by DecodeId.
	EncodeId(ctx context.Context, id string) (encodedId string)
}
