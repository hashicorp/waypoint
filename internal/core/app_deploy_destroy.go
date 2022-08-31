package core

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-argmapper"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/opaqueany"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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

// DestroyDeploy destroys a specific deployment.
func (a *App) DestroyDeploy(ctx context.Context, d *pb.Deployment) error {
	return a.destroyDeploy(ctx, d, nil)
}

// destroyAllDeploys will destroy all non-destroyed releases.
func (a *App) destroyAllDeploys(ctx context.Context) error {
	resp, err := a.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
		Application:   a.ref,
		Workspace:     a.workspace,
		PhysicalState: pb.Operation_CREATED,
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Desc:  true,
		},
	})
	if err != nil {
		return nil
	}

	results := resp.Deployments
	if len(results) == 0 {
		return nil
	}

	// current deploy is the latest deploy
	currentDeploy := results[0]

	a.UI.Output(fmt.Sprintf("Destroying deployments for application '%s'...", a.config.Name), terminal.WithHeaderStyle())
	for _, v := range results {
		err := a.destroyDeploy(ctx, v, currentDeploy)
		if err != nil {
			return err
		}
	}

	return nil
}

// destroyDeploy destroys a specific deployment. "d" is the deployment
// to destroy. "configD" is the deployment to use to render the configuration.
// If configD is nil, then "d" is used.
//
// We separate configD and d because when destroying multiple deployments, we
// only have access to the current config, which we can only render using
// the latest deploy typically. This lets callers make that determination.
func (a *App) destroyDeploy(
	ctx context.Context,
	d *pb.Deployment,
	configD *pb.Deployment,
) error {
	// If the deploy is destroyed already then do nothing.
	if d.State == pb.Operation_DESTROYED {
		a.logger.Info("deployment already destroyed, doing nothing", "id", d.Id)
		return nil
	}

	if configD == nil {
		configD = d
	}

	// We need to get the pushed artifact if it isn't loaded.
	artifact, err := a.deployArtifact(ctx, configD)
	if err != nil {
		return err
	}

	// Add our build to our config
	var evalCtx hcl.EvalContext
	if _, err := a.deployEvalContext(ctx, &evalCtx); err != nil {
		a.logger.Warn("failed to prepare entrypoint variables, will not be available",
			"err", err)
	}
	if err := evalCtxTemplateProto(&evalCtx, "artifact", artifact); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, &evalCtx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, destroyment, err := a.doOperation(ctx, a.logger.Named("deploy"), &deployDestroyOperation{
		Component:  c,
		Deployment: d,
	})

	destroyProto, ok := destroyment.(*pb.Deployment)
	if !ok {
		return errors.New("failed to convert destroyment proto to a Deployment")
	}

	var message string
	if len(destroyProto.DeclaredResources) > 0 {
		message = message + fmt.Sprintf("These resources were not destroyed for app %q:\n", a.ref.Application)
		for _, resource := range destroyProto.DeclaredResources {
			message = message + "- " + resource.Name + "\n"
		}
		a.UI.Output(message, terminal.WithWarningStyle())
		message = ""
		if len(destroyProto.DestroyedResources) > 0 {
			message = message + fmt.Sprintf("These resources were destroyed for app %q:\n", a.ref.Application)
			for _, resource := range destroyProto.DestroyedResources {
				message = message + "- " + resource.Name + "\n"
			}
			a.UI.Output(message, terminal.WithSuccessStyle())
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

	// We need to get the pushed artifact if it isn't loaded.
	artifact, err := a.deployArtifact(ctx, results[0])
	if err != nil {
		return err
	}

	// Add our build to our config
	var evalCtx hcl.EvalContext
	if _, err := a.deployEvalContext(ctx, &evalCtx); err != nil {
		a.logger.Warn("failed to prepare entrypoint variables, will not be available",
			"err", err)
	}
	if err := evalCtxTemplateProto(&evalCtx, "artifact", artifact); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, &evalCtx)
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

	a.UI.Output(fmt.Sprintf("Destroying shared deploy resources for application '%s'...", a.config.Name), terminal.WithHeaderStyle())
	_, err = a.callDynamicFunc(ctx,
		log,
		nil,
		c,
		d.DestroyWorkspaceFunc(),
		plugin.ArgNamedAny("deployment", results[0].Deployment),
	)
	return err
}

type deployDestroyOperation struct {
	Component  *Component
	Deployment *pb.Deployment
}

// Name returns the name of the operation
func (op *deployDestroyOperation) Name() string {
	return "deployment destroy"
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
		return nil, errors.Wrapf(err, "failed upserting deployment destroy operation")
	}

	return resp.Deployment, nil
}

func (op *deployDestroyOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	destroy, ok := msg.(*pb.Deployment)
	if !ok {
		return nil, errors.New("failed to cast deploy destroy operation proto to a Deployment")
	}

	destroyer, ok := op.Component.Value.(component.Destroyer)
	if !ok || destroyer.DestroyFunc() == nil {
		return nil, nil
	}

	if op.Deployment.Deployment == nil {
		log.Error("Unable to destroy the Deployment as the proto message Deployment returned from the plugin's DeployFunc is nil. This situation occurs when the deployment process is interrupted by the user.", "deployment", op.Deployment)
		return nil, nil // Fail silently for now, this will be fixed in v0.2
	}

	baseArgs := []argmapper.Arg{plugin.ArgNamedAny("deployment", op.Deployment.Deployment)}
	declaredResourcesResp := &component.DeclaredResourcesResp{}
	destroyedResourcesResp := &component.DestroyedResourcesResp{}
	args := append(baseArgs, argmapper.Typed(declaredResourcesResp), argmapper.Typed(destroyedResourcesResp))

	// We don't need the result, we just need the declared and destroyed resources
	// which we can access without the result since they were passed by reference
	_, err := app.callDynamicFunc(ctx,
		log,
		nil,
		op.Component,
		destroyer.DestroyFunc(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	declaredResources := make([]*pb.DeclaredResource, len(declaredResourcesResp.DeclaredResources))
	if len(declaredResourcesResp.DeclaredResources) > 0 {
		for i, pluginDeclaredResource := range declaredResourcesResp.DeclaredResources {
			var serverDeclaredResource pb.DeclaredResource
			if err := mapstructure.Decode(pluginDeclaredResource, &serverDeclaredResource); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to decode plugin declared resource named %q: %s", pluginDeclaredResource.Name, err)
			}
			declaredResources[i] = &serverDeclaredResource
		}
	}
	destroy.DeclaredResources = declaredResources

	destroyedResources := make([]*pb.DestroyedResource, len(destroyedResourcesResp.DestroyedResources))
	if len(destroyedResourcesResp.DestroyedResources) > 0 {
		for i, pluginDestroyedResource := range destroyedResourcesResp.DestroyedResources {
			var serverDestroyedResource pb.DestroyedResource
			if err := mapstructure.Decode(pluginDestroyedResource, &serverDestroyedResource); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to decode plugin declared resource named %q: %s", pluginDestroyedResource.Name, err)
			}
			destroyedResources[i] = &serverDestroyedResource
		}
	}
	destroy.DestroyedResources = destroyedResources

	return nil, err
}

func (op *deployDestroyOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *deployDestroyOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	return nil, nil
}

var _ operation = (*deployDestroyOperation)(nil)
