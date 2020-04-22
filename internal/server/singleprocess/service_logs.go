package singleprocess

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// TODO: test
func (s *service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Devflow_GetLogStreamServer,
) error {
	log := hclog.FromContext(srv.Context()).With("deployment_id", req.DeploymentId)
	ws := memdb.NewWatchSet()

	// Get all our initial records
	records, err := s.instancesByDeployment(req.DeploymentId, ws)
	if err != nil {
		return err
	}
	log.Trace("instances for deployment", "len", len(records))

	// For each record, start a goroutine that reads the log entries and sends them.
	for _, record := range records {
		instanceId := record.Id
		r := record.LogBuffer.Reader()

		instanceLog := log.With("instance_id", instanceId)
		instanceLog.Trace("instance log stream starting")
		go r.CloseContext(srv.Context())
		go func() {
			defer instanceLog.Debug("instance log stream ending")

			for {
				entries := r.Read(64)
				if entries == nil {
					return
				}

				instanceLog.Trace("sending instance log data", "entries", len(entries))
				srv.Send(&pb.LogBatch{
					DeploymentId: req.DeploymentId,
					InstanceId:   instanceId,
					Lines:        entries,
				})
			}
		}()
	}

	// TODO: use the watchset to detect new instances and add them to the
	// listening list.

	// Wait until we're done
	<-srv.Context().Done()
	return nil
}
