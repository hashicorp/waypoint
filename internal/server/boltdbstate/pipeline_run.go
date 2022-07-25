package boltdbstate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	//"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

var pipelineRunBucket = []byte("pipeline_run")

func init() {
	dbBuckets = append(dbBuckets, pipelineRunBucket)
	dbIndexers = append(dbIndexers, (*State).pipelineRunIndexInit)
	schemas = append(schemas, pipelineRunIndexSchema)
}

// PipelineRunPut creates or updates the given PipelineRun.
func (s *State) PipelineRunPut(p *pb.PipelineRun) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if p.Pipeline == nil {
			return status.Error(codes.FailedPrecondition,
				"a Pipeline ref for the pipeline run is required")
		}

		if p.Id == "" {
			id, err := ulid()
			if err != nil {
				return err
			}

			p.Id = id
		}

		return s.pipelineRunPut(dbTxn, memTxn, p)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

func (s *State) pipelineRunPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.PipelineRun,
) error {
	// The data should be validated before this, but it is a critical
	// issue if there are validation errors so we test again.

	// TODO:XX add validation later
	//if err := ptypes.ValidatePipelineRun(value); err != nil {
	//	return status.Errorf(codes.FailedPrecondition, err.Error())
	//}
	//

	// Get the global bucket and write the value to it.
	id := s.pipelineRunId(value)
	b := dbTxn.Bucket(pipelineRunBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.pipelineRunIndexSet(memTxn, id, value)
}

// PipelineRunGet gets a PipelineRun by sequence or UUID.
func (s *State) PipelineRunGet(ref *pb.Ref_Pipeline, seq uint64) (*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		p, err := s.pipelineGet(dbTxn, memTxn, ref)
		result, err = s.pipelineRunGet(dbTxn, memTxn, p, fmt.Sprint(seq))
		return err
	})

	return result, err
}

func (s *State) pipelineRunGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	p *pb.Pipeline,
	seq string,
) (*pb.PipelineRun, error) {
	var result pb.PipelineRun
	b := dbTxn.Bucket(pipelineRunBucket)

	// Look up the first instance of the pipeline run where the sequence number and pipeline ref match
	raw, err := memTxn.First(pipelineRunIndexTableName,
		pipelineRunIndexIdBySeq, p.Id, seq)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound,
			"pipeline run v%q could not found for pipeline name: %q, uuid: %q",
			seq, p.Name, p.Id)
	}
	// set the id to be looked up for GET
	idx, ok := raw.(*pipelineRunIndexRecord)
	if !ok {
		// This shouldn't happen, but guard against it...
		return nil, status.Error(codes.Internal,
			"failed to decode raw result to *pipelineRunIndexRecord!")
	}
	var pipelineRunId string
	pipelineRunId = idx.Id

	return &result, dbGet(b, []byte(strings.ToLower(pipelineRunId)), &result)
}

//////// ---------------------------- ////
//func (s *State) PipelineRunList(pRef *pb.Ref_Project) ([]*pb.PipelineRun, error) {
//	memTxn := s.inmem.Txn(false)
//	defer memTxn.Abort()
//
//	refs, err := s.pipelineRunList(memTxn, pRef)
//	if err != nil {
//		return nil, err
//	}
//
//	var out []*pb.PipelineRun
//	err = s.db.View(func(dbTxn *bolt.Tx) error {
//		for _, ref := range refs {
//			val, err := s.pipelineRunGet(dbTxn, memTxn, ref)
//			if err != nil {
//				return err
//			}
//
//			out = append(out, val)
//		}
//
//		return nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return out, nil
//}

//func (s *State) pipelineRunList(
//	memTxn *memdb.Txn,
//	ref *pb.Ref_Project,
//) ([]*pb.Ref_PipelineRun, error) {
//	iter, err := memTxn.Get(pipelineRunIndexTableName, pipelineRunIndexId+"_prefix", "")
//	if err != nil {
//		return nil, err
//	}
//
//	var result []*pb.Ref_pipelineRun
//	for {
//		next := iter.Next()
//		if next == nil {
//			break
//		}
//		idx := next.(*pipelineRunIndexRecord)
//
//		result = append(result, &pb.Ref_pipelineRun{
//			Ref: &pb.Ref_pipelineRun_Id{
//				Id: &pb.Ref_pipelineRunId{
//					Id: idx.Id,
//				},
//			},
//		})
//	}
//
//	return result, nil
//}

// pipelineRunIndexSet writes an index record for a single pipelineRun.
func (s *State) pipelineRunIndexSet(txn *memdb.Txn, id []byte, value *pb.PipelineRun) error {
	record := &pipelineRunIndexRecord{
		Id:         string(id),
		Sequence:   fmt.Sprint(value.Sequence),
		PipelineId: value.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id,
	}

	// Insert the index
	return txn.Insert(pipelineRunIndexTableName, record)
}

// pipelineRunIndexInit initializes the pipelineRun index from persisted data.
func (s *State) pipelineRunIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(pipelineRunBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.PipelineRun
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		if err := s.pipelineRunIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) pipelineRunId(p *pb.PipelineRun) []byte {
	return []byte(strings.ToLower(p.Id))
}

func pipelineRunIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: pipelineRunIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			pipelineRunIndexId: {
				Name:         pipelineIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			pipelineRunIndexIdBySeq: {
				Name:         pipelineRunIndexIdBySeq,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "PipelineId",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Sequence",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

const (
	pipelineRunIndexTableName = "pipelineRun-index"
	pipelineRunIndexId        = "id"
	pipelineRunIndexIdBySeq   = "seq"
)

type pipelineRunIndexRecord struct {
	Id         string
	Sequence   string
	PipelineId string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *pipelineRunIndexRecord) Copy() *pipelineRunIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
