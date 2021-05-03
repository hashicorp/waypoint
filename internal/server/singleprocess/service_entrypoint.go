package singleprocess

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	"github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Waypoint_EntrypointConfigServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Fetch the deployment info so we can calculate the config variables to send.
	// This also verifies this deployment exists.
	deployment, err := s.GetDeployment(srv.Context(), &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{Id: req.DeploymentId},
		},
	})
	if err != nil {
		return err
	}

	// Create our record
	log = log.With("deployment_id", req.DeploymentId, "instance_id", req.InstanceId)
	log.Trace("registering entrypoint")
	record := &state.Instance{
		Id:           req.InstanceId,
		DeploymentId: req.DeploymentId,
		Project:      deployment.Application.Project,
		Application:  deployment.Application.Application,
		Workspace:    deployment.Workspace.Workspace,
		LogBuffer:    logbuffer.New(),
		Type:         req.Type,
	}
	if err := s.state.InstanceCreate(record); err != nil {
		return err
	}

	// Defer deleting this.
	// TODO(mitchellh): this is too aggressive and we want to have some grace
	// period for reconnecting clients. We should clean this up.
	defer func() {
		// We want to close all our readers at the end of this
		defer record.LogBuffer.Close()

		// Delete the entrypoint first
		log.Trace("deleting entrypoint")
		if err := s.state.InstanceDelete(record.Id); err != nil {
			log.Error("failed to delete instance data. This should not happen.", "err", err)
		}

		// Delete any active but unconnected exec requests. This can happen
		// if the entrypoint crashed after an exec was assigned to the entrypoint.
		log.Trace("closing any unconnected exec requests")
		execs, err := s.state.InstanceExecListByInstanceId(record.Id, nil)
		if err != nil {
			log.Error("failed to query instance exec list. This should not happen.", "err", err)
		} else {
			for _, exec := range execs {
				if atomic.CompareAndSwapUint32(&exec.Connected, 0, 1) {
					close(exec.EntrypointEventCh)
				}
			}
		}
	}()

	// Build our config in a loop.
	for {
		ws := memdb.NewWatchSet()
		execs, err := s.state.InstanceExecListByInstanceId(req.InstanceId, ws)
		if err != nil {
			return err
		}

		// Build our config
		config := &pb.EntrypointConfig{}
		for _, exec := range execs {
			config.Exec = append(config.Exec, &pb.EntrypointConfig_Exec{
				Index: exec.Id,
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
		vars, err := s.state.ConfigGetWatch(&pb.ConfigGetRequest{
			Scope: &pb.ConfigGetRequest_Application{
				Application: deployment.Application,
			},
		}, ws)
		if err != nil {
			return err
		}
		config.EnvVars = vars

		config.FileChangeSignal, err = s.state.GetFileChangeSignal(
			deployment.Application,
		)
		if err != nil {
			return err
		}

		// Get the config sources we need for our vars. We only do this if
		// at least one var has a dynamic value.
		if varContainsDynamic(vars) {
			// NOTE(mitchellh): For now we query all the types and always send it
			// all down. In the future we may want to consider filtering this
			// by only the types we actually need above.
			sources, err := s.state.ConfigSourceGetWatch(&pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},
			}, ws)
			if err != nil {
				return err
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
			return err
		}

		// Nil out the stuff we used so that if we're waiting awhile we can GC
		config = nil
		execs = nil

		// Wait for any changes
		if err := ws.WatchCtx(srv.Context()); err != nil {
			return err
		}
	}
}

// TODO: test
func (s *service) EntrypointLogStream(
	server pb.Waypoint_EntrypointLogStreamServer,
) error {
	log := hclog.FromContext(server.Context())

	var buf *logbuffer.Buffer
	for {
		// Read the next log entry
		batch, err := server.Recv()
		if err != nil {
			return err
		}

		// If we haven't initialized our buffer yet, do that
		if buf == nil {
			log = log.With("instance_id", batch.InstanceId)

			// Read our instance record
			instance, err := s.state.InstanceById(batch.InstanceId)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					// See if we have a instance logs entry to use instead.
					// These are used by the logs plugin functionality to provide a place
					// to rendezvous logs sent by the plugin with the waiting client
					// without generating a full Instance.

					log.Info("no Instance found, attempting to lookup InstanceLogs record instead")
					il, err := s.state.InstanceLogsByInstanceId(batch.InstanceId)
					if err != nil {
						return err
					}

					buf = il.LogBuffer
				} else {
					return err
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

func (s *service) EntrypointExecStream(
	server pb.Waypoint_EntrypointExecStreamServer,
) error {
	log := hclog.FromContext(server.Context())

	// Receive our opening message so we can determine the exec stream.
	req, err := server.Recv()
	if err != nil {
		return err
	}
	open, ok := req.Event.(*pb.EntrypointExecRequest_Open_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"first message must be open type")
	}

	// Get our instance and look for this exec index
	exec, err := s.state.InstanceExecById(open.Open.Index)
	if err != nil {
		return err
	}
	log = log.With("instance_id", exec.InstanceId, "index", open.Open.Index)

	// Mark we're connected
	if !atomic.CompareAndSwapUint32(&exec.Connected, 0, 1) {
		return status.Errorf(codes.FailedPrecondition,
			"exec session is already open for this index")
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

			// Send the event along. We if the reciever is gone (ie it's context is Done())
			// then we don't bother. We don't depend on server.Context() here because
			// that is done within server.Recv() and we want to allow it to returned buffered
			// messages even if it's internal context is Done().
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
		return err
	}

	// Loop through our receive loop
	for {
		select {

		// Wait on the above goroutine.
		case <-ctx.Done():

			// Double check and see if there was an error and if so, return it.
			select {
			case err = <-errCh:
				return err
			default:
				return nil
			}

		// The above goroutine has errored.
		case err := <-errCh:
			return err

		// The client goroutine has finished.
		case req, active := <-exec.ClientEventCh:
			if !active {
				log.Debug("client event channel closed, exiting")
				return nil
			}

			if err := s.handleClientExecRequest(log, server, req); err != nil {
				return err
			}
		}
	}
}

func (s *service) handleClientExecRequest(
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
