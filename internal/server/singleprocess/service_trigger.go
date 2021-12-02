package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func (s *service) UpsertTrigger(
	ctx context.Context,
	req *pb.UpsertTriggerRequest,
) (*pb.UpsertTriggerResponse, error) {
	if err := serverptypes.ValidateUpsertTriggerRequest(req); err != nil {
		return nil, err
	}

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
	if err := serverptypes.ValidateGetTriggerRequest(req); err != nil {
		return nil, err
	}

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
	if err := serverptypes.ValidateDeleteTriggerRequest(req); err != nil {
		return nil, err
	}

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
	// NOTE: no ptype validation at the moment, as all Ref fields are optional

	result, err := s.state.TriggerList(req.Workspace, req.Project, req.Application)
	if err != nil {
		return nil, err
	}

	return &pb.ListTriggerResponse{Triggers: result}, nil
}
