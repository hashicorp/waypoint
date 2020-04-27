package singleprocess

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TODO: test
func (s *service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Devflow_GetLogStreamServer,
) error {
	log := hclog.FromContext(srv.Context()).With("deployment_id", req.DeploymentId)

	// We keep track of what instances we already have readers for here.
	var instanceSetLock sync.Mutex
	instanceSet := make(map[string]struct{})

	// We loop forever so that we can automatically get any new instances that
	// join as we have an open log stream.
	for {
		// Get all our records
		ws := memdb.NewWatchSet()
		records, err := s.state.InstancesByDeployment(req.DeploymentId, ws)
		if err != nil {
			return err
		}
		log.Trace("instances for deployment", "len", len(records))

		// For each record, start a goroutine that reads the log entries and sends them.
		for _, record := range records {
			instanceId := record.Id

			// If we already have a reader for this, then do nothing.
			instanceSetLock.Lock()
			_, exit := instanceSet[instanceId]
			instanceSet[instanceId] = struct{}{}
			instanceSetLock.Unlock()
			if exit {
				continue
			}

			// Start our reader up
			r := record.LogBuffer.Reader()
			instanceLog := log.With("instance_id", instanceId)
			instanceLog.Trace("instance log stream starting")
			go r.CloseContext(srv.Context())
			go func() {
				defer instanceLog.Debug("instance log stream ending")
				defer func() {
					instanceSetLock.Lock()
					defer instanceSetLock.Unlock()
					delete(instanceSet, instanceId)
				}()

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

		// Wait for changes or to be done
		if err := ws.WatchCtx(srv.Context()); err != nil {
			// If our context ended, exit with that
			if err := srv.Context().Err(); err != nil {
				return err
			}

			return err
		}
	}
}
