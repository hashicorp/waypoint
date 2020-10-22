package core

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/config2"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// CanDestroyDeploy returns true if this app supports destroying deployments.
func (a *App) CanDestroyDeploy() bool {
	c, err := componentCreatorMap[component.PlatformType].Create(context.Background(), a, nil)
	if err != nil {
		return false
	}
	defer c.Close()

	_, ok := c.Value.(component.Destroyer)
	return ok
}

// DestroyDeploy destroyes a specific deployment.
func (a *App) DestroyDeploy(ctx context.Context, d *pb.Deployment) error {
	// If the deploy is destroyed already then do nothing.
	if d.State == pb.Operation_DESTROYED {
		a.logger.Info("deployment already destroyed, doing nothing", "id", d.Id)
		return nil
	}

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	_, _, err = a.doOperation(ctx, a.logger.Named("deploy"), &deployDestroyOperation{
		Component:  c,
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

	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, nil)
	if status.Code(err) == codes.Unimplemented {
		return status.Errorf(codes.FailedPrecondition,
			"Created deployments must be destroyed but no deployment plugin is configured! "+
				"Please configure a deployment plugin in your Waypoint configuration.")
	}
	if err != nil {
		return err
	}
	defer c.Close()

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

	// If we have no operations, we don't call the hook.
	results := resp.Deployments
	if len(results) == 0 {
		return nil
	}

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	// Call the hook
	d, ok := c.Value.(component.WorkspaceDestroyer)
	if !ok || d.DestroyWorkspaceFunc() == nil {
		// Workspace deletion is optional.
		return nil
	}

	a.UI.Output("Destroying shared deploy resources...", terminal.WithHeaderStyle())
	_, err = a.callDynamicFunc(ctx,
		log,
		nil,
		c,
		d.DestroyWorkspaceFunc(),
		argNamedAny("deployment", results[0].Deployment),
	)
	return err
}

type deployDestroyOperation struct {
	Component  *Component
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
	destroyer, ok := op.Component.Value.(component.Destroyer)
	if !ok || destroyer.DestroyFunc() == nil {
		return nil, nil
	}

	if op.Deployment.Deployment == nil {
		log.Error("Unable to destroy the Deployment as the proto message Deployment returned from the plugin's DeployFunc is nil. This situation occurs when the deployment process is interupted by the user.", "deployment", op.Deployment)
		return nil, nil // Fail silently for now, this will be fixed in v0.2
	}

	return app.callDynamicFunc(ctx,
		log,
		nil,
		op.Component,
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
