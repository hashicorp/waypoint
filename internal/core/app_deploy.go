package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
)

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, push *pb.PushedArtifact) (component.Deployment, error) {
	result, _, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployOperation{
		Push: push,
	})
	if err != nil {
		return nil, err
	}

	return result.(component.Deployment), nil
}

type deployOperation struct {
	Push *pb.PushedArtifact
}

func (op *deployOperation) Init(app *App) (proto.Message, error) {
	return &pb.Deployment{
		Component:  app.components[app.Platform],
		ArtifactId: op.Push.Id,
	}, nil
}

func (op *deployOperation) Upsert(
	ctx context.Context,
	client pb.DevflowClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: msg.(*pb.Deployment),
	})
	if err != nil {
		return nil, err
	}

	return resp.Deployment, nil
}

func (op *deployOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	return app.callDynamicFunc(ctx,
		log,
		(*component.Deployment)(nil),
		app.Platform,
		app.Platform.DeployFunc(),
		op.Push.Artifact.Artifact,
	)
}

func (op *deployOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Deployment).Status)
}

func (op *deployOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.Deployment).Deployment)
}

var _ operation = (*deployOperation)(nil)
