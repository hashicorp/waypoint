package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) QueueJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) (*pb.QueueJobResponse, error) {
	job := req.Job

	// Validation
	if job == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "job must be set")
	}
	if job.Id != "" {
		return nil, status.Errorf(codes.FailedPrecondition, "id must not be set")
	}

	// Get the next id
	id, err := server.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	// Queue the job
	if err := s.state.JobCreate(job); err != nil {
		return nil, err
	}

	return &pb.QueueJobResponse{JobId: job.Id}, nil
}
