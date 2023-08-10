// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/hashstructure/v2"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serversort "github.com/hashicorp/waypoint/pkg/server/sort"
)

var configSourceBucket = []byte("config_source")

func init() {
	dbBuckets = append(dbBuckets, configSourceBucket)
	dbIndexers = append(dbIndexers, (*State).configSourceIndexInit)
	schemas = append(schemas, configSourceIndexSchema)
}

// ConfigSourceDelete deletes the config sourcer value from the db
func (s *State) ConfigSourceDelete(ctx context.Context, vs ...*pb.ConfigSource) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		for _, v := range vs {
			v.Delete = true
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

// ConfigSourceSet writes a set of config source values to the database.
func (s *State) ConfigSourceSet(ctx context.Context, vs ...*pb.ConfigSource) error {
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
func (s *State) ConfigSourceGet(ctx context.Context, req *pb.GetConfigSourceRequest) ([]*pb.ConfigSource, error) {
	return s.ConfigSourceGetWatch(ctx, req, nil)
}

// ConfigSourceGetWatch gets all the configuration sources for the given request.
// If a non-nil WatchSet is given, this can be watched for potential changes
// in the config source settings.
func (s *State) ConfigSourceGetWatch(ctx context.Context, req *pb.GetConfigSourceRequest, ws memdb.WatchSet) ([]*pb.ConfigSource, error) {
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
	// The scope and type of a config source is used to establish a unique record
	// in the config sources table.
	idHash, err := hashstructure.Hash(map[string]interface{}{
		"scope":     value.Scope,
		"type":      value.Type,
		"workspace": value.Workspace,
	}, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	id := s.configSourceId(idHash)

	// Write the hashed value of the config source. We use a map here so
	// that it is easy for us to add more keys to the hash.
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
	sources, err := s.configSourceGetExact(dbTxn, memTxn, ws, &pb.Ref_Global{}, req.Type)
	if err != nil {
		return nil, err
	}

	switch scope := req.Scope.(type) {
	case *pb.GetConfigSourceRequest_Global:
		// Do nothing, because we already got the global config sources

	case *pb.GetConfigSourceRequest_Project:
		// Project scope, grab our project scope vars and only those
		projectSources, err := s.configSourceGetExact(dbTxn, memTxn, ws, scope.Project, req.Type)
		if err != nil {
			return nil, err
		}

		sources = append(sources, projectSources...)

	case *pb.GetConfigSourceRequest_Application:
		projectSources, err := s.configSourceGetExact(dbTxn, memTxn, ws, &pb.Ref_Project{
			Project: scope.Application.Project,
		}, req.Type)
		if err != nil {
			return nil, err
		}

		sources = append(sources, projectSources...)

		appSources, err := s.configSourceGetExact(dbTxn, memTxn, ws, scope.Application, req.Type)
		if err != nil {
			return nil, err
		}

		sources = append(sources, appSources...)

	case *pb.GetConfigSourceRequest_All:
		sources, err = s.configSourceGetExact(dbTxn, memTxn, ws, nil, req.Type)
		if err != nil {
			return nil, err
		}
		goto SORT
	default:
		panic("unknown scope")
	}

	// Filter based on the workspace if we have it set.
	if req.Workspace != nil {
		for key, source := range sources {
			if source.Workspace != nil &&
				!strings.EqualFold(source.Workspace.Workspace, req.Workspace.Workspace) {
				sources[key] = nil
			}
		}
	}

SORT:
	var result []*pb.ConfigSource
	for _, source := range sources {
		if source != nil {
			result = append(result, source)
		}
	}
	sort.Sort(serversort.ConfigSource(result))

	return result, nil
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
	var err error
	if ref != nil {
		switch ref := ref.(type) {
		case *pb.Ref_Global:
			iter, err = memTxn.Get(
				configSourceIndexTableName,
				configSourceIndexGlobalIndexName+"_prefix",
				true,
				typeVal,
			)
			if err != nil {
				return nil, err
			}

		case *pb.Ref_Project:
			iter, err = memTxn.Get(
				configSourceIndexTableName,
				configIndexProjectIndexName+"_prefix",
				ref.Project,
				true,
				typeVal,
			)
			if err != nil {
				return nil, err
			}

		case *pb.Ref_Application:
			iter, err = memTxn.Get(
				configSourceIndexTableName,
				configSourceIndexApplicationIndexName+"_prefix",
				ref.Project,
				ref.Application,
				typeVal,
			)
			if err != nil {
				return nil, err
			}

		default:
			panic("unknown scope")
		}
	} else {
		// if no scope was passed in here, we're either getting all config
		// sources, or getting all config sources of a certain type
		if typeVal != "" {
			iter, err = memTxn.Get(
				configSourceIndexTableName,
				configSourceIndexTypeIndexName+"_prefix",
				typeVal)
		} else {
			iter, err = memTxn.Get(
				configSourceIndexTableName,
				configSourceIndexIdIndexName+"_prefix",
				"")
		}
		if err != nil {
			return nil, err
		}
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
	var projectName, applicationName string
	global := false
	project := false

	switch scope := value.Scope.(type) {
	case *pb.ConfigSource_Application:
		projectName = scope.Application.Project
		applicationName = scope.Application.Application

	case *pb.ConfigSource_Project:
		projectName = scope.Project.Project
		project = true

	case *pb.ConfigSource_Global:
		global = true

	default:
		panic("unknown scope")
	}

	record := &configSourceIndexRecord{
		Id:              string(id),
		ProjectName:     projectName,
		ApplicationName: applicationName,
		Type:            value.Type,
		Global:          global,
		Project:         project,
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

	// Regex to match for the names of any config sourcer plugins v0.10.5 and earlier
	re := regexp.MustCompile(`aws-ssm|consul|kubernetes|null|packer|terraform-cloud|vault`)

	return bucket.ForEach(func(k, v []byte) error {
		var value pb.ConfigSource
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		// If there are any records whose ID is the name of a plugin, pre-v0.10.4
		// behavior, then we delete that record from the database, since it is
		// about to be saved with a hashed ID. This is not very elegant, but
		// is the simplest way for users to upgrade since custom config sourcer
		// plugins aren't yet supported, as of 1/10/2023.
		key := string(k)
		if re.MatchString(key) || !utf8.Valid(k) {
			if err := bucket.Delete(k); err != nil {
				return err
			}
			// configSourceSet will create a new record in Bolt DB, AND update
			// Mem DB, with the ID value hashed as per the v0.10.4 hashing logic.
			if err := s.configSourceSet(dbTxn, memTxn, &value); err != nil {
				return err
			}
		} else {
			if err := s.configSourceIndexSet(memTxn, k, &value); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *State) configSourceId(idHash uint64) []byte {
	// Convert uint64 to string before conversion to byte slice. BoltDB
	// operations are OK with uint64, but server snapshots (a client of BoltDB)
	// expect that all keys and values are UTF-8 encoded, which uint64 is not.
	return []byte(fmt.Sprintf("%d", idHash))
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

			configSourceIndexTypeIndexName: {
				Name:         configSourceIndexTypeIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Type",
					Lowercase: true,
				},
			},

			configSourceIndexGlobalIndexName: {
				Name:         configSourceIndexGlobalIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "Global",
						},

						&memdb.StringFieldIndex{
							Field:     "Type",
							Lowercase: true,
						},
					},
				},
			},

			configSourceIndexProjectIndexName: {
				Name:         configSourceIndexProjectIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "ProjectName",
							Lowercase: true,
						},

						&memdb.BoolFieldIndex{
							Field: "Project",
						},

						&memdb.StringFieldIndex{
							Field:     "Type",
							Lowercase: true,
						},
					},
				},
			},

			configSourceIndexApplicationIndexName: {
				Name:         configSourceIndexApplicationIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "ProjectName",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "ApplicationName",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Type",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

const (
	configSourceIndexTableName            = "config-source-index"
	configSourceIndexIdIndexName          = "id"
	configSourceIndexTypeIndexName        = "type"
	configSourceIndexGlobalIndexName      = "global"
	configSourceIndexProjectIndexName     = "project"
	configSourceIndexApplicationIndexName = "application"
)

type configSourceIndexRecord struct {
	Id              string
	ProjectName     string
	ApplicationName string
	Type            string
	Global          bool
	Project         bool
}

// isConfigSourceDelete returns true if the config var represents a deletion.
func isConfigSourceDelete(value *pb.ConfigSource) bool {
	return value.Delete
}
