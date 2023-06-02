package singleprocess

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Waypoint_EntrypointConfigServer,
) error {
	log := hclog.FromContext(srv.Context())
	ctx := srv.Context()

	// Fetch the deployment info so we can calculate the config variables to send.
	// This also verifies this deployment exists.
	deployment, err := s.GetDeployment(srv.Context(), &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{Id: req.DeploymentId},
		},
	})
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to get deployment in entrypoint config",
			"deployment_id",
			req.DeploymentId,
		)
	}

	if tok := s.decodedTokenFromContext(ctx); tok != nil {
		if tok.UnusedEntrypoint != nil && tok.UnusedEntrypoint.DeploymentId != req.DeploymentId {
			return hcerr.Externalize(
				log,
				err,
				"entrypoint token invalid for this deployment ID", "deployment_id",
				req.DeploymentId,
			)
		}
	}

	// Create our record
	log = log.With("deployment_id", req.DeploymentId, "instance_id", req.InstanceId)
	log.Trace("registering entrypoint")
	record := &serverstate.Instance{
		Id:           req.InstanceId,
		DeploymentId: req.DeploymentId,
		Project:      deployment.Application.Project,
		Application:  deployment.Application.Application,
		Workspace:    deployment.Workspace.Workspace,
		LogBuffer:    logbuffer.New(),
		Type:         req.Type,
		DisableExec:  req.DisableExec,
	}
	if err := s.state(ctx).InstanceCreate(ctx, record); err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to get create an instance in entrypoint config",
			"deployment_id",
			req.DeploymentId,
		)
	}

	// Handling exec requests is optional so we check if state supports
	// them and only if so do we add them to the list of things we'll check
	// for.
	iexec, ok := s.state(ctx).(serverstate.InstanceExecHandler)
	if !ok {
		iexec = nil
	}

	// Defer deleting this.
	// TODO(mitchellh): this is too aggressive and we want to have some grace
	// period for reconnecting clients. We should clean this up.
	defer func() {
		// We want to close all our readers at the end of this
		defer record.LogBuffer.Close()

		// Delete the entrypoint first
		log.Trace("deleting entrypoint")
		if err := s.state(ctx).InstanceDelete(ctx, record.Id); err != nil {
			log.Error("failed to delete instance data. This should not happen.", "err", err)
		}

		if iexec != nil {
			// Delete any active but unconnected exec requests. This can happen
			// if the entrypoint crashed after an exec was assigned to the entrypoint.
			log.Trace("closing any unconnected exec requests")
			execs, err := iexec.InstanceExecListByInstanceId(ctx, record.Id, nil)
			if err != nil {
				log.Error("failed to query instance exec list. This should not happen.", "err", err)
			} else {
				for _, exec := range execs {
					if atomic.CompareAndSwapUint32(&exec.Connected, 0, 1) {
						close(exec.EntrypointEventCh)
					}
				}
			}
		}
	}()

	// Build our config in a loop.
	for {
		ws := memdb.NewWatchSet()

		// Get our exec requests
		var execs []*serverstate.InstanceExec
		if iexec != nil {
			execs, err = iexec.InstanceExecListByInstanceId(ctx, req.InstanceId, ws)
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to find an instance for exec",
				)
			}
		}

		// Refresh the application record in case it's been deleted
		_, err = s.state(ctx).AppGet(ctx, deployment.Application)
		if err != nil {
			log.Warn("detected removed application in entrypoint",
				"project", deployment.Application.Project,
				"application", deployment.Application.Application,
				"lookup-error", err,
			)

			log.Warn("exiting EntrypointConfig because project was deleted")
			return status.Error(codes.Unavailable, "project has been deleted")
		}

		// Build our config
		config := &pb.EntrypointConfig{}
		for _, exec := range execs {
			config.Exec = append(config.Exec, &pb.EntrypointConfig_Exec{
				Index: int64(exec.Id),
				Args:  exec.Args,
				Pty:   exec.Pty,
			})
		}

		// Write our deployment info
		config.Deployment = &pb.EntrypointConfig_DeploymentInfo{
			Component: deployment.Component,
			Labels:    deployment.Labels,
		}

		// Get the config vars in use
		vars, err := s.state(ctx).ConfigGetWatch(ctx, &pb.ConfigGetRequest{
			Scope: &pb.ConfigGetRequest_Application{
				Application: deployment.Application,
			},
			Workspace: deployment.Workspace,
			Labels:    deployment.Labels,
		}, ws)
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to watch config in entrypoint config",
			)
		}
		config.EnvVars = vars

		config.FileChangeSignal, err = s.state(ctx).GetFileChangeSignal(
			ctx, deployment.Application,
		)
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to get file change signal in entrypoint config",
			)
		}

		// Get the config sources we need for our vars. We only do this if
		// at least one var has a dynamic value.
		if varContainsDynamic(vars) {
			// NOTE(mitchellh): For now we query all the types and always send it
			// all down. In the future we may want to consider filtering this
			// by only the types we actually need above.
			sources, err := s.state(ctx).ConfigSourceGetWatch(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Application{
					Application: deployment.Application,
				},
				Workspace: deployment.Workspace,
			}, ws)
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to watch config source in entrypoint config",
				)
			}

			config.ConfigSources = sources
		}

		// If we have the URL service setup, note that
		if v := s.urlConfig; v != nil {
			// Get our base config
			s.urlCEBMu.RLock()
			pbVal := proto.Clone(s.urlCEB).(*pb.EntrypointConfig_URLService)
			ws.Add(s.urlCEBWatchCh) // Watch for changes
			s.urlCEBMu.RUnlock()

			var flatLabels []string
			for k, v := range deployment.Labels {
				flatLabels = append(flatLabels, fmt.Sprintf("%s=%s", k, v))
			}

			// Determine our URL fragment
			pd := &ptypes.Deployment{Deployment: deployment}

			// We always have these default labels for the URL service.
			flatLabels = append(flatLabels,
				hznLabelApp+"="+deployment.Application.Application,
				hznLabelProject+"="+deployment.Application.Project,
				hznLabelWorkspace+"="+deployment.Workspace.Workspace,
				hznLabelInstance+"="+record.Id,

				":deployment="+pd.URLFragment(),
				":deployment-order="+strings.ToLower(deployment.Id),
			)
			pbVal.Labels = strings.Join(flatLabels, ",")

			// If the token is empty, don't send anything down to the
			// entrypoint yet since the entrypoint doesn't yet handle changes
			// to the URL config.
			// NOTE(mitchellh): when ceb supports changes to URL config, remove this
			if pbVal.Token != "" {
				config.UrlService = pbVal
			}
		}

		// Send new config
		if err := srv.Send(&pb.EntrypointConfigResponse{
			Config: config,
		}); err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to send config in entrypoint config",
			)
		}

		// Nil out the stuff we used so that if we're waiting awhile we can GC
		config = nil
		execs = nil

		// Wait for any changes
		if err := ws.WatchCtx(srv.Context()); err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to watch for changes in entrypoint config",
			)
		}
	}
}

