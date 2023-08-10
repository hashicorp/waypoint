// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"reflect"

	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// This file has the methods related to tracking assigned jobs. I pulled this
// out into a separate file since the job queueing logic is already quite long.

// blockOps are the operation types that can NOT be run in parallel. Anything
// in this list will block if any other operation in this list is running
// for the app and workspace.
var blockOps = map[reflect.Type]struct{}{
	reflect.TypeOf((*pb.Job_Deploy)(nil)):  {},
	reflect.TypeOf((*pb.Job_Destroy)(nil)): {},
	reflect.TypeOf((*pb.Job_Release)(nil)): {},
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
			jobAssignedIdIndexName: {
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

// jobIsBlocked will return true if the given job is currently blocked because:
//   - a job with the same project/app/workspace is executing.
//   - a dependent job has not completed
//
// If ws is set then a watch will be added for any changes in assigned jobs.
// The watchset value is only modified in situations we are already blocked.
// If the job is not blocked, ws will not trigger for scenarios we might become
// blocked.
//
// Note that this trigger doesn't mean that the blocking is necessarily gone
// but something changed to warrant rechecking.
func (s *State) jobIsBlocked(memTxn *memdb.Txn, idx *jobIndex, ws memdb.WatchSet) (bool, error) {
	ok, err := s.jobIsBlockedByApp(memTxn, idx, ws)
	if err != nil {
		return ok, err
	}
	if ok {
		// If we're already blocked by this reason, just early exit.
		return true, nil
	}

	return s.jobIsBlockedByDeps(memTxn, idx, ws)
}

// jobIsBlockedByDeps returns true if a job can't be run because dependent
// jobs haven't completed yet.
func (s *State) jobIsBlockedByDeps(
	memTxn *memdb.Txn, idx *jobIndex, ws memdb.WatchSet) (bool, error) {
	// Can't be blocked without dependencies
	if len(idx.DependsOn) == 0 {
		return false, nil
	}

	// Go through each and determine if it is completed. This should be pretty
	// quick because all of this are in-memory indexed lookups.
	for _, id := range idx.DependsOn {
		watchCh, raw, err := memTxn.FirstWatch(jobTableName, jobIdIndexName, id)
		if err != nil {
			return false, err
		}

		// If the job no longer exists, we consider that it is complete.
		// We do this because on create, we verified all jobs we depended on
		// existed. Further, on erroneous terminal state, jobs also error all
		// dependents. Therefore, the only scenario a job suddenly disappears
		// is that it was pruned after being complete.
		if raw == nil {
			continue
		}

		// We only check if the job state is terminal. Any other state blocks us.
		// We allow errors here because errors should cascade to failing
		// this job earlier.
		depIdx := raw.(*jobIndex)
		if depIdx.State != pb.Job_SUCCESS && depIdx.State != pb.Job_ERROR {
			// Add this to our watch status, because we should recheck if
			// this job ever changes.
			ws.Add(watchCh)
			return true, nil
		}
	}

	// Not blocked
	return false, nil
}

// jobIsBlockedByApp returns true if a job can't be run because another
// job is already assigned to run with the same project/app/workspace.
func (s *State) jobIsBlockedByApp(memTxn *memdb.Txn, idx *jobIndex, ws memdb.WatchSet) (bool, error) {
	// If this job represents a parallelizable operation type, then allow it.
	if _, ok := blockOps[idx.OpType]; !ok {
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

// jobAssignedSet records the given job as assigned.
func (s *State) jobAssignedSet(memTxn *memdb.Txn, idx *jobIndex, assigned bool) error {
	// If this job represents a parallelizable operation type, then do nothing.
	if _, ok := blockOps[idx.OpType]; !ok {
		return nil
	}

	rec := &jobAssignedIndex{
		Project:     idx.Application.Project,
		Application: idx.Application.Application,
		Workspace:   idx.Workspace.Workspace,
	}

	if assigned {
		return memTxn.Insert(jobAssignedTableName, rec)
	}

	return memTxn.Delete(jobAssignedTableName, rec)
}
