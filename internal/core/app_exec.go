package core

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/ceb"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Exec launches an exec plugin
func (a *App) Exec(ctx context.Context, id string, d *pb.Deployment) error {
	// We need to get the pushed artifact if it isn't loaded.
	// TODO(evanphx) I don't understand why we do this. We're mimicing what the
	// deploy destroy code does and they load this, so we do to?
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
		}

		if err != nil {
			a.logger.Error("error querying artifact",
				"artifact_id", d.ArtifactId,
				"error", err)
			return err
		}

		artifact = resp
	}

	// Add our build to our config
	var evalCtx hcl.EvalContext
	if err := evalCtxTemplateProto(&evalCtx, "artifact", artifact); err != nil {
		a.logger.Warn("failed to prepare template variables, will not be available",
			"err", err)
	}

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, &evalCtx)
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return err
	}
	defer c.Close()

	a.logger.Debug("spooling exec operation")

	_, _, err = a.doOperation(ctx, a.logger.Named("exec"), &execOperation{
		InstanceId: id,
		Component:  c,
		Deployment: d,
		Client:     a.client,
	})
	return err
}

// execOperation is a simple barebones operation value to conform to the
// doOperation doOperation interface.
type execOperation struct {
	Log        hclog.Logger
	InstanceId string
	Component  *Component
	Deployment *pb.Deployment
	Client     pb.WaypointClient
}

func (op *execOperation) Init(app *App) (proto.Message, error) {
	return op.Deployment, nil
}

func (op *execOperation) Hooks(app *App) map[string][]*config.Hook {
	return nil
}

func (op *execOperation) Labels(app *App) map[string]string {
	return op.Deployment.Labels
}

func (op *execOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	return msg, nil
}

// pluginExecVirtHandler is an implementation of ceb.VirtualExecHandler
// that hands off the exec session info to the plugin's Exec function.
type pluginExecVirtHandler struct {
	app    *App
	log    hclog.Logger
	op     *execOperation
	execer component.Execer
	value  *interface{}

	// Set in CreateSession
	info *ceb.VirtualExecInfo

	// Any window size updates that we get from the virtual CEB
	wsUpdates chan component.WindowSize

	// wired up to the context running the CEB to allow us the ability
	// to cancel it from another goroutine
	cancel func()
}

// CreateSession just returns itself because we only use one per virtual
// ceb instance.
func (p *pluginExecVirtHandler) CreateSession(
	ctx context.Context,
	sess *ceb.VirtualExecInfo,
) (ceb.VirtualExecSession, error) {

	p.log.Info("creating plugin virt handler session")

	p.info = sess

	return p, nil
}

// Run translates the session info set in CreateSession into the
// equiv component types and calls the Exec function.
func (p *pluginExecVirtHandler) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	defer func() {
		p.cancel = nil
	}()

	esi := &component.ExecSessionInfo{
		Input:     p.info.Input,
		Output:    p.info.Output,
		Error:     p.info.Error,
		Arguments: p.info.Arguments,
	}

	if p.info.PTY != nil && p.info.PTY.Enable {
		p.wsUpdates = make(chan component.WindowSize, 1)

		esi.WindowSizeUpdates = p.wsUpdates
		esi.IsTTY = true
		esi.Term = p.info.PTY.Term
		esi.InitialWindowSize = component.WindowSize{
			Height: int(p.info.PTY.WindowSize.Height),
			Width:  int(p.info.PTY.WindowSize.Width),
		}
	}

	p.log.Debug("calling plugin with session-info", "arguments", esi.Arguments)

	val, err := p.app.callDynamicFunc(ctx,
		p.log,
		nil,
		p.op.Component,
		p.execer.ExecFunc(),
		argmapper.Named("deployment", p.op.Deployment),
		argmapper.Named("exec_info", esi),
	)
	if err != nil {
		p.log.Error("error executing plugin function", "error", err)
		return err
	}

	// This is to shuffle the actual return value back for app.Exec()
	*p.value = val
	return nil
}

func (p *pluginExecVirtHandler) Close() error {
	if p.cancel == nil {
		return nil
	}

	p.cancel()
	return nil
}

func (p *pluginExecVirtHandler) PTYResize(winch *pb.ExecStreamRequest_WindowSize) error {
	if p.wsUpdates == nil {
		return nil
	}

	p.wsUpdates <- component.WindowSize{
		Height: int(winch.Height),
		Width:  int(winch.Width),
	}

	return nil
}

func (op *execOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	execer, ok := op.Component.Value.(component.Execer)
	if !ok || execer.ExecFunc() == nil {
		log.Debug("component is not an Execer or has no ExecFunc()")
		return nil, nil
	}

	if op.Deployment == nil {
		return nil, fmt.Errorf("no deployment given to exec operation")
	}

	log.Debug("spawn virtual ceb to handle exec")

	virt, err := ceb.NewVirtual(log, ceb.VirtualConfig{
		DeploymentId: op.Deployment.Id,
		InstanceId:   op.InstanceId,
		Client:       op.Client,
	})

	if err != nil {
		return nil, err
	}

	var value interface{}

	return value, virt.RunExec(ctx, &pluginExecVirtHandler{
		app:    app,
		log:    log,
		op:     op,
		execer: execer,
		value:  &value,
	}, 1)
}

func (op *execOperation) StatusPtr(msg proto.Message) **pb.Status {
	return nil
}

func (op *execOperation) ValuePtr(msg proto.Message) **any.Any {
	return nil
}

var _ operation = (*execOperation)(nil)
