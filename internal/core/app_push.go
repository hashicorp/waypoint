package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Push pushes the given build to a registry.
// TODO(mitchellh): test
func (a *App) PushBuild(ctx context.Context, build *pb.Build) (*pb.PushedArtifact, error) {
	_, msg, err := a.doOperation(ctx, a.logger.Named("push"), &pushBuildOperation{
		Build: build,
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.PushedArtifact), nil
}

type pushBuildOperation struct {
	Build *pb.Build
}

func (op *pushBuildOperation) Init(app *App) (proto.Message, error) {
	return &pb.PushedArtifact{
		Component: app.components[app.Registry],
		BuildId:   op.Build.Id,
	}, nil
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

func (op *pushBuildOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		app.Registry,
		app.Registry.PushFunc(),
		op.Build.Artifact.Artifact,
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
