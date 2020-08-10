package state

import (
	"reflect"

	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// This file has the methods related to tracking assigned jobs. I pulled this
// out into a separate file since the job queueing logic is already quite long.

// parallelOps are the operation types that can be run in parallel. Anything
// NOT in this list will be queued by app and workspace.
var parallelOps = map[reflect.Type]struct{}{
	reflect.TypeOf((*pb.Job_Noop)(nil)):     struct{}{},
	reflect.TypeOf((*pb.Job_Validate)(nil)): struct{}{},
	reflect.TypeOf((*pb.Job_Auth)(nil)):     struct{}{},
}

func init() {
	schemas = append(schemas, jobAssignedSchema)
}

const (
	jobAssignedTableName   = "jobs-assigned"
	jobAssignedIdIndexName = "id"
)

func jobAssignedSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobAssignedTableName,
		Indexes: map[string]*memdb.IndexSchema{
			jobAssignedIdIndexName: &memdb.IndexSchema{
				Name:         jobAssignedIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Application",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Workspace",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

type jobAssignedIndex struct {
	Project     string
	Application string
	Workspace   string
}

// jobIsBlocked will return true if the given job is currently blocked because
// a job with the same project/app/workspace is executing.
//
// If ws is set then a watch will be added for any changes in assigned jobs.
// Note that this trigger doesn't mean that the blocking is necessarily gone
// but something changed to warrant rechecking.
func (s *State) jobIsBlocked(memTxn *memdb.Txn, idx *jobIndex, ws memdb.WatchSet) (bool, error) {
	// If this job represents a parallelizable operation type, then allow it.
	if _, ok := parallelOps[idx.OpType]; ok {
		return false, nil
	}

	// Look for this project/app/ws combo
	watchCh, value, err := memTxn.FirstWatch(
		jobAssignedTableName,
		jobAssignedIdIndexName,
		idx.Application.Project,
		idx.Application.Application,
		idx.Workspace.Workspace,
	)
	if err != nil {
		return false, err
	}
	if ws != nil {
		ws.Add(watchCh)
	}

	// Blocked if we have a record
	return value != nil, nil
}
