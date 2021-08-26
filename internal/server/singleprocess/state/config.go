package state

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/hashstructure/v2"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serversort "github.com/hashicorp/waypoint/internal/server/sort"
)

var configBucket = []byte("config")

func init() {
	dbBuckets = append(dbBuckets, configBucket)
	dbIndexers = append(dbIndexers, (*State).configIndexInit)
	schemas = append(schemas, configIndexSchema)
}

// ConfigSet writes a configuration variable to the data store. Deletes
// are always ordered before writes so that you can't delete a new value
// in a single ConfigSet (you should never want to do this).
func (s *State) ConfigSet(vs ...*pb.ConfigVar) error {
	// Sort the variables so that deletes are handled before writes.
	// TODO: test this
	sort.Slice(vs, func(i, j int) bool {
		// i < j if i is a delete request.
		return isConfigVarDelete(vs[i])
	})

	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		for _, v := range vs {
			if err := s.configSet(dbTxn, memTxn, v); err != nil {
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

// ConfigGet gets all the configuration for the given request.
func (s *State) ConfigGet(req *pb.ConfigGetRequest) ([]*pb.ConfigVar, error) {
	return s.ConfigGetWatch(req, nil)
}

// ConfigGetWatch gets all the configuration for the given request. If a non-nil
// WatchSet is given, this can be watched for potential changes in the config.
func (s *State) ConfigGetWatch(req *pb.ConfigGetRequest, ws memdb.WatchSet) ([]*pb.ConfigVar, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result []*pb.ConfigVar
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.configGetMerged(dbTxn, memTxn, ws, req)
		return err
	})

	return result, err
}

func (s *State) configSet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.ConfigVar,
) error {
	s.configVarNormalize(value)
	id := s.configVarId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(configBucket)

	if isConfigVarDelete(value) {
		if err := b.Delete(id); err != nil {
			return err
		}
	} else {
		// If this is a runner, we don't support dynamic values currently.
		if value.Target.Runner != nil {
			if _, ok := value.Value.(*pb.ConfigVar_Static); !ok {
				return status.Errorf(codes.FailedPrecondition,
					"runner-scoped configuration must be static")
			}
		}

		if err := dbPut(b, id, value); err != nil {
			return err
		}
	}

	// Create our index value and write that.
	return s.configIndexSet(memTxn, id, value)
}

func (s *State) configGetMerged(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
	req *pb.ConfigGetRequest,
) ([]*pb.ConfigVar, error) {
	var mergeSet [][]*pb.ConfigVar
	switch scope := req.Scope.(type) {
	case *pb.ConfigGetRequest_Project:
		// For project scope, we just return the project scoped values.
		return s.configGetExact(dbTxn, memTxn, ws, scope.Project, req.Prefix)

	case *pb.ConfigGetRequest_Application:
		// Application scope, we have to get the project scope first
		projectVars, err := s.configGetExact(dbTxn, memTxn, ws, &pb.Ref_Project{
			Project: scope.Application.Project,
		}, req.Prefix)
		if err != nil {
			return nil, err
		}

		// Then the application scope
		appVars, err := s.configGetExact(dbTxn, memTxn, ws, scope.Application, req.Prefix)
		if err != nil {
			return nil, err
		}

		// Build our merge set
		mergeSet = append(mergeSet, projectVars, appVars)

	case *pb.ConfigGetRequest_Runner:
		var err error
		mergeSet, err = s.configGetRunner(dbTxn, memTxn, ws, scope.Runner, req.Prefix)
		if err != nil {
			return nil, err
		}

	default:
		panic("unknown scope")
	}

	// Merge our merge set
	merged := make(map[string]*pb.ConfigVar)
	for _, set := range mergeSet {
		for _, v := range set {
			merged[v.Name] = v
		}
	}

	result := make([]*pb.ConfigVar, 0, len(merged))
	for _, v := range merged {
		result = append(result, v)
	}

	sort.Sort(serversort.ConfigName(result))

	return result, nil
}

