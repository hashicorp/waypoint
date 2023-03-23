// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package virtualceb

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"

	"github.com/hashicorp/waypoint/internal/appconfig"
	"github.com/hashicorp/waypoint/internal/ceb/execwriter"
	"github.com/hashicorp/waypoint/internal/plugin"
)

// ExecInfo contains values to run an exec session.
type ExecInfo struct {
	Input  io.Reader // stdin
	Output io.Writer // stdout
	Error  io.Writer // stderr

	// Command line arguments
	Arguments []string

	// The environment variables to set in the exec context
	Environment []string

	// Specifies if we and how we should allocate a pty to handle
	// the command.
	PTY *pb.ExecStreamRequest_PTY
}

// ExecSession represents a running exec, spawned by ExecHandler.
type ExecSession interface {
	// Called to start the session. Should block until the session is finished.
	Run(ctx context.Context) error

	// Close the running session down. Called concurrently to Run.
	Close() error

	// Resize the PTY according to the given window information.
	// Called concurrently to Run.
	PTYResize(*pb.ExecStreamRequest_WindowSize) error
}

// ExecHandler represents the ability to spawn exec sessions.
// It is the abstraction layer Virtual uses for creating exec sessions.
type ExecHandler interface {
	CreateSession(ctx context.Context, sess *ExecInfo) (ExecSession, error)
}

// Config is the configuration of the CEB Virtual value
type Config struct {
	// The deployment id that this virtual session is for. The server
	// will validate this value.
	DeploymentId string

	// The instance id for this virtual instance.
	InstanceId string

	// How to connect back to the server. Because Virtual is usually used in the context
	// of a Runner, this can be the same Client the Runner is using.
	Client pb.WaypointClient

	// Support Dynamic Config
	EnableDynamicConfig bool
}

// Virtual represents a virtual CEB instance. It is used to manifest an instance that
// performs exec operations via a Go interface rather than just running a command.
type Virtual struct {
	cfg Config
	log hclog.Logger
}

// New creates a new Virtual value based on the logger and config.
func New(log hclog.Logger, cfg Config) (*Virtual, error) {
	virt := &Virtual{
		cfg: cfg,
		log: log,
	}
	return virt, nil
}

// RunExec connects to the server and handles any inbound Exec requests via the
// ExecHandler. The count parameter indicates how many exec sessions to handle
// before returning. If count is less than 0, it handles sessions forever.
func (v *Virtual) RunExec(ctx context.Context, h ExecHandler, count int) error {
	v.log.Info("connecting virtual instance")

	// A quick heads up: gRPC provides to ability to let the client of a recieve stream
	// tell the remote side "hey, I'm done now, don't send me anything else.". So instead
	// we wire up a context to the call and cancel it when we are done. This cancelation
	// will propogate to the server and they'll see that we have gone away.
	//
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// On the server side, EntrypointConfig is what registers an instance. Keeping
	// the epclient alive is what controls if the instance is registered or not.
	//
	epclient, err := v.cfg.Client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: v.cfg.DeploymentId,
		InstanceId:   v.cfg.InstanceId,
		Type:         pb.Instance_VIRTUAL,
	})
	if err != nil {
		return err
	}

	// We never send anything
	epclient.CloseSend()

	v.log.Info("virtual instance connected")

	var highestExec int64

	// They can be used for config sources that we might be sent.
	configPlugins := plugin.ConfigSourcers

	// Setup our config watcher
	w, err := appconfig.NewWatcher(
		appconfig.WithLogger(v.log),
		appconfig.WithPlugins(configPlugins),
		appconfig.WithDynamicEnabled(v.cfg.EnableDynamicConfig),
	)
	if err != nil {
		return err
	}
	defer w.Close()

	// This is much more paired down than the version in the official CEB because the
	// expectation is that a virtual instance is used for a single operation and then
	// exits. So we only need to see a single view of the config variables before we
	// can continue on.

	for {
		msg, err := epclient.Recv()
		if err != nil {
			return err
		}

		// Always update our config vars even if we have no exec so that
		// while we're waiting for exec we're at least gathering config.
		w.UpdateSources(ctx, msg.Config.ConfigSources)
		w.UpdateVars(ctx, msg.Config.EnvVars)

		if msg.Config.Exec == nil {
			continue
		}

		var env []string
		if len(msg.Config.EnvVars) > 0 {
			// Get the config values, we use a iterator value of 0 here so that
			// we will wait for at least one set of config values to resolve.
			// In most virtual CEB cases we are only open for a single exec
			// anyways so this resolves perfectly to our expected values.
			uenv, _, err := w.Next(ctx, 0)
			if err != nil {
				// we drop the error here (only log it don't return) because
				// that is what we did prior to this change too
				v.log.Warn("error retrieving config values", "err", err)
			}

			env = uenv.EnvVars
		}

		idx := highestExec
		for _, exec := range msg.Config.Exec {
			// Skip sessions we already know about. Normal CEB does this too, I guess beacuse
			// the server can resend exec info.
			if exec.Index <= highestExec {
				continue
			}

			if exec.Index > idx {
				idx = exec.Index
			}

			if err := v.startExec(ctx, h, exec, env); err != nil {
				return err
			}
			if count > 0 {
				count--
				if count == 0 {
					v.log.Info("virtual instance stopping")
					return nil
				}
			}
		}

		highestExec = idx
	}
}

