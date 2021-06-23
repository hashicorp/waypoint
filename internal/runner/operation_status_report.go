package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeStatusReportOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_StatusReport)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	var statusReportResult *pb.StatusReport

	switch t := op.StatusReport.Target.(type) {
	case *pb.Job_StatusReportOp_Deployment:
		statusReportResult, err = app.DeploymentStatusReport(ctx, t.Deployment)
	case *pb.Job_StatusReportOp_Release:
		statusReportResult, err = app.ReleaseStatusReport(ctx, t.Release)
	default:
		err = fmt.Errorf("unknown status report target: %T", op.StatusReport.Target)
	}

	if err != nil {
		return nil, err
	}

	if statusReportResult != nil {
		err = r.enableApplicationPoll(ctx, log, job.Application)

		if err != nil {
			return nil, err
		}
	}

	return &pb.Job_Result{
		StatusReport: &pb.Job_StatusReportResult{
			StatusReport: statusReportResult,
		},
	}, nil
}

func (r *Runner) enableApplicationPoll(
	ctx context.Context,
	log hclog.Logger,
	appRef *pb.Ref_Application,
) error {
	log.Trace("calling GetProject to determine app polling status")
	resp, err := r.client.GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: appRef.Project,
		},
	})
	if err != nil {
		return err
	}
	project := resp.Project

	for _, a := range project.Applications {
		if a.Name == appRef.Application {
			// check that polling isn't enabled for app then do, otherwise break and return

			// get project client and upsert update to app
			_, err := r.client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: project.Name},
				Name:    appRef.Application,
				Poll:    true,
			})

			if err != nil {
				return err
			} else {
				break // we found the app we were trying to update
			}
		}
	}

	return nil
}
