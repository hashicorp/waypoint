package state

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var (
	runnerBucket = []byte("runners")
)

const (
	runnerTableName   = "runners"
	runnerIdIndexName = "id"
)

func init() {
	// Note: there is no dbIndexer for runners because we never have to
	// reinit data from disk. When a runner registers, we always fetch the
	// data from disk. We don't need to prefetch it.
	dbBuckets = append(dbBuckets, runnerBucket)
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
}

func (s *State) RunnerCreate(r *pb.Runner) error {
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
func (s *State) RunnerDelete(id string) error {
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

func (s *State) RunnerById(id string) (*pb.Runner, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(runnerTableName, runnerIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "runner ID not found: %s", id)
	}

	return raw.(*runnerIndex).Runner, nil
}

func (s *State) RunnerList() ([]*pb.Runner, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(runnerTableName, runnerIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Runner
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		record := next.(*runnerIndex)
		result = append(result, record.Runner)
	}

	return result, nil
}

// runnerCreate creates the runner record and inserts it into the database.
// This operation is an upsert; it will update information if this runner
// has been seen before.
func (s *State) runnerCreate(dbTxn *bolt.Tx, memTxn *memdb.Txn, runnerpb *pb.Runner) error {
	now, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return err
	}

	// Zero out all the records that are server side set. These will be
	// replaced with real values if we have them.
	runnerpb.Online = false
	runnerpb.FirstSeen = now
	runnerpb.LastSeen = now
	runnerpb.AdoptionState = pb.Runner_NEW

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
	bucket := dbTxn.Bucket(triggerBucket)
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

// runnerIndexSet writes an index record for a single runner.
func (s *State) runnerIndexSet(txn *memdb.Txn, id []byte, runnerpb *pb.Runner) (*runnerIndex, error) {
	rec := &runnerIndex{
		Runner: runnerpb,
		Id:     runnerpb.Id,
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
