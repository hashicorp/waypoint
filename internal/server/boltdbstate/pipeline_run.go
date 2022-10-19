package boltdbstate

import (
	"context"
	"math"
	"strings"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

var pipelineRunBucket = []byte("pipeline_run")

func init() {
	dbBuckets = append(dbBuckets, pipelineRunBucket)
	dbIndexers = append(dbIndexers, (*State).pipelineRunIndexInit)
	schemas = append(schemas, pipelineRunIndexSchema)
}

// PipelineRunPut creates or updates the given PipelineRun.
func (s *State) PipelineRunPut(ctx context.Context, pr *pb.PipelineRun) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if pr.Pipeline == nil {
			return status.Error(codes.FailedPrecondition,
				"A pipeline ref for the pipeline run is required")
		}

		if pr.Id == "" {
			id, err := ulid()
			if err != nil {
				return err
			}
			pr.Id = id
		}

		// only alter sequence if this is a new pipeline run
		if pr.Sequence == 0 {
			pId := pr.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id
			iter, err := memTxn.ReverseLowerBound(
				pipelineRunIndexTableName,
				pipelineRunIndexPIdBySeq,
				pId, uint(math.MaxInt))
			if err != nil {
				return err
			}

			raw := iter.Next()

			// increment sequence if this is not the first run
			if raw != nil {
				idx := raw.(*pipelineRunIndexRecord)
				pr.Sequence = idx.Sequence + 1
			} else {
				pr.Sequence = 1
			}
		}

		return s.pipelineRunPut(dbTxn, memTxn, pr)
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
	// The data should be validated before this, but since it is a critical
	// issue if there are validation errors, we test again.

	if err := ptypes.ValidatePipelineRun(value); err != nil {
		return status.Errorf(codes.FailedPrecondition, err.Error())
	}

	// Get the global bucket and write the value to it.
	id := s.pipelineRunId(value)
	b := dbTxn.Bucket(pipelineRunBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.pipelineRunIndexSet(memTxn, id, value)
}

func (s *State) PipelineRunGetByJobId(ctx context.Context, jobId string) (*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		job, err := s.jobById(dbTxn, jobId)
		if job.Pipeline == nil {
			err = status.Errorf(codes.FailedPrecondition, "no pipeline run associated with job %q", job)
			return err
		}
		ref := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: job.Pipeline.PipelineId,
			},
		}
		p, err := s.pipelineGet(dbTxn, memTxn, ref)
		result, err = s.pipelineRunGet(dbTxn, memTxn, p.Id, job.Pipeline.RunSequence)
		return err
	})

	if result != nil && len(result.Jobs) < 1 {
		err = status.Errorf(codes.FailedPrecondition, "no jobs queued for pipeline run %q", result)
	}
	return result, err
}

// PipelineRunGet gets a PipelineRun by pipeline and sequence.
func (s *State) PipelineRunGet(ctx context.Context, ref *pb.Ref_Pipeline, seq uint64) (*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		p, err := s.pipelineGet(dbTxn, memTxn, ref)
		result, err = s.pipelineRunGet(dbTxn, memTxn, p.Id, seq)
		return err
	})

	return result, err
}

func (s *State) pipelineRunGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	pId string,
	seq uint64,
) (*pb.PipelineRun, error) {
	var result pb.PipelineRun
	b := dbTxn.Bucket(pipelineRunBucket)

	// Look up the first instance of the pipeline run where the sequence number and pipeline match
	raw, err := memTxn.First(pipelineRunIndexTableName,
		pipelineRunIndexPIdBySeq, pId, seq)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound,
			"pipeline run v%q could not be found for pipeline Id: %q",
			seq, pId)
	}
	// set the id to be looked up for GET
	idx, ok := raw.(*pipelineRunIndexRecord)
	if !ok {
		// This shouldn't happen, but guard against it...
		return nil, status.Error(codes.Internal,
			"failed to decode raw result to *pipelineRunIndexRecord!")
	}

	return &result, dbGet(b, []byte(strings.ToLower(idx.Id)), &result)
}

