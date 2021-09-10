package singleprocess

import (
	"context"
	"sort"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) UI_GetProject(
	ctx context.Context,
	req *pb.UI_GetProjectRequest,
) (*pb.UI_GetProjectResponse, error) {
	project, err := s.state.ProjectGet(req.Project)

	if err != nil {
		return nil, err
	}

	jobs, err := s.state.JobList()

	if err != nil {
		return nil, err
	}

	// Sort jobs by queue time (descending)
	sort.Slice(jobs, func(i int, j int) bool {
		ti := jobs[i].QueueTime.AsTime()
		tj := jobs[j].QueueTime.AsTime()

		return !ti.After(tj)
	})

	var latestInitJob *pb.Job

	for _, job := range jobs {
		switch job.Operation.(type) {
		case *pb.Job_Init:
			if job.Application.Project == project.Name {
				latestInitJob = job
				break
			}
		}

	}

	return &pb.UI_GetProjectResponse{
		Project:       project,
		LatestInitJob: latestInitJob,
	}, nil
}
