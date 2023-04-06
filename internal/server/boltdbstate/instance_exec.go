// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

const (
	instanceExecTableName           = "instance-execs"
	instanceExecIdIndexName         = "id"
	instanceExecInstanceIdIndexName = "deployment-id"
)

func init() {
	schemas = append(schemas, instanceExecSchema)
}

var instanceExecId int64

func instanceExecSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: instanceExecTableName,
		Indexes: map[string]*memdb.IndexSchema{
			instanceExecIdIndexName: {
				Name:         instanceExecIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.IntFieldIndex{
					Field: "Id",
				},
			},

			instanceExecInstanceIdIndexName: {
				Name:         instanceExecInstanceIdIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "InstanceId",
					Lowercase: true,
				},
			},
		},
	}
}

var _ serverstate.InstanceExecHandler = (*State)(nil)

func (s *State) InstanceExecCreateByTargetedInstance(ctx context.Context, id string, exec *serverstate.InstanceExec) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// If the caller specified an instance id already, then just validate it.
	if id == "" {
		return status.Errorf(codes.NotFound, "No instance id given")
	}

	raw, err := txn.First(instanceTableName, instanceIdIndexName, id)
	if err != nil {
		return err
	}

	if raw == nil {
		return status.Errorf(codes.NotFound, "No instance by given id: %s", id)
	}

	instance := raw.(*serverstate.Instance)
	if instance.DisableExec {
		return status.Errorf(codes.PermissionDenied,
			"The requested instance (id: %s) does not support exec.", id)
	}

	// Set our ID
	exec.Id = atomic.AddInt64(&instanceExecId, 1)
	exec.InstanceId = id

	// Insert
	if err := txn.Insert(instanceExecTableName, exec); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

// InstanceExecCreateForVirtualInstance registers the given InstanceExec record on
// the instance specified. The instance does not yet have to be known (as it may
// not yet have connected to the server) so this code will use memdb watchers
// to detect the instance as it connects and then register the exec.
func (s *State) InstanceExecCreateForVirtualInstance(ctx context.Context, id string, exec *serverstate.InstanceExec) error {
	// If the caller specified an instance id already, then just validate it.
	if id == "" {
		return status.Errorf(codes.NotFound, "No instance id given")
	}

	_, err := s.instanceByIdWaiting(ctx, id)
	if err != nil {
		return err
	}

	// Now make a new write txn. We don't want to hold a write txn above (plus we have
	// to create a one per loop)
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Set our ID
	exec.Id = atomic.AddInt64(&instanceExecId, 1)
	exec.InstanceId = id

	// Insert
	if err := txn.Insert(instanceExecTableName, exec); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceExecCreateByDeployment(ctx context.Context, did string, exec *serverstate.InstanceExec) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Find all the instances by deployment
	iter, err := txn.Get(instanceTableName, instanceDeploymentIdIndexName, did)
	if err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	// Go through each to try to find the least loaded. Most likely there
	// will be an instance with no exec sessions and we prefer that.
	var min *serverstate.Instance
	minCount := 0
	empty := true
	disabled := false
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		rec := raw.(*serverstate.Instance)

		// When looking through all the instances for an exec capable instance
		// we only consider LONG_RUNNING type instances. These are the only ones
		// that it makes sense to send random exec sessions to.
		if rec.Type != pb.Instance_LONG_RUNNING {
			continue
		}

		// We saw at least one instance
		empty = false

		// If this instance doesn't support exec then we ignore it.
		if rec.DisableExec {
			// We saw at least one disabled
			disabled = true
			continue
		}

		execs, err := s.instanceExecListByInstanceId(txn, rec.Id, nil)
		if err != nil {
			return err
		}

		// Zero length exec means we take it right away
		if len(execs) == 0 {
			min = rec
			break
		}

		// Otherwise we keep track of the lowest "load" exec which we just
		// choose by the minimum number of registered sessions.
		if min == nil || len(execs) < minCount {
			min = rec
			minCount = len(execs)
		}
	}

	if min == nil {
		// If we have no running instances then show a special error.
		if empty {
			return status.Errorf(codes.ResourceExhausted, strings.TrimSpace(errNoRunningInstances))
		}

		// If we have at least one disabled, that means all have to be disabled
		// to get to this error.
		if disabled {
			return status.Errorf(codes.ResourceExhausted, strings.TrimSpace(errExecAllDisabled))
		}

		// This SHOULD be impossible since right now we'll always assign
		// an instance exec if we're non-empty and have any non-disabled.
		// Therefore, we'll keep this error somewhat vague. This should not
		// happen.
		return status.Errorf(codes.ResourceExhausted, "No available instances for exec.")
	}

	// Set the instance ID that we'll be using
	exec.InstanceId = min.Id

	// Set our ID
	exec.Id = atomic.AddInt64(&instanceExecId, 1)

	// Insert
	if err := txn.Insert(instanceExecTableName, exec); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceExecDelete(ctx context.Context, id int64) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	if _, err := txn.DeleteAll(instanceExecTableName, instanceExecIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceExecById(ctx context.Context, id int64) (*serverstate.InstanceExec, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceExecTableName, instanceExecIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance exec ID not found")
	}

	return raw.(*serverstate.InstanceExec), nil
}

func (s *State) InstanceExecConnect(ctx context.Context, id int64) (*serverstate.InstanceExec, error) {
	return s.InstanceExecById(ctx, id)
}

func (s *State) InstanceExecWaitConnected(ctx context.Context, exec *serverstate.InstanceExec) error {
	return nil
}

func (s *State) InstanceExecListByInstanceId(ctx context.Context, id string, ws memdb.WatchSet) ([]*serverstate.InstanceExec, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()
	return s.instanceExecListByInstanceId(txn, id, ws)
}

func (s *State) instanceExecListByInstanceId(
	txn *memdb.Txn, id string, ws memdb.WatchSet,
) ([]*serverstate.InstanceExec, error) {
	// Find all the exec sessions
	iter, err := txn.Get(instanceExecTableName, instanceExecInstanceIdIndexName, id)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	var result []*serverstate.InstanceExec
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*serverstate.InstanceExec))
	}

	ws.Add(iter.WatchCh())

	return result, nil
}

const (
	errNoRunningInstances = `
No running instances found for exec!

If you just recently deployed, the instances could still be starting up.
Otherwise, please diagnose the issue by inspecting your application logs.
If application logs are not available, the application may have failed to start.
`

	errExecAllDisabled = `
Exec is not available for any running instance. Every instance has exec
explicitly disabled. This is only possible by disabling exec at deploy time.
It is not possible to re-enable Waypoint exec for this deployment using
"waypoint config".
`
)
