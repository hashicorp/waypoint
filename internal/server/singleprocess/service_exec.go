package singleprocess

import (
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) StartExecStream(
	srv pb.Devflow_StartExecStreamServer,
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
	eventCh := make(chan *pb.EntrypointExecRequest)
	execRec := &state.InstanceExec{
		Args: start.Start.Args,
		Reader: &grpc_net_conn.Conn{
			Stream:   srv,
			Response: &pb.ExecStreamRequest{},
			Decode: grpc_net_conn.SimpleDecoder(func(msg proto.Message) *[]byte {
				return &msg.(*pb.ExecStreamRequest).Event.(*pb.ExecStreamRequest_Input_).Input.Data
			}),
		},
		EventCh: eventCh,
	}

	// Register the exec session
	err = s.state.InstanceExecCreateByDeployment(start.Start.DeploymentId, execRec)
	if err != nil {
		return err
	}

	// Make sure we always deregister it
	defer s.state.InstanceExecDelete(execRec.Id)

	// Loop through and read events
	for {
		var entryReq *pb.EntrypointExecRequest
		var closed bool
		select {
		case <-srv.Context().Done():
			// TODO: we need to notify the entrypoint side that we're over
			return nil

		case entryReq, closed = <-eventCh:
			// We got an event, exit out of the select and determine our action
		}

		if closed {
			return nil
		}

		log.Trace("event received", "event", entryReq.Event)
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
	}
}
