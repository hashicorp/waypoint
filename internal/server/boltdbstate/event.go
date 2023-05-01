// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EventPut puts an event based on the proto passed in, currently this is not implemented in Waypoint
func (s *State) EventPut(ctx context.Context, value *serverstate.Event) error {
	return status.Errorf(codes.Unimplemented, "method EventPut not implemented")
}

// EventListBundles returns the list of events
func (s *State) EventListBundles(ctx context.Context, eventReq *pb.UI_ListEventsRequest) ([]*pb.UI_EventBundle, *pb.PaginationResponse, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.eventListBundles(memTxn, eventReq.Pagination)
}

// eventListBundles returns a list of event bundles
func (s *State) eventListBundles(
	memTxn *memdb.Txn,
	paginationRequest *pb.PaginationRequest,
) ([]*pb.UI_EventBundle, *pb.PaginationResponse, error) {
	return nil, &pb.PaginationResponse{}, status.Errorf(codes.Unimplemented, "method EventListBundles not implemented")
}
