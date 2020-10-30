package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	c, err := componentCreatorMap[component.ReleaseManagerType].Create(ctx, a, nil)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer c.Close()

	result, releasepb, err := a.doOperation(ctx, a.logger.Named("release"), &releaseOperation{
		Component: c,
		Target:    target,
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
	Component *Component
	Target    *pb.Deployment

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

	if v := op.Component; v != nil {
		release.Component = v.Info
		release.Labels = v.labels
	}

	if op.result != nil {
		release.Url = op.result.URL()
	}

	return release, nil
}

func (op *releaseOperation) Hooks(app *App) map[string][]*config.Hook {
	if op.Component == nil {
		return nil
	}

	return op.Component.hooks
}

func (op *releaseOperation) Labels(app *App) map[string]string {
	if op.Component == nil {
		return nil
	}

	return op.Component.labels
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
	if op.Component == nil {
		return nil, nil
	}

	result, err := app.callDynamicFunc(ctx,
		log,
		(*component.Release)(nil),
		op.Component,
		op.Component.Value.(component.ReleaseManager).ReleaseFunc(),
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
