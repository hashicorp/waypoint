// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"

	"github.com/hashicorp/opaqueany"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (a *App) DeploymentStatusReport(
	ctx context.Context,
	deployTarget *pb.Deployment,
) (*pb.StatusReport, error) {
	var evalCtx hcl.EvalContext
	// Load the deployment variables context
	if err := a.deployStatusReportEvalContext(ctx, deployTarget, &evalCtx); err != nil {
		return nil, err
	}

	// Load variables from deploy
	hclCtx := evalCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, err
	}

	// Get the previous status report (if any)
	lastResp, err := a.client.GetLatestStatusReport(ctx, &pb.GetLatestStatusReportRequest{
		Application: a.ref,
		Workspace:   a.workspace,
		Target: &pb.GetLatestStatusReportRequest_DeploymentId{
			DeploymentId: deployTarget.Id,
		},
	})
	if status.Code(err) == codes.NotFound {
		err = nil
		lastResp = nil
	}
	if err != nil {
		return nil, err
	}

	c, err := a.createStatusReporter(ctx, hclCtx, component.PlatformType)
	if status.Code(err) == codes.Unimplemented {
		a.logger.Debug("status report is not implemented in plugin, cannot report on status")
		return &pb.StatusReport{}, nil
	}
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return nil, err
	}

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil
	}
	defer c.Close()

	// Create the initial report. Assuming no error, we'll query ListInstances
	// to find how many active instance connections are registered for this
	// deployment. If there is a non-nil error, just return immediately.
	report, rErr := a.statusReport(ctx, "deploy_statusreport", c, deployTarget, lastResp)
	if rErr != nil {
		return report, errors.Wrapf(rErr, "failed generating status report")
	}

	// Currently the statusReport() call above will only return a nil report if
	// the error is NOT nil, but it is not an exported method so we guard here
	// in case that changes.
	if report != nil {
		// check if we have any connected instances
		resp, err := a.client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: deployTarget.Id,
			},
		})
		if err != nil {
			a.logger.Warn("error retrieving connected instances", "error", err)
			// we intentionally do not return the error from ListInstances, and
			// instead simply log the error and return the original report and
			// report error (if any)
			return report, rErr
		}

		// Modify the status report with the active instance count. The
		// statusReport() method called above both creates the report and saves
		// it to state, so modifying it here requires us to re-upsert the report
		// with the updated count.
		report.InstancesCount = uint32(len(resp.Instances))
		newReport, err := a.client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
			StatusReport: report,
		})
		if err != nil {
			a.logger.Warn("error upserting updated status report", "error", err)
			// we intentionally do not return the error from UpsertStatusReport,
			// and instead simply log the error and return the original report
			// and report error (if any)
			return report, rErr
		}

		if newReport != nil && newReport.StatusReport != nil {
			report = newReport.StatusReport
		}
	}

	return report, rErr
}

func (a *App) ReleaseStatusReport(
	ctx context.Context,
	releaseTarget *pb.Release,
) (*pb.StatusReport, error) {
	var evalCtx hcl.EvalContext
	// Load the release variables context
	if err := a.releaseStatusReportEvalContext(ctx, releaseTarget, &evalCtx); err != nil {
		return nil, err
	}

	// Load variables from deploy
	hclCtx := evalCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, err
	}

	// Get the previous status report (if any)
	lastResp, err := a.client.GetLatestStatusReport(ctx, &pb.GetLatestStatusReportRequest{
		Application: a.ref,
		Workspace:   a.workspace,
		Target: &pb.GetLatestStatusReportRequest_ReleaseId{
			ReleaseId: releaseTarget.Id,
		},
	})
	if status.Code(err) == codes.NotFound {
		err = nil
		lastResp = nil
	}
	if err != nil {
		return nil, err
	}

	c, err := a.createStatusReporter(ctx, &evalCtx, component.ReleaseManagerType)
	if status.Code(err) == codes.Unimplemented {
		a.logger.Debug("status report is not implemented in plugin, cannot report on status")
		return &pb.StatusReport{}, nil
	}
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return nil, err
	}

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil
	}
	defer c.Close()

	return a.statusReport(ctx, "release_statusreport", c, releaseTarget, lastResp)
}

