package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Builds a status report on the given deployment
// TODO(briancain): test
func (a *App) StatusReport(ctx context.Context, target *pb.Deployment) (*pb.StatusReport, error) {
	var evalCtx hcl.EvalContext
	if err := evalCtxTemplateProto(&evalCtx, "deploy", target); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	c, err := a.createStatusReporter(ctx, &evalCtx)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return nil, err
	}
	defer c.Close()

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil
	}

	_, msg, err := a.doOperation(ctx, a.logger.Named("statusreport"), &statusReportOperation{
		Component: c,
		Target:    target,
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.StatusReport), nil
}

func (a *App) createStatusReporter(
	ctx context.Context,
	hclCtx *hcl.EvalContext,
) (*Component, error) {
	log := a.logger

	log.Debug("initializing status report plugin")
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, hclCtx)
	if err == nil {
		// We have a releaser configured, use that.
		return c, nil
	}

	// If we received Unimplemented, we just don't have a status report. Otherwise
	// we want to return the error we got.
	if status.Code(err) != codes.Unimplemented {
		return nil, err
	}

	// TODO remove this
	return nil, err
}

type statusReportOperation struct {
	Component *Component
	Target    *pb.Deployment

	result component.Status
}

func (op *statusReportOperation) Init(app *App) (proto.Message, error) {
	// TODO: Maybe more is needed here... deployment id?
	return &pb.StatusReport{
		Application: app.ref,
		Workspace:   app.workspace,
	}, nil
}

func (op *statusReportOperation) Hooks(app *App) map[string][]*config.Hook {
	if op.Component == nil {
		return nil
	}

	return op.Component.hooks
}

func (op *statusReportOperation) Labels(app *App) map[string]string {
	if op.Component == nil {
		return nil
	}

	return op.Component.labels
}

func (op *statusReportOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: msg.(*pb.StatusReport),
	})
	if err != nil {
		return nil, err
	}

	return resp.StatusReport, nil
}

func (op *statusReportOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	// If we have no statusRreport, we do nothing since we just update the
	// blank status report metadata.
	if op.Component == nil {
		return nil, nil
	}

	result, err := app.callDynamicFunc(ctx,
		log,
		(*component.Status)(nil),
		op.Component,
		op.Component.Value.(component.Status).StatusFunc(),
		argNamedAny("target", op.Target.Deployment),
	)
	if err != nil {
		// TODO: Look at deployments LoadDetails and how it matches the plugins
		// Deploy message proto. Need to do the same for StatusReport
		return nil, err
	}

	op.result = result.(component.Status)

	return result, nil
}

func (op *statusReportOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.StatusReport).Status)
}

func (op *statusReportOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.StatusReport).StatusReport)
}

var _ operation = (*statusReportOperation)(nil)
