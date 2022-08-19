package boltdbstate

import (
	"strings"

	"github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var onDemandRunnerBucket = []byte("ondemandRunner")

func init() {
	dbBuckets = append(dbBuckets, onDemandRunnerBucket)
	dbIndexers = append(dbIndexers, (*State).onDemandRunnerIndexInit)
	schemas = append(schemas, onDemandRunnerIndexSchema)
}

// OnDemandRunnerConfigPut creates or updates the given ondemandRunner.
//
// Application changes will be ignored, you must use the Application APIs.
func (s *State) OnDemandRunnerConfigPut(o *pb.OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error) {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {

		existingConfig, err := s.onDemandRunnerGet(dbTxn, memTxn, &pb.Ref_OnDemandRunnerConfig{Id: o.Id})
		if err != nil && status.Code(err) != codes.NotFound {
			return errors.Wrapf(err, "failed to check for existing on-demand runner config")
		}

		if status.Code(err) == codes.NotFound && o.Id != "" {
			return status.Errorf(codes.InvalidArgument, "cannot set the ID of a new odr profile")
		}

		if existingConfig != nil && o.Name != "" && existingConfig.Id != "" && o.Id != existingConfig.Id {
			return status.Errorf(codes.InvalidArgument, "cannot update the ID id existing runner profile with name %q", o.Name)
		}

		if o.Id != "" && status.Code(err) == codes.NotFound {
			return err
		}

		id, err := ulid()
		if err != nil {
			return err
		}
		o.Id = id

		// If no name was given, set it to the id.
		if o.Name == "" {
			o.Name = o.Id
		}

		return s.ondemandRunnerPut(dbTxn, memTxn, o)
	})
	if err == nil {
		memTxn.Commit()
	} else {
		return nil, err
	}

	ret, err := s.OnDemandRunnerConfigGet(&pb.Ref_OnDemandRunnerConfig{Id: o.Id, Name: o.Name})
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting on-demand runner config after setting it.")
	}
	return ret, nil
}

// OnDemandRunnerConfigGet gets a ondemandRunner by reference.
func (s *State) OnDemandRunnerConfigGet(ref *pb.Ref_OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.OnDemandRunnerConfig
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.onDemandRunnerGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// OnDemandRunnerConfigDelete deletes a ondemandRunner by reference. This is a complete data
// delete. This will delete all operations associated with this ondemandRunner
// as well.
func (s *State) OnDemandRunnerConfigDelete(ref *pb.Ref_OnDemandRunnerConfig) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.onDemandRunnerDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// OnDemandRunnerConfigList returns the list of ondemandRunners.
func (s *State) OnDemandRunnerConfigList() ([]*pb.OnDemandRunnerConfig, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.onDemandRunnerList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.OnDemandRunnerConfig

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.onDemandRunnerGet(dbTxn, memTxn, ref)
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

// OnDemandRunnerConfigDefault returns the list of ondemandRunners that are defaults.
func (s *State) OnDemandRunnerConfigDefault() ([]*pb.Ref_OnDemandRunnerConfig, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.onDemandRunnerDefaultRefs(memTxn)
}

func (s *State) ondemandRunnerPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.OnDemandRunnerConfig,
) error {
	// This is to prevent mistakes or abuse. Realistically a waypoint.hcl
	// file should be MUCH smaller than this so this catches the really big
	// mistakes.
	if len(value.PluginConfig) > projectWaypointHclMaxSize {
		return status.Errorf(codes.FailedPrecondition,
			"ondemandRunner 'waypoint_hcl' exceeds maximum size (5MB)",
		)
	}

	id := s.onDemandRunnerId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(onDemandRunnerBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.onDemandRunnerIndexSet(memTxn, id, value)
}

func (s *State) onDemandRunnerGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_OnDemandRunnerConfig,
) (*pb.OnDemandRunnerConfig, error) {
	var result pb.OnDemandRunnerConfig
	b := dbTxn.Bucket(onDemandRunnerBucket)

	if ref.Id != "" {
		s.log.Info("looking up ondemand runner config by id", "id", ref.Id)
		return &result, dbGet(b, []byte(strings.ToLower(ref.Id)), &result)
	}

	// Look for one by name if possible.
	if ref.Name != "" {
		s.log.Info("looking up ondemand runner config by name", "name", ref.Name)
		iter, err := memTxn.Get(
			onDemandRunnerIndexTableName,
			onDemandRunnerIndexName+"_prefix",
			ref.Name,
		)
		if err != nil {
			return nil, err
		}

		next := iter.Next()
		if next == nil {
			// Indicates that there isn't one of the given name.
			return nil, status.Errorf(codes.NotFound, "ondemand runner config not found")
		}

		idx := next.(*onDemandRunnerIndexRecord)

		return &result, dbGet(b, []byte(strings.ToLower(idx.Id)), &result)
	}

	return nil, nil
}

func (s *State) onDemandRunnerList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_OnDemandRunnerConfig, error) {
	iter, err := memTxn.Get(onDemandRunnerIndexTableName, onDemandRunnerIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_OnDemandRunnerConfig
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*onDemandRunnerIndexRecord)

		result = append(result, &pb.Ref_OnDemandRunnerConfig{
			Id: idx.Id,
		})
	}

	return result, nil
}

