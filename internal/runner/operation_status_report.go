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

	log = log.With("app", job.Application.Application)

	log.Trace("generating status report")
	var statusReportResult *pb.StatusReport

	switch t := op.StatusReport.Target.(type) {
	case *pb.Job_StatusReportOp_Deployment:
		log.Trace("starting a status report against a deployment")
		statusReportResult, err = app.DeploymentStatusReport(ctx, t.Deployment)
	case *pb.Job_StatusReportOp_Release:
		log.Trace("starting a status report against a release")
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

// enableApplicationPoll will switch on its queued poll handler. This means
// on the defined interval, Waypoint will generate a status report for this
// application. When an application is initially inserted, like on a `waypoint init`,
// it won't enable the poller to generate a status report. This method switches
// it on after the first status report is generated. Each time after it should
// do nothinng when its polling has been enabled.
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

	// We never turn on application polling for status reports if the applications
	// project is not configured with a remote data source via git. This is because
	// the runner needs access to the project to generate a status report, and if
	// the project source is local (i.e. a local waypoint up), the remote runner
	// has no way to access the projects code. For now, we only enable application
	// polling for continuous status reports if the project has a data source configured.
	if project.DataSource == nil {
		log.Trace("cannot use status report polling if there is not a data source configured")
		return nil
	}

	if project.StatusReportPoll != nil && project.StatusReportPoll.Enabled {
		// Status report polling is already enabled
		log.Trace("application polling for status reports already enabled")
		return nil
	} else {
		log.Info("enabling application polling")
		// get project client and upsert update to app
		project.StatusReportPoll = &pb.Project_AppPoll{
			Enabled: true,
		}
		_, err := r.client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: project,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
