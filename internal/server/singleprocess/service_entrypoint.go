package singleprocess

import (
	"sync"

	"github.com/armon/circbuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/pkg/circbufsync"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// For now we just store logs in memory in circular buffers, one per
// instance of an application. This is NOT what we want to do long term
// probably but it was easiest to get started.
var (
	logBuffers     = make(map[string]*circbufsync.Buffer)
	logBuffersLock sync.Mutex
	logBufferSize  int64 = 1024 * 1024 * 4 // 4 MB
)

// TODO: test
func (s *service) EntrypointConfig(
	req *pb.EntrypointConfigRequest,
	srv pb.Devflow_EntrypointConfigServer,
) error {
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

	var buf *circbufsync.Buffer
	for {
		// Read the next log entry
		batch, err := server.Recv()
		if err != nil {
			return err
		}

		// If we haven't initialized our buffer yet, do that
		if buf == nil {
			buf, err = s.initLogBuffer(batch.InstanceId)
			if err != nil {
				return err
			}

			log = log.With("instance_id", batch.InstanceId)
		}

		// Log that we received data in trace mode
		if log.IsTrace() {
			log.Trace("received data", "len", len(batch.Data))
		}

		// Write our log data to the circular buffer
		if _, err := buf.Write(batch.Data); err != nil {
			return err
		}
	}

	return server.SendAndClose(&empty.Empty{})
}

// initLogBuffer initializes the circular buffer for an entrypoint ID.
func (s *service) initLogBuffer(id string) (*circbufsync.Buffer, error) {
	logBuffersLock.Lock()
	defer logBuffersLock.Unlock()

	buf, ok := logBuffers[id]
	if ok {
		return buf, nil
	}

	cbuf, err := circbuf.NewBuffer(logBufferSize)
	if err != nil {
		return nil, err
	}

	buf = circbufsync.New(cbuf)
	logBuffers[id] = buf
	return buf, nil
}
