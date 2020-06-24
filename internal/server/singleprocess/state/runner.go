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

type runnerRecord struct {
	// The full Runner. All other fiels are derivatives of this.
	Runner *pb.Runner

	// Id of the runner
	Id string

	// Components are the components that this runner has access to.
	Components map[pb.Component_Type]map[string]*pb.Component
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
		return nil, status.Errorf(codes.NotFound, "runner ID not found")
	}

	return raw.(*runnerRecord).Runner, nil
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
		Runner:     r,
		Id:         r.Id,
		Components: make(map[pb.Component_Type]map[string]*pb.Component),
	}

	for _, c := range r.Components {
		m, ok := rec.Components[c.Type]
		if !ok {
			m = make(map[string]*pb.Component)
			rec.Components[c.Type] = m
		}

		m[c.Name] = c
	}

	return rec
}

// MatchComponentRefs tests if the given references are satisfied by this
// runner.
func (r *runnerRecord) MatchComponentRefs(refs []*pb.Ref_Component) bool {
	for _, ref := range refs {
		m, ok := r.Components[ref.Type]
		if !ok {
			return false
		}

		_, ok = m[ref.Name]
		if !ok {
			return false
		}

		// NOTE(mitchellh): In the future I imagine we'll do more checks here
		// such as versioning.
	}

	return true
}
