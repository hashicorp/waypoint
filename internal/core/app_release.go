// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/opaqueany"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Release releases a set of deploys.
// TODO(mitchellh): test
func (a *App) Release(ctx context.Context, target *pb.Deployment) (
	*pb.Release,
	component.Release,
	error,
) {
	// Query the artifact
	var artifact *pb.PushedArtifact
	if target.Preload != nil && target.Preload.Artifact != nil {
		artifact = target.Preload.Artifact
	}
	if artifact == nil {
		a.logger.Debug("querying artifact", "artifact_id", target.ArtifactId)
		resp, err := a.client.GetPushedArtifact(ctx, &pb.GetPushedArtifactRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: target.ArtifactId,
				},
			},
		})
		if err != nil {
			a.logger.Error("error querying artifact",
				"artifact_id", target.ArtifactId,
				"error", err)
			return nil, nil, err
		}

		artifact = resp
	}

	// Add our config
	var evalCtx hcl.EvalContext
	if err := evalCtxTemplateProto(&evalCtx, "artifact", artifact); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}
	if err := evalCtxTemplateProto(&evalCtx, "deploy", target); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	c, err := a.createReleaser(ctx, &evalCtx)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer c.Close()

	result, releasepb, err := a.doOperation(ctx, a.logger.Named("release"), &releaseOperation{
		Component: c,
		Target:    target,
	})
	if err != nil {
		return nil, nil, err
	}

	var release component.Release
	if result != nil {
		release = result.(component.Release)
	}

	releaseResult, ok := releasepb.(*pb.Release)
	if !ok {
		return nil, nil, status.Error(codes.Internal, "app_release failed to convert the operation message into a Release proto")
	}

	return releaseResult, release, nil
}

// createReleaser creates the releaser component instance by trying to
// first load the explicit releaser, but falling back to the default releaser
// if available.
func (a *App) createReleaser(ctx context.Context, hclCtx *hcl.EvalContext) (*Component, error) {
	log := a.logger

	log.Debug("initializing release manager plugin")
	c, err := componentCreatorMap[component.ReleaseManagerType].Create(ctx, a, hclCtx)
	if err == nil {
		// We have a releaser configured, use that.
		return c, nil
	}

	// If we received Unimplemented, we just don't have a releaser. Otherwise
	// we want to return the error we got.
	if status.Code(err) != codes.Unimplemented {
		return nil, err
	}

	// No releaser. Let's try a default releaser if we can. We first
	// initialize the platform. We need to configure the eval context to
	// match a deployment.
	hclCtx = hclCtx.NewChild()
	if _, err := a.deployEvalContext(ctx, hclCtx); err != nil {
		return nil, err
	}

	log.Debug("no release manager plugin, initializing platform to check for default releaser")
	platformC, err := componentCreatorMap[component.PlatformType].Create(ctx, a, hclCtx)
	if err != nil {
		return nil, err
	}

	// Then check if the platform has a default releaser.
	pr, ok := platformC.Value.(component.PlatformReleaser)
	if !ok || pr.DefaultReleaserFunc() == nil {
		log.Debug("default releaser not supported by platform, stopping")
		platformC.Close()
		return nil, status.Errorf(codes.Unimplemented, "no releaser is supported by the requested platform")
	}

	// It does! Initialize it.
	log.Debug("default releaser supported! initializing...")
	raw, err := a.callDynamicFunc(
		ctx,
		log,
		(*component.ReleaseManager)(nil),
		platformC,
		pr.DefaultReleaserFunc(),
	)
	if err != nil {
		platformC.Close()
		return nil, err
	}

	// Set the value
	platformC.Value = raw

	// Do NOT close the platformC here, since the platform component
	// is the plugin instance that also is holding our default releaser.
	// Return to the user and let them close it.

	return platformC, nil
}

type releaseOperation struct {
	Component *Component
	Target    *pb.Deployment

	result component.Release
}

func (op *releaseOperation) Init(app *App) (proto.Message, error) {
	release := &pb.Release{
		Application:   app.ref,
		Workspace:     app.workspace,
		DeploymentId:  op.Target.Id,
		State:         pb.Operation_CREATED,
		Component:     op.Target.Component,
		Unimplemented: true,
	}

	if v := op.Component; v != nil {
		release.Component = v.Info
		release.Labels = v.labels
	}

	if op.result != nil {
		release.Url = op.result.URL()
	}

	return release, nil
}

func (op *releaseOperation) Hooks(app *App) map[string][]*config.Hook {
	if op.Component == nil {
		return nil
	}

	return op.Component.hooks
}

func (op *releaseOperation) Labels(app *App) map[string]string {
	if op.Component == nil {
		return nil
	}

	return op.Component.labels
}

func (op *releaseOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: msg.(*pb.Release),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed upserting release operation")
	}

	return resp.Release, nil
}

// Name returns the name of the operation
func (op *releaseOperation) Name() string {
	return "release"
}

func (op *releaseOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	// If we have no releaser, we do nothing since we just update the
	// blank release metadata.
	if op.Component == nil {
		return nil, nil
	}

	declaredResourcesResp := &component.DeclaredResourcesResp{}

	result, err := app.callDynamicFunc(ctx,
		log,
		(*component.Release)(nil),
		op.Component,
		op.Component.Value.(component.ReleaseManager).ReleaseFunc(),
		plugin.ArgNamedAny("target", op.Target.Deployment),
		argmapper.Typed(declaredResourcesResp),
	)
	if err != nil {
		return nil, err
	}

	op.result = result.(component.Release)

	rm := msg.(*pb.Release)
	rm.Unimplemented = false
	rm.Url = op.result.URL()

	// Convert from the plugin declaredResources to server declaredResources. Should be identical.
	declaredResources := make([]*pb.DeclaredResource, len(declaredResourcesResp.DeclaredResources))
	for i, pluginDeclaredResource := range declaredResourcesResp.DeclaredResources {
		var serverDeclaredResource pb.DeclaredResource
		if err := mapstructure.Decode(pluginDeclaredResource, &serverDeclaredResource); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode plugin declared resource named %q: %s", pluginDeclaredResource.Name, err)
		}
		declaredResources[i] = &serverDeclaredResource
	}

	rm.DeclaredResources = declaredResources

	return result, nil
}

func (op *releaseOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Release).Status)
}

func (op *releaseOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	return &(msg.(*pb.Release).Release), &(msg.(*pb.Release).ReleaseJson)
}

var _ operation = (*releaseOperation)(nil)
