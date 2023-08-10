// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// NewStatus returns a new Status message with the given initial state.
func NewStatus(init pb.Status_State) *pb.Status {
	return &pb.Status{
		State:     init,
		StartTime: timestamppb.Now(),
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
	s.CompleteTime = timestamppb.Now()
}

// StatusSetSuccess sets state of the status to success and marks the
// completion time.
func StatusSetSuccess(s *pb.Status) {
	s.State = pb.Status_SUCCESS
	s.CompleteTime = timestamppb.Now()
}
