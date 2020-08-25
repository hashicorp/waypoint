package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// CanDestroyRelease returns true if this app supports destroying releases.
func (a *App) CanDestroyRelease() bool {
	d, ok := a.Releaser.(component.Destroyer)
	return ok && d.DestroyFunc() != nil
}

// DestroyRelease destroyes a specific release.
func (a *App) DestroyRelease(ctx context.Context, d *pb.Release) error {
	// If the release is destroyed already then do nothing.
	if d.State == pb.Operation_DESTROYED {
		a.logger.Info("release already destroyed, doing nothing", "id", d.Id)
		return nil
	}

	_, _, err := a.doOperation(ctx, a.logger.Named("release"), &releaseDestroyOperation{
		Release: d,
	})
	return err
}

type releaseDestroyOperation struct {
	Release *pb.Release
}

func (op *releaseDestroyOperation) Init(app *App) (proto.Message, error) {
	return op.Release, nil
}

func (op *releaseDestroyOperation) Hooks(app *App) map[string][]*config.Hook {
	return nil
}

func (op *releaseDestroyOperation) Labels(app *App) map[string]string {
	return op.Release.Labels
}

func (op *releaseDestroyOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	d := msg.(*pb.Release)
	d.State = pb.Operation_DESTROYED

	resp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: d,
	})
	if err != nil {
		return nil, err
	}

	return resp.Release, nil
}

func (op *releaseDestroyOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	destroyer := app.Releaser.(component.Destroyer)

	return app.callDynamicFunc(ctx,
		log,
		nil,
		destroyer,
		destroyer.DestroyFunc(),
		argNamedAny("release", op.Release.Release),
	)
}

func (op *releaseDestroyOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *releaseDestroyOperation) ValuePtr(msg proto.Message) **any.Any {
	return nil
}

var _ operation = (*releaseDestroyOperation)(nil)
