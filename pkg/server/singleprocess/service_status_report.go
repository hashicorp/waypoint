package singleprocess

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertStatusReport(
	ctx context.Context,
	req *pb.UpsertStatusReportRequest,
) (*pb.UpsertStatusReportResponse, error) {
	if err := serverptypes.ValidateUpsertStatusReportRequest(req); err != nil {
		return nil, err
	}

	result := req.StatusReport

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"failed to generate a uuid while upserting a status report",
			)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state(ctx).StatusReportPut(ctx, !insert, result); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error upserting status report",
		)
	}

	return &pb.UpsertStatusReportResponse{StatusReport: result}, nil
}

func (s *Service) ListStatusReports(
	ctx context.Context,
	req *pb.ListStatusReportsRequest,
) (*pb.ListStatusReportsResponse, error) {
	if err := serverptypes.ValidateListStatusReportsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).StatusReportList(ctx, req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing status reports",
		)
	}

	var response *pb.ListStatusReportsResponse
	switch req.Target.(type) {
	case *pb.ListStatusReportsRequest_Deployment:
		var r []*pb.StatusReport

		for _, sr := range result {
			if _, ok := sr.TargetId.(*pb.StatusReport_DeploymentId); ok {
				r = append(r, sr)
			}
		}

		response = &pb.ListStatusReportsResponse{StatusReports: r}
	case *pb.ListStatusReportsRequest_Release:
		var r []*pb.StatusReport

		for _, sr := range result {
			if _, ok := sr.TargetId.(*pb.StatusReport_ReleaseId); ok {
				r = append(r, sr)
			}
		}
		response = &pb.ListStatusReportsResponse{StatusReports: r}
	default:
		response = &pb.ListStatusReportsResponse{StatusReports: result}
	}

	return response, nil
}

func (s *Service) GetLatestStatusReport(
	ctx context.Context,
	req *pb.GetLatestStatusReportRequest,
) (*pb.StatusReport, error) {
	if err := serverptypes.ValidateGetLatestStatusReportRequest(req); err != nil {
		return nil, err
	}

	filter := func(r *pb.StatusReport) (bool, error) {
		switch target := req.Target.(type) {
		case *pb.GetLatestStatusReportRequest_Any:
			return true, nil

		case *pb.GetLatestStatusReportRequest_DeploymentAny:
			_, ok := r.TargetId.(*pb.StatusReport_DeploymentId)
			return ok, nil

		case *pb.GetLatestStatusReportRequest_ReleaseAny:
			_, ok := r.TargetId.(*pb.StatusReport_ReleaseId)
			return ok, nil

		case *pb.GetLatestStatusReportRequest_DeploymentId:
			id, ok := r.TargetId.(*pb.StatusReport_DeploymentId)
			return ok && id.DeploymentId == target.DeploymentId, nil

		case *pb.GetLatestStatusReportRequest_ReleaseId:
			id, ok := r.TargetId.(*pb.StatusReport_ReleaseId)
			return ok && id.ReleaseId == target.ReleaseId, nil

		case nil:
			// Nil is allowed for backwards compatibility before we had
			// Target and is equal to Any.
			return true, nil

		default:
			// This shouldn't happen for valid proto clients.
			return false, status.Errorf(codes.FailedPrecondition,
				"invalid target type for request")
		}
	}

	r, err := s.state(ctx).StatusReportLatest(ctx, req.Application, req.Workspace, filter)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error retrieving latest status reports",
			"application",
			req.Application.GetApplication(),
		)
	}

	return r, nil
}

// GetStatusReport returns a StatusReport based on ID
func (s *Service) GetStatusReport(
	ctx context.Context,
	req *pb.GetStatusReportRequest,
) (*pb.StatusReport, error) {
	if err := serverptypes.ValidateGetStatusReportRequest(req); err != nil {
		return nil, err
	}

	r, err := s.state(ctx).StatusReportGet(ctx, req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting status report",
		)
	}

	return r, nil
}

// Builds a status report job, queues it, and returns the job ID
func (s *Service) ExpediteStatusReport(
	ctx context.Context,
	req *pb.ExpediteStatusReportRequest,
) (*pb.ExpediteStatusReportResponse, error) {
	statusReportJob := &pb.Job_StatusReport{
		StatusReport: &pb.Job_StatusReportOp{},
	}

	// Get target from request
	var applicationRef *pb.Ref_Application
	switch target := req.Target.(type) {
	case *pb.ExpediteStatusReportRequest_Deployment:
		d, err := s.state(ctx).DeploymentGet(ctx, target.Deployment)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error getting deployment for expedited status report",
			)
		}

		applicationRef = d.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Deployment{
			Deployment: d,
		}
	case *pb.ExpediteStatusReportRequest_Release:
		r, err := s.state(ctx).ReleaseGet(ctx, target.Release)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error getting release in expedited status report",
			)
		}

		applicationRef = r.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Release{
			Release: r,
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unknown status report target: %T", req.Target)
	}

	// Status report jobs need the parent project to obtain its datasource
	project, err := s.state(ctx).ProjectGet(ctx, &pb.Ref_Project{Project: applicationRef.Project})
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting project in expedited status report",
			"project",
			applicationRef.GetProject(),
		)
	}

	var workspace string
	if req.Workspace == nil {
		workspace = "default"
	} else {
		workspace = req.Workspace.Workspace
	}

	// build job
	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one on demand
			// status report any time queued per application.
			SingletonId: fmt.Sprintf("status-report-ondemand/%s", applicationRef.Application),

			Application: applicationRef,

			// Status reports requires a data source to be configured for the project
			// Otherwise a status report can't properly eval the projects hcl context
			// needed to query the deploy or release
			DataSource: project.DataSource,

			Workspace: &pb.Ref_Workspace{Workspace: workspace},

			// Generate a status report
			Operation: statusReportJob,

			// Any runner is fine for polling.
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		},
	}

	queueJobResponse, err := s.QueueJob(ctx, jobRequest)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error queueing job for expedited status report",
		)
	}
	jobID := queueJobResponse.JobId

	return &pb.ExpediteStatusReportResponse{
		JobId: jobID,
	}, nil
}
