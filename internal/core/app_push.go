package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	_, msg, err := a.doOperation(ctx, a.logger.Named("push"), &pushBuildOperation{
		Build: opts.Build,
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.PushedArtifact), nil
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
	Build *pb.Build
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
	var component interface{} = app.Registry
	if component == nil {
		component = app.Builder
	}

	return &pb.PushedArtifact{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   app.components[component].Info,
		Labels:      app.components[component].Labels,
		BuildId:     op.Build.Id,
	}, nil
}

func (op *pushBuildOperation) Hooks(app *App) map[string][]*config.Hook {
	if app.Registry == nil {
		return nil
	}

	return app.components[app.Registry].Hooks
}

func (op *pushBuildOperation) Labels(app *App) map[string]string {
	if app.Registry == nil {
		return nil
	}

	return app.components[app.Registry].Labels
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
		return nil, err
	}

	return resp.Artifact, nil
}

func (op *pushBuildOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	// If we have no registry, we just push the local build.
	if app.Registry == nil {
		return op.Build.Artifact.Artifact, nil
	}

	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		app.Registry,
		app.Registry.PushFunc(),
		argNamedAny("artifact", op.Build.Artifact.Artifact),
	)
}

func (op *pushBuildOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.PushedArtifact).Status)
}

func (op *pushBuildOperation) ValuePtr(msg proto.Message) **any.Any {
	v := msg.(*pb.PushedArtifact)
	if v.Artifact == nil {
		v.Artifact = &pb.Artifact{}
	}

	return &v.Artifact.Artifact
}

var _ operation = (*pushBuildOperation)(nil)
