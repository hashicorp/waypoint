package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// CanDestroyDeploy returns true if this app supports destroying deployments.
func (a *App) CanDestroyDeploy() bool {
	_, ok := a.Platform.(component.Destroyer)
	return ok
}

// DestroyDeploy destroyes a specific deployment.
func (a *App) DestroyDeploy(ctx context.Context, d *pb.Deployment) error {
	_, _, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployDestroyOperation{
		Deployment: d,
	})
	return err
}

type deployDestroyOperation struct {
	Deployment *pb.Deployment
}

func (op *deployDestroyOperation) Init(app *App) (proto.Message, error) {
	return op.Deployment, nil
}

func (op *deployDestroyOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	// We don't interact with the server
	return op.Deployment, nil
}

func (op *deployDestroyOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	destroyer := app.Platform.(component.Destroyer)

	return app.callDynamicFunc(ctx,
		log,
		nil,
		destroyer,
		destroyer.DestroyFunc(),
		op.Deployment.Deployment,
	)
}

func (op *deployDestroyOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *deployDestroyOperation) ValuePtr(msg proto.Message) **any.Any {
	return nil
}

var _ operation = (*deployDestroyOperation)(nil)
