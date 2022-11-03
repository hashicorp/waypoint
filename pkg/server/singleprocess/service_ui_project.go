package singleprocess

import (
	"context"
	"sort"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UI_GetProject(
	ctx context.Context,
	req *pb.UI_GetProjectRequest,
) (*pb.UI_GetProjectResponse, error) {
	if err := serverptypes.ValidateUIGetProjectRequest(req); err != nil {
		return nil, err
	}

	project, err := s.state(ctx).ProjectGet(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting project",
		)
	}

	jobs, err := s.state(ctx).JobList(ctx, &pb.ListJobsRequest{})
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing jobs",
		)
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
