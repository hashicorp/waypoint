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

	result := req.Trigger
	if err := s.state.TriggerPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertTriggerResponse{Trigger: result}, nil
}

// GetTrigger returns a Trigger based on ID
func (s *service) GetTrigger(
	ctx context.Context,
	req *pb.GetTriggerRequest,
) (*pb.GetTriggerResponse, error) {
	// TODO: Validate request with ptypes

	t, err := s.state.TriggerGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return &pb.GetTriggerResponse{Trigger: t}, nil
}

// DeleteTrigger deletes a Trigger based on ID
func (s *service) DeleteTrigger(
	ctx context.Context,
	req *pb.DeleteTriggerRequest,
) (*empty.Empty, error) {
	// TODO: Validate request with ptypes

	err := s.state.TriggerDelete(req.Ref)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) ListTriggers(
	ctx context.Context,
	req *pb.ListTriggerRequest,
) (*pb.ListTriggerResponse, error) {
	// TODO: Validate request with ptypes

	result, err := s.state.TriggerList()
	if err != nil {
		return nil, err
	}

	return &pb.ListTriggerResponse{Triggers: result}, nil
}
