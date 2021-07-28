package singleprocess

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertStatusReport(
	ctx context.Context,
	req *pb.UpsertStatusReportRequest,
) (*pb.UpsertStatusReportResponse, error) {
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
	result, err := s.state.StatusReportList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListStatusReportsResponse{StatusReports: result}, nil
}

func (s *service) GetLatestStatusReport(
	ctx context.Context,
	req *pb.GetLatestStatusReportRequest,
) (*pb.StatusReport, error) {
	r, err := s.state.StatusReportLatest(req.Application, req.Workspace)
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

	var applicationRef *pb.Ref_Application
	switch target := req.Target.(type) {
	case *pb.ExpediteStatusReportRequest_Deployment:
		applicationRef = target.Deployment.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Deployment{
			Deployment: target.Deployment,
		}
	case *pb.ExpediteStatusReportRequest_Release:
		applicationRef = target.Release.Application
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Release{
			Release: target.Release,
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unknown status report target: %T", req.Target)
	}

	// Status report jobs need the parent project to obtain its datasource
	project, err := s.state.ProjectGet(&pb.Ref_Project{Project: applicationRef.Project})
	if err != nil {
		return nil, err
	}

	// Get target from request

	// build job
	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one poll operation at
			// any time queued per application.
			SingletonId: fmt.Sprintf("status-report-ondemand/%s", applicationRef.Application),

			Application: applicationRef,

			// Application polling requires a data source to be configured for the project
			// Otherwise a status report can't properly eval the projects hcl context
			// needed to query the deploy or release
			DataSource: project.DataSource,

			Workspace: &pb.Ref_Workspace{Workspace: "default"},

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
		Id: jobID,
	}, nil
}
