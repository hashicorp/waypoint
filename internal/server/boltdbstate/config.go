package boltdbstate

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/mitchellh/pointerstructure"
	"github.com/zclconf/go-cty/cty"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/waypoint/pkg/config/funcs"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serversort "github.com/hashicorp/waypoint/pkg/server/sort"
)

var (
	configBucket = []byte("config_v2")

	// configBucketOld is the name of the config bucket in WP versions
	// 0.5 and earlier. We changed the bucket name so that we can safely
	// upgrade the data. We should not use the "config" bucket name again
	// unless we are sure we've migrated all data.
	configBucketOld = []byte("config")
)

func init() {
	dbBuckets = append(dbBuckets, configBucket)
	dbIndexers = append(dbIndexers, (*State).configIndexInit)
	schemas = append(schemas, configIndexSchema)
}

// ConfigDelete looks at each passed in ConfigVar and checks to see if it has
// been properly set for deletion either through ConfigVar_Unset or when its
// static value i.e. ConfigVar_Static is empty string. It then calls through
// to configSet to delete it.
func (s *State) ConfigDelete(ctx context.Context, vs ...*pb.ConfigVar) error {
	for i, v := range vs {
		if !isConfigVarDelete(vs[i]) {
			return fmt.Errorf("config var is not set to delete: %s. "+
				"Its value must be set to empty string if static or ConfigVar_Unset", v.Name)
		}
	}

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

// ConfigSet writes a configuration variable to the data store. Deletes
// are always ordered before writes so that you can't delete a new value
// in a single ConfigSet (you should never want to do this).
func (s *State) ConfigSet(ctx context.Context, vs ...*pb.ConfigVar) error {
	// Sort the variables so that deletes are handled before writes.
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
func (s *State) ConfigGet(ctx context.Context, req *pb.ConfigGetRequest) ([]*pb.ConfigVar, error) {
	return s.ConfigGetWatch(ctx, req, nil)
}

// ConfigGetWatch gets all the configuration for the given request. If a non-nil
// WatchSet is given, this can be watched for potential changes in the config.
func (s *State) ConfigGetWatch(ctx context.Context, req *pb.ConfigGetRequest, ws memdb.WatchSet) ([]*pb.ConfigVar, error) {
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
	// Always get our global values
	globalVars, err := s.configGetExact(dbTxn, memTxn, ws, &pb.Ref_Global{}, req.Prefix)
	if err != nil {
		return nil, err
	}

	// mergeSet is the set of variables we'll merge and resolve.
	mergeSet := [][]*pb.ConfigVar{globalVars}
	merge := true

	switch scope := req.Scope.(type) {
	case *pb.ConfigGetRequest_Project:
		// Project scope, grab our project scope vars and only those
		projectVars, err := s.configGetExact(dbTxn, memTxn, ws, scope.Project, req.Prefix)
		if err != nil {
			return nil, err
		}

		// Project scope never resolves so we just return it as is. We
		// set the "merge = false" flag to notify that we don't want to do
		// resolution.
		mergeSet = append(mergeSet, projectVars)
		merge = false

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

	case nil:
		// For nil scope we just look at the global variables.
		//
		// Note this is IMPORTANT for backwards compatibility, since when we
		// introduced app-scoped runner vars, we pulled the "runner" out of the
		// oneof, allowing old clients to send a nil scope. So we must handle
		// this case so long as WP 0.5 and earlier clients exist, at least.

	default:
		panic("unknown scope")
	}

	// Sort all of our merge sets by the resolution rules
	for _, set := range mergeSet {
		sort.Sort(serversort.ConfigResolution(set))
	}

	// If we have a runner set, then we want to filter all our config vars
	// by runner. This is more complex than that though, because tighter
	// scoped runner refs should overwrite weaker scoped (i.e. ID-ref overwrites
	// Any-ref). So we have to split our merge set from <X, Y> to
	// <X_any, X_id, Y_any, Y_id> so it merges properly later.
	if req.Runner != nil {
		var newMergeSet [][]*pb.ConfigVar
		for _, set := range mergeSet {
			splitSets, err := s.configRunnerSet(set, req.Runner)
			if err != nil {
				return nil, err
			}

			newMergeSet = append(newMergeSet, splitSets...)
		}

		mergeSet = newMergeSet
	} else {
		// If runner isn't set, then we want to ensure we're not getting
		// any runner env vars.
		for _, set := range mergeSet {
			for i, v := range set {
				if v == nil {
					continue
				}

				if v.Target.Runner != nil {
					set[i] = nil
				}
			}
		}
	}

	// Filter based on the workspace if we have it set.
	if req.Workspace != nil {
		for _, set := range mergeSet {
			for i, v := range set {
				if v == nil {
					continue
				}

				if v.Target.Workspace != nil &&
					!strings.EqualFold(v.Target.Workspace.Workspace, req.Workspace.Workspace) {
					set[i] = nil
				}
			}
		}
	}

	// Filter by labels
	ctyMap := cty.MapValEmpty(cty.String)
	if len(req.Labels) > 0 {
		mapValues := map[string]cty.Value{}
		for k, v := range req.Labels {
			mapValues[k] = cty.StringVal(v)
		}
		ctyMap = cty.MapVal(mapValues)
	}

	for _, set := range mergeSet {
		for i, v := range set {
			if v == nil {
				continue
			}

			// If there is no selector, ignore.
			if v.Target.LabelSelector == "" {
				continue
			}

			// Use our selectormatch HCL function for equal logic
			result, err := funcs.SelectorMatch(ctyMap, cty.StringVal(v.Target.LabelSelector))
			if errors.Is(err, pointerstructure.ErrNotFound) {
				// this means that the label selector contains a label
				// that isn't set, this means we do not match.
				err = nil
				result = cty.BoolVal(false)
			}
			if err != nil {
				return nil, err
			}

			if result.False() {
				set[i] = nil
			}
		}
	}

	// If we aren't merging, then we're done. We just flatten the list.
	if !merge {
		var result []*pb.ConfigVar
		for _, set := range mergeSet {
			for _, v := range set {
				if v != nil {
					result = append(result, v)
				}
			}
		}
		sort.Sort(serversort.ConfigName(result))
		return result, nil
	}

	// Merge our merge set
	merged := make(map[string]*pb.ConfigVar)
	for _, set := range mergeSet {
		for _, v := range set {
			// Ignore nil since those are filtered out values.
			if v == nil {
				continue
			}

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
	case *pb.Ref_Global:
		var err error
		iter, err = memTxn.Get(
			configIndexTableName,
			configIndexGlobalIndexName+"_prefix",
			true,
			prefix,
		)
		if err != nil {
			return nil, err
		}

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

// configRunnerSet splits a set of config vars into a merge set depending
// on priority to match a runner.
func (s *State) configRunnerSet(
	set []*pb.ConfigVar,
	req *pb.Ref_RunnerId,
) ([][]*pb.ConfigVar, error) {
	// Results go into two buckets
	result := make([][]*pb.ConfigVar, 2)
	const (
		idxAny = 0
		idxId  = 1
	)

	// Go through the iterator and accumulate the results
	for _, current := range set {
		if current.Target.Runner == nil {
			// We are not a config for a runner.
			continue
		}

		idx := -1
		switch ref := current.Target.Runner.Target.(type) {
		case *pb.Ref_Runner_Any:
			idx = idxAny

		case *pb.Ref_Runner_Id:
			idx = idxId

			// We need to match this ID
			if ref.Id.Id != req.Id {
				continue
			}

		default:
			return nil, fmt.Errorf("config has unknown target type: %T", current.Target.Runner.Target)
		}

		result[idx] = append(result[idx], current)
	}

	return result, nil
}

// configIndexSet writes an index record for a single config var.
func (s *State) configIndexSet(txn *memdb.Txn, id []byte, value *pb.ConfigVar) error {
	var project, application string
	global := false
	switch scope := value.Target.AppScope.(type) {
	case *pb.ConfigVar_Target_Application:
		project = scope.Application.Project
		application = scope.Application.Application

	case *pb.ConfigVar_Target_Project:
		project = scope.Project.Project

	case *pb.ConfigVar_Target_Global:
		global = true

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
		Global:      global,
	}

	// If we have no value, we delete from the memdb index
	if isConfigVarDelete(value) {
		err := txn.Delete(configIndexTableName, record)
		if errors.Is(err, memdb.ErrNotFound) {
			// If it doesn't exist that is okay
			err = nil
		}

		return err
	}

	// Insert the index
	return txn.Insert(configIndexTableName, record)
}

// configIndexInit initializes the config index from persisted data.
func (s *State) configIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	// If we have the old (WP 0.5 and earlier) bucket, then perform an
	// upgrade over all the items in this bucket.
	oldBucket := dbTxn.Bucket(configBucketOld)
	if oldBucket != nil {
		err := oldBucket.ForEach(func(k, v []byte) error {
			var value pb.ConfigVar
			if err := proto.Unmarshal(v, &value); err != nil {
				return err
			}

			// configSet normalizes the old to new format and writes it
			// into our new bucket.
			if err := s.configSet(dbTxn, memTxn, &value); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		// Delete the old bucket. We shouldn't fail since we just verified
		// this key exists and it is a bucket.
		if err := dbTxn.DeleteBucket(configBucketOld); err != nil {
			return err
		}
	}

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

			configIndexGlobalIndexName: {
				Name:         configIndexGlobalIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "Global",
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
	configIndexGlobalIndexName      = "global"

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
	Global      bool
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
