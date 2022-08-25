package core

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/opaqueany"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Deploy deploys the given artifact.
// TODO(mitchellh): test
func (a *App) Deploy(ctx context.Context, push *pb.PushedArtifact) (*pb.Deployment, error) {
	// Add our build to our config
	var evalCtx hcl.EvalContext
	evalCtx.Variables = map[string]cty.Value{}
	if err := evalCtxTemplateProto(&evalCtx, "artifact", push); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}
	deployConfig, err := a.deployEvalContext(ctx, &evalCtx)
	if err != nil {
		return nil, err
	}

	op := &deployOperation{
		ComponentFactory: componentCreatorMap[component.PlatformType].Create,
		EvalContext:      &evalCtx,
		Push:             push,
		DeploymentConfig: deployConfig,
	}
	defer op.Close()

	_, msg, err := a.doOperation(ctx, a.logger.Named("deploy"), op)
	if err != nil {
		return nil, err
	}

	result, ok := msg.(*pb.Deployment)
	if !ok {
		return nil, status.Error(codes.Internal, "app_deploy failed to convert the operation message into a Deployment proto")
	}

	return result, nil
}

// deployEvalContext sets the HCL evaluation context for `deploy` blocks.
//
// Note that the eval context set won't be entirely identical to the eval
// context that exists during a real deploy operation. During a real deploy
// operation, the `entrypoint.env` map will contain an auth token that is
// generated just-in-time for the deploy. We omit this for non-deploy operations
// such as destroys or default releasers.
func (a *App) deployEvalContext(
	ctx context.Context,
	evalCtx *hcl.EvalContext,
) (*component.DeploymentConfig, error) {
	if evalCtx.Variables == nil {
		evalCtx.Variables = map[string]cty.Value{}
	}

	// Get the deployment config
	resp, err := a.client.RunnerGetDeploymentConfig(ctx, &pb.RunnerGetDeploymentConfigRequest{})
	if err != nil {
		return nil, err
	}

	// Build our deployment config and expose the env we need to the config
	deployConfig := &component.DeploymentConfig{
		ServerAddr:          resp.ServerAddr,
		ServerTls:           resp.ServerTls,
		ServerTlsSkipVerify: resp.ServerTlsSkipVerify,
	}
	deployEnv := map[string]cty.Value{}
	for k, v := range deployConfig.Env() {
		deployEnv[k] = cty.StringVal(v)
	}
	evalCtx.Variables["entrypoint"] = cty.ObjectVal(map[string]cty.Value{
		"env": cty.MapVal(deployEnv),
	})

	return deployConfig, nil
}

// deployArtifact loads the pushed artifact for a deployment.
func (a *App) deployArtifact(
	ctx context.Context,
	d *pb.Deployment,
) (*pb.PushedArtifact, error) {
	var artifact *pb.PushedArtifact
	if d.Preload != nil && d.Preload.Artifact != nil {
		artifact = d.Preload.Artifact
	}

	if artifact == nil {
		a.logger.Debug("querying artifact", "artifact_id", d.ArtifactId)
		resp, err := a.client.GetPushedArtifact(ctx, &pb.GetPushedArtifactRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: d.ArtifactId,
				},
			},
		})
		if status.Code(err) == codes.NotFound {
			resp = nil
			err = nil
			a.logger.Warn("artifact not found, will attempt destroy regardless",
				"artifact_id", d.ArtifactId)
		}
		if err != nil {
			a.logger.Error("error querying artifact",
				"artifact_id", d.ArtifactId,
				"error", err)
			return nil, err
		}

		artifact = resp
	}

	return artifact, nil
}

type deployOperation struct {
	ComponentFactory func(context.Context, *App, *hcl.EvalContext) (*Component, error)
	EvalContext      *hcl.EvalContext
	Push             *pb.PushedArtifact
	DeploymentConfig *component.DeploymentConfig

	// Set by init
	autoHostname pb.UpsertDeploymentRequest_Tristate

	// component is initialized and recycled at various points in the deployment
	// operation as we have more templating information for the config.
	component *Component

	// id is populated with the deployment id on Upsert
	id string

	// sequence is the monotonically incrementing version number of
	// the deployment
	sequence uint64

	// cebToken is the token to set for the deployment to auth
	cebToken string

	// result is either a component.Deployment or component.DeploymentWithUrl
	result interface{}
}

// Name returns the name of the operation
func (op *deployOperation) Name() string {
	return "deploy"
}

func (op *deployOperation) Close() error {
	if op.component != nil {
		return op.component.Close()
	}

	return nil
}

func (op *deployOperation) reinitComponent(app *App) error {
	// Shut down any previously running plugin if we have one
	if op.component != nil {
		if err := op.component.Close(); err != nil {
			return err
		}
		op.component = nil
	}

	// Recalculate our entrypoint env since this is the one
	// dynamic element that changes while the operation runs.
	deployEnv := map[string]cty.Value{}
	for k, v := range op.DeploymentConfig.Env() {
		deployEnv[k] = cty.StringVal(v)
	}
	op.EvalContext.Variables["entrypoint"] = cty.ObjectVal(map[string]cty.Value{
		"env": cty.MapVal(deployEnv),
	})

	c, err := op.ComponentFactory(context.Background(), app, op.EvalContext)
	if err != nil {
		return err
	}

	op.component = c
	return nil
}

