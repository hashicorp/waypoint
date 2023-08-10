// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"context"
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
func (s *State) PipelinePut(ctx context.Context, p *pb.Pipeline) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if p.Id == "" {
			// look up if pipeline already exists by Owner and Name
			pipelineProjRef, ok := p.Owner.(*pb.Pipeline_Project)
			if !ok {
				return status.Error(codes.FailedPrecondition,
					"unsupported pipeline project owner, Pipeline_Project is expected")
			}

			pipe, err := s.pipelineGet(dbTxn, memTxn, &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      pipelineProjRef.Project,
						PipelineName: p.Name,
					},
				},
			})
			if err != nil && status.Code(err) != codes.NotFound {
				// if not found, that's ok. We're creating a new entry
				return err
			}

			if pipe != nil {
				p.Id = pipe.Id
			} else {
				id, err := ulid()
				if err != nil {
					return err
				}

				p.Id = id
			}
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
func (s *State) PipelineGet(ctx context.Context, ref *pb.Ref_Pipeline) (*pb.Pipeline, error) {
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
		s.log.Trace("looking up pipeline by id", "id", r.Id)
		pipelineId = r.Id
	case *pb.Ref_Pipeline_Owner:
		s.log.Trace("looking up pipeline by owner and name",
			"owner", r.Owner.Project, "name", r.Owner.PipelineName)

		// Look up the first instance of the pipeline owners pipeline name
		raw, err := memTxn.First(pipelineIndexTableName,
			pipelineIndexName, r.Owner.Project.Project, r.Owner.PipelineName)
		if err != nil {
			return nil, err
		}
		if raw == nil {
			return nil, status.Errorf(codes.NotFound,
				"pipeline %q could not found for owner %q",
				r.Owner.PipelineName, r.Owner.Project.Project)
		}

		// set the id to be looked up for GET
		idx, ok := raw.(*pipelineIndexRecord)
		if !ok {
			// This shouldn't happen, but guard against it...
			return nil, status.Error(codes.Internal,
				"failed to decode raw result to *pipelineIndexRecord!")
		}

		pipelineId = idx.Id
	default:
		return nil, status.Error(
			codes.FailedPrecondition,
			"No valid ref provided to pipelineGet",
		)
	}

	return &result, dbGet(b, []byte(strings.ToLower(pipelineId)), &result)
}

// PipelineDelete deletes a pipeline by reference.
func (s *State) PipelineDelete(ctx context.Context, ref *pb.Ref_Pipeline) error {
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

func (s *State) PipelineList(ctx context.Context, pRef *pb.Ref_Project) ([]*pb.Pipeline, error) {
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
	iter, err := memTxn.Get(pipelineIndexTableName, pipelineIndexProjectId, ref.Project)
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
				Id: idx.Id,
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
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "ProjectId",
					Lowercase: true,
				},
			},
			pipelineIndexName: {
				Name:         pipelineIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "ProjectId",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Name",
							Lowercase: true,
						},
					},
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
