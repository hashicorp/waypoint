package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config2"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, push *pb.PushedArtifact) (*pb.Deployment, error) {
	// Get the deployment config
	resp, err := a.client.RunnerGetDeploymentConfig(ctx, &pb.RunnerGetDeploymentConfigRequest{})
	if err != nil {
		return nil, err
	}

	// Render the config
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	_, msg, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployOperation{
		Component: c,
		Push:      push,
		DeploymentConfig: &component.DeploymentConfig{
			ServerAddr:          resp.ServerAddr,
			ServerTls:           resp.ServerTls,
			ServerTlsSkipVerify: resp.ServerTlsSkipVerify,
		},
	})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.Deployment), nil
}

type deployOperation struct {
	Component        *Component
	Push             *pb.PushedArtifact
	DeploymentConfig *component.DeploymentConfig

	// Set by init
	autoHostname pb.UpsertDeploymentRequest_Tristate

	// id is populated with the deployment id on Upsert
	id string

	// cebToken is the token to set for the deployment to auth
	cebToken string
}

func (op *deployOperation) Init(app *App) (proto.Message, error) {
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
		Component:   op.Component.Info,
		Labels:      op.Component.labels,
		ArtifactId:  op.Push.Id,
		State:       pb.Operation_CREATED,
		HasEntrypointConfig: op.DeploymentConfig != nil &&
			op.DeploymentConfig.ServerAddr != "",
	}, nil
}

func (op *deployOperation) Hooks(app *App) map[string][]*config.Hook {
	return op.Component.hooks
}

func (op *deployOperation) Labels(app *App) map[string]string {
	return op.Component.labels
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

	if op.id == "" {
		// Set our internal ID for the Do step
		op.id = resp.Deployment.Id

		// We need to get our token that we'll give this deployment
		resp, err := client.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			// TODO: this needs to be configurable. For now we set it
			// to long enough that it should be forever.
			Duration: "87600h", // 10 years

			// This is an entrypoint token specifically for this deployment
			Entrypoint: &pb.Token_Entrypoint{
				DeploymentId: op.id,
			},
		})
		if err != nil {
			return nil, err
		}

		// Set our token up
		op.cebToken = resp.Token
	}

	return resp.Deployment, nil
}

func (op *deployOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	dconfig := *op.DeploymentConfig
	dconfig.Id = op.id
	dconfig.EntrypointInviteToken = op.cebToken

	return app.callDynamicFunc(ctx,
		log,
		(*component.Deployment)(nil),
		op.Component,
		op.Component.Value.(component.Platform).DeployFunc(),
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
