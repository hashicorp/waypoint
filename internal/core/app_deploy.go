package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, push *pb.PushedArtifact) (*pb.Deployment, error) {
	_, msg, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployOperation{
		Push:             push,
		DeploymentConfig: &a.dconfig,
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.Deployment), nil
}

type deployOperation struct {
	Push             *pb.PushedArtifact
	DeploymentConfig *component.DeploymentConfig

	// id is populated with the deployment id on Upsert
	id string
}

func (op *deployOperation) Init(app *App) (proto.Message, error) {
	return &pb.Deployment{
		Component:  app.components[app.Platform].Info,
		Labels:     app.components[app.Platform].Labels,
		ArtifactId: op.Push.Id,
		State:      pb.Deployment_DEPLOY,
	}, nil
}

func (op *deployOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: msg.(*pb.Deployment),
	})
	if err != nil {
		return nil, err
	}

	// Set our internal ID for the Do step
	op.id = resp.Deployment.Id

	return resp.Deployment, nil
}

func (op *deployOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	dconfig := *op.DeploymentConfig
	dconfig.Id = op.id

	return app.callDynamicFunc(ctx,
		log,
		(*component.Deployment)(nil),
		app.Platform,
		app.Platform.DeployFunc(),
		argNamedAny("artifact", op.Push.Artifact.Artifact),
		argmapper.Typed(&dconfig),
	)
}

func (op *deployOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Deployment).Status)
}

func (op *deployOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.Deployment).Deployment)
}

var _ operation = (*deployOperation)(nil)
