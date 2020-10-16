package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Release releases a set of deploys.
// TODO(mitchellh): test
func (a *App) Release(ctx context.Context, target *pb.Deployment) (
	*pb.Release,
	component.Release,
	error,
) {
	result, releasepb, err := a.doOperation(ctx, a.logger.Named("release"), &releaseOperation{
		Target: target,
	})
	if err != nil {
		return nil, nil, err
	}

	var release component.Release
	if result != nil {
		release = result.(component.Release)
	}

	return releasepb.(*pb.Release), release, nil
}

type releaseOperation struct {
	Target *pb.Deployment

	result component.Release
}

func (op *releaseOperation) Init(app *App) (proto.Message, error) {
	release := &pb.Release{
		Application:  app.ref,
		Workspace:    app.workspace,
		DeploymentId: op.Target.Id,
		State:        pb.Operation_CREATED,
		Component:    op.Target.Component,
	}

	if app.Releaser != nil {
		release.Component = app.components[app.Releaser].Info
		release.Labels = app.components[app.Releaser].Labels
	}

	if op.result != nil {
		release.Url = op.result.URL()
	}

	return release, nil
}

func (op *releaseOperation) Hooks(app *App) map[string][]*config.Hook {
	if app.Releaser == nil {
		return nil
	}

	return app.components[app.Releaser].Hooks
}

func (op *releaseOperation) Labels(app *App) map[string]string {
	if app.Releaser == nil {
		return nil
	}

	return app.components[app.Releaser].Labels
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
		return nil, err
	}

	return resp.Release, nil
}

func (op *releaseOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	// If we have no releaser, we do nothing since we just update the
	// blank release metadata.
	if app.Releaser == nil {
		return nil, nil
	}

	result, err := app.callDynamicFunc(ctx,
		log,
		(*component.Release)(nil),
		app.Releaser,
		app.Releaser.ReleaseFunc(),
		argNamedAny("target", op.Target.Deployment),
	)
	if err != nil {
		return nil, err
	}

	op.result = result.(component.Release)

	rm := msg.(*pb.Release)
	rm.Url = op.result.URL()

	return result, nil
}

func (op *releaseOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Release).Status)
}

func (op *releaseOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.Release).Release)
}

var _ operation = (*releaseOperation)(nil)