// configGetExact returns the list of config variables for a scope
// exactly. By "exactly" we mean without any merging logic: if you request
// app-scoped variables, you'll get app-scoped variables. If a project-scoped
// variable matches, it will not be merged in.
func (s *State) configGetExact(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
	ref interface{}, // should be one of the *pb.Ref_ values.
	prefix string,
) ([]*pb.ConfigVar, error) {
	// We have to get the correct iterator based on the scope. We check the
	// scope and use the proper index to get the iterator here.
	var iter memdb.ResultIterator
	switch ref := ref.(type) {
	case *pb.Ref_Application:
		var err error
		iter, err = memTxn.Get(
			configIndexTableName,
			configIndexApplicationIndexName+"_prefix",
			ref.Project,
			ref.Application,
			prefix,
		)
		if err != nil {
			return nil, err
		}

	case *pb.Ref_Project:
		var err error
		iter, err = memTxn.Get(
			configIndexTableName,
			configIndexProjectIndexName+"_prefix",
			ref.Project,
			prefix,
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
	var result []*pb.ConfigVar
	b := dbTxn.Bucket(configBucket)
	for {
		current := iter.Next()
		if current == nil {
			break
		}

		var value pb.ConfigVar
		record := current.(*configIndexRecord)
		if err := dbGet(b, []byte(record.Id), &value); err != nil {
			return nil, err
		}

		result = append(result, &value)
	}

	return result, nil
}

// configGetRunner gets the config vars for a runner.
func (s *State) configGetRunner(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
	req *pb.Ref_RunnerId,
	prefix string,
) ([][]*pb.ConfigVar, error) {
	iter, err := memTxn.Get(
		configIndexTableName,
		configIndexRunnerIndexName+"_prefix",
		true,
		prefix,
	)
	if err != nil {
		return nil, err
	}

	// Add to our watch set
	ws.Add(iter.WatchCh())

	// Results go into two buckets
	result := make([][]*pb.ConfigVar, 2)
	const (
		idxAny = 0
		idxId  = 1
	)

	// Go through the iterator and accumulate the results
	b := dbTxn.Bucket(configBucket)
	for {
		current := iter.Next()
		if current == nil {
			break
		}
		record := current.(*configIndexRecord)

		idx := -1
		switch ref := record.RunnerRef.Target.(type) {
		case *pb.Ref_Runner_Any:
			idx = idxAny

		case *pb.Ref_Runner_Id:
			idx = idxId

			// We need to match this ID
			if ref.Id.Id != req.Id {
				continue
			}

		default:
			return nil, fmt.Errorf("config has unknown target type: %T", record.RunnerRef.Target)
		}

		var value pb.ConfigVar
		if err := dbGet(b, []byte(record.Id), &value); err != nil {
			return nil, err
		}

		result[idx] = append(result[idx], &value)
	}

	return result, nil
}

// configIndexSet writes an index record for a single config var.
func (s *State) configIndexSet(txn *memdb.Txn, id []byte, value *pb.ConfigVar) error {
	var project, application string
	switch scope := value.Target.AppScope.(type) {
	case *pb.ConfigVar_Target_Application:
		project = scope.Application.Project
		application = scope.Application.Application

	case *pb.ConfigVar_Target_Project:
		project = scope.Project.Project

	case *pb.ConfigVar_Target_Global:
		// Global scope set nothing

	default:
		panic("unknown scope")
	}

	record := &configIndexRecord{
		Id:          string(id),
		Project:     project,
		Application: application,
		Name:        value.Name,
		Runner:      value.Target.Runner != nil,
		RunnerRef:   value.Target.Runner,
	}

	// If we have no value, we delete from the memdb index
	if isConfigVarDelete(value) {
		return txn.Delete(configIndexTableName, record)
	}

	// Insert the index
	return txn.Insert(configIndexTableName, record)
}

// configIndexInit initializes the config index from persisted data.
func (s *State) configIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(configBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.ConfigVar
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.configIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) configVarId(v *pb.ConfigVar) []byte {
	// It is very important the variable is normalized before we get the ID.
	s.configVarNormalize(v)

	// Hash the ConfigVar Target to get a unique ID for this variable.
	// This works because hashstructure ignores unexported fields.
	hash, err := hashstructure.Hash(v.Target, hashstructure.FormatV2, nil)
	if err != nil {
		// This should never happen so we panic so we can report a bug.
		// This _may_ crash the server which is not ideal, but we usually
		// do per-request recover() and this usually isn't in a goroutine.
		// And again, it shouldn't happen.
		panic(err)
	}

	return []byte(strings.ToLower(fmt.Sprintf("%s:%d/%s", configIdPrefix, hash, v.Name)))
}

// configVarNormalize takes a ConfigVar and "normalizes" it by applying
// transforms to ensure config vars have the same structure. For example,
// this auto-upgrades older fields in the protobuf.
//
// This is safe to call multiple times.
func (s *State) configVarNormalize(v *pb.ConfigVar) {
	// UnusedScope is from WP 0.5 and earlier (and was named "scope" then).
	// We upgrade it to target. If its set we always overwrite the target.
	if v.UnusedScope != nil {
		v.Target = &pb.ConfigVar_Target{}

		switch scope := v.UnusedScope.(type) {
		case *pb.ConfigVar_Application:
			v.Target.AppScope = &pb.ConfigVar_Target_Application{
				Application: scope.Application,
			}

		case *pb.ConfigVar_Project:
			v.Target.AppScope = &pb.ConfigVar_Target_Project{
				Project: scope.Project,
			}

		case *pb.ConfigVar_Runner:
			// In WP 0.5 and earlier, a scope of runner implied a global
			// runner variable no matter where the runner lived.
			v.Target.AppScope = &pb.ConfigVar_Target_Global{
				Global: &pb.Ref_Global{},
			}
			v.Target.Runner = scope.Runner
		}

		// Set our scope to nil so that we don't try to use it.
		v.UnusedScope = nil
	}
}

func configIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: configIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			configIndexIdIndexName: {
				Name:         configIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			configIndexProjectIndexName: {
				Name:         configIndexProjectIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Name",
							Lowercase: true,
						},
					},
				},
			},

			configIndexApplicationIndexName: {
				Name:         configIndexApplicationIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Application",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Name",
							Lowercase: true,
						},
					},
				},
			},

			configIndexRunnerIndexName: {
				Name:         configIndexRunnerIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "Runner",
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
	configIndexTableName            = "config-index"
	configIndexIdIndexName          = "id"
	configIndexProjectIndexName     = "project"
	configIndexApplicationIndexName = "application"
	configIndexRunnerIndexName      = "runner"

	// configIdPrefix prefixes our ID so that we can change the
	// ID hashing in the future if we want to.
	configIdPrefix = "v2:"
)

type configIndexRecord struct {
	Id          string
	Project     string
	Application string
	Name        string
	Runner      bool // true if this is a runner config
	RunnerRef   *pb.Ref_Runner
}

// isConfigVarDelete returns true if the config var represents a deletion.
func isConfigVarDelete(value *pb.ConfigVar) bool {
	switch v := value.Value.(type) {
	case *pb.ConfigVar_Unset:
		return true

	case *pb.ConfigVar_Static:
		return v.Static == ""

	case nil:
		return true
	}

	return false
}
