package singleprocess

import (
	"io"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) StartExecStream(
	srv pb.Waypoint_StartExecStreamServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Read our first event which must be a Start event.
	log.Trace("waiting for Start message")
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	start, ok := req.Event.(*pb.ExecStreamRequest_Start_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"first message must be start type")
	}
	log = log.With("deployment_id", start.Start.DeploymentId)
	log.Debug("exec requested", "args", start.Start.Args)

	// Create our exec. We have to populate everything here first because
	// once we register, this will trigger any watchers to be notified of
	// a change and the instance should try to connect to us.
	clientEventCh := make(chan *pb.ExecStreamRequest)
	eventCh := make(chan *pb.EntrypointExecRequest)
	execRec := &state.InstanceExec{
		Args:              start.Start.Args,
		Pty:               start.Start.Pty,
		ClientEventCh:     clientEventCh,
		EntrypointEventCh: eventCh,
	}

	// Register the exec session
	err = s.state.InstanceExecCreateByDeployment(start.Start.DeploymentId, execRec)
	if err != nil {
		return err
	}

	// Make sure we always deregister it
	defer s.state.InstanceExecDelete(execRec.Id)

	// Always send the open message. In the future we'll send some metadata here.
	if err := srv.Send(&pb.ExecStreamResponse{
		Event: &pb.ExecStreamResponse_Open_{
			Open: &pb.ExecStreamResponse_Open{},
		},
	}); err != nil {
		return err
	}

	// Start our receive loop to read data from the client
	clientCloseCh := make(chan error, 1)
	go func() {
		defer close(clientEventCh)
		defer close(clientCloseCh)
		for {
			resp, err := srv.Recv()
			if err == io.EOF {
				// This means our client closed the stream. if the client
				// closed the stream, we want to end the exec stream completely.
				return
			}

			if err != nil {
				// Non EOF errors we will just send the error down and exit.
				clientCloseCh <- err
				return
			}

			clientEventCh <- resp
		}
	}()

	// Loop through and read events
	for {
		select {
		case <-srv.Context().Done():
			// The context was closed so we just exit. This will trigger
			// the EOF in the recv goroutine which will end the entrypoint
			// side as well.
			return nil

		case err := <-clientCloseCh:
			// The client closed the connection so we want to exit the stream.
			return err

		case entryReq, active := <-eventCh:
			// We got an event, exit out of the select and determine our action
			if !active {
				log.Debug("event channel closed, exiting")
				return nil
			}

			exit, err := s.handleEntrypointExecRequest(log, srv, entryReq)
			if exit || err != nil {
				return err
			}
		}
	}
}

func (s *service) handleEntrypointExecRequest(
	log hclog.Logger,
	srv pb.Waypoint_StartExecStreamServer,
	entryReq *pb.EntrypointExecRequest,
) (bool, error) {
	log.Trace("event received from entrypoint", "event", entryReq.Event)
	var send *pb.ExecStreamResponse
	exit := false
	switch event := entryReq.Event.(type) {
	case *pb.EntrypointExecRequest_Output_:
		send = &pb.ExecStreamResponse{
			Event: &pb.ExecStreamResponse_Output_{
				Output: &pb.ExecStreamResponse_Output{
					Channel: pb.ExecStreamResponse_Output_Channel(event.Output.Channel),
					Data:    event.Output.Data,
				},
			},
		}

	case *pb.EntrypointExecRequest_Exit_:
		exit = true
		send = &pb.ExecStreamResponse{
			Event: &pb.ExecStreamResponse_Exit_{
				Exit: &pb.ExecStreamResponse_Exit{
					Code: event.Exit.Code,
				},
			},
		}
	}

	// Send our response
	if send != nil {
		if err := srv.Send(send); err != nil {
			log.Warn("stream error", "err", err)
			return false, err
		}
	}

	return exit, nil
}
