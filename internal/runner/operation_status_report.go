package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeStatusReportOp(
	ctx context.Context,
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
		statusReportResult, _, err = app.StatusReport(ctx, t.Deployment, nil)
	case *pb.Job_StatusReportOp_Release:
		statusReportResult, _, err = app.StatusReport(ctx, nil, t.Release)
	default:
		err = fmt.Errorf("unknown destruction target: %T", op.StatusReport.Target)
	}

	if err != nil {
		return nil, err
	}

	// Update to the latest deployment in order to get all the preload data.
	var statusReport *pb.StatusReport

	if statusReportResult != nil {
		statusReport, err = r.client.GetStatusReport(ctx, &pb.GetStatusReportRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: statusReportResult.Id},
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.Job_Result{
		StatusReport: &pb.Job_StatusReportResult{
			StatusReport: statusReport,
		},
	}, nil
}
