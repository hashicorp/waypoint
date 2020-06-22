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
			runnerIdIndexName: &memdb.IndexSchema{
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

func (s *State) RunnerCreate(rec *pb.Runner) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Create our runner
	if err := txn.Insert(runnerTableName, rec); err != nil {
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
		return nil, status.Errorf(codes.NotFound, "runner ID not found")
	}

	return raw.(*pb.Runner), nil
}
