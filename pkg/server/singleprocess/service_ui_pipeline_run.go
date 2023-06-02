// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/telemetry/metrics"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

// UI_ListPipelineRuns returns pipeline runs for a given pipeline. While paginating is
// part of the request, this doesn't yet support pagination and will return
// everything every time.
func (s *Service) UI_ListPipelineRuns(
	ctx context.Context,
	req *pb.UI_ListPipelineRunsRequest,
) (*pb.UI_ListPipelineRunsResponse, error) {
	log := hclog.FromContext(ctx)

	if err := serverptypes.ValidateUIListPipelineRunsRequest(req); err != nil {
		return nil, err
	}

	// Create uninitialized array of pipeline run bundles
	var allPipelineRuns []*pb.UI_PipelineRunBundle

	// Get list of all pipeline runs
	pipelineRunListResponse, err := s.state(ctx).PipelineRunList(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error listing piplines",
		)
	}

	// Create bundles
	var pipelineRunBundle *pb.UI_PipelineRunBundle
	for _, pipelineRun := range pipelineRunListResponse {
		if len(pipelineRun.Jobs) != 0 {
			job, err := s.GetJob(ctx, &pb.GetJobRequest{
				JobId: pipelineRun.Jobs[0].Id,
			})
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to get first job for latest pipeline run",
				)
			}
			pipelineRunBundle = &pb.UI_PipelineRunBundle{
				PipelineRun:   pipelineRun,
				QueueTime:     job.QueueTime,
				Application:   job.Application,
				DataSourceRef: job.DataSourceRef,
			}
		}

		// Add pipeline bundle to uninitialized array
		allPipelineRuns = append(allPipelineRuns, pipelineRunBundle)

	}

	// Return the array
	return &pb.UI_ListPipelineRunsResponse{
		PipelineRunBundles: allPipelineRuns,
		Pagination:         &pb.PaginationResponse{},
	}, nil
}

func (s *Service) UI_GetPipelineRun(
	ctx context.Context,
	req *pb.UI_GetPipelineRunRequest,
) (*pb.UI_GetPipelineRunResponse, error) {
	log := hclog.FromContext(ctx)

	if err := serverptypes.ValidateUIGetPipelineRunRequest(req); err != nil {
		return nil, err
	}

	runResp, err := s.GetPipelineRun(ctx, &pb.GetPipelineRunRequest{
		Pipeline: req.Pipeline,
		Sequence: req.Sequence,
	})
	if err != nil {
		return nil, err
	}
	run := runResp.PipelineRun

	// Fetch full jobs
	start := time.Now()
	var jobs []*pb.Job
	for _, ref := range run.Jobs {
		job, err := s.GetJob(ctx, &pb.GetJobRequest{JobId: ref.Id})
		if err != nil {
			return nil, hcerr.Externalize(
				log,
				err,
				"failed to get jobs for all pipeline steps",
			)
		}
		jobs = append(jobs, job)
	}
	metrics.MeasureOperation(ctx, start, "fetch_jobs_for_ui_get_pipeline_run")

	// Fetch latest status report for every deployment and release
	start = time.Now()
	var statusReports []*pb.StatusReport
	for _, job := range jobs {
		if d := job.Result.GetDeploy().GetDeployment(); d != nil {
			sr, err := s.GetLatestStatusReport(ctx, &pb.GetLatestStatusReportRequest{
				Application: d.Application,
				Workspace:   d.Workspace,
				Target: &pb.GetLatestStatusReportRequest_DeploymentId{
					DeploymentId: d.Id,
				},
			})
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to get latest status report for deployment %q",
					d.Id,
				)
			}
			if sr != nil {
				statusReports = append(statusReports, sr)
			}
		}
		if r := job.Result.GetRelease().GetRelease(); r != nil {
			sr, err := s.GetLatestStatusReport(ctx, &pb.GetLatestStatusReportRequest{
				Application: r.Application,
				Workspace:   r.Workspace,
				Target: &pb.GetLatestStatusReportRequest_ReleaseId{
					ReleaseId: r.Id,
				},
			})
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to get latest status report for release %q",
					r.Id,
				)
			}
			if sr != nil {
				statusReports = append(statusReports, sr)
			}
		}
	}
	metrics.MeasureOperation(ctx, start, "fetch_latest_status_reports_for_ui_get_pipeline_run")

	runBundle := &pb.UI_PipelineRunBundle{
		PipelineRun: run,
	}
	if len(jobs) > 0 {
		j := jobs[0]
		runBundle.Application = j.Application
		runBundle.DataSourceRef = j.DataSourceRef
		runBundle.QueueTime = j.QueueTime
	}

	rootNode, err := serverptypes.UI_PipelineRunTreeFromJobs(jobs, statusReports)
	if err != nil {
		return nil, err
	}

	return &pb.UI_GetPipelineRunResponse{
		PipelineRunBundle: runBundle,
		RootTreeNode:      rootNode,
	}, nil
}
