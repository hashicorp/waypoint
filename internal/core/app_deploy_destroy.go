package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// CanDestroyDeploy returns true if this app supports destroying deployments.
func (a *App) CanDestroyDeploy() bool {
	_, ok := a.Platform.(component.Destroyer)
	return ok
}

// DestroyDeploy destroyes a specific deployment.
func (a *App) DestroyDeploy(ctx context.Context, d *pb.Deployment) error {
	// If the deploy is destroyed already then do nothing.
	if d.State == pb.Operation_DESTROYED {
		a.logger.Info("deployment already destroyed, doing nothing", "id", d.Id)
		return nil
	}

	_, _, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployDestroyOperation{
		Deployment: d,
	})
	return err
}

// destroyAllDeploys will destroy all non-destroyed releases.
func (a *App) destroyAllDeploys(ctx context.Context) error {
	resp, err := a.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
		Application:   a.ref,
		Workspace:     a.workspace,
		PhysicalState: pb.Operation_CREATED,
	})
	if err != nil {
		return nil
	}

	results := resp.Deployments
	if len(results) == 0 {
		return nil
	}

	if a.Platform == nil {
		return status.Errorf(codes.FailedPrecondition,
			"Created deployments must be destroyed but no deployment plugin is configured! "+
				"Please configure a deployment plugin in your Waypoint configuration.")
	}

	a.UI.Output("Destroying deployments...", terminal.WithHeaderStyle())
	for _, v := range results {
		err := a.DestroyDeploy(ctx, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// destroyDeployWorkspace will call the DestroyWorkspace hook if there
// are any valid operations. This expects all operations of this type to
// already be destroyed and will error if they are not.
func (a *App) destroyDeployWorkspace(ctx context.Context) error {
	log := a.logger

	// Get the last destroyed value.
	resp, err := a.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
		Application:   a.ref,
		Workspace:     a.workspace,
		PhysicalState: pb.Operation_DESTROYED,
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Limit: 1,
		},
	})
	if err != nil {
		return nil
	}

	// If we have no opeartions, we don't call the hook.
	results := resp.Deployments
	if len(results) == 0 {
		return nil
	}

	// Call the hook
	d, ok := a.Platform.(component.WorkspaceDestroyer)
	if !ok || d.DestroyWorkspaceFunc() == nil {
		return status.Errorf(codes.FailedPrecondition,
			"Created deployments must be destroyed but no deployment plugin is configured! "+
				"Please configure a deployment plugin in your Waypoint configuration.")
	}

	a.UI.Output("Destroying shared deploy resources...", terminal.WithHeaderStyle())
	_, err = a.callDynamicFunc(ctx,
		log,
		nil,
		d,
		d.DestroyWorkspaceFunc(),
		argNamedAny("deployment", results[0].Deployment),
	)
	return err
}

type deployDestroyOperation struct {
	Deployment *pb.Deployment
}

func (op *deployDestroyOperation) Init(app *App) (proto.Message, error) {
	return op.Deployment, nil
}

func (op *deployDestroyOperation) Hooks(app *App) map[string][]*config.Hook {
	return nil
}

func (op *deployDestroyOperation) Labels(app *App) map[string]string {
	return op.Deployment.Labels
}

func (op *deployDestroyOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	d := msg.(*pb.Deployment)
	d.State = pb.Operation_DESTROYED

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: d,
	})
	if err != nil {
		return nil, err
	}

	return resp.Deployment, nil
}

func (op *deployDestroyOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	destroyer := app.Platform.(component.Destroyer)

	return app.callDynamicFunc(ctx,
		log,
		nil,
		destroyer,
		destroyer.DestroyFunc(),
		argNamedAny("deployment", op.Deployment.Deployment),
	)
}

func (op *deployDestroyOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *deployDestroyOperation) ValuePtr(msg proto.Message) **any.Any {
	return nil
}

var _ operation = (*deployDestroyOperation)(nil)
