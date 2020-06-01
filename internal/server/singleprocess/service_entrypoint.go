package singleprocess

import (
	"context"
	"io"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Waypoint_EntrypointConfigServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Create our record
	log = log.With("deployment_id", req.DeploymentId, "instance_id", req.InstanceId)
	log.Trace("registering entrypoint")
	record := &state.Instance{
		Id:           req.InstanceId,
		DeploymentId: req.DeploymentId,
		LogBuffer:    logbuffer.New(),
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

		log.Trace("deleting entrypoint")
		if err := s.state.InstanceDelete(record.Id); err != nil {
			log.Error("failed to delete instance data. This should not happen.", "err", err)
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
				return err
			}

			// Get our log buffer
			buf = instance.LogBuffer
		}

		// Log that we received data in trace mode
		if log.IsTrace() {
			log.Trace("received data", "lines", len(batch.Lines))
		}

		// Strip any trailing whitespace
		for _, entry := range batch.Lines {
			entry.Line = strings.TrimSuffix(entry.Line, "\n")
		}

		// Write our log data to the circular buffer
		buf.Write(batch.Lines...)
	}
}

// TODO: test
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

	// Always close the event channel which signals to the reader end that
	// we are done.
	defer close(exec.EntrypointEventCh)

	// Note to the caller that we're opened
	if err := server.Send(&pb.EntrypointExecResponse{
		Event: &pb.EntrypointExecResponse_Opened{
			Opened: true,
		},
	}); err != nil {
		return err
	}

	// Create a context we can use to cancel
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	// Create a goroutine that just waits for events from the entrypoint
	// and sends them along to the client side.
	errCh := make(chan error, 1)
	go func() {
		defer cancel()

		for {
			log.Trace("waiting for entrypoint exec event")
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
			log.Trace("entrypoint event received", "event", req.Event)

			// Send the event along
			select {
			case exec.EntrypointEventCh <- req:
			case <-ctx.Done():
				return
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

	// Loop through our receive loop
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errCh:
			return err

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
	log.Trace("event received from client", "event", req.Event)
	var send *pb.EntrypointExecResponse
	switch event := req.Event.(type) {
	case *pb.ExecStreamRequest_Input_:
		send = &pb.EntrypointExecResponse{
			Event: &pb.EntrypointExecResponse_Input{
				Input: event.Input.Data,
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
