package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Builds a status report on the given deployment
func (a *App) StatusReport(
	ctx context.Context,
	deployTarget *pb.Deployment,
	releaseTarget *pb.Release,
) (*pb.StatusReport, *sdk.StatusReport, error) {
	var evalCtx hcl.EvalContext
	if err := evalCtxTemplateProto(&evalCtx, "deploy", deployTarget); err != nil {
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
		return nil, nil, err
	}
	defer c.Close()

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil, nil
	}

	result, msg, err := a.doOperation(ctx, a.logger.Named("statusreport"), &statusReportOperation{
		Component:     c,
		DeployTarget:  deployTarget,
		ReleaseTarget: releaseTarget,
	})
	if err != nil {
		return nil, nil, err
	}
	var status *sdk.StatusReport
	if result != nil {
		status = result.(*sdk.StatusReport)
	}

	return msg.(*pb.StatusReport), status, nil
}

func (a *App) createStatusReporter(
	ctx context.Context,
	hclCtx *hcl.EvalContext,
) (*Component, error) {
	log := a.logger

	// Load variables from deploy
	hclCtx = hclCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, err
	}

	log.Debug("initializing status report plugin")
	// Potential bug here with k8s apply plugin
	// Works with docker and k8s ok
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, hclCtx)
	if err == nil {
		// We have a status reporter configured, use that.
		return c, nil
	}

	// If we received Unimplemented, we just don't have a status report. Otherwise
	// we want to return the error we got.
	if status.Code(err) != codes.Unimplemented {
		c.Close()
		return nil, err
	}

	return nil, err
}

type statusReportOperation struct {
	Component     *Component
	DeployTarget  *pb.Deployment
	ReleaseTarget *pb.Release

	result *sdk.StatusReport
}

func (op *statusReportOperation) Init(app *App) (proto.Message, error) {
	return &pb.StatusReport{
		Application:   app.ref,
		Workspace:     app.workspace,
		HealthStatus:  "UNKNOWN",
		HealthMessage: "Unknown health status",
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

func (op *statusReportOperation) Do(
	ctx context.Context,
	log hclog.Logger,
	app *App,
	msg proto.Message,
) (interface{}, error) {
	// If we have no statusReport, we do nothing since we just update the
	// blank status report metadata.
	if op.Component == nil {
		return nil, nil
	}

	// Call func on deployment _or_ release target
	var args []argmapper.Arg
	if op.DeployTarget != nil && op.DeployTarget.Deployment != nil {
		args = append(args, argNamedAny("target", op.DeployTarget.Deployment))
	} else if op.ReleaseTarget != nil && op.ReleaseTarget.Release != nil {
		args = append(args, argNamedAny("target", op.ReleaseTarget.Release))
	} else {
		return nil, status.Errorf(codes.FailedPrecondition, "unsupported status report target given")
	}

	result, err := app.callDynamicFunc(ctx,
		log,
		nil,
		op.Component,
		op.Component.Value.(component.Status).StatusFunc(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	op.result = result.(*sdk.StatusReport)

	return result, nil
}

func (op *statusReportOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.StatusReport).Status)
}

func (op *statusReportOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.StatusReport).StatusReport)
}

var _ operation = (*statusReportOperation)(nil)