// PipelineRunGetLatest gets the latest PipelineRun by pipeline ID.
func (s *State) PipelineRunGetLatest(ctx context.Context, pId string) (*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.pipelineRunGetLatest(dbTxn, memTxn, pId)
		return err
	})

	return result, err
}

func (s *State) pipelineRunGetLatest(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	pId string,
) (*pb.PipelineRun, error) {
	var result pb.PipelineRun
	b := dbTxn.Bucket(pipelineRunBucket)

	// Look up the last instance of the pipeline run where the pipeline matches
	iter, err := memTxn.ReverseLowerBound(
		pipelineRunIndexTableName,
		pipelineRunIndexPIdBySeq,
		pId, uint(math.MaxInt))
	if err != nil {
		return nil, err
	}

	raw := iter.Next()

	if raw == nil {
		return nil, status.Errorf(codes.NotFound,
			"pipeline run could not be found for pipeline Id: %q", pId)
	}
	// set the id to be looked up for GET
	idx, ok := raw.(*pipelineRunIndexRecord)
	if !ok {
		// This shouldn't happen, but guard against it...
		return nil, status.Error(codes.Internal,
			"failed to decode raw result to *pipelineRunIndexRecord!")
	}

	return &result, dbGet(b, []byte(strings.ToLower(idx.Id)), &result)
}

// PipelineRunGetById gets a PipelineRun by pipeline run ID.
func (s *State) PipelineRunGetById(ctx context.Context, id string) (*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.pipelineRunGetById(dbTxn, memTxn, id)
		return err
	})

	return result, err
}

func (s *State) pipelineRunGetById(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	Id string,
) (*pb.PipelineRun, error) {
	var result pb.PipelineRun
	b := dbTxn.Bucket(pipelineRunBucket)

	return &result, dbGet(b, []byte(strings.ToLower(Id)), &result)
}

func (s *State) PipelineRunList(ctx context.Context, pRef *pb.Ref_Pipeline) ([]*pb.PipelineRun, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	pId, ok := pRef.Ref.(*pb.Ref_Pipeline_Id)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			"could not convert Ref %t to pipeline ID", pRef.Ref)
	}

	var out []*pb.PipelineRun
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		rrs, err := s.pipelineRunList(memTxn, pId.Id)
		if err != nil {
			return err
		}

		for _, idx := range rrs {
			var pr *pb.PipelineRun
			pr, err = s.pipelineRunGet(dbTxn, memTxn, idx.PipelineId, idx.Sequence)
			if err != nil {
				return err
			}
			out = append(out, pr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *State) pipelineRunList(
	memTxn *memdb.Txn,
	pId string,
) ([]*pipelineRunIndexRecord, error) {
	iter, err := memTxn.Get(pipelineRunIndexTableName, pipelineRunIndexPId+"_prefix", pId)
	if err != nil {
		return nil, err
	}

	var result []*pipelineRunIndexRecord
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*pipelineRunIndexRecord)

		result = append(result, idx)
	}

	return result, nil
}

// pipelineRunIndexSet writes an index record for a single pipelineRun.
func (s *State) pipelineRunIndexSet(txn *memdb.Txn, id []byte, value *pb.PipelineRun) error {
	record := &pipelineRunIndexRecord{
		Id:         string(id),
		Sequence:   value.Sequence,
		PipelineId: value.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id,
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
			pipelineRunIndexPIdBySeq: {
				Name:         pipelineRunIndexPIdBySeq,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "PipelineId",
							Lowercase: true,
						},

						&memdb.UintFieldIndex{
							Field: "Sequence",
						},
					},
				},
			},
			pipelineRunIndexPId: {
				Name:         pipelineRunIndexPId,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "PipelineId",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	pipelineRunIndexTableName = "pipelineRun-index"
	pipelineRunIndexId        = "id"
	pipelineRunIndexPId       = "pipeline-id"
	pipelineRunIndexPIdBySeq  = "seq"
)

type pipelineRunIndexRecord struct {
	Id         string
	Sequence   uint64
	PipelineId string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *pipelineRunIndexRecord) Copy() *pipelineRunIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
