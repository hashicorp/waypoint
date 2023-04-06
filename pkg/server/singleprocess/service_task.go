// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertTask(
	ctx context.Context,
	req *pb.UpsertTaskRequest,
) (*pb.UpsertTaskResponse, error) {
	if err := serverptypes.ValidateUpsertTaskRequest(req); err != nil {
		return nil, err
	}

	result := req.Task
	if err := s.state(ctx).TaskPut(ctx, result); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to upsert task", "id", req.Task.Id)
	}

	return &pb.UpsertTaskResponse{Task: result}, nil
}

// GetTask returns a Task based on ID
func (s *Service) GetTask(
	ctx context.Context,
	req *pb.GetTaskRequest,
) (*pb.GetTaskResponse, error) {
	log := hclog.FromContext(ctx)
	if err := serverptypes.ValidateGetTaskRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state(ctx).TaskGet(ctx, req.Ref)
	if err != nil {
		var refArgs []interface{}
		switch r := req.Ref.Ref.(type) {
		case *pb.Ref_Task_Id:
			refArgs = append(refArgs, "id", r.Id)
		case *pb.Ref_Task_JobId:
			refArgs = append(refArgs, "job_id", r.JobId)
		}
		return nil, hcerr.Externalize(log, err, "failed to get task", refArgs...)
	}

	// Get the Start, Run, and Stop jobs
	startJob, taskJob, stopJob, watchJob, err := s.state(ctx).JobsByTaskRef(ctx, t)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to get jobs for task", "id", t.Id)
	}

	return &pb.GetTaskResponse{
		Task:     t,
		TaskJob:  taskJob,
		StartJob: startJob,
		StopJob:  stopJob,
		WatchJob: watchJob,
	}, nil
}

func (s *Service) ListTask(
	ctx context.Context,
	req *pb.ListTaskRequest,
) (*pb.ListTaskResponse, error) {
	// NOTE: no ptype validation at the moment, request params are optional

	log := hclog.FromContext(ctx)
	result, err := s.state(ctx).TaskList(ctx, req)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to list tasks")
	}

	var tasks []*pb.GetTaskResponse
	for _, t := range result {
		startJob, taskJob, stopJob, watchJob, err := s.state(ctx).JobsByTaskRef(ctx, t)
		if err != nil {
			return nil, hcerr.Externalize(log, err, "failed to get jobs for task", "id", t.Id)
		}
		tsk := &pb.GetTaskResponse{Task: t,
			TaskJob: taskJob, StartJob: startJob, StopJob: stopJob, WatchJob: watchJob}

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

	if err := s.state(ctx).TaskCancel(ctx, req.Ref); err != nil {
		var refArgs []interface{}
		switch r := req.Ref.Ref.(type) {
		case *pb.Ref_Task_Id:
			refArgs = append(refArgs, "id", r.Id)
		case *pb.Ref_Task_JobId:
			refArgs = append(refArgs, "job_id", r.JobId)
		}
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to cancel task", refArgs...)
	}

	return &empty.Empty{}, nil
}
