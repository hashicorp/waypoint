package singleprocess

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/serverstate"
)

func (s *service) UpsertStatusReport(
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
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state.StatusReportPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertStatusReportResponse{StatusReport: result}, nil
}

func (s *service) ListStatusReports(
	ctx context.Context,
	req *pb.ListStatusReportsRequest,
) (*pb.ListStatusReportsResponse, error) {
	if err := serverptypes.ValidateListStatusReportsRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state.StatusReportList(req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
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

func (s *service) GetLatestStatusReport(
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

	r, err := s.state.StatusReportLatest(req.Application, req.Workspace, filter)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetStatusReport returns a StatusReport based on ID
func (s *service) GetStatusReport(
	ctx context.Context,
	req *pb.GetStatusReportRequest,
) (*pb.StatusReport, error) {
	if err := serverptypes.ValidateGetStatusReportRequest(req); err != nil {
		return nil, err
	}

	r, err := s.state.StatusReportGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Builds a status report job, queues it, and returns the job ID
func (s *service) ExpediteStatusReport(
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
		d, err := s.state.DeploymentGet(target.Deployment)
		if err != nil {
			return nil, err
		}

		applicationRef = d.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Deployment{
			Deployment: d,
		}
	case *pb.ExpediteStatusReportRequest_Release:
		r, err := s.state.ReleaseGet(target.Release)
		if err != nil {
			return nil, err
		}

		applicationRef = r.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Release{
			Release: r,
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unknown status report target: %T", req.Target)
	}

	// Status report jobs need the parent project to obtain its datasource
	project, err := s.state.ProjectGet(&pb.Ref_Project{Project: applicationRef.Project})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	jobID := queueJobResponse.JobId

	return &pb.ExpediteStatusReportResponse{
		JobId: jobID,
	}, nil
}