// A generic status report func that takes an already setup component and a
// specific target to execute the report on.
func (a *App) statusReport(
	ctx context.Context,
	loggerName string,
	component *Component,
	target interface{},
	last *pb.StatusReport,
) (*pb.StatusReport, error) {
	if loggerName == "" {
		loggerName = "statusreport"
	}

	_, msg, err := a.doOperation(ctx, a.logger.Named(loggerName), &statusReportOperation{
		Component: component,
		Target:    target,
		Last:      last,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed performing status report op")
	}

	reportResp, ok := msg.(*pb.StatusReport)
	if !ok {
		return nil,
			status.Errorf(codes.FailedPrecondition, "unsupported status report response returned from plugin")
	}

	return reportResp, nil
}

// Sets up the eval context for a status report for deployments
func (a *App) deployStatusReportEvalContext(
	ctx context.Context,
	d *pb.Deployment,
	evalCtx *hcl.EvalContext,
) error {
	// Query the artifact
	var artifact *pb.PushedArtifact
	if d != nil && d.Preload != nil {
		artifact = d.Preload.Artifact
	}
	if artifact == nil && d != nil {
		a.logger.Debug("querying artifact", "artifact_id", d.ArtifactId)
		resp, err := a.client.GetPushedArtifact(ctx, &pb.GetPushedArtifactRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: d.ArtifactId,
				},
			},
		})
		if err != nil {
			a.logger.Error("error querying artifact",
				"artifact_id", d.ArtifactId,
				"error", err)
			return err
		}

		artifact = resp
	}

	evalCtx.Variables = map[string]cty.Value{}
	if err := evalCtxTemplateProto(evalCtx, "artifact", artifact); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}
	return nil
}

// Sets up the eval context for a status report for deployments
func (a *App) releaseStatusReportEvalContext(
	ctx context.Context,
	r *pb.Release,
	evalCtx *hcl.EvalContext,
) error {
	// Query the deployment
	a.logger.Debug("querying deployment", "deployment_id", r.DeploymentId)
	resp, err := a.client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: r.DeploymentId,
			},
		},

		LoadDetails: pb.Deployment_ARTIFACT,
	})
	if status.Code(err) == codes.NotFound {
		resp = nil
		err = nil
		a.logger.Warn("deployment not found, will attempt status report regardless",
			"deployment_id", r.DeploymentId)
	}
	if err != nil {
		a.logger.Error("error querying deployment",
			"deployment_id", r.DeploymentId,
			"error", err)
		return err
	}
	deploy := resp

	evalCtx.Variables = map[string]cty.Value{}
	// Add our build to our config if deploy found
	if deploy != nil {
		if err := evalCtxTemplateProto(evalCtx, "artifact", deploy.Preload.Artifact); err != nil {
			a.logger.Warn("failed to prepare template variables, will not be available",
				"err", err)
		}
		if err := evalCtxTemplateProto(evalCtx, "deploy", deploy); err != nil {
			a.logger.Warn("failed to prepare template variables, will not be available",
				"err", err)
		}
	}
	return nil
}

func (a *App) createStatusReporter(
	ctx context.Context,
	hclCtx *hcl.EvalContext,
	componentType component.Type,
) (*Component, error) {
	log := a.logger

	log.Debug("initializing status report plugin")
	c, err := componentCreatorMap[componentType].Create(ctx, a, hclCtx)
	if err != nil {
		if status.Code(err) != codes.Unimplemented {
			c.Close()
		}
		return nil, err
	}

	// We have a status reporter configured, use that.
	return c, nil
}

type statusReportOperation struct {
	Component *Component
	Target    interface{}      // Target to run a Status Report against
	Last      *pb.StatusReport // The last status report

	result *sdk.StatusReport
}

