package server

import (
	"github.com/golang/protobuf/ptypes"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// NewStatus returns a new Status message with the given initial state.
func NewStatus(init pb.Status_State) *pb.Status {
	return &pb.Status{
		State:     init,
		StartTime: ptypes.TimestampNow(),
	}
}
