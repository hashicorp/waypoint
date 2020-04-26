package singleprocess

import (
	"io"
	"strings"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/internal/server/logbuffer"
	"github.com/mitchellh/devflow/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Devflow_EntrypointConfigServer,
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

	return nil
}

// TODO: test
func (s *service) EntrypointLogStream(
	server pb.Devflow_EntrypointLogStreamServer,
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

	return server.SendAndClose(&empty.Empty{})
}

// TODO: test
func (s *service) EntrypointExecStream(
	server pb.Devflow_EntrypointExecStreamServer,
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
	defer close(exec.EventCh)

	// Connect the reader to send data down
	go io.Copy(&grpc_net_conn.Conn{
		Stream:  server,
		Request: &pb.EntrypointExecResponse{},
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			return &msg.(*pb.EntrypointExecResponse).Data
		}),
	}, exec.Reader)

	// Loop through our receive loop
	for {
		req, err := server.Recv()
		if err != nil {
			// TODO: error handling
			return err
		}

		// Send the event
		exec.EventCh <- req

		// If this is an exit or error event then we also exit this loop now.
		switch event := req.Event.(type) {
		case *pb.EntrypointExecRequest_Exit_:
			log.Debug("exec stream exiting due to exit message", "code", event.Exit.Code)
			return nil
		case *pb.EntrypointExecRequest_Error_:
			log.Debug("exec stream exiting due to client error",
				"error", event.Error.Error.Message)
			return nil
		}
	}

	return nil
}
