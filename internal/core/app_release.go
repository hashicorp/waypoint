package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Release releases a set of deploys.
// TODO(mitchellh): test
func (a *App) Release(ctx context.Context, targets []component.ReleaseTarget) (component.Release, error) {
	result, _, err := a.doOperation(ctx, a.logger.Named("release"), &releaseOperation{
		Targets: targets,
	})
	if err != nil {
		return nil, err
	}

	return result.(component.Release), nil
}

type releaseOperation struct {
	Targets []component.ReleaseTarget
}

func (op *releaseOperation) Init(app *App) (proto.Message, error) {
	release := &pb.Release{
		Component:    app.components[app.Releaser].Info,
		Labels:       app.components[app.Releaser].Labels,
		TrafficSplit: &pb.Release_Split{},
	}

	// Create our splits for the release
	for _, target := range op.Targets {
		release.TrafficSplit.Targets = append(release.TrafficSplit.Targets, &pb.Release_SplitTarget{
			DeploymentId: target.DeploymentId,
			Percent:      int32(target.Percent),
		})
	}

	return release, nil
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

func (op *releaseOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	return app.callDynamicFunc(ctx,
		log,
		(*component.Release)(nil),
		app.Releaser,
		app.Releaser.ReleaseFunc(),
		op.Targets,
	)
}

func (op *releaseOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Release).Status)
}

func (op *releaseOperation) ValuePtr(msg proto.Message) **any.Any {
	return &(msg.(*pb.Release).Release)
}

var _ operation = (*releaseOperation)(nil)
