package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *service) UpsertTask(
	ctx context.Context,
	req *pb.UpsertTaskRequest,
) (*pb.UpsertTaskResponse, error) {
	if err := serverptypes.ValidateUpsertTaskRequest(req); err != nil {
		return nil, err
	}

	result := req.Task
	if err := s.state.TaskPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertTaskResponse{Task: result}, nil
}

// GetTask returns a Task based on ID
func (s *service) GetTask(
	ctx context.Context,
	req *pb.GetTaskRequest,
) (*pb.GetTaskResponse, error) {
	if err := serverptypes.ValidateGetTaskRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state.TaskGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskResponse{Task: t}, nil
}

// DeleteTask deletes a Task based on ID
func (s *service) DeleteTask(
	ctx context.Context,
	req *pb.DeleteTaskRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteTaskRequest(req); err != nil {
		return nil, err
	}

	err := s.state.TaskDelete(req.Ref)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) ListTask(
	ctx context.Context,
	req *pb.ListTaskRequest,
) (*pb.ListTaskResponse, error) {
	// NOTE: no ptype validation at the moment, there are no request params

	result, err := s.state.TaskList()
	if err != nil {
		return nil, err
	}

	return &pb.ListTaskResponse{Tasks: result}, nil
}
