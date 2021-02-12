package state

import (
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/hashstructure/v2"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var configSourceBucket = []byte("config_source")

func init() {
	dbBuckets = append(dbBuckets, configSourceBucket)
	dbIndexers = append(dbIndexers, (*State).configSourceIndexInit)
	schemas = append(schemas, configSourceIndexSchema)
}

// ConfigSourceSet writes a set of config source values to the database.
func (s *State) ConfigSourceSet(vs ...*pb.ConfigSource) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		for _, v := range vs {
			if err := s.configSourceSet(dbTxn, memTxn, v); err != nil {
				return err
			}
		}

		return nil
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ConfigSourceGet gets all the configuration sources for the given request.
func (s *State) ConfigSourceGet(req *pb.GetConfigSourceRequest) ([]*pb.ConfigSource, error) {
	return s.ConfigSourceGetWatch(req, nil)
}

// ConfigSourceGetWatch gets all the configuration sources for the given request.
// If a non-nil WatchSet is given, this can be watched for potential changes
// in the config source settings.
func (s *State) ConfigSourceGetWatch(req *pb.GetConfigSourceRequest, ws memdb.WatchSet) ([]*pb.ConfigSource, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result []*pb.ConfigSource
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.configSourceGetMerged(dbTxn, memTxn, ws, req)
		return err
	})

	return result, err
}

func (s *State) configSourceSet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.ConfigSource,
) error {
	id := s.configSourceId(value)

	// Write the hashed value of the config source. We use a map here so
	// that it is easy for us to add more keys to the hash.
	var err error
	value.Hash, err = hashstructure.Hash(map[string]interface{}{
		"config": value.Config,
	}, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(configSourceBucket)

	if isConfigSourceDelete(value) {
		if err := b.Delete(id); err != nil {
			return err
		}
	} else {
		if err := dbPut(b, id, value); err != nil {
			return err
		}
	}

	// Create our index value and write that.
	return s.configSourceIndexSet(memTxn, id, value)
}

func (s *State) configSourceGetMerged(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
	req *pb.GetConfigSourceRequest,
) ([]*pb.ConfigSource, error) {
	switch scope := req.Scope.(type) {
	case *pb.GetConfigSourceRequest_Global:
		return s.configSourceGetExact(dbTxn, memTxn, ws, scope.Global, req.Type)

	default:
		panic("unknown scope")
	}

	// In the future, when we implement scoping, this is where we'd do merging.
	// See the logic in config.go for more details.
}

// configSourceGetExact returns the list of config sources for a scope
// exactly. By "exactly" we mean without any merging logic.
func (s *State) configSourceGetExact(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
	ref interface{}, // should be one of the *pb.Ref_ values.
	typeVal string,
) ([]*pb.ConfigSource, error) {
	// We have to get the correct iterator based on the scope. We check the
	// scope and use the proper index to get the iterator here.
	var iter memdb.ResultIterator
	switch ref.(type) {
	case *pb.Ref_Global:
		var err error
		iter, err = memTxn.Get(
			configSourceIndexTableName,
			configSourceIndexIdIndexName+"_prefix",
			typeVal,
		)
		if err != nil {
			return nil, err
		}

	default:
		panic("unknown scope")
	}

	// Add to our watchset
	ws.Add(iter.WatchCh())

	// Go through the iterator and accumulate the results
	var result []*pb.ConfigSource
	b := dbTxn.Bucket(configSourceBucket)
	for {
		current := iter.Next()
		if current == nil {
			break
		}

		var value pb.ConfigSource
		record := current.(*configSourceIndexRecord)
		if err := dbGet(b, []byte(record.Id), &value); err != nil {
			return nil, err
		}

		result = append(result, &value)
	}

	return result, nil
}

// configSourceIndexSet writes an index record for a single config var.
func (s *State) configSourceIndexSet(txn *memdb.Txn, id []byte, value *pb.ConfigSource) error {
	record := &configSourceIndexRecord{
		Id:   string(id),
		Type: value.Type,
	}

	// If we have no value, we delete from the memdb index
	if isConfigSourceDelete(value) {
		return txn.Delete(configSourceIndexTableName, record)
	}

	// Insert the index
	return txn.Insert(configSourceIndexTableName, record)
}

// configSourceIndexInit initializes the config index from persisted data.
func (s *State) configSourceIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(configSourceBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.ConfigSource
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.configSourceIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) configSourceId(v *pb.ConfigSource) []byte {
	// For now the type is a unique ID. In the future when we introduce
	// scoping we'll have to do something different. This ID is only used
	// for the in-memory index so when we change this the server just needs
	// to be restarted for the fix to stick.
	return []byte(strings.ToLower(v.Type))
}

func configSourceIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: configSourceIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			configSourceIndexIdIndexName: {
				Name:         configSourceIndexIdIndexName,
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

const (
	configSourceIndexTableName   = "config-source-index"
	configSourceIndexIdIndexName = "id"
)

type configSourceIndexRecord struct {
	Id   string
	Type string
}

// isConfigSourceDelete returns true if the config var represents a deletion.
func isConfigSourceDelete(value *pb.ConfigSource) bool {
	return value.Delete
}
