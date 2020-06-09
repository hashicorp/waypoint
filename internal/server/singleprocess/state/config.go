package state

import (
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serversort "github.com/hashicorp/waypoint/internal/server/sort"
)

var configBucket = []byte("config")

func init() {
	dbBuckets = append(dbBuckets, configBucket)
	dbIndexers = append(dbIndexers, (*State).configIndexInit)
	schemas = append(schemas, configIndexSchema)
}

// ConfigSet writes a configuration variable to the data store.
func (s *State) ConfigSet(vs ...*pb.ConfigVar) error {
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
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result []*pb.ConfigVar
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.configGetMerged(dbTxn, memTxn, req)
		return err
	})

	return result, err
}

func (s *State) configSet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.ConfigVar,
) error {
	id := s.configVarId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(configBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.configIndexSet(memTxn, id, value)
}

func (s *State) configGetMerged(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	req *pb.ConfigGetRequest,
) ([]*pb.ConfigVar, error) {
	var mergeSet [][]*pb.ConfigVar
	switch scope := req.Scope.(type) {
	case *pb.ConfigGetRequest_Project:
		// For project scope, we just return the project scoped values.
		return s.configGetExact(dbTxn, memTxn, scope.Project, req.Prefix)

	case *pb.ConfigGetRequest_Application:
		// Application scope, we have to get the project scope first
		projectVars, err := s.configGetExact(dbTxn, memTxn, &pb.Ref_Project{
			Project: scope.Application.Project,
		}, req.Prefix)
		if err != nil {
			return nil, err
		}

		// Then the application scope
		appVars, err := s.configGetExact(dbTxn, memTxn, scope.Application, req.Prefix)
		if err != nil {
			return nil, err
		}

		// Build our merge set
		mergeSet = append(mergeSet, projectVars, appVars)

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

// configIndexSet writes an index record for a single config var.
func (s *State) configIndexSet(txn *memdb.Txn, id []byte, value *pb.ConfigVar) error {
	var project, application string
	switch scope := value.Scope.(type) {
	case *pb.ConfigVar_Application:
		project = scope.Application.Project
		application = scope.Application.Application

	case *pb.ConfigVar_Project:
		project = scope.Project.Project

	default:
		panic("unknown scope")
	}

	return txn.Insert(configIndexTableName, &configIndexRecord{
		Id:          string(id),
		Project:     project,
		Application: application,
		Name:        value.Name,
	})
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
	switch scope := v.Scope.(type) {
	case *pb.ConfigVar_Application:
		return []byte(fmt.Sprintf("%s/%s/%s",
			scope.Application.Project,
			scope.Application.Application,
			v.Name,
		))

	case *pb.ConfigVar_Project:
		return []byte(fmt.Sprintf("%s/%s/%s",
			scope.Project.Project,
			"",
			v.Name,
		))

	default:
		panic("unknown scope")
	}
}

func configIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: configIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			configIndexIdIndexName: &memdb.IndexSchema{
				Name:         configIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			configIndexProjectIndexName: &memdb.IndexSchema{
				Name:         configIndexProjectIndexName,
				AllowMissing: false,
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

			configIndexApplicationIndexName: &memdb.IndexSchema{
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
		},
	}
}

const (
	configIndexTableName            = "config-index"
	configIndexIdIndexName          = "id"
	configIndexProjectIndexName     = "project"
	configIndexApplicationIndexName = "application"
)

type configIndexRecord struct {
	Id          string
	Project     string
	Application string
	Name        string
}
