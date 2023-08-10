// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"context"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

var (
	runnerBucket = []byte("runners")
)

const (
	runnerTableName   = "runners"
	runnerIdIndexName = "id"
)

func init() {
	dbBuckets = append(dbBuckets, runnerBucket)
	dbIndexers = append(dbIndexers, (*State).runnerIndexInit)
	schemas = append(schemas, runnerSchema)
}

func runnerSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: runnerTableName,
		Indexes: map[string]*memdb.IndexSchema{
			runnerIdIndexName: {
				Name:         runnerIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
		},
	}
}

type runnerIndex struct {
	// The full Runner. All other fiels are derivatives of this.
	Runner *pb.Runner

	// Id of the runner
	Id string

	// State of adoption for this runner
	AdoptionState pb.Runner_AdoptionState
}

func (s *State) RunnerCreate(ctx context.Context, r *pb.Runner) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.runnerCreate(dbTxn, txn, r)
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

// RunnerDelete permanently deletes the runner record including any
// on-disk state such as adoption state. This effectively "forgets" the
// runner.
func (s *State) RunnerDelete(ctx context.Context, id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.runnerDelete(dbTxn, txn, id)
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

func (s *State) RunnerById(ctx context.Context, id string, ws memdb.WatchSet) (*pb.Runner, error) {
	// We only grab a read txn to memdb if we want a watch. Otherwise,
	// we just load the runner from disk.
	if ws != nil {
		memTxn := s.inmem.Txn(false)
		defer memTxn.Abort()

		watchCh, _, err := memTxn.FirstWatch(runnerTableName, runnerIdIndexName, id)
		if err != nil {
			return nil, err
		}
		ws.Add(watchCh)
	}

	// Get our value
	var result *pb.Runner
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.runnerById(dbTxn, id)
		return err
	})
	if err != nil {
		result = nil
	}

	return result, err
}

// RunnerList lists all the runners. This isn't a list of only online runners;
// this is ALL runners the database currently knows about.
func (s *State) RunnerList(ctx context.Context) ([]*pb.Runner, error) {
	var result []*pb.Runner
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		bucket := dbTxn.Bucket(runnerBucket)
		c := bucket.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var value pb.Runner
			if err := proto.Unmarshal(v, &value); err != nil {
				return err
			}

			result = append(result, &value)
		}

		return nil
	})

	return result, err
}

