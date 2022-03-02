package server

import (
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// NewStatus returns a new Status message with the given initial state.
func NewStatus(init pb.Status_State) *pb.Status {
	return &pb.Status{
		State:     init,
		StartTime: ptypes.TimestampNow(),
	}
}

// StatusSetError sets the error state on the status and marks the
// completion time.
func StatusSetError(s *pb.Status, err error) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.Newf(codes.Internal, "Non-status error %T: %s", err, err)
	}

	s.State = pb.Status_ERROR
	s.Error = st.Proto()
	s.CompleteTime = ptypes.TimestampNow()
}

// StatusSetSuccess sets state of the status to success and marks the
// completion time.
func StatusSetSuccess(s *pb.Status) {
	s.State = pb.Status_SUCCESS
	s.CompleteTime = ptypes.TimestampNow()
}