func (s *State) onDemandRunnerDefaultRefs(
	memTxn *memdb.Txn,
) ([]*pb.Ref_OnDemandRunnerConfig, error) {
	iter, err := memTxn.Get(
		onDemandRunnerIndexTableName,
		onDemandRunnerIndexDefault+"_prefix",
		true,
		"",
	)
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_OnDemandRunnerConfig
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*onDemandRunnerIndexRecord)

		result = append(result, &pb.Ref_OnDemandRunnerConfig{
			Id: idx.Id,
		})
	}

	return result, nil
}

func (s *State) onDemandRunnerDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_OnDemandRunnerConfig,
) error {
	// Get the ondemandRunner. If it doesn't exist then we're successful.
	_, err := s.onDemandRunnerGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id := s.onDemandRunnerIdByRef(ref)
	if err := dbTxn.Bucket(onDemandRunnerBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(onDemandRunnerIndexTableName, &onDemandRunnerIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// onDemandRunnerIndexSet writes an index record for a single ondemandRunner.
func (s *State) onDemandRunnerIndexSet(txn *memdb.Txn, id []byte, value *pb.OnDemandRunnerConfig) error {
	record := &onDemandRunnerIndexRecord{
		Id:      string(id),
		Name:    value.Name,
		Default: value.Default,
	}

	// Insert the index
	return txn.Insert(onDemandRunnerIndexTableName, record)
}

// onDemandRunnerIndexInit initializes the ondemandRunner index from persisted data.
func (s *State) onDemandRunnerIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(onDemandRunnerBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.OnDemandRunnerConfig
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		// Do a minor upgrade, namely set a name if there is no name.
		if value.Name == "" {
			value.Name = value.Id
		}

		if err := s.onDemandRunnerIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) onDemandRunnerId(p *pb.OnDemandRunnerConfig) []byte {
	return []byte(strings.ToLower(p.Id))
}

func (s *State) onDemandRunnerIdByRef(ref *pb.Ref_OnDemandRunnerConfig) []byte {
	return []byte(strings.ToLower(ref.Id))
}

func onDemandRunnerIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: onDemandRunnerIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			onDemandRunnerIndexId: {
				Name:         onDemandRunnerIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			onDemandRunnerIndexName: {
				Name:         onDemandRunnerIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Name",
					Lowercase: true,
				},
			},
			onDemandRunnerIndexDefault: {
				Name:         onDemandRunnerIndexDefault,
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
	onDemandRunnerIndexTableName = "ondemandRunner-index"
	onDemandRunnerIndexId        = "id"
	onDemandRunnerIndexName      = "name"
	onDemandRunnerIndexDefault   = "default"

	onDemandRunnerWaypointHclMaxSize = 5 * 1024 // 5 MB
)

type onDemandRunnerIndexRecord struct {
	Id      string
	Name    string
	Default bool
}

// Copy should be called prior to any modifications to an existing record.
func (idx *onDemandRunnerIndexRecord) Copy() *onDemandRunnerIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