// startExec launches and manages a ExecStream for the given exec config. It will
// spawn an exec session from the handler and then shuffle the data between gRPC
// and the ExecSession and ExecInfo interfaces.
func (v *Virtual) startExec(
	ctx context.Context,
	h ExecHandler,
	execConfig *pb.EntrypointConfig_Exec,
	env []string,
) error {
	v.log.Info("starting exec stream", "args", execConfig.Args)
	defer v.log.Info("exec stream finished")

	client, err := v.cfg.Client.EntrypointExecStream(ctx)
	if err != nil {
		v.log.Warn("error opening exec stream", "err", err)
		return nil
	}

	defer client.CloseSend()

	// Send our open message
	v.log.Trace("sending open message")
	if err := client.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: v.cfg.InstanceId,
				Index:      execConfig.Index,
			},
		},
	}); err != nil {
		v.log.Warn("error opening exec stream", "err", err)
		return nil
	}

	// Create our pipe for stdin so that we can send data
	stdinR, stdinW := io.Pipe()
	defer stdinW.Close()

	// We need to modify our command so the input/output is all over gRPC
	stdout := execwriter.Writer(client, pb.EntrypointExecRequest_Output_STDOUT)
	stderr := execwriter.Writer(client, pb.EntrypointExecRequest_Output_STDERR)

	done := make(chan error, 1)

	// Spawn a new exec session for this exec config.
	xsess, err := h.CreateSession(ctx, &ExecInfo{
		Input:       stdinR,
		Output:      stdout,
		Error:       stderr,
		Arguments:   execConfig.Args,
		Environment: env,
		PTY:         execConfig.Pty,
	})
	if err != nil {
		return err
	}

	// Start our receive data loop. We use this loop so that our
	// main processing loop can select on multiple channels.

	respCh := make(chan *pb.EntrypointExecResponse)
	go func() {
		defer close(respCh)

		for {
			resp, err := client.Recv()
			if err != nil {
				if err == io.EOF || err == context.Canceled || status.Code(err) == codes.Canceled {
					v.log.Info("exec stream ended by client")
				} else {
					v.log.Warn("error receiving from server stream", "err", err)
				}
				return
			}

			respCh <- resp
		}
	}()

	// We don't block on Run in the main goroutine so that we can just shuffle
	// and orchestrate the data there.
	go func() {
		defer v.log.Debug("virtual ceb exec session has ended")
		done <- xsess.Run(ctx)
	}()

	for {
		select {

		// Run has finished.
		case err := <-done:
			v.log.Info("command has finished", "error", err)
			exitCode := 0
			if err != nil {
				// Following in the path of exec.Command and ssh.ExitError, try to
				// detect if the error has a exit status associated with it and pass
				// it back to the client.
				if es, ok := err.(interface{ ExitStatus() int }); ok {
					exitCode = es.ExitStatus()
				} else {
					exitCode = 1
				}
			}

			if err := client.Send(&pb.EntrypointExecRequest{
				Event: &pb.EntrypointExecRequest_Exit_{
					Exit: &pb.EntrypointExecRequest_Exit{
						Code: int32(exitCode),
					},
				},
			}); err != nil {
				v.log.Warn("error sending exit message", "err", err)
			}

			// We don't return here, instead we wait for the remote side to see
			// our exit message and close the stream. That will be observed
			// as the above go routine closing respCh and the below case
			// seeing ok = false.

		// The server sent new info
		case resp, ok := <-respCh:
			if !ok {
				// channel is closed, we should terminate our child process.
				v.log.Info("exec recv stream closed")
				return xsess.Close()
			}

			switch event := resp.Event.(type) {
			case *pb.EntrypointExecResponse_Input:
				// Copy the input to stdin
				v.log.Trace("input received", "data", event.Input)
				io.Copy(stdinW, bytes.NewReader(event.Input))

			case *pb.EntrypointExecResponse_InputEof:
				v.log.Trace("input EOF, closing stdin")
				stdinW.Close()

			case *pb.EntrypointExecResponse_Winch:
				v.log.Debug("window size change event, changing")

				if err := xsess.PTYResize(event.Winch); err != nil {
					v.log.Warn("error changing window size, this doesn't quit the stream",
						"err", err)
				}
			}
		}
	}
}

