package handlers

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

type Service interface {
	pb.WaypointServer

	State(ctx context.Context) serverstate.Interface

	// SuperUser - if true, forces all API actions to behave as if a superuser
	// made them. This is usually turned on for local mode only.
	SuperUser() bool

	DecodeId(id string) (decodedId string, err error)

	EncodeId(ctx context.Context, id string) string
}