// RunnerOffline marks that a runner has gone offline. This is the preferred
// approach to deregistering a runner, since this will keep the adoption state
// around.
func (s *State) RunnerOffline(ctx context.Context, id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.runnerOffline(dbTxn, txn, id)
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

// RunnerAdopt marks a runner as adopted.
//
// If "implicit" is true, then this runner is implicitly adopted and the
// state goes to "PREADOPTED". This means that the runner instance itself
// was never explicitly adopted, but it already has a valid token so it is
// accepted.
func (s *State) RunnerAdopt(ctx context.Context, id string, implicit bool) error {
	state := pb.Runner_ADOPTED
	if implicit {
		state = pb.Runner_PREADOPTED
	}

	txn := s.inmem.Txn(true)
	defer txn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.runnerSetAdoptionState(dbTxn, txn, id, state)
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

// RunnerReject marks a runner as rejected.
func (s *State) RunnerReject(ctx context.Context, id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.runnerSetAdoptionState(dbTxn, txn, id, pb.Runner_REJECTED)
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

// runnerCreate creates the runner record and inserts it into the database.
// This operation is an upsert; it will update information if this runner
// has been seen before.
func (s *State) runnerCreate(dbTxn *bolt.Tx, memTxn *memdb.Txn, runnerpb *pb.Runner) error {
	now := timestamppb.Now()

	// Zero out all the records that are server side set. These will be
	// replaced with real values if we have them.
	runnerpb.FirstSeen = now
	runnerpb.LastSeen = now
	runnerpb.AdoptionState = pb.Runner_PENDING

	// Look up the runner in our database. If it exists, override the
	// values that are persistently stored.
	runnerOld, err := s.runnerById(dbTxn, runnerpb.Id)
	if status.Code(err) == codes.NotFound {
		runnerOld = nil
		err = nil
	}
	if err != nil {
		return err
	}
	if runnerOld != nil {
		runnerpb.FirstSeen = runnerOld.FirstSeen
		runnerpb.AdoptionState = runnerOld.AdoptionState

		// If we have non-matching labels, then we reset the adoption state to
		// new. This prevents a runner from being adopted for one environment
		// such as dev and then reregistering with labels for prod and retaining
		// the adoption status.
		hash1, err := serverptypes.RunnerLabelHash(runnerpb.Labels)
		if err != nil {
			return err
		}
		hash2, err := serverptypes.RunnerLabelHash(runnerOld.Labels)
		if err != nil {
			return err
		}

		// NOTE(mitchellh): in the future, we may want to have a setting that
		// allows a runner to change their labels without affecting adoption.
		// For now, we do not support this.
		if hash1 != hash2 {
			runnerpb.AdoptionState = pb.Runner_PENDING
		}
	}

	// If the runner has no first seen set, then set that.
	if runnerpb.FirstSeen == nil {
		runnerpb.FirstSeen = now
	}

	// Every time we see a runner, we update the last seen.
	runnerpb.LastSeen = now

	// Insert into bolt
	id := []byte(runnerpb.Id)
	if err := dbPut(dbTxn.Bucket(runnerBucket), id, runnerpb); err != nil {
		return err
	}

	// Insert into memdb
	idx := newRunnerIndex(runnerpb)
	return memTxn.Insert(runnerTableName, idx)
}

func (s *State) runnerDelete(dbTxn *bolt.Tx, memTxn *memdb.Txn, id string) error {
	// Delete from the database
	bucket := dbTxn.Bucket(runnerBucket)
	if err := bucket.Delete([]byte(id)); err != nil {
		return err
	}

	// Delete from memory
	err := memTxn.Delete(runnerTableName, &runnerIndex{Id: id})
	if err == memdb.ErrNotFound {
		err = nil
	}

	return nil
}

func (s *State) runnerSetAdoptionState(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	id string,
	state pb.Runner_AdoptionState,
) error {
	// Get our runner
	r, err := s.runnerById(dbTxn, id)
	if err != nil {
		return err
	}

	// Set our state
	r.AdoptionState = state

	// Insert into bolt
	if err := dbPut(dbTxn.Bucket(runnerBucket), []byte(id), r); err != nil {
		return err
	}

	// Insert into memdb, this will replace the old value.
	idx := newRunnerIndex(r)
	return memTxn.Insert(runnerTableName, idx)
}

func (s *State) runnerOffline(dbTxn *bolt.Tx, memTxn *memdb.Txn, id string) error {
	r, err := s.runnerById(dbTxn, id)
	if status.Code(err) == codes.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	// Determine if we need to delete this runner from the persisted DB
	//
	// NOTE(mitchellh): One day, we may want to keep around old runners for
	// awhile for historical purposes, to enable functions like "job history by runner",
	// etc. That's not part of this initial scope of work I'm doing around
	// adoption so we just delete in most cases. But, if we ever implemented that,
	// this is where you change it.
	del := false
	switch r.Kind.(type) {
	case *pb.Runner_Remote_:
		// Delete if the state is new, because there's no reason to keep around
		// a pending record for a non-adopted/rejected runner. But we DO want
		// to keep around states like ADOPTED or REJECTED so that if the runner
		// comes back, we know exactly how to handle them.
		//
		// We also delete PREADOPTED because if the runner comes back, we expect
		// it'll have a valid token to set it back to the PREADOPTED state.
		del = r.AdoptionState == pb.Runner_PENDING || r.AdoptionState == pb.Runner_PREADOPTED

	default:
		// All other runner types like ODR and Local we don't keep records of.
		del = true
	}
	if del {
		return s.runnerDelete(dbTxn, memTxn, id)
	}

	return nil
}

func (s *State) runnerById(dbTxn *bolt.Tx, id string) (*pb.Runner, error) {
	var result pb.Runner
	b := dbTxn.Bucket(runnerBucket)
	return &result, dbGet(b, []byte(id), &result)
}

// runnerEmpty returns true if there are no runners registered.
func (s *State) runnerEmpty(memTxn *memdb.Txn) (bool, error) {
	iter, err := memTxn.LowerBound(runnerTableName, runnerIdIndexName, "")
	if err != nil {
		return false, err
	}

	return iter.Next() == nil, nil
}

// runnerIndexInit initializes the config index from persisted data.
func (s *State) runnerIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(runnerBucket)
	c := bucket.Cursor()

	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		var value pb.Runner
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		_, err := s.runnerIndexSet(memTxn, k, &value)
		if err != nil {
			return err
		}
	}

	return nil
}

// runnerIndexSet writes an index record for a single runner.
func (s *State) runnerIndexSet(txn *memdb.Txn, id []byte, runnerpb *pb.Runner) (*runnerIndex, error) {
	rec := &runnerIndex{
		Runner:        runnerpb,
		Id:            runnerpb.Id,
		AdoptionState: runnerpb.AdoptionState,
	}

	// Insert the index
	return rec, txn.Insert(runnerTableName, rec)
}

// newRunnerIndex creates a runnerIndex from a runner.
func newRunnerIndex(r *pb.Runner) *runnerIndex {
	rec := &runnerIndex{
		Runner: r,
		Id:     r.Id,
	}

	return rec
}

// Copy should be called prior to any modifications to an existing runnerIndex.
func (idx *runnerIndex) Copy() *runnerIndex {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