// TODO: test
func (s *Service) EntrypointLogStream(
	server pb.Waypoint_EntrypointLogStreamServer,
) error {
	ctx := server.Context()

	log := hclog.FromContext(server.Context())

	// TODO(mitchellh): We only support logs if we're using the in-memory
	// state store. We will add support for our other stores later.
	inmemstate, ok := s.state(ctx).(*boltdbstate.State)
	if !ok {
		return hcerr.Externalize(
			hclog.FromContext(ctx),
			status.Errorf(codes.Unimplemented,
				"state storage doesn't support log streaming"),
			"state storage doesn't support log streaming",
		)
	}

	var buf *logbuffer.Buffer
	for {
		// Read the next log entry
		batch, err := server.Recv()
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to receive entrypoint log",
			)
		}

		// If we haven't initialized our buffer yet, do that
		if buf == nil {
			log = log.With("instance_id", batch.InstanceId)

			// Read our instance record
			instance, err := s.state(ctx).InstanceById(ctx, batch.InstanceId)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					// See if we have a instance logs entry to use instead.
					// These are used by the logs plugin functionality to provide a place
					// to rendezvous logs sent by the plugin with the waiting client
					// without generating a full Instance.

					log.Info("no Instance found, attempting to lookup InstanceLogs record instead")
					il, err := inmemstate.InstanceLogsByInstanceId(ctx, batch.InstanceId)
					if err != nil {
						return hcerr.Externalize(
							log,
							err,
							"failed to find log by instance ID in entrypoint",
						)
					}

					log.Info("using InstanceLogs record")
					buf = il.LogBuffer
				} else {
					return hcerr.Externalize(
						log,
						err,
						"failed to find log by instance ID in entrypoint",
					)
				}
			} else {
				// Get our log buffer
				buf = instance.LogBuffer
			}
		}

		// Log that we received data in trace mode
		if log.IsTrace() {
			log.Trace("received data", "lines", len(batch.Lines))
		}

		// Strip any trailing whitespace
		for _, entry := range batch.Lines {
			entry.Line = strings.TrimSuffix(entry.Line, "\n")
		}

		// Convert the lines to an interface{} for our buffer
		entries := make([]logbuffer.Entry, len(batch.Lines))
		for i, l := range batch.Lines {
			entries[i] = l
		}

		// Write our log data to the circular buffer
		buf.Write(entries...)
	}
}

