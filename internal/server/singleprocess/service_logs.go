package singleprocess

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Waypoint_GetLogStreamServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Default the limit
	if req.LimitBacklog == 0 {
		req.LimitBacklog = 50
	}

	var instanceFunc func(ws memdb.WatchSet) ([]*state.Instance, error)
	switch scope := req.Scope.(type) {
	case *pb.GetLogStreamRequest_DeploymentId:
		log = log.With("deployment_id", scope.DeploymentId)
		instanceFunc = func(ws memdb.WatchSet) ([]*state.Instance, error) {
			return s.state.InstancesByDeployment(scope.DeploymentId, ws)
		}

	case *pb.GetLogStreamRequest_Application_:
		if scope.Application == nil ||
			scope.Application.Application == nil ||
			scope.Application.Workspace == nil {
			return status.Errorf(
				codes.FailedPrecondition,
				"application scope requires the application and workspace fields to be set",
			)
		}

		log = log.With(
			"project", scope.Application.Application.Project,
			"application", scope.Application.Application.Application,
			"workspace", scope.Application.Workspace.Workspace,
		)
		instanceFunc = func(ws memdb.WatchSet) ([]*state.Instance, error) {
			return s.state.InstancesByApp(
				scope.Application.Application,
				scope.Application.Workspace,
				ws,
			)
		}

	default:
		return status.Errorf(
			codes.FailedPrecondition,
			"invalid scope supplied: %T",
			req.Scope,
		)
	}

	// We keep track of what instances we already have readers for here.
	var instanceSetLock sync.Mutex
	instanceSet := make(map[string]struct{})

	// We loop forever so that we can automatically get any new instances that
	// join as we have an open log stream.
	for {
		// Get all our records
		ws := memdb.NewWatchSet()
		records, err := instanceFunc(ws)
		if err != nil {
			return err
		}
		log.Trace("instances loaded", "len", len(records))

		// For each record, start a goroutine that reads the log entries and sends them.
		for _, record := range records {
			instanceId := record.Id
			deploymentId := record.DeploymentId

			// If we already have a reader for this, then do nothing.
			instanceSetLock.Lock()
			_, exit := instanceSet[instanceId]
			instanceSet[instanceId] = struct{}{}
			instanceSetLock.Unlock()
			if exit {
				continue
			}

			// Start our reader up
			r := record.LogBuffer.Reader(req.LimitBacklog)
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
					entries := r.Read(64, true)
					if entries == nil {
						return
					}

					lines := make([]*pb.LogBatch_Entry, len(entries))
					for i, v := range entries {
						lines[i] = v.(*pb.LogBatch_Entry)
					}

					instanceLog.Trace("sending instance log data", "entries", len(entries))
					srv.Send(&pb.LogBatch{
						DeploymentId: deploymentId,
						InstanceId:   instanceId,
						Lines:        lines,
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
