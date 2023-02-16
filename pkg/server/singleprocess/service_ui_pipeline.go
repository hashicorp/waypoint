package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UI_ListPipelines(
	ctx context.Context,
	req *pb.UI_ListPipelinesRequest,
) (*pb.UI_ListPipelinesResponse, error) {
	if err := serverptypes.ValidateUIListPipelinesRequest(req); err != nil {
		return nil, err
	}

	// create uninitialized array of pipeline bundles
	var allPipelines []*pb.UI_PipelineBundle
	// get list of all pipelines
	pipelineListResponse, err := s.state(ctx).PipelineList(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing piplines",
		)
	}

	// for each pipeline,

	// we will put this in a state function- right now its new db query every time. if i did it in a db function- could do it once.
	// see line 30 in service_ui_project
	// do this after i get it running
	for _, pipeline := range pipelineListResponse {
		ref := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		}
		// get total runs for the pipeline
		pipelineRunListResponse, err := s.state(ctx).PipelineRunList(ctx, ref)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"failed to count pipeline runs",
			)
		}
		// get the last run
		pipelineLastRun, err := s.state(ctx).PipelineRunGetLatest(ctx, pipeline.Id)
		if err != nil && status.Code(err) != codes.NotFound {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"failed to find latest pipeline run",
			)
		}

		pipelineBundle := &pb.UI_PipelineBundle{
			TotalRuns: uint64(len(pipelineRunListResponse)),
			LastRun:   pipelineLastRun,
		}
		// add pipeline bundle to uninitialzed array
		allPipelines = append(allPipelines, pipelineBundle)

	}

	// return the array
	return &pb.UI_ListPipelinesResponse{
		Pipelines:  allPipelines,
		Pagination: nil, // TODO: must be done before PR
	}, nil
}