func (op *statusReportOperation) Init(app *App) (proto.Message, error) {
	return &pb.StatusReport{
		Application: app.ref,
		Workspace:   app.workspace,
		Health: &pb.StatusReport_Health{
			HealthStatus:  "UNKNOWN",
			HealthMessage: "Unknown health status",
		},
	}, nil
}

func (op *statusReportOperation) Hooks(app *App) map[string][]*config.Hook {
	return nil
}

func (op *statusReportOperation) Labels(app *App) map[string]string {
	return nil
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
		return nil, errors.Wrapf(err, "failed upserting status report operation")
	}

	return resp.StatusReport, nil
}

// Name returns the name of the operation
func (op *statusReportOperation) Name() string {
	return "status report"
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

	args, err := op.argsStatusReport()
	if err != nil {
		return nil, err
	}

	// realMsg is a ref to the actual proto message that gets Upserted into the DB
	realMsg, ok := msg.(*pb.StatusReport)
	if !ok {
		return nil,
			status.Errorf(codes.FailedPrecondition, "unsupported status report response returned from plugin")
	}

	// Call func on deployment _or_ release target
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

	// Populate message with results before Upsert is called
	report := result.(*sdk.StatusReport)

	// Load Status Report message compiled by the plugin into the overall generated report
	reportAny, err := opaqueany.New(report)
	if err != nil {
		return nil, err
	}
	realMsg.StatusReport = reportAny // marshall to Any

	realMsg.Health = &pb.StatusReport_Health{
		HealthStatus:  report.Health.String(),
		HealthMessage: report.HealthMessage,
	}

	// Add our resources
	var resources []*pb.StatusReport_Resource
	for _, sdkResource := range report.Resources {
		var resource pb.StatusReport_Resource
		if err := mapstructure.Decode(sdkResource, &resource); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode status report resource named %q: %s", sdkResource.Name, err)
		}
		resources = append(resources, &resource)
	}
	realMsg.Resources = resources

	// Add the deployment/release ID to the report.
	// TODO: this is a stopgap solution - we should wire resource ID information into here in a more generic way
	// rather than continue to append to this switch case.
	switch target := op.Target.(type) {
	case *pb.Deployment:
		realMsg.TargetId = &pb.StatusReport_DeploymentId{DeploymentId: target.Id}
	case *pb.Release:
		realMsg.TargetId = &pb.StatusReport_ReleaseId{ReleaseId: target.Id}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unsupported status operation type")
	}

	realMsg.GeneratedTime = report.GeneratedTime
	realMsg.External = report.External

	op.result = result.(*sdk.StatusReport)

	return result, nil
}

// args returns the args we send to the Status function call
func (op *statusReportOperation) argsStatusReport() ([]argmapper.Arg, error) {
	var args []argmapper.Arg

	switch t := op.Target.(type) {
	case *pb.Deployment:
		args = append(args, plugin.ArgNamedAny("target", t.Deployment))
	case *pb.Release:
		args = append(args, plugin.ArgNamedAny("target", t.Release))
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unsupported status report target given")
	}

	// Get our last status report. This can be nil if this is the first report.
	var last *sdk.StatusReport
	if op.Last != nil {
		var err error
		last, err = serverToSDKStatusReport(op.Last)
		if err != nil {
			return nil, err
		}
	}
	args = append(args, argmapper.Named("last_report", last))

	return args, nil
}

func (op *statusReportOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.StatusReport).Status)
}

func (op *statusReportOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	return &(msg.(*pb.StatusReport).StatusReport), &(msg.(*pb.StatusReport).StatusReportJson)
}

// serverToSDKStatusReport converts a pb.StatusReport to an SDK one.
func serverToSDKStatusReport(from *pb.StatusReport) (*sdk.StatusReport, error) {
	// We embed the status report in the message so we just unmarshal
	// the original.
	var result sdk.StatusReport
	return &result, opaqueany.UnmarshalTo(
		from.StatusReport, &result, proto.UnmarshalOptions{})
}

var _ operation = (*statusReportOperation)(nil)
