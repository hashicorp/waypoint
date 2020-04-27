package singleprocess

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
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

	// Start our receive loop to read data from the client
	go func() {
		defer close(clientEventCh)
		for {
			resp, err := srv.Recv()
			if err != nil {
				// TODO
				continue
			}

			clientEventCh <- resp
		}
	}()

	// Loop through and read events
	for {
		select {
		case <-srv.Context().Done():
			// TODO: we need to notify the entrypoint side that we're over
			log.Debug("context ended")
			return nil

		case entryReq, active := <-eventCh:
			// We got an event, exit out of the select and determine our action
			if !active {
				log.Debug("event channel closed, exiting")
				return nil
			}

			if err := s.handleEntrypointExecRequest(log, srv, entryReq); err != nil {
				return err
			}
		}
	}
}

func (s *service) handleEntrypointExecRequest(
	log hclog.Logger,
	srv pb.Waypoint_StartExecStreamServer,
	entryReq *pb.EntrypointExecRequest,
) error {
	log.Trace("event received from entrypoint", "event", entryReq.Event)
	var send *pb.ExecStreamResponse
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
			return err
		}
	}

	return nil
}
