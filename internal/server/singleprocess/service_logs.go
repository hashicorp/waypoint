package singleprocess

import (
	"github.com/hashicorp/go-memdb"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// TODO: test
func (s *service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Devflow_GetLogStreamServer,
) error {
	ws := memdb.NewWatchSet()

	// Get all our initial records
	records, err := s.instancesByDeployment(req.DeploymentId, ws)
	if err != nil {
		return err
	}

	// For each record, start a goroutine that reads the log entries and sends them.
	for _, record := range records {
		instanceId := record.Id
		r := record.LogBuffer.Reader()
		go r.CloseContext(srv.Context())
		go func() {
			entries := r.Read(64)
			if entries == nil {
				return
			}

			srv.Send(&pb.LogBatch{
				DeploymentId: req.DeploymentId,
				InstanceId:   instanceId,
				Lines:        entries,
			})
		}()
	}

	// TODO: use the watchset to detect new instances and add them to the
	// listening list.

	// Wait until we're done
	<-srv.Context().Done()
	return nil
}