// RunLogs connects to the server and streams logs from the function passed
// back to the server. The log entries will be associated with the Instance Id
// the virtual CEB is using.
func (v *Virtual) RunLogs(
	ctx context.Context,
	startTime time.Time,
	limit int,
	f func(ctx context.Context, lv *component.LogViewer) error,
) error {
	// Open our log stream
	client, err := v.cfg.Client.EntrypointLogStream(ctx)
	if err != nil {
		return status.Errorf(codes.Aborted,
			"failed to open a log stream: %s", err)
	}

	defer client.CloseAndRecv()

	v.log.Info("connecting virtual instance log stream")

	logs := make(chan component.LogEvent, 10)

	// A quick heads up: gRPC provides to ability to let the client of a recieve stream
	// tell the remote side "hey, I'm done now, don't send me anything else.". So instead
	// we wire up a context to the call and cancel it when we are done. This cancelation
	// will propogate to the server and they'll see that we have gone away.
	//
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ls := &component.LogViewer{
		Output:     logs,
		StartingAt: startTime,
		Limit:      limit,
	}

	done := make(chan error, 1)

	go func() {
		done <- f(ctx, ls)
	}()

	v.log.Info("processing logs for virtual instance")
	defer v.log.Info("finished logs for virtual instance")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-done:
			return err
		case lev, ok := <-logs:
			if !ok {
				// the logs channel was closed, so the plugin must have exited
				return nil
			}

			ts := timestamppb.New(lev.Timestamp)

			entry := &pb.LogBatch_Entry{
				Source:    pb.LogBatch_Entry_APP,
				Timestamp: ts,
				Line:      lev.Message,
			}

			v.log.Debug("sending log entry via EntrypointLogBatch")

			// TODO(evanphx) Do something with lev.Partition
			err = client.Send(&pb.EntrypointLogBatch{
				InstanceId: v.cfg.InstanceId,
				Lines:      []*pb.LogBatch_Entry{entry},
			})

			if err != nil {
				v.log.Error("error sending entrypoint log batch", "error", err)
				return nil
			}
		}
	}
}
