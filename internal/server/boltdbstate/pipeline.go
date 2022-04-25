package boltdbstate

import (
	"strings"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var pipelineBucket = []byte("pipeline")

func init() {
	dbBuckets = append(dbBuckets, pipelineBucket)
	dbIndexers = append(dbIndexers, (*State).pipelineIndexInit)
	schemas = append(schemas, pipelineIndexSchema)
}

// PipelinePut creates or updates the given Pipeline.
func (s *State) PipelinePut(p *pb.Pipeline) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if p.Id == "" {
			id, err := ulid()
			if err != nil {
				return err
			}

			p.Id = id
		}

		return s.pipelinePut(dbTxn, memTxn, p)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

func (s *State) pipelinePut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Pipeline,
) error {
	// Basic preconditions for state storage. The data is expected
	// to be fully validated prior to this step so we just validate
	// what we need to ensure data consistency.
	if value.Owner == nil {
		return status.Error(codes.FailedPrecondition,
			"an owner must be set on the project")
	}

	if len(value.Steps) == 0 {
		return status.Error(codes.FailedPrecondition,
			"a pipeline requires at least one step")
	}

	// Verify we have exactly one root
	rootCount := 0
	for _, step := range value.Steps {
		if len(step.DependsOn) == 0 {
			rootCount++
		}
	}
	if rootCount != 1 {
		return status.Error(codes.FailedPrecondition,
			"a pipeline requires exactly one root step")
	}

	// Verify there are no cycles in the steps
	var stepGraph graph.Graph
	for _, step := range value.Steps {
		// Add our job
		stepGraph.Add(step.Name)

		// Add any dependencies
		for _, dep := range step.DependsOn {
			stepGraph.Add(dep)
			stepGraph.AddEdge(dep, step.Name)

			if _, ok := value.Steps[dep]; !ok {
				return status.Errorf(codes.FailedPrecondition,
					"Step %q depends on non-existent step %q", step, dep)
			}
		}
	}
	if cycles := stepGraph.Cycles(); len(cycles) > 0 {
		return status.Errorf(codes.FailedPrecondition,
			"Step dependencies contain one or more cycles: %#v", cycles)
	}

	// Get the global bucket and write the value to it.
	id := s.pipelineId(value)
	b := dbTxn.Bucket(pipelineBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.pipelineIndexSet(memTxn, id, value)
}

// PipelineGet gets a pipeline by reference.
func (s *State) PipelineGet(ref *pb.Ref_Pipeline) (*pb.Pipeline, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Pipeline
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.pipelineGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

func (s *State) pipelineGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Pipeline,
) (*pb.Pipeline, error) {
	var result pb.Pipeline
	b := dbTxn.Bucket(pipelineBucket)

	var pipelineId string
	switch r := ref.Ref.(type) {
	case *pb.Ref_Pipeline_Id:
		s.log.Info("looking up pipeline by id", "id", r.Id)
		pipelineId = r.Id.Id
	default:
		return nil, status.Error(
			codes.FailedPrecondition,
			"No valid ref provided to pipelineGet",
		)
	}

	return &result, dbGet(b, []byte(strings.ToLower(pipelineId)), &result)
}

// pipelineIndexSet writes an index record for a single pipeline.
func (s *State) pipelineIndexSet(txn *memdb.Txn, id []byte, value *pb.Pipeline) error {
	record := &pipelineIndexRecord{
		Id:        string(id),
		ProjectId: value.Owner.(*pb.Pipeline_Project).Project.Project,
	}

	// Insert the index
	return txn.Insert(pipelineIndexTableName, record)
}

// pipelineIndexInit initializes the pipeline index from persisted data.
func (s *State) pipelineIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(pipelineBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Pipeline
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		if err := s.pipelineIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) pipelineId(p *pb.Pipeline) []byte {
	return []byte(strings.ToLower(p.Id))
}

func pipelineIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: pipelineIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			pipelineIndexId: {
				Name:         pipelineIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			pipelineIndexProjectId: {
				Name:         pipelineIndexProjectId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "ProjectId",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	pipelineIndexTableName = "pipeline-index"
	pipelineIndexId        = "id"
	pipelineIndexProjectId = "projectid"
)

type pipelineIndexRecord struct {
	Id        string
	ProjectId string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *pipelineIndexRecord) Copy() *pipelineIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
