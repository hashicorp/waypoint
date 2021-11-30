package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) UpsertTrigger(
	ctx context.Context,
	req *pb.UpsertTriggerRequest,
) (*pb.UpsertTriggerResponse, error) {
	// TODO: Validate request with ptypes
	return nil, nil
}

// GetTrigger returns a Trigger based on ID
func (s *service) GetTrigger(
	ctx context.Context,
	req *pb.GetTriggerRequest,
) (*pb.GetTriggerResponse, error) {
	// TODO: Validate request with ptypes
	return nil, nil
}

// DeleteTrigger deletes a Trigger based on ID
func (s *service) DeleteTrigger(
	ctx context.Context,
	req *pb.DeleteTriggerRequest,
) (*empty.Empty, error) {
	// TODO: Validate request with ptypes
	return nil, nil
}

func (s *service) ListTriggers(
	ctx context.Context,
	req *pb.ListTriggerRequest,
) (*pb.ListTriggerResponse, error) {
	// TODO: Validate request with ptypes
	return nil, nil
}
