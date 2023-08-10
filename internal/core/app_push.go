// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/opaqueany"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Push pushes the given build to the configured registry. This requires
// that the build artifact be available, which we leave up to the caller.
// Therefore, please note that this generally can't be called on separate
// machines, long after a build is done (because the person may have deleted
// the physical artifact, etc.).
//
// TODO(mitchellh): test
func (a *App) PushBuild(ctx context.Context, optFuncs ...PushBuildOption) (*pb.PushedArtifact, error) {
	opts, err := newPushBuildOptions(optFuncs...)
	if err != nil {
		return nil, err
	}

	// Add our build to our config
	var evalCtx hcl.EvalContext
	if err := evalCtxTemplateProto(&evalCtx, "artifact", opts.Build); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	// Make our registry
	cr, err := componentCreatorMap[component.RegistryType].Create(ctx, a, &evalCtx)
	if status.Code(err) == codes.Unimplemented {
		cr = nil
		err = nil
	}
	if err != nil {
		return nil, err
	}
	defer cr.Close()

	cb, err := componentCreatorMap[component.BuilderType].Create(ctx, a, nil)
	if err != nil {
		return nil, err
	}
	defer cb.Close()

	_, msg, err := a.doOperation(ctx, a.logger.Named("push"), &pushBuildOperation{
		ComponentRegistry: cr,
		ComponentBuilder:  cb,
		Build:             opts.Build,
	})
	if err != nil {
		return nil, err
	}

	result, ok := msg.(*pb.PushedArtifact)
	if !ok {
		return nil, status.Error(codes.Internal, "app_push failed to convert the operation message into a PushedArtifact proto")
	}

	return result, nil
}

// PushBuildOption is used to configure a Build
type PushBuildOption func(*pushBuildOptions) error

// BuildWithPush sets whether or not the build will push. The default
// is for the build to push.
func PushWithBuild(b *pb.Build) PushBuildOption {
	return func(opts *pushBuildOptions) error {
		opts.Build = b
		return nil
	}
}

type pushBuildOptions struct {
	Build *pb.Build
}

func newPushBuildOptions(opts ...PushBuildOption) (*pushBuildOptions, error) {
	def := &pushBuildOptions{}
	for _, f := range opts {
		if err := f(def); err != nil {
			return nil, err
		}
	}

	return def, def.Validate()
}

type pushBuildOperation struct {
	ComponentRegistry *Component
	ComponentBuilder  *Component
	Build             *pb.Build
}

func (opts *pushBuildOptions) Validate() error {
	return validation.ValidateStruct(opts,
		validation.Field(&opts.Build, validation.Required),
	)
}

func (op *pushBuildOperation) Init(app *App) (proto.Message, error) {
	// Our component is typically the registry but if we don't have
	// one configured, then we specify the component as our builder since
	// that is what is creating the pushed artifact.
	component := op.ComponentRegistry
	if component == nil {
		component = op.ComponentBuilder
	}

	return &pb.PushedArtifact{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   component.Info,
		Labels:      component.labels,
		BuildId:     op.Build.Id,
	}, nil
}

func (op *pushBuildOperation) Hooks(app *App) map[string][]*config.Hook {
	if op.ComponentRegistry == nil {
		return nil
	}

	return op.ComponentRegistry.hooks
}

func (op *pushBuildOperation) Labels(app *App) map[string]string {
	if op.ComponentRegistry == nil {
		return nil
	}

	return op.ComponentRegistry.labels
}

func (op *pushBuildOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: msg.(*pb.PushedArtifact),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed upserting pushed artifact operation")
	}

	return resp.Artifact, nil
}

// Name returns the name of the operation
func (op *pushBuildOperation) Name() string {
	return "push build"
}

func (op *pushBuildOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	// If we have no registry, we just push the local build.
	if op.ComponentRegistry == nil {
		return op.Build.Artifact.Artifact, nil
	}

	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		op.ComponentRegistry,
		op.ComponentRegistry.Value.(component.Registry).PushFunc(),
		plugin.ArgNamedAny("artifact", op.Build.Artifact.Artifact),
	)
}

func (op *pushBuildOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.PushedArtifact).Status)
}

func (op *pushBuildOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	v := msg.(*pb.PushedArtifact)
	if v.Artifact == nil {
		v.Artifact = &pb.Artifact{}
	}

	return &v.Artifact.Artifact, &v.Artifact.ArtifactJson
}

var _ operation = (*pushBuildOperation)(nil)
