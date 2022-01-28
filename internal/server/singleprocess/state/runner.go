package state

import (
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const (
	runnerTableName   = "runners"
	runnerIdIndexName = "id"
)

func init() {
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

type runnerRecord struct {
	// The full Runner. All other fiels are derivatives of this.
	Runner *pb.Runner

	// Id of the runner
	Id string
}

func (s *State) RunnerCreate(r *pb.Runner) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Create our runner
	if err := txn.Insert(runnerTableName, newRunnerRecord(r)); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	txn.Commit()

	return nil
}

func (s *State) RunnerDelete(id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll(runnerTableName, runnerIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
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

	return raw.(*runnerRecord).Runner, nil
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
		record := next.(*runnerRecord)
		result = append(result, record.Runner)
	}

	return result, nil
}

// runnerEmpty returns true if there are no runners registered.
func (s *State) runnerEmpty(memTxn *memdb.Txn) (bool, error) {
	iter, err := memTxn.LowerBound(runnerTableName, runnerIdIndexName, "")
	if err != nil {
		return false, err
	}

	return iter.Next() == nil, nil
}

// newRunnerRecord creates a runnerRecord from a runner.
func newRunnerRecord(r *pb.Runner) *runnerRecord {
	rec := &runnerRecord{
		Runner: r,
		Id:     r.Id,
	}

	return rec
}
