package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/opaqueany"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// CanDestroyRelease returns true if this app supports destroying releases.
func (a *App) CanDestroyRelease() bool {
	c, err := componentCreatorMap[component.ReleaseManagerType].Create(context.Background(), a, nil)
	if status.Code(err) == codes.Unimplemented {
		// The if statement below catches this too but I want to just explicitly
		// state that we want to ensure this is false in case we ever refactor
		// the below to return an error.
		return false
	}
	if err != nil {
		return false
	}
	defer c.Close()

	d, ok := c.Value.(component.Destroyer)
	return ok && d.DestroyFunc() != nil
}

// DestroyRelease destroys a specific release.
func (a *App) DestroyRelease(ctx context.Context, d *pb.Release) error {
	// If the release is destroyed already then do nothing.
	if d.State == pb.Operation_DESTROYED {
		a.logger.Info("release already destroyed, doing nothing", "id", d.Id)
		return nil
	}

	// Setup our context
	var evalCtx hcl.EvalContext
	if err := a.releaserEvalContext(ctx, d, &evalCtx); err != nil {
		return err
	}

	c, err := a.createReleaser(ctx, &evalCtx)
	if status.Code(err) == codes.Unimplemented {
		c = nil
		err = nil
	}
	if err != nil {
		return err
	}
	defer c.Close()

	_, _, err = a.doOperation(ctx, a.logger.Named("release"), &releaseDestroyOperation{
		Component: c,
		Release:   d,
	})
	return err
}

// releaserEvalContext populates the typical HCL context for rendering
// the releaser.
func (a *App) releaserEvalContext(
	ctx context.Context,
	r *pb.Release,
	evalCtx *hcl.EvalContext,
) error {
	// Query the deployment
	a.logger.Debug("querying deployment", "deployment_id", r.DeploymentId)
	resp, err := a.client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: r.DeploymentId,
			},
		},

		LoadDetails: pb.Deployment_ARTIFACT,
	})
	if status.Code(err) == codes.NotFound {
		resp = nil
		err = nil
		a.logger.Warn("deployment not found, will attempt destroy regardless",
			"deployment_id", r.DeploymentId)
	}
	if err != nil {
		a.logger.Error("error querying deployment",
			"deployment_id", r.DeploymentId,
			"error", err)
		return err
	}
	deploy := resp

	// Add our build to our config
	if deploy != nil {
		if err := evalCtxTemplateProto(evalCtx, "artifact", deploy.Preload.Artifact); err != nil {
			a.logger.Warn("failed to prepare template variables, will not be available",
				"err", err)
		}
		if err := evalCtxTemplateProto(evalCtx, "deploy", deploy); err != nil {
			a.logger.Warn("failed to prepare template variables, will not be available",
				"err", err)
		}
	}

	return nil
}

// destroyAllReleases will destroy all non-destroyed releases.
func (a *App) destroyAllReleases(ctx context.Context) error {
	resp, err := a.client.ListReleases(ctx, &pb.ListReleasesRequest{
		Application:   a.ref,
		Workspace:     a.workspace,
		PhysicalState: pb.Operation_CREATED,
	})
	if err != nil {
		return nil
	}

	rels := resp.Releases
	if len(rels) == 0 {
		return nil
	}

	a.UI.Output(fmt.Sprintf("Destroying releases for application '%s'...", a.config.Name), terminal.WithHeaderStyle())
	for _, rel := range rels {
		err := a.DestroyRelease(ctx, rel)
		if err != nil {
			return err
		}
	}

	return nil
}

// destroyReleaseWorkspace will call the DestroyWorkspace hook if there
// are any valid operations. This expects all operations of this type to
// already be destroyed and will error if they are not.
func (a *App) destroyReleaseWorkspace(ctx context.Context) error {
	log := a.logger

	// Get the last destroyed value.
	resp, err := a.client.ListReleases(ctx, &pb.ListReleasesRequest{
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
	results := resp.Releases
	if len(results) == 0 {
		return nil
	}

	// Setup our context
	var evalCtx hcl.EvalContext
	if err := a.releaserEvalContext(ctx, results[0], &evalCtx); err != nil {
		return err
	}

	// Start the plugin
	c, err := a.createReleaser(ctx, &evalCtx)
	if status.Code(err) == codes.Unimplemented {
		return nil
	}
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

	a.UI.Output(fmt.Sprintf("Destroying shared release resources for application '%s'...", a.config.Name), terminal.WithHeaderStyle())
	_, err = a.callDynamicFunc(ctx,
		log,
		nil,
		c,
		d.DestroyWorkspaceFunc(),
		plugin.ArgNamedAny("release", results[0].Release),
	)
	return err
}

type releaseDestroyOperation struct {
	Component *Component
	Release   *pb.Release
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
		return nil, errors.Wrapf(err, "failed upserting release destroy operation")
	}

	return resp.Release, nil
}

// Name returns the name of the operation
func (op *releaseDestroyOperation) Name() string {
	return "release destroy"
}

func (op *releaseDestroyOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	// If we have no releaser then we're done.
	if op.Component == nil {
		return nil, nil
	}

	destroy, ok := msg.(*pb.Release)
	if !ok {
		return nil, errors.New("failed to cast release destroy operation proto to a Release")
	}

	// If we don't implement the destroy plugin we just mark it as destroyed.
	destroyer, ok := op.Component.Value.(component.Destroyer)
	if !ok || destroyer.DestroyFunc() == nil {
		return nil, nil
	}

	if op.Release.Release == nil {
		log.Error("Unable to destroy the Release as the proto message Release returned from the plugin's ReleaseFunc is nil. This situation occurs when the release process is interupted by the user.", "release", op.Release)
		return nil, nil // Fail silently for now, this will be fixed in v0.2
	}

	declaredResourcesResp := &component.DeclaredResourcesResp{}
	destroyedResourcesResp := &component.DestroyedResourcesResp{}

	// We don't need the result, we just need the declared and destroyed resources
	// which we can access without the result since they were passed by reference
	_, err := app.callDynamicFunc(ctx,
		log,
		nil,
		op.Component,
		destroyer.DestroyFunc(),
		plugin.ArgNamedAny("release", op.Release.Release),
		argmapper.Typed(declaredResourcesResp),
		argmapper.Typed(destroyedResourcesResp),
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

func (op *releaseDestroyOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *releaseDestroyOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	return nil, nil
}

var _ operation = (*releaseDestroyOperation)(nil)
