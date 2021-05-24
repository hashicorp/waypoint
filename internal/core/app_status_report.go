package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/zclconf/go-cty/cty"
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

func (a *App) DeploymentStatusReport(
	ctx context.Context,
	deployTarget *pb.Deployment,
) (*pb.StatusReport, *sdk.StatusReport, error) {
	var evalCtx hcl.EvalContext
	// Load the deployment variables context
	if err := a.deployStatusReportEvalContext(ctx, deployTarget, &evalCtx); err != nil {
		return nil, nil, err
	}

	// Load variables from deploy
	hclCtx := evalCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, nil, err
	}

	c, err := a.createStatusReporter(ctx, hclCtx, component.PlatformType)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return nil, nil, err
	}

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil, nil
	}
	defer c.Close()

	return a.statusReport(ctx, "deploy_statusreport", c, deployTarget)
}

func (a *App) ReleaseStatusReport(
	ctx context.Context,
	releaseTarget *pb.Release,
) (*pb.StatusReport, *sdk.StatusReport, error) {
	var evalCtx hcl.EvalContext
	// Load the deployment variables context
	if err := a.releaseStatusReportEvalContext(ctx, releaseTarget, &evalCtx); err != nil {
		return nil, nil, err
	}

	// Load variables from deploy
	hclCtx := evalCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, nil, err
	}

	c, err := a.createStatusReporter(ctx, &evalCtx, component.ReleaseManagerType)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return nil, nil, err
	}

	a.logger.Debug("starting status report operation")
	statusReporter, ok := c.Value.(component.Status)

	if !ok || statusReporter.StatusFunc() == nil {
		a.logger.Debug("component is not a Status or has no StatusFunc()")
		return nil, nil, nil
	}
	defer c.Close()

	return a.statusReport(ctx, "release_statusreport", c, releaseTarget)
}

// A generic status report func that takes an already setup component and a
// specific target to execute the report on.
func (a *App) statusReport(
	ctx context.Context,
	loggerName string,
	component *Component,
	target interface{},
) (*pb.StatusReport, *sdk.StatusReport, error) {
	if loggerName == "" {
		loggerName = "statusreport"
	}

	result, msg, err := a.doOperation(ctx, a.logger.Named(loggerName), &statusReportOperation{
		Component: component,
		Target:    target,
	})
	if err != nil {
		return nil, nil, err
	}

	var statusReport *sdk.StatusReport
	if result != nil {
		statusReport = result.(*sdk.StatusReport)
	}

	reportResp, ok := msg.(*pb.StatusReport)
	if !ok {
		return nil, nil,
			status.Errorf(codes.FailedPrecondition, "unsupported status report response returned from plugin")
	}

	return reportResp, statusReport, nil
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
	Target    interface{} // Target to run a Status Report against

	result *sdk.StatusReport
}

func (op *statusReportOperation) Init(app *App) (proto.Message, error) {
	return &pb.StatusReport{
		Application: app.ref,
		Workspace:   app.workspace,
		ResourcesHealth: []*pb.StatusReport_Health{
			{
				HealthStatus:  "UNKNOWN",
				HealthMessage: "Unknown health status",
			},
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

	args, err := op.argsStatusReport()
	if err != nil {
		return nil, err
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

	op.result = result.(*sdk.StatusReport)

	return result, nil
}

// args returns the args we send to the Status function call
func (op *statusReportOperation) argsStatusReport() ([]argmapper.Arg, error) {
	var args []argmapper.Arg

	switch t := op.Target.(type) {
	case *pb.Deployment:
		args = append(args, argNamedAny("target", t.Deployment))
	case *pb.Release:
		args = append(args, argNamedAny("target", t.Release))
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "unsupported status report target given")
	}

	return args, nil
}

func (op *statusReportOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.StatusReport).Status)
}

func (op *statusReportOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.StatusReport).StatusReport)
}

var _ operation = (*statusReportOperation)(nil)
