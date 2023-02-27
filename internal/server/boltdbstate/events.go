package boltdbstate

import (
	"context"
	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// EventListBundles returns the list of events
func (s *State) EventListBundles(ctx context.Context, paginationRequest *pb.PaginationRequest) ([]*pb.UI_EventBundle, *pb.PaginationResponse, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.eventListBundles(memTxn, paginationRequest)
}

// eventListBundles returns a list of event bundles
func (s *State) eventListBundles(
	memTxn *memdb.Txn,
	paginationRequest *pb.PaginationRequest,
) ([]*pb.UI_EventBundle, *pb.PaginationResponse, error) {
	return nil, &pb.PaginationResponse{}, nil
}
