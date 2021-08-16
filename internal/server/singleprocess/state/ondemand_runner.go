package state

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var ondemandRunnerBucket = []byte("ondemandRunner")

func init() {
	dbBuckets = append(dbBuckets, ondemandRunnerBucket)
	dbIndexers = append(dbIndexers, (*State).ondemandRunnerIndexInit)
	schemas = append(schemas, ondemandRunnerIndexSchema)
}

// OndemandRunnerPut creates or updates the given ondemandRunner.
//
// Application changes will be ignored, you must use the Application APIs.
func (s *State) OndemandRunnerPut(o *pb.OndemandRunner) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	if o.Id == "" {
		id, err := ulid()
		if err != nil {
			return err
		}

		o.Id = id
	}

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.ondemandRunnerPut(dbTxn, memTxn, o)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// OndemandRunnerGet gets a ondemandRunner by reference.
func (s *State) OndemandRunnerGet(ref *pb.Ref_OndemandRunner) (*pb.OndemandRunner, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.OndemandRunner
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.ondemandRunnerGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// OndemandRunnerDelete deletes a ondemandRunner by reference. This is a complete data
// delete. This will delete all operations associated with this ondemandRunner
// as well.
func (s *State) OndemandRunnerDelete(ref *pb.Ref_OndemandRunner) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.ondemandRunnerDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// OndemandRunnerList returns the list of ondemandRunners.
func (s *State) OndemandRunnerList() ([]*pb.OndemandRunner, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.ondemandRunnerList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.OndemandRunner

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.ondemandRunnerGet(dbTxn, memTxn, ref)
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

// OndemandRunnerDefault returns the list of ondemandRunners that are defaults.
func (s *State) OndemandRunnerDefault() ([]*pb.Ref_OndemandRunner, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.ondemandRunnerDefaultRefs(memTxn)
}

func (s *State) ondemandRunnerPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.OndemandRunner,
) error {
	// This is to prevent mistakes or abuse. Realistically a waypoint.hcl
	// file should be MUCH smaller than this so this catches the really big
	// mistakes.
	if len(value.PluginConfig) > projectWaypointHclMaxSize {
		return status.Errorf(codes.FailedPrecondition,
			"ondemandRunner 'waypoint_hcl' exceeds maximum size (5MB)",
		)
	}

	id := s.ondemandRunnerId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(ondemandRunnerBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.ondemandRunnerIndexSet(memTxn, id, value)
}

func (s *State) ondemandRunnerGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_OndemandRunner,
) (*pb.OndemandRunner, error) {
	var result pb.OndemandRunner
	b := dbTxn.Bucket(ondemandRunnerBucket)

	return &result, dbGet(b, []byte(strings.ToLower(ref.Id)), &result)
}

func (s *State) ondemandRunnerList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_OndemandRunner, error) {
	iter, err := memTxn.Get(ondemandRunnerIndexTableName, ondemandRunnerIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_OndemandRunner
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*ondemandRunnerIndexRecord)

		result = append(result, &pb.Ref_OndemandRunner{
			Id: idx.Id,
		})
	}

	return result, nil
}

func (s *State) ondemandRunnerDefaultRefs(
	memTxn *memdb.Txn,
) ([]*pb.Ref_OndemandRunner, error) {
	iter, err := memTxn.Get(
		ondemandRunnerIndexTableName,
		ondemandRunnerIndexDefault+"_prefix",
		true,
		"",
	)
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_OndemandRunner
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*ondemandRunnerIndexRecord)

		result = append(result, &pb.Ref_OndemandRunner{
			Id: idx.Id,
		})
	}

	return result, nil
}

func (s *State) ondemandRunnerDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_OndemandRunner,
) error {
	// Get the ondemandRunner. If it doesn't exist then we're successful.
	_, err := s.ondemandRunnerGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id := s.ondemandRunnerIdByRef(ref)
	if err := dbTxn.Bucket(ondemandRunnerBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(ondemandRunnerIndexTableName, &ondemandRunnerIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// ondemandRunnerIndexSet writes an index record for a single ondemandRunner.
func (s *State) ondemandRunnerIndexSet(txn *memdb.Txn, id []byte, value *pb.OndemandRunner) error {
	record := &ondemandRunnerIndexRecord{
		Id:      string(id),
		Default: value.Default,
	}

	// Insert the index
	return txn.Insert(ondemandRunnerIndexTableName, record)
}

// ondemandRunnerIndexInit initializes the ondemandRunner index from persisted data.
func (s *State) ondemandRunnerIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(ondemandRunnerBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.OndemandRunner
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.ondemandRunnerIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) ondemandRunnerId(p *pb.OndemandRunner) []byte {
	return []byte(strings.ToLower(p.Id))
}

func (s *State) ondemandRunnerIdByRef(ref *pb.Ref_OndemandRunner) []byte {
	return []byte(strings.ToLower(ref.Id))
}

func ondemandRunnerIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: ondemandRunnerIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			ondemandRunnerIndexId: {
				Name:         ondemandRunnerIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			ondemandRunnerIndexDefault: {
				Name:         ondemandRunnerIndexDefault,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "Default",
						},
						&memdb.StringFieldIndex{
							Field:     "Id",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

const (
	ondemandRunnerIndexTableName = "ondemandRunner-index"
	ondemandRunnerIndexId        = "id"
	ondemandRunnerIndexDefault   = "default"

	ondemandRunnerWaypointHclMaxSize = 5 * 1024 // 5 MB
)

type ondemandRunnerIndexRecord struct {
	Id      string
	Default bool
}

// Copy should be called prior to any modifications to an existing record.
func (idx *ondemandRunnerIndexRecord) Copy() *ondemandRunnerIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
