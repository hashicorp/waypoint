package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) UpsertTask(
	ctx context.Context,
	req *pb.UpsertTaskRequest,
) (*pb.UpsertTaskResponse, error) {
	if err := serverptypes.ValidateUpsertTaskRequest(req); err != nil {
		return nil, err
	}

	result := req.Task
	if err := s.state(ctx).TaskPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertTaskResponse{Task: result}, nil
}

// GetTask returns a Task based on ID
func (s *Service) GetTask(
	ctx context.Context,
	req *pb.GetTaskRequest,
) (*pb.GetTaskResponse, error) {
	if err := serverptypes.ValidateGetTaskRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state(ctx).TaskGet(req.Ref)
	if err != nil {
		return nil, err
	}

	// Get the Start, Run, and Stop jobs
	startJob, taskJob, stopJob, err := s.state(ctx).JobsByTaskRef(t)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskResponse{
		Task:     t,
		TaskJob:  taskJob,
		StartJob: startJob,
		StopJob:  stopJob,
	}, nil
}

func (s *Service) ListTask(
	ctx context.Context,
	req *pb.ListTaskRequest,
) (*pb.ListTaskResponse, error) {
	// NOTE: no ptype validation at the moment, request params are optional

	result, err := s.state(ctx).TaskList(req)
	if err != nil {
		return nil, err
	}

	var tasks []*pb.GetTaskResponse
	for _, t := range result {
		startJob, taskJob, stopJob, err := s.state(ctx).JobsByTaskRef(t)
		if err != nil {
			return nil, err
		}
		tsk := &pb.GetTaskResponse{Task: t, TaskJob: taskJob, StartJob: startJob, StopJob: stopJob}

		tasks = append(tasks, tsk)
	}

	return &pb.ListTaskResponse{Tasks: tasks}, nil
}

func (s *Service) CancelTask(
	ctx context.Context,
	req *pb.CancelTaskRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateCancelTaskRequest(req); err != nil {
		return nil, err
	}

	if err := s.state(ctx).TaskCancel(req.Ref); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