func (op *deployOperation) Init(app *App) (proto.Message, error) {
	// Initialize our component
	if err := op.reinitComponent(app); err != nil {
		return nil, err
	}

	if v := app.config.URL; v != nil {
		if v.AutoHostname != nil {
			if *v.AutoHostname {
				op.autoHostname = pb.UpsertDeploymentRequest_TRUE
			} else {
				op.autoHostname = pb.UpsertDeploymentRequest_FALSE
			}
		}
	}

	// If the deployment plugin supports creating a generation ID, then
	// get that ID up front and set it.
	var generationId string
	if g, ok := op.component.Value.(component.Generation); ok {
		if f := g.GenerationFunc(); f != nil {
			// Get the ID from the plugin.
			idBytesRaw, err := app.callDynamicFunc(context.Background(),
				app.logger,
				nil,
				op.component,
				g.GenerationFunc(),
				op.args()...,
			)
			if err != nil {
				return nil, err
			}
			idBytes := idBytesRaw.([]byte)

			// The plugin can return an empty ID for us to default it to random.
			// If it isn't empty, then we SHA-1 (Version 5 UUID) the bytes
			// to create the actual generation.
			if len(idBytes) > 0 {
				generationId = strings.Replace(
					uuid.NewSHA1(uuid.NameSpaceDNS, idBytes).String(),
					"-", "", -1,
				)
			}
		}
	}

	deployment := &pb.Deployment{
		Generation:  &pb.Generation{Id: generationId},
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   op.component.Info,
		Labels:      op.component.labels,
		ArtifactId:  op.Push.Id,
		State:       pb.Operation_CREATED,
		HasEntrypointConfig: op.DeploymentConfig != nil &&
			op.DeploymentConfig.ServerAddr != "",
	}

	return deployment, nil
}

func (op *deployOperation) Hooks(app *App) map[string][]*config.Hook {
	return op.component.hooks
}

func (op *deployOperation) Labels(app *App) map[string]string {
	return op.component.labels
}

func (op *deployOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment:   msg.(*pb.Deployment),
		AutoHostname: op.autoHostname,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed upserting deployment operation")
	}

	if op.id == "" {
		// Set our internal ID for the Do step
		op.id = resp.Deployment.Id
		op.sequence = resp.Deployment.Sequence

		// We need to get our token that we'll give this deployment
		resp, err := client.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			// TODO: this needs to be configurable. For now we set it
			// to long enough that it should be forever.
			Duration: "87600h", // 10 years

			// This is an entrypoint token specifically for this deployment.
			// NOTE: we purposely use the old "unused" version here until
			// we have an account system.
			UnusedEntrypoint: &pb.Token_Entrypoint{
				DeploymentId: op.id,
			},
		})
		if err != nil {
			return nil, err
		}

		// Set our token up
		op.cebToken = resp.Token
	}

	// Set the new values on our deployment config
	dconfig := *op.DeploymentConfig
	dconfig.Id = op.id
	dconfig.Sequence = op.sequence
	dconfig.EntrypointInviteToken = op.cebToken
	op.DeploymentConfig = &dconfig

	return resp.Deployment, nil
}

func (op *deployOperation) Do(ctx context.Context, log hclog.Logger, app *App, msg proto.Message) (interface{}, error) {
	deploy := msg.(*pb.Deployment)

	// Reinitialize our plugin so that we can render the configuration
	// with the entrypoint token. We need to do this because we have a
	// chicken/egg problem: we need to initialize the config first to
	// get to this point, but we need this env token for the entrypoint
	// to function properly.
	if err := op.reinitComponent(app); err != nil {
		return nil, err
	}

	// Sync our config first
	if err := app.ConfigSync(ctx); err != nil {
		return nil, err
	}

	baseArgs := op.args()

	// Add an outparameter for declared resources, which will be populated by the dynamic func
	declaredResourcesResp := &component.DeclaredResourcesResp{}
	args := append(baseArgs, argmapper.Typed(declaredResourcesResp))

	result, err := app.callDynamicFunc(ctx,
		log,
		(*component.Deployment)(nil),
		op.component,
		op.component.Value.(component.Platform).DeployFunc(),
		args...,
	)

	if err != nil {
		return nil, err
	}

	// Convert from the plugin declaredResources to server declaredResources. Should be identical.
	declaredResources := make([]*pb.DeclaredResource, len(declaredResourcesResp.DeclaredResources))
	for i, pluginDeclaredResource := range declaredResourcesResp.DeclaredResources {
		var serverDeclaredResource pb.DeclaredResource
		if err := mapstructure.Decode(pluginDeclaredResource, &serverDeclaredResource); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode plugin declared resource named %q: %s", pluginDeclaredResource.Name, err)
		}
		declaredResources[i] = &serverDeclaredResource
	}

	deploy.DeclaredResources = declaredResources

	if ep, ok := op.component.Value.(component.Execer); ok && ep.ExecFunc() != nil {
		log.Debug("detected deployment uses an exec plugin, decorating deployment with info")
		deploy.HasExecPlugin = true
	} else {
		log.Debug("no exec plugin detected on platform component")
	}

	if ep, ok := op.component.Value.(component.LogPlatform); ok && ep.LogsFunc() != nil {
		log.Debug("detected deployment uses a logs plugin, decorating deployment with info")
		deploy.HasLogsPlugin = true
	} else {
		log.Debug("no logs plugin detected on platform component")
	}
	return result, err
}

// args returns the args we send to the Deploy function call
func (op *deployOperation) args() []argmapper.Arg {
	var args []argmapper.Arg

	if v := op.Push.Artifact.Artifact; v != nil {
		// This should always be non-nil but for tests we sometimes set this to nil.
		// If we don't do this we will panic so its best to protect against this
		// anyways.
		args = append(args,
			plugin.ArgNamedAny("artifact", op.Push.Artifact.Artifact),
		)
	}

	args = append(args, argmapper.Typed(op.DeploymentConfig))
	return args
}

func (op *deployOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Deployment).Status)
}

func (op *deployOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	return &(msg.(*pb.Deployment).Deployment), &(msg.(*pb.Deployment).DeploymentJson)
}

var _ operation = (*deployOperation)(nil)
