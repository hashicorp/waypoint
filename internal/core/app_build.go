// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
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

// Build builds the artifact from source for this app.
// TODO(mitchellh): test
func (a *App) Build(ctx context.Context, optFuncs ...BuildOption) (
	*pb.Build,
	*pb.PushedArtifact,
	error,
) {
	opts, err := newBuildOptions(optFuncs...)
	if err != nil {
		return nil, nil, err
	}

	// Render the config
	c, err := componentCreatorMap[component.BuilderType].Create(ctx, a, nil)
	if err != nil {
		return nil, nil, err
	}
	defer c.Close()

	cr, err := componentCreatorMap[component.RegistryType].Create(ctx, a, nil)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			cr = nil
			err = nil
		} else {
			return nil, nil, err
		}
	}

	if cr != nil {
		defer cr.Close()
	}

	// First we do the build
	_, msg, err := a.doOperation(ctx, a.logger.Named("build"), &buildOperation{
		Component:   c,
		Registry:    cr,
		HasRegistry: cr != nil,
	})
	if err != nil {
		return nil, nil, err
	}
	build, ok := msg.(*pb.Build)
	if !ok {
		return nil, nil, status.Error(codes.Internal,
			"app_build failed to convert the operation message into a Build proto message")
	}

	// If we're not pushing, then we're done!
	if !opts.Push {
		return build, nil, nil
	}

	// We're also pushing to a registry, so invoke that.
	artifact, err := a.PushBuild(ctx, PushWithBuild(build))
	return build, artifact, err
}

// Name returns the name of the operation
func (op *buildOperation) Name() string {
	return "build"
}

// BuildOption is used to configure a Build
type BuildOption func(*buildOptions) error

// BuildWithPush sets whether or not the build will push. The default
// is for the build to push.
func BuildWithPush(v bool) BuildOption {
	return func(opts *buildOptions) error {
		opts.Push = v
		return nil
	}
}

type buildOptions struct {
	Push bool
}

func defaultBuildOptions() *buildOptions {
	return &buildOptions{
		Push: true,
	}
}

func newBuildOptions(opts ...BuildOption) (*buildOptions, error) {
	def := defaultBuildOptions()
	for _, f := range opts {
		if err := f(def); err != nil {
			return nil, err
		}
	}

	return def, def.Validate()
}

func (opts *buildOptions) Validate() error {
	return nil
}

// buildOperation implements the operation interface.
type buildOperation struct {
	Component *Component
	Registry  *Component
	Build     *pb.Build

	HasRegistry bool
}

func (op *buildOperation) Init(app *App) (proto.Message, error) {
	return &pb.Build{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   op.Component.Info,
	}, nil
}

func (op *buildOperation) Hooks(app *App) map[string][]*config.Hook {
	return op.Component.hooks
}

func (op *buildOperation) Labels(app *App) map[string]string {
	return op.Component.labels
}

func (op *buildOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertBuild(ctx, &pb.UpsertBuildRequest{
		Build: msg.(*pb.Build),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed upserting build operation")
	}

	return resp.Build, nil
}

func (op *buildOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	args := []argmapper.Arg{
		argmapper.Named("HasRegistry", op.HasRegistry),
	}

	// If there is a registry defined and it implements RegistryAccess...
	if op.Registry != nil {
		if ra, ok := op.Registry.Value.(component.RegistryAccess); ok && ra.AccessInfoFunc() != nil {
			raw, err := app.callDynamicFunc(ctx, log, nil, op.Component, ra.AccessInfoFunc())
			if err == nil {
				args = append(args, argmapper.Typed(raw))

				if pm, ok := raw.(interface {
					TypedAny() *opaqueany.Any
				}); ok {
					any := pm.TypedAny()

					// ... which we make available to build plugin.
					args = append(args, plugin.ArgNamedAny("access_info", any))
					log.Debug("injected access info")
				} else {
					log.Error("unexpected response type from callDynamicFunc", "type", hclog.Fmt("%T", raw))
					return nil, errors.New("AccessInfoFunc didn't provide a typed any")
				}
			} else {
				log.Error("error calling dynamic func", "error", err)
				return nil, err
			}
		} else {
			if ok && ra != nil && ra.AccessInfoFunc() == nil {
				return nil, status.Error(codes.Internal, "The plugin requested does not "+
					"define an AccessInfoFunc() in its Registry plugin. This is an internal "+
					"error and should be reported to the author of the plugin.")
			}
		}
	}

	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		op.Component,
		op.Component.Value.(component.Builder).BuildFunc(),
		args...,
	)
}

func (op *buildOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Build).Status)
}

func (op *buildOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	v := msg.(*pb.Build)
	if v.Artifact == nil {
		v.Artifact = &pb.Artifact{}
	}

	return &v.Artifact.Artifact, &v.Artifact.ArtifactJson
}

var _ operation = (*buildOperation)(nil)
