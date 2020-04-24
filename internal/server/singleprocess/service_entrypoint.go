package singleprocess

import (
	"context"
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
)

func init() {
	memdbSchema.Tables["instances"] = &memdb.TableSchema{
		Name: "instances",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			"deployment-id": &memdb.IndexSchema{
				Name:         "deployment-id",
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "DeploymentId",
					Lowercase: true,
				},
			},
		},
	}
}

// TODO: test
func (s *service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Devflow_EntrypointConfigServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Create our record
	log = log.With("deployment_id", req.DeploymentId, "instance_id", req.InstanceId)
	log.Trace("registering entrypoint")
	record := &instanceRecord{
		Id:           req.InstanceId,
		DeploymentId: req.DeploymentId,
		LogBuffer:    logbuffer.New(),
	}
	if err := s.instancesCreate(record); err != nil {
		return err
	}

	// Defer deleting this.
	// TODO(mitchellh): this is too aggressive and we want to have some grace
	// period for reconnecting clients. We should clean this up.
	defer func() {
		// We want to close all our readers at the end of this
		defer record.LogBuffer.Close()

		log.Trace("deleting entrypoint")
		tx := s.inmem.Txn(true)
		if err := tx.Delete("instances", record); err != nil {
			log.Error("failed to delete instance data. This should not happen.", "err", err)
		}
		tx.Commit()
	}()

	// Send initial config
	if err := srv.Send(&pb.EntrypointConfigResponse{}); err != nil {
		return err
	}

	// TODO(mitchellh): loop, send down any changes in configuration.
	<-srv.Context().Done()

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
			instance, err := s.instanceById(batch.InstanceId)
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
	rec, err := s.instanceById(open.Open.InstanceId)
	if err != nil {
		return err
	}
	exec, ok := rec.Exec[int32(open.Open.Index)]
	if !ok {
		return status.Errorf(codes.NotFound,
			"exec session index not found")
	}
	log = log.With("instance_id", open.Open.InstanceId, "index", open.Open.Index)

	// Mark we're connected
	if !atomic.CompareAndSwapUint32(&exec.Connected, 0, 1) {
		return status.Errorf(codes.FailedPrecondition,
			"exec session is already open for this index")
	}
	log.Debug("exec stream open")

	// We want to ensure we close the exec stream ALWAYS in case we don't
	// gracefully close it below.
	defer exec.Cancel()

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

func (s *service) instancesCreate(record *instanceRecord) error {
	// Insert this mapping into our memdb
	tx := s.inmem.Txn(true)
	defer tx.Abort()
	if err := tx.Insert("instances", record); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	tx.Commit()

	return nil
}

func (s *service) instanceById(id string) (*instanceRecord, error) {
	tx := s.inmem.Txn(false)
	raw, err := tx.First("instances", "id", id)
	tx.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"instance ID not found, please call EntrypointConfig first")
	}

	return raw.(*instanceRecord), nil
}

func (s *service) instancesByDeployment(id string, ws memdb.WatchSet) ([]*instanceRecord, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()

	// Find all the instances
	iter, err := txn.Get("instances", "deployment-id", id)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// If we're tracking changes, add that
	if ws != nil {
		ws.Add(iter.WatchCh())
	}

	var result []*instanceRecord
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*instanceRecord))
	}

	return result, nil
}

// instanceRecord is the record type that we'll insert into memdb
type instanceRecord struct {
	Id           string
	DeploymentId string
	LogBuffer    *logbuffer.Buffer
	Exec         map[int32]*instanceExec
}

type instanceExec struct {
	Args      []string
	Reader    io.Reader
	EventCh   chan<- *pb.EntrypointExecRequest
	Cancel    context.CancelFunc
	Connected uint32
}
