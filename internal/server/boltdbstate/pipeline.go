package boltdbstate

import (
	"strings"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
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
	// The data should be validated before this, but it is a critical
	// issue if there are validation errors so we test again.
	if err := ptypes.ValidatePipeline(value); err != nil {
		return status.Errorf(codes.FailedPrecondition, err.Error())
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
	case *pb.Ref_Pipeline_Owner:
		s.log.Info("looking up pipeline by owner and name",
			"owner", r.Owner.Project, "name", r.Owner.PipelineName)

		// NOTE(briancain): This query doesn't seem to work as I'd expect it to.
		// It returns the "last inserted pipeline" rather than all pipelines
		// by project name. :thinking:
		iter, err := memTxn.Get(pipelineIndexTableName, pipelineIndexProjectId, r.Owner.Project.Project)
		if err != nil {
			return nil, err
		}

		// Look up if there's a pipeline name that exists by project ID owner
		for {
			raw := iter.Next()
			if raw == nil {
				// We're out of candidates and we found none.
				break
			}

			pipeIndex := raw.(*pipelineIndexRecord)
			// TODO:  delete this, used for debugging tests
			s.log.Info("looking at", "name", pipeIndex.Name)

			if pipeIndex.ProjectId == r.Owner.Project.Project &&
				pipeIndex.Name == r.Owner.PipelineName {
				pipelineId = pipeIndex.Id
				break
			}
		}

		if pipelineId == "" {
			// better error message here
			return nil, status.Errorf(codes.NotFound, "pipeline %q not found", r.Owner.PipelineName)
		}
	default:
		return nil, status.Error(
			codes.FailedPrecondition,
			"No valid ref provided to pipelineGet",
		)
	}

	return &result, dbGet(b, []byte(strings.ToLower(pipelineId)), &result)
}

// PipelineDelete deletes a pipeline by reference.
func (s *State) PipelineDelete(ref *pb.Ref_Pipeline) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.pipelineDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

func (s *State) pipelineDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Pipeline,
) error {
	// Get the pipeline. If it doesn't exist then we're successful.
	p, err := s.pipelineGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	id := s.pipelineId(p)
	if err := dbTxn.Bucket(pipelineBucket).Delete([]byte(p.Id)); err != nil {
		return err
	}

	// Delete from memdb
	if _, err := memTxn.DeleteAll(pipelineIndexTableName, pipelineIndexId, string(id)); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	return nil
}

func (s *State) PipelineList(pRef *pb.Ref_Project) ([]*pb.Pipeline, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.pipelineList(memTxn, pRef)
	if err != nil {
		return nil, err
	}

	var out []*pb.Pipeline
	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.pipelineGet(dbTxn, memTxn, ref)
			if err != nil {
				return err
			}

			out = append(out, val)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *State) pipelineList(
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) ([]*pb.Ref_Pipeline, error) {
	iter, err := memTxn.Get(pipelineIndexTableName, pipelineIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_Pipeline
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*pipelineIndexRecord)

		result = append(result, &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{
					Id: idx.Id,
				},
			},
		})
	}

	return result, nil
}

// pipelineIndexSet writes an index record for a single pipeline.
func (s *State) pipelineIndexSet(txn *memdb.Txn, id []byte, value *pb.Pipeline) error {
	record := &pipelineIndexRecord{
		Id:        string(id),
		ProjectId: value.Owner.(*pb.Pipeline_Project).Project.Project,
		Name:      value.Name,
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
			pipelineIndexName: {
				Name:         pipelineIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Name",
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
	pipelineIndexName      = "name"
)

type pipelineIndexRecord struct {
	Id        string
	ProjectId string
	Name      string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *pipelineIndexRecord) Copy() *pipelineIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
