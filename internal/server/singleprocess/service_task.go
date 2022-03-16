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

	// Get the Start, Run, and Stop jobs
	resp, err := s.getJobsByTaskRef(ctx, t)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

	var tasks []*pb.GetTaskResponse
	for _, t := range result {
		tsk, err := s.getJobsByTaskRef(ctx, t)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, tsk)
	}

	return &pb.ListTaskResponse{Tasks: tasks}, nil
}

func (s *service) getJobsByTaskRef(
	ctx context.Context,
	t *pb.Task,
) (*pb.GetTaskResponse, error) {
	var taskJob, startJob, stopJob *pb.Job

	if t.TaskJob != nil {
		var err error
		taskJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.TaskJob.Id})
		if err != nil {
			return nil, err
		}
	} else {
	}

	if t.StartJob == nil {
		var err error
		startJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.StartJob.Id})
		if err != nil {
			return nil, err
		}
	} else {
	}

	if t.StopJob == nil {
		var err error
		stopJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.StopJob.Id})
		if err != nil {
			return nil, err
		}
	} else {
	}

	return &pb.GetTaskResponse{Task: t, TaskJob: taskJob, StartJob: startJob, StopJob: stopJob}, nil
}
