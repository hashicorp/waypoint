package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/ceb"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Exec launches an exec plugin. Exec plugins are only used if the plugin's
// platforms plugin wishes to implement the ExecFunc protocol. And even then, we
// only trigger this code path if there are no long running instances associated
// with the given Deployment.
// Under traditional platform scenarios, we don't need to run a exec plugin, instead
// the exec command can connect directly to a long running instance to provide the
// exec session.
// The result of running this task is that the platform plugin is called
// and made available as a virtual instance with the given id.
// enableDynConfig controls if exec jobs will attempt to read from any dynamic config sources.
// Reading from those sources requires the runner to have credentials to those sources.
func (a *App) Exec(ctx context.Context, id string, d *pb.Deployment, enableDynConfig bool) error {
	// We need to get the pushed artifact if it isn't loaded.
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

	execer, ok := c.Value.(component.Execer)
	if !ok || execer.ExecFunc() == nil {
		a.logger.Debug("component is not an Execer or has no ExecFunc()")
		return nil
	}

	a.logger.Debug("spawn virtual ceb to handle exec")

	virt, err := ceb.NewVirtual(a.logger, ceb.VirtualConfig{
		DeploymentId:        d.Id,
		InstanceId:          id,
		Client:              a.client,
		EnableDynamicConfig: enableDynConfig,
	})

	if err != nil {
		return err
	}

	return virt.RunExec(ctx, &pluginExecVirtHandler{
		app:        a,
		log:        a.logger,
		component:  c,
		deployment: d,
		execer:     execer,
		artifact:   artifact,
	}, 1)
}

// pluginExecVirtHandler is an implementation of ceb.VirtualExecHandler
// that hands off the exec session info to the plugin's Exec function.
type pluginExecVirtHandler struct {
	app        *App
	log        hclog.Logger
	component  *Component
	deployment *pb.Deployment
	execer     component.Execer
	artifact   *pb.PushedArtifact

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
	info *ceb.VirtualExecInfo,
) (ceb.VirtualExecSession, error) {

	p.log.Info("creating plugin virt handler session")

	p.info = info

	return p, nil
}

type exitStatusError struct {
	code int
}

func (e *exitStatusError) Error() string {
	return fmt.Sprintf("command exited with status %d", e.code)
}

func (e *exitStatusError) ExitStatus() int {
	return e.code
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
		Input:       p.info.Input,
		Output:      p.info.Output,
		Error:       p.info.Error,
		Arguments:   p.info.Arguments,
		Environment: p.info.Environment,
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

	result, err := p.app.callDynamicFunc(ctx,
		p.log,
		nil,
		p.component,
		p.execer.ExecFunc(),
		argNamedAny("deployment", p.deployment.Deployment),
		argNamedAny("image", p.artifact.Artifact.Artifact),
		argmapper.Typed(esi),
	)
	if err != nil {
		p.log.Error("error executing plugin function", "error", err)
		return err
	}

	p.log.Info("exec finished", "result", hclog.Fmt("%#v", result))

	if ec, ok := result.(*component.ExecResult); ok {
		if ec.ExitCode != 0 {
			return &exitStatusError{ec.ExitCode}
		}
	}

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
