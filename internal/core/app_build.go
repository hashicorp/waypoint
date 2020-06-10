package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Build builds the artifact from source for this app.
// TODO(mitchellh): test
func (a *App) Build(ctx context.Context) (*pb.Build, error) {
	_, msg, err := a.doOperation(ctx, a.logger.Named("build"), &buildOperation{})
	if err != nil {
		return nil, err
	}

	return msg.(*pb.Build), nil
}

type buildOperation struct {
	Build *pb.Build
}

func (op *buildOperation) Init(app *App) (proto.Message, error) {
	return &pb.Build{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   app.components[app.Builder].Info,
		Labels:      app.components[app.Builder].Labels,
	}, nil
}

func (op *buildOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertBuild(ctx, &pb.UpsertBuildRequest{
		Build: msg.(*pb.Build),
	})
	if err != nil {
		return nil, err
	}

	return resp.Build, nil
}

func (op *buildOperation) Do(ctx context.Context, log hclog.Logger, app *App) (interface{}, error) {
	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		app.Builder,
		app.Builder.BuildFunc(),
	)
}

func (op *buildOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Build).Status)
}

func (op *buildOperation) ValuePtr(msg proto.Message) **any.Any {
	v := msg.(*pb.Build)
	if v.Artifact == nil {
		v.Artifact = &pb.Artifact{}
	}

	return &v.Artifact.Artifact
}

var _ operation = (*buildOperation)(nil)
