package state

import (
	"context"
	"math/rand"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

type InstanceExec struct {
	Id         int64
	InstanceId string

	Args []string
	Pty  *pb.ExecStreamRequest_PTY

	ClientEventCh     <-chan *pb.ExecStreamRequest
	EntrypointEventCh chan<- *pb.EntrypointExecRequest
	Connected         uint32

	// This is the context that the client side is running inside.
	// It is used by the entrypoint side to detect if the client is still
	// around or not.
	Context context.Context
}

func (s *State) InstanceExecCreateByTargetedInstance(id string, exec *InstanceExec) error {
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
func (s *State) InstanceExecCreateForVirtualInstance(ctx context.Context, id string, exec *InstanceExec) error {
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

// CalculateInstanceExecByDeployment considers all the instances registered
// to the given deployment and finds the one that is least loaded. If there
// are no instances, returns a ResourceExhausted error. Calls to this
// function in quick succession will return could return the same instance,
// which is why a simple random sampling is done on all prospective instances.
func (s *State) CalculateInstanceExecByDeployment(did string) (*Instance, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()

	// Find all the instances by deployment
	iter, err := txn.Get(instanceTableName, instanceDeploymentIdIndexName, did)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// Go through each to try to find the least loaded. Most likely there
	// will be an instance with no exec sessions and we prefer that.
	var min []*Instance
	minCount := 0
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		rec := raw.(*Instance)

		// When looking through all the instances for an exec capable instance
		// we only consider LONG_RUNNING type instances. These are the only ones
		// that it makes sense to send random exec sessions to.
		if rec.Type != pb.Instance_LONG_RUNNING {
			continue
		}

		execs, err := s.instanceExecListByInstanceId(txn, rec.Id, nil)
		if err != nil {
			return nil, err
		}

		// Otherwise we keep track of the lowest "load" exec which we just
		// choose by the minimum number of registered sessions.
		if len(execs) < minCount {
			// If we're less than the min count we've seen before, then
			// we reset the slice of candidates because we only want
			// candidates with this count.
			min = nil
		}
		if min == nil || len(execs) <= minCount {
			min = append(min, rec)
			minCount = len(execs)
		}
	}

	if min == nil {
		return nil, status.Errorf(codes.ResourceExhausted,
			"No available instances for exec.")
	}

	if len(min) == 1 {
		return min[0], nil
	}

	// To avoid callers always picking the first one if there are multiple
	// canditates, pick a random one. This helps even out the load.

	tgt := rand.Intn(len(min))

	return min[tgt], nil
}

func (s *State) InstanceExecCreateByDeployment(did string, exec *InstanceExec) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Find all the instances by deployment
	iter, err := txn.Get(instanceTableName, instanceDeploymentIdIndexName, did)
	if err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	// Go through each to try to find the least loaded. Most likely there
	// will be an instance with no exec sessions and we prefer that.
	var min *Instance
	minCount := 0
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		rec := raw.(*Instance)

		// When looking through all the instances for an exec capable instance
		// we only consider LONG_RUNNING type instances. These are the only ones
		// that it makes sense to send random exec sessions to.
		if rec.Type != pb.Instance_LONG_RUNNING {
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
		return status.Errorf(codes.ResourceExhausted,
			"No available instances for exec.")
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

func (s *State) InstanceExecDelete(id int64) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll(instanceExecTableName, instanceExecIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceExecById(id int64) (*InstanceExec, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceExecTableName, instanceExecIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance exec ID not found")
	}

	return raw.(*InstanceExec), nil
}

func (s *State) InstanceExecListByInstanceId(id string, ws memdb.WatchSet) ([]*InstanceExec, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()
	return s.instanceExecListByInstanceId(txn, id, ws)
}

func (s *State) instanceExecListByInstanceId(
	txn *memdb.Txn, id string, ws memdb.WatchSet,
) ([]*InstanceExec, error) {
	// Find all the exec sessions
	iter, err := txn.Get(instanceExecTableName, instanceExecInstanceIdIndexName, id)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	var result []*InstanceExec
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*InstanceExec))
	}

	ws.Add(iter.WatchCh())

	return result, nil
}