func (s *Service) EntrypointExecStream(
	server pb.Waypoint_EntrypointExecStreamServer,
) error {
	ctx := server.Context()
	log := hclog.FromContext(server.Context())

	// Exec support is optional for a state interface, so we need to check that
	// it's supported first.
	iexec, ok := s.state(ctx).(serverstate.InstanceExecHandler)
	if !ok {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.Unimplemented,
				"state storage doesn't support exec streaming"),
			"state storage doesn't support exec streaming",
		)
	}

	// Receive our opening message so we can determine the exec stream.
	req, err := server.Recv()
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to receive entrypoint exex stream",
		)
	}
	open, ok := req.Event.(*pb.EntrypointExecRequest_Open_)
	if !ok {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.FailedPrecondition,
				"error reading entrypoint exec stream, first message must be open type"),
			"error reading entrypoint exec stream, first message must be open type",
		)
	}

	// Get our instance and look for this exec index
	exec, err := iexec.InstanceExecConnect(ctx, open.Open.Index)
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to connect to instance for exec",
		)
	}
	log = log.With("instance_id", exec.InstanceId, "index", open.Open.Index)

	// Mark we're connected
	if !atomic.CompareAndSwapUint32(&exec.Connected, 0, 1) {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.FailedPrecondition,
				"exec session is already open for this index"),
			"exec session is already open for this index",
		)
	}
	log.Debug("exec stream open")

	// Create a new context that we'll use manage the below go routine.
	// We use a new context rather than server.Context() so that we don't
	// exit too early because we see the context Done(). Specificly we want
	// to allow server.Recv() to return any buffered messages even if the
	// server.Context() is Done().
	// This fresh context will only be canceled we are no longer able to
	// process any client messages.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a goroutine that just waits for events from the entrypoint
	// and sends them along to the client side.
	errCh := make(chan error, 1)
	go func() {
		defer cancel()

		// backgroundRetry indicates that we detected a context cancellation
		// while holding a message and are handling that message in another
		// goroutine first.
		//
		// sawExit indicates that we observed either an Exit or Error message
		// flow through. We use this to control if we should synthesizes a Error
		// message about an improper closure.
		var sawExit bool

		// Close the event channel we send to. This signals to the receiving
		// side in StartExecStream (service_exec.go) that the entrypoint
		// exited and it should also exit the client stream.
		defer func() {
			// We defer the close again here so that if sendExit fails, we always
			// manage to get the channel closed.
			defer close(exec.EntrypointEventCh)

			// If we observed a Exit or Error event, we don't need to do this.
			if sawExit {
				return
			}

			req := &pb.EntrypointExecRequest{
				Event: &pb.EntrypointExecRequest_Error_{
					Error: &pb.EntrypointExecRequest_Error{
						Error: status.New(codes.Aborted, "server side exited unexpectedly").Proto(),
					},
				},
			}

			// We're again careful here to not block forever, but depend on the
			// client's context to decide if we should give up.
			select {
			case exec.EntrypointEventCh <- req:
				log.Debug("entrypoint event for improper closure dispatched")
			case <-exec.Context.Done():
			}
		}()

		for {
			log.Debug("waiting for entrypoint exec event")
			req, err := server.Recv()
			if err == io.EOF {
				// On EOF, this means the client closed their write side.
				// In this case, we assume we have exited and exit accordingly.
				return
			}

			if err != nil {
				// For any other error, we send the error along and exit the
				// read loop. The sent error will be picked up and sent back
				// as a result to the client.
				errCh <- err
				return
			}
			log.Trace("entrypoint event received", "event", hclog.Fmt("%#v", req.Event))

			// If this is an exit or error event track that so we don't synthesize one.
			// We synthesize an Error value if this loop is going to returning without
			// having sent one of these types to tell the client that something happened.
			switch req.Event.(type) {
			case *pb.EntrypointExecRequest_Exit_, *pb.EntrypointExecRequest_Error_:
				sawExit = true
			}

			// Send the event along. We if the reciever is gone (ie its context is Done())
			// then we don't bother. We don't depend on server.Context() here because
			// that is done within server.Recv() and we want to allow it to returned buffered
			// messages even if its internal context is Done().
			select {
			case exec.EntrypointEventCh <- req:
			// ok
			case <-exec.Context.Done():
				// oops, guess they went away. Keep reading req's so we don't
				// block though.
			}

			// If this is an exit or error event then we also exit this loop now.
			switch event := req.Event.(type) {
			case *pb.EntrypointExecRequest_Exit_:
				log.Debug("exec stream exiting due to exit message", "code", event.Exit.Code)
				return
			case *pb.EntrypointExecRequest_Error_:
				log.Debug("exec stream exiting due to client error",
					"error", event.Error.Error.Message)
				return
			}
		}
	}()

	// Note to the caller that we're opened. It is very important that
	// we call this AFTER the goroutine is started above so that we
	// properly close the exec.EntrypointEventCh channel. This channel has
	// to be closed when this function exits so that the client side properly
	// exits.
	if err := server.Send(&pb.EntrypointExecResponse{
		Event: &pb.EntrypointExecResponse_Opened{
			Opened: true,
		},
	}); err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failure to send entrypoint exec response",
		)
	}

	// Loop through our receive loop
	for {
		select {

		// Wait on the above goroutine.
		case <-ctx.Done():

			// Double check and see if there was an error and if so, return it.
			select {
			case err = <-errCh:
				return hcerr.Externalize(
					log,
					err,
					"error in entrypoint exec stream",
				)
			default:
				return nil
			}

		// The above goroutine has errored.
		case err := <-errCh:
			return hcerr.Externalize(
				log,
				err,
				"error in entrypoint exec stream",
			)

		// The client goroutine has finished.
		case req, active := <-exec.ClientEventCh:
			if !active {
				log.Debug("client event channel closed, exiting")
				return nil
			}

			if err := s.handleClientExecRequest(log, server, req); err != nil {
				return hcerr.Externalize(
					log,
					err,
					"error handling client exec request",
				)
			}
		}
	}
}

