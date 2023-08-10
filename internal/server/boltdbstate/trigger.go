// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"context"
	"strings"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var triggerBucket = []byte("trigger")

func init() {
	dbBuckets = append(dbBuckets, triggerBucket)
	dbIndexers = append(dbIndexers, (*State).triggerIndexInit)
	schemas = append(schemas, triggerIndexSchema)
}

// TriggerPut creates or updates the given Trigger.
func (s *State) TriggerPut(ctx context.Context, t *pb.Trigger) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if t.Project == nil {
			return status.Error(codes.FailedPrecondition, "A Project is required to create a trigger")
		}

		if t.Workspace == nil || t.Workspace.Workspace == "" {
			t.Workspace = &pb.Ref_Workspace{Workspace: "default"}
		}

		if t.Id == "" {
			id, err := ulid()
			if err != nil {
				return err
			}

			t.Id = id
		}

		return s.triggerPut(dbTxn, memTxn, t)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TriggerGet gets a trigger by reference.
func (s *State) TriggerGet(ctx context.Context, ref *pb.Ref_Trigger) (*pb.Trigger, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Trigger
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.triggerGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// TriggerDelete deletes a trigger by reference.
func (s *State) TriggerDelete(ctx context.Context, ref *pb.Ref_Trigger) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.triggerDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TriggerList returns the list of triggers.
func (s *State) TriggerList(
	ctx context.Context,
	refws *pb.Ref_Workspace,
	refproj *pb.Ref_Project,
	refapp *pb.Ref_Application,
	tagFilter []string,
) ([]*pb.Trigger, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	if refproj != nil && refws == nil {
		return nil, status.Error(codes.FailedPrecondition,
			"Workspace Ref is required when filtering on a Project")
	}

	if refapp != nil {
		if refproj == nil {
			return nil, status.Error(codes.FailedPrecondition,
				"Project Ref is required when filtering on an Application")
		}
		if refws == nil {
			return nil, status.Error(codes.FailedPrecondition,
				"Workspace Ref is required when filtering on an Application")
		}
	}

	refs, err := s.triggerList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.Trigger

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.triggerGet(dbTxn, memTxn, ref)
			if err != nil {
				return err
			}

			// filter out triggers on request
			if len(tagFilter) > 0 {
				if len(val.Tags) == 0 {
					// the trigger has no tags, so it's not a match
					continue
				}

				tagMatch := false
			EXIT:
				for _, f := range tagFilter {
					for _, tf := range val.Tags {
						if tf == f {
							tagMatch = true // we found a matching tag on this value
							// break to continue to compare ws, proj, app on value
							break EXIT
						}
					}
				}

				if !tagMatch {
					continue
				}
			}

			if refws != nil {
				if val.Workspace.Workspace != refws.Workspace {
					continue
				}

				if refproj != nil && val.Project != nil {
					if val.Project.Project != refproj.Project {
						continue
					}

					if refapp != nil && val.Application != nil {
						if val.Application.Application != refapp.Application {
							continue
						}
					}
				}
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

func (s *State) triggerPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Trigger,
) error {
	// If no name was given, set it to the id.
	if value.Name == "" {
		value.Name = value.Id
	}

	id := s.triggerId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(triggerBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.triggerIndexSet(memTxn, id, value)
}

func (s *State) triggerGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Trigger,
) (*pb.Trigger, error) {
	var result pb.Trigger
	b := dbTxn.Bucket(triggerBucket)

	if ref.Id != "" {
		s.log.Info("looking up trigger by id", "id", ref.Id)
		return &result, dbGet(b, []byte(strings.ToLower(ref.Id)), &result)
	} else {
		return nil, status.Error(codes.FailedPrecondition, "No id provided in Trigger ref to triggerGet")
	}
}

func (s *State) triggerList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_Trigger, error) {
	iter, err := memTxn.Get(triggerIndexTableName, triggerIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_Trigger
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*triggerIndexRecord)

		result = append(result, &pb.Ref_Trigger{
			Id: idx.Id,
		})
	}

	return result, nil
}

func (s *State) triggerDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Trigger,
) error {
	// Get the trigger. If it doesn't exist then we're successful.
	_, err := s.triggerGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id := s.triggerIdByRef(ref)
	if err := dbTxn.Bucket(triggerBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(triggerIndexTableName, &triggerIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// triggerIndexSet writes an index record for a single trigger.
func (s *State) triggerIndexSet(txn *memdb.Txn, id []byte, value *pb.Trigger) error {
	record := &triggerIndexRecord{
		Id:   string(id),
		Name: value.Name,
	}

	// Insert the index
	return txn.Insert(triggerIndexTableName, record)
}

// triggerIndexInit initializes the trigger index from persisted data.
func (s *State) triggerIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(triggerBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Trigger
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		// Do a minor upgrade, namely set a name if there is no name.
		if value.Name == "" {
			value.Name = value.Id
		}

		if err := s.triggerIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) triggerId(t *pb.Trigger) []byte {
	return []byte(strings.ToLower(t.Id))
}

func (s *State) triggerIdByRef(ref *pb.Ref_Trigger) []byte {
	return []byte(strings.ToLower(ref.Id))
}

func triggerIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: triggerIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			triggerIndexId: {
				Name:         triggerIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			triggerIndexName: {
				Name:         triggerIndexName,
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
	triggerIndexTableName = "trigger-index"
	triggerIndexId        = "id"
	triggerIndexName      = "name"
)

type triggerIndexRecord struct {
	Id   string
	Name string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *triggerIndexRecord) Copy() *triggerIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
