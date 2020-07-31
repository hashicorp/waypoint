package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, push *pb.PushedArtifact) (*pb.Deployment, error) {
	// Get the deployment config
	resp, err := a.client.RunnerGetDeploymentConfig(ctx, &pb.RunnerGetDeploymentConfigRequest{})
	if err != nil {
		return nil, err
	}

	_, msg, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployOperation{
		Push: push,
		DeploymentConfig: &component.DeploymentConfig{
			ServerAddr:     resp.ServerAddr,
			ServerInsecure: resp.ServerInsecure,
		},
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.Deployment), nil
}

type deployOperation struct {
	Push             *pb.PushedArtifact
	DeploymentConfig *component.DeploymentConfig

	// Set by init
	autoHostname pb.UpsertDeploymentRequest_Tristate

	// id is populated with the deployment id on Upsert
	id string
}

func (op *deployOperation) Init(app *App) (proto.Message, error) {
	if app.components[app.Platform] == nil {
		return nil, status.Error(codes.NotFound, "no deployment configured")
	}

	if v := app.config.URL; v != nil {
		if v.AutoHostname != nil {
			if *v.AutoHostname {
				op.autoHostname = pb.UpsertDeploymentRequest_TRUE
			} else {
				op.autoHostname = pb.UpsertDeploymentRequest_FALSE
			}
		}
	}

	return &pb.Deployment{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   app.components[app.Platform].Info,
		Labels:      app.components[app.Platform].Labels,
		ArtifactId:  op.Push.Id,
		State:       pb.Deployment_DEPLOYED,
	}, nil
}

func (op *deployOperation) Hooks(app *App) map[string][]*config.Hook {
	return app.components[app.Platform].Hooks
}

func (op *deployOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment:   msg.(*pb.Deployment),
		AutoHostname: op.autoHostname,
	})
	if err != nil {
		return nil, err
	}

	// Set our internal ID for the Do step
	op.id = resp.Deployment.Id

	return resp.Deployment, nil
}

func (op *deployOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
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