func (s *Service) handleClientExecRequest(
	log hclog.Logger,
	srv pb.Waypoint_EntrypointExecStreamServer,
	req *pb.ExecStreamRequest,
) error {
	log.Debug("event received from client", "event", req.Event)
	var send *pb.EntrypointExecResponse
	switch event := req.Event.(type) {
	case *pb.ExecStreamRequest_Input_:
		send = &pb.EntrypointExecResponse{
			Event: &pb.EntrypointExecResponse_Input{
				Input: event.Input.Data,
			},
		}

	case *pb.ExecStreamRequest_InputEof:
		send = &pb.EntrypointExecResponse{
			Event: &pb.EntrypointExecResponse_InputEof{
				InputEof: &empty.Empty{},
			},
		}

	case *pb.ExecStreamRequest_Winch:
		send = &pb.EntrypointExecResponse{
			Event: &pb.EntrypointExecResponse_Winch{
				Winch: event.Winch,
			},
		}
	}

	// Send our response
	if send != nil {
		if err := srv.Send(send); err != nil {
			log.Warn("stream error", "err", err)
			return err
		}
	}

	return nil
}

// varContainsDynamic returns true if there are any dynamic values in the list.
func varContainsDynamic(vars []*pb.ConfigVar) bool {
	for _, v := range vars {
		if _, ok := v.Value.(*pb.ConfigVar_Dynamic); ok {
			return true
		}
	}

	return false
}
