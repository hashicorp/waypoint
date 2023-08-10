// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// appOperation is an abstraction on any "operation" that may happen to
// an app such as a build, push, etc. This allows uniform API calls on
// top of operations at a basic level.
type appOperation struct {
	// Struct is the record structure used for this operation. Struct is
	// expected to have the following fields with the following types. The
	// names and types must match exactly.
	//
	//   - Id string
	//   - Status *pb.Status
	//   - Application *pb.Ref_Application
	//
	// It may also have the special field "Preload". If this field exists,
	// it is automatically set to nil on disk and set to empty on read. This
	// field is expected to be used for just-in-time data loading that is not
	// persisted.
	//
	Struct interface{}

	// Bucket is the global bucket for all records of this operation.
	Bucket []byte

	// The number of records that should be indexed off disk. This allows
	// dormant records to remain on disk but not indexed.
	MaximumIndexedRecords int

	// This guards indexedRecords for manipulation during pruning
	pruneMu sync.Mutex

	// Holds how many records we've indexed at runtime.
	indexedRecords int

	// seq is the previous sequence number to set. This is initialized by the
	// index init on server boot and `sync/atomic` should be used to increment
	// it on each use.
	//
	// NOTE(mitchellh): we currently never prune this map. We should do that
	// once we implement Delete.
	seq map[string]map[string]*uint64
}

// Test validates that the operation struct is setup properly. This
// is expected to be called in a unit test.
func (op *appOperation) Test(t testing.T) {
	require := require.New(t)

	// Validate the struct is a struct
	typ := reflect.TypeOf(op.Struct)
	require.Equal(reflect.Ptr, typ.Kind())
	require.Equal(reflect.Struct, typ.Elem().Kind())

	// Fields we need
	v := op.newStruct()
	{
		field := op.valueField(v, "Id")
		require.NotNil(field)
		require.IsType("", field)
	}
	{
		field := op.valueField(v, "Status")
		require.IsType((*pb.Status)(nil), field)
	}
	{
		field := op.valueField(v, "Application")
		require.IsType((*pb.Ref_Application)(nil), field)
	}
	{
		field := op.valueField(v, "Workspace")
		require.IsType((*pb.Ref_Workspace)(nil), field)
	}
}

// register should be called in init() to register this operation with
// all the proper global variables to setup the state for this operation.
func (op *appOperation) register() {
	dbBuckets = append(dbBuckets, op.Bucket)
	dbIndexers = append(dbIndexers, op.indexInit)
	schemas = append(schemas, op.memSchema)

	if op.MaximumIndexedRecords > 0 {
		pruneFns = append(pruneFns, func(memTxn *memdb.Txn) (string, int, error) {
			cnt, err := op.pruneOld(memTxn, op.MaximumIndexedRecords)
			return op.memTableName(), cnt, err
		})
	}
}

// Put inserts or updates an operation record.
func (op *appOperation) Put(s *State, update bool, value proto.Message) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return op.dbPut(s, dbTxn, memTxn, update, value)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// Get gets an operation record by reference.
func (op *appOperation) Get(s *State, ref *pb.Ref_Operation) (interface{}, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	result := op.newStruct()
	err := s.db.View(func(tx *bolt.Tx) error {
		var id string
		switch t := ref.Target.(type) {
		case *pb.Ref_Operation_Id:
			id = t.Id

		case *pb.Ref_Operation_Sequence:
			var err error
			id, err = op.getIdForSeq(s, tx, memTxn, t.Sequence)
			if err != nil {
				return err
			}

		default:
			return status.Errorf(codes.FailedPrecondition,
				"unknown operation reference type: %T", ref.Target)
		}

		// Check if we are tracking this value in the indexes before returning
		// it. When pruning, we leave the values on disk but remove them
		// from the indexes.
		raw, err := memTxn.First(
			op.memTableName(),
			opIdIndexName,
			id,
		)
		if err != nil {
			return err
		}

		if raw == nil {
			return status.Errorf(codes.NotFound,
				"value with given id not found: %s", id)
		}

		return op.get(s, tx, []byte(id), result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes an operation record by reference
func (op *appOperation) Delete(s *State, value proto.Message) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return op.delete(dbTxn, memTxn, value)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

func (op *appOperation) delete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value proto.Message,
) error {
	id := op.valueField(value, "Id").(string)
	// Delete the value from the bucket
	if err := op.dbDelete(dbTxn, value, []byte(id)); err != nil {
		return err
	}
	if err := memTxn.Delete(op.memTableName(), &operationIndexRecord{Id: string(id)}); err != nil {
		return err
	}
	return nil
}

// dbDelete deletes the value from the database
func (op *appOperation) dbDelete(
	dbTxn *bolt.Tx,
	value proto.Message,
	id []byte,
) error {
	// Get our application
	appRef := op.valueField(value, "Application").(*pb.Ref_Application)
	if appRef == nil {
		return status.Errorf(codes.Internal, "state: Application must be set on value %T", value)
	}

	var seq uint64
	if f := op.valueFieldReflect(value, "Sequence"); f.IsValid() {
		// Subtract 1 from the sequence # because we're deleting an operation
		seq = atomic.AddUint64(op.appSeq(appRef), ^uint64(1-1))
		f.Set(reflect.ValueOf(seq))
	}

	b := dbTxn.Bucket(op.Bucket)
	return b.Delete(id)
}

func (op *appOperation) getIdForSeq(
	s *State,
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_OperationSeq,
) (string, error) {
	raw, err := memTxn.First(
		op.memTableName(),
		opSeqIndexName,
		ref.Application.Project,
		ref.Application.Application,
		ref.Number,
	)
	if err != nil {
		return "", err
	}
	if raw == nil {
		return "", status.Errorf(codes.NotFound,
			"not found for sequence id number %d, application %q, and project %q", ref.Number, ref.Application.Application, ref.Application.Project)
	}

	idx := raw.(*operationIndexRecord)
	return idx.Id, nil
}

// List lists all the records.
func (op *appOperation) List(
	s *State, opts *serverstate.ListOperationOptions) ([]interface{}, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	// Set the proper index for our ordering
	idx := opStartTimeIndexName
	if opts.Order != nil {
		switch opts.Order.Order {
		case pb.OperationOrder_COMPLETE_TIME:
			idx = opCompleteTimeIndexName
		}
	}

	// We need an application to do this, so return an error
	// to avoid a panic. This should not commonly be a user-facing
	// error
	if opts.Application == nil {
		return nil, errors.New("must provide an Application.Ref to List")
	}

	// Get the iterator for lower-bound based querying
	iter, err := memTxn.LowerBound(
		op.memTableName(),
		idx,
		opts.Application.Project,
		opts.Application.Application,
		indexTimeLatest{},
	)
	if err != nil {
		return nil, err
	}

	var result []interface{}

	s.db.View(func(tx *bolt.Tx) error {
		for {
			current := iter.Next()
			if current == nil {
				return nil
			}

			record := current.(*operationIndexRecord)
			if !record.MatchRef(opts.Application) {
				return nil
			}

			// If our workspace doesn't match then continue to the next result.
			if opts.Workspace != nil && record.Workspace != opts.Workspace.Workspace {
				continue
			}

			value := op.newStruct()
			if err := op.get(s, tx, []byte(record.Id), value); err != nil {
				return err
			}

			if opts.PhysicalState > 0 {
				if raw := op.valueField(value, "State"); raw != nil {
					state := raw.(pb.Operation_PhysicalState)
					if state != opts.PhysicalState {
						continue
					}
				}
			}

			if len(opts.Status) > 0 {
				// Get our status field
				status := op.valueField(value, "Status").(*pb.Status)

				// Filter. If we don't match the filter, then ignore this result.
				if !statusFilterMatch(opts.Status, status) {
					continue
				}
			}

			result = append(result, value)

			// If we have a limit, check that now
			if o := opts.Order; o != nil && o.Limit > 0 && len(result) >= int(o.Limit) {
				return nil
			}
		}
	})

	// We can derive a watchch by doing a get on the "first" prefix, discarding the result, and
	// just use the watchch.
	if opts.WatchSet != nil {
		watchCh, _, err := memTxn.FirstWatch(
			op.memTableName(),
			opSeqIndexName,
			opts.Application.Project,
			opts.Application.Application,
			uint(0),
		)
		if err != nil {
			return nil, err
		}

		opts.WatchSet.Add(watchCh)
	}

	return result, nil
}

// LatestFilter gets the latest operation that was completed successfully
// and matches the given filter. This works by iterating over the operations
// in most-recently-completed order, so if you specify a filter that rarely is
// true, this may require effectively a table scan.
func (op *appOperation) LatestFilter(
	s *State,
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
	filter func(interface{}) (bool, error),
) (interface{}, error) {
	// If we have no filter, create a filter that always returns true.
	if filter == nil {
		filter = func(interface{}) (bool, error) { return true, nil }
	}

	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.LowerBound(
		op.memTableName(),
		opCompleteTimeIndexName,
		ref.Project,
		ref.Application,
		indexTimeLatest{},
	)
	if err != nil {
		return nil, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}

		record := raw.(*operationIndexRecord)
		if !record.MatchRef(ref) {
			break
		}

		// If our workspace doesn't match then continue to the next result.
		if ws != nil && record.Workspace != ws.Workspace {
			continue
		}

		v, err := op.Get(s, &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{Id: record.Id},
		})
		if err != nil {
			return nil, err
		}

		// Shouldn't happen but if it does, return nothing.
		st := op.valueField(v, "Status")
		if st == nil {
			break
		}

		// State must be success.
		if st.(*pb.Status).State != pb.Status_SUCCESS {
			continue
		}

		// If we have no filter, return it
		filterResult, err := filter(v)
		if err != nil {
			return nil, err
		}

		if filterResult {
			return v, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "No application named %q is available, or application has no successful operations", ref.Application)
}

// Latest gets the latest operation that was completed successfully.
func (op *appOperation) Latest(
	s *State,
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (interface{}, error) {
	return op.LatestFilter(s, ref, ws, nil)
}

// get reads the value from the database. This populates any computed fields
// as necessary, unlike dbGet which just pulls the raw value. This should
// be used instead of dbGet in most cases.
func (op *appOperation) get(
	s *State,
	dbTxn *bolt.Tx,
	id []byte,
	result proto.Message,
) error {
	// Get the value
	if err := op.dbGet(dbTxn, id, result); err != nil {
		return err
	}

	// If we have a preload field then we check if we have to prepopulate.
	if f := op.valueFieldReflect(result, "Preload"); f.IsValid() {
		preloadIface := f.Interface()

		// If we have a job data source ref field, then we attempt to
		// load the job and populate this.
		jobIdRaw := op.valueField(result, "JobId")
		dsrefF := op.valueFieldReflect(preloadIface, "JobDataSourceRef")
		if jobIdRaw != nil && dsrefF.IsValid() {
			job, err := s.jobById(dbTxn, jobIdRaw.(string))

			// We ignore not found by simply not populating the field.
			if status.Code(err) == codes.NotFound {
				err = nil
				job = nil
			}

			// Any other error we return back to the user.
			if err != nil {
				return err
			}

			// If we found the job, we set it.
			if job != nil {
				dsrefF.Set(reflect.ValueOf(job.DataSourceRef))
			}
		}
	}

	return nil
}

// dbGet reads the value from the database.
func (op *appOperation) dbGet(
	dbTxn *bolt.Tx,
	id []byte,
	result proto.Message,
) error {
	// Read the value
	if err := dbGet(dbTxn.Bucket(op.Bucket), []byte(id), result); err != nil {
		return err
	}

	// If there is a preload field, we want to set that to non-nil.
	if f := op.valueFieldReflect(result, "Preload"); f.IsValid() {
		f.Set(reflect.New(f.Type().Elem()))
	}

	return nil
}

// dbPut wites the value to the database and also sets up any index records.
// It expects to hold a write transaction to both bolt and memdb.
func (op *appOperation) dbPut(
	s *State,
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	update bool,
	value proto.Message,
) error {
	// Get our application and ensure it is created
	appRef := op.valueField(value, "Application").(*pb.Ref_Application)
	if appRef == nil {
		return status.Errorf(codes.Internal, "state: Application must be set on value %T", value)
	}
	if err := s.appPut(dbTxn, memTxn, s.appDefaultForRef(appRef)); err != nil {
		return err
	}

	// Get our workspace reference and ensure it is created
	wsRef := op.valueField(value, "Workspace").(*pb.Ref_Workspace)
	if wsRef == nil {
		return status.Errorf(codes.Internal, "state: Workspace must be set on value %T", value)
	}

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(op.Bucket)

	id := []byte(op.valueField(value, "Id").(string))
	if update {
		// Load the value so that we can retain the values that are read-only.
		// At the same time we verify it exists
		existing := op.newStruct()
		err := op.dbGet(dbTxn, []byte(id), existing)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return status.Errorf(codes.NotFound, "record with ID %q not found for update", string(id))
			}

			return err
		}

		// Next, ensure that the fields we want to match are matched.
		matchFields := []string{"Generation", "Sequence"}
		for _, name := range matchFields {
			f := op.valueFieldReflect(value, name)
			if !f.IsValid() {
				continue
			}

			fOld := op.valueFieldReflect(existing, name)
			if !fOld.IsValid() {
				continue
			}

			f.Set(fOld)
		}
	}

	if !update {
		// If we're not updating, then set the sequence number up if we have one.
		var seq uint64
		if f := op.valueFieldReflect(value, "Sequence"); f.IsValid() {
			seq = atomic.AddUint64(op.appSeq(appRef), 1)
			f.Set(reflect.ValueOf(seq))
		}

		if f := op.valueFieldReflect(value, "Generation"); f.IsValid() {
			gen := f.Interface().(*pb.Generation)

			// Default the generation to a new ULID if it isn't set.
			if gen == nil || gen.Id == "" {
				v, err := ulid()
				if err != nil {
					return err
				}

				gen = &pb.Generation{Id: v}
				f.Set(reflect.ValueOf(gen))
			}

			// Our initial sequence is always our current to start. But
			// if we can find an older version, we will update it.
			gen.InitialSequence = seq

			// Set our initial sequence number by searching the history
			// to the first operation that used this generation.
			iter, err := memTxn.LowerBound(
				op.memTableName(),
				opGenIndexName,
				appRef.Project,
				appRef.Application,
				gen.Id,
				uint64(0),
			)
			if err != nil {
				return err
			}
			if raw := iter.Next(); raw != nil {
				idx := raw.(*operationIndexRecord)
				if idx.MatchRef(appRef) &&
					idx.Generation == gen.Id {
					gen.InitialSequence = idx.Sequence
				}
			}
		}
	}

	// If there is a preload field, we want to set that to nil.
	if f := op.valueFieldReflect(value, "Preload"); f.IsValid() {
		f.Set(reflect.New(f.Type().Elem()))
	}

	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	indexRec, err := op.indexPut(s, memTxn, value)
	if err != nil {
		return err
	}

	// Ensure the workspace record is updated.
	touchTime := indexRec.StartTime
	if v := indexRec.CompleteTime; v.After(touchTime) {
		touchTime = v
	}
	if err := s.workspaceTouchApp(dbTxn, memTxn, wsRef, appRef, touchTime); err != nil {
		return err
	}

	return nil
}

// appSeq gets the pointer to the sequence number for the given application.
// This can only safely be called while holding the memdb write transaction.
func (op *appOperation) appSeq(ref *pb.Ref_Application) *uint64 {
	if op.seq == nil {
		op.seq = map[string]map[string]*uint64{}
	}

	// Get our apps
	k := strings.ToLower(ref.Project)
	apps, ok := op.seq[k]
	if !ok {
		apps = map[string]*uint64{}
		op.seq[k] = apps
	}

	// Get our app
	k = strings.ToLower(ref.Application)
	seq, ok := apps[k]
	if !ok {
		var value uint64
		seq = &value
		apps[k] = seq
	}

	return seq
}

// indexInit initializes the index table in memdb from all the records
// persisted on disk.
func (op *appOperation) indexInit(s *State, dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(op.Bucket)
	c := bucket.Cursor()

	var cnt int

	// This algorithm depends on boltdb's iteration order. Specificly that the keys are
	// lexically order AND because we're using ULID's for the keys in production, the newest
	// records will have the higher lexical value and thusly be at the end of the database.
	//
	// So we just start at the end and insert records until we hit the maximum, knowing
	// we'll be inserted the newest records.

	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		result := op.newStruct()
		if err := proto.Unmarshal(v, result); err != nil {
			return err
		}
		if _, err := op.indexPut(s, memTxn, result); err != nil {
			return err
		}

		// Check if this has a bigger sequence number
		if v := op.valueField(result, "Sequence"); v != nil {
			seq := v.(uint64)

			appRef := op.valueField(result, "Application").(*pb.Ref_Application)
			current := op.appSeq(appRef)
			if seq > *current {
				*current = seq
			}
		}

		cnt++

		if op.MaximumIndexedRecords > 0 && cnt >= op.MaximumIndexedRecords {
			break
		}
	}

	return nil
}

// indexPut writes an index record for a single operation record.
func (op *appOperation) indexPut(
	s *State,
	txn *memdb.Txn,
	value proto.Message,
) (*operationIndexRecord, error) {
	var startTime, completeTime time.Time

	statusRaw := op.valueField(value, "Status")
	if statusRaw != nil {
		statusVal := statusRaw.(*pb.Status)
		if statusVal != nil {
			if t := statusVal.StartTime; t != nil {
				startTime = t.AsTime()
			}

			if t := statusVal.CompleteTime; t != nil {
				completeTime = statusVal.CompleteTime.AsTime()
			}
		}
	}

	var sequence uint64
	if v := op.valueField(value, "Sequence"); v != nil {
		sequence = v.(uint64)
	}

	var generation string
	if v := op.valueField(value, "Generation"); v != nil && v.(*pb.Generation) != nil {
		generation = v.(*pb.Generation).Id
	}

	// Get our refs
	ref := op.valueField(value, "Application").(*pb.Ref_Application)
	wsRef := op.valueField(value, "Workspace").(*pb.Ref_Workspace)

	rec := &operationIndexRecord{
		Id:           op.valueField(value, "Id").(string),
		Project:      ref.Project,
		App:          ref.Application,
		Workspace:    wsRef.Workspace,
		Sequence:     sequence,
		Generation:   generation,
		StartTime:    startTime,
		CompleteTime: completeTime,
	}

	// If there is no maximum, don't track the record count.
	if op.MaximumIndexedRecords != 0 {
		op.pruneMu.Lock()
		op.indexedRecords++
		op.pruneMu.Unlock()
	}

	return rec, txn.Insert(op.memTableName(), rec)
}

func (op *appOperation) pruneOld(memTxn *memdb.Txn, max int) (int, error) {
	return pruneOld(memTxn, pruneOp{
		lock:      &op.pruneMu,
		table:     op.memTableName(),
		index:     opIdIndexName,
		indexArgs: []interface{}{""},
		max:       max,
		cur:       &op.indexedRecords,
	})
}

func (op *appOperation) valueField(value interface{}, field string) interface{} {
	fv := op.valueFieldReflect(value, field)
	if !fv.IsValid() {
		return nil
	}

	return fv.Interface()
}

func (op *appOperation) valueFieldReflect(value interface{}, field string) reflect.Value {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	return v.FieldByName(field)
}

// newStruct creates a pointer to a new value of the type of op.Struct.
// The value of op.Struct is usually itself a pointer so the result of this
// is a pointer to a pointer.
func (op *appOperation) newStruct() proto.Message {
	return reflect.New(reflect.TypeOf(op.Struct).Elem()).Interface().(proto.Message)
}

// workspaceResource returns the resource to use for workspaces.
func (op *appOperation) workspaceResource() string {
	value := reflect.TypeOf(op.Struct).Elem().String()
	if idx := strings.Index(value, "."); idx >= 0 {
		value = value[idx+1:]
	}

	return value
}

func (op *appOperation) memTableName() string {
	return strings.ToLower(string(op.Bucket))
}

// memSchema is the memdb schema for this operation.
func (op *appOperation) memSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: op.memTableName(),
		Indexes: map[string]*memdb.IndexSchema{
			opIdIndexName: {
				Name:         opIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Id",
				},
			},

			opStartTimeIndexName: {
				Name:         opStartTimeIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&IndexTime{
							Field: "StartTime",
						},
					},
				},
			},

			opCompleteTimeIndexName: {
				Name:         opCompleteTimeIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&IndexTime{
							Field: "CompleteTime",
						},
					},
				},
			},

			opSeqIndexName: {
				Name:         opSeqIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&memdb.UintFieldIndex{
							Field: "Sequence",
						},
					},
				},
			},

			opGenIndexName: {
				Name:   opGenIndexName,
				Unique: false,

				// Allow missing since not every app operation has a
				// generation field.
				AllowMissing: true,

				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Generation",
							Lowercase: true,
						},

						&memdb.UintFieldIndex{
							Field: "Sequence",
						},
					},
				},
			},
		},
	}
}

// operationIndexRecord is the record we store in MemDB to perform
// indexed lookup operations by project, app, time, etc.
type operationIndexRecord struct {
	Id           string
	Project      string
	App          string
	Workspace    string
	Sequence     uint64
	Generation   string
	StartTime    time.Time
	CompleteTime time.Time
}

// MatchRef checks if a record matches the ref value. We have to provide
// this because we use LowerBound lookups in memdb and this may return
// a non-matching value at a certain point after iteration.
func (rec *operationIndexRecord) MatchRef(ref *pb.Ref_Application) bool {
	return rec.Project == ref.Project && rec.App == ref.Application
}

const (
	opIdIndexName           = "id"            // id index name
	opStartTimeIndexName    = "start-time"    // start time index
	opCompleteTimeIndexName = "complete-time" // complete time index
	opSeqIndexName          = "seq"           // sequence number index
	opGenIndexName          = "generation"    // generation index
)

// statusFilterMatch is a helper that compares a pb.Status to a set of
// StatusFilters. This returns true if the filters match.
func statusFilterMatch(
	filters []*pb.StatusFilter,
	status *pb.Status,
) bool {
	if len(filters) == 0 {
		return true
	}

NEXT_FILTER:
	for _, group := range filters {
		for _, filter := range group.Filters {
			if !statusFilterMatchSingle(filter, status) {
				continue NEXT_FILTER
			}
		}

		// If any match we match (OR)
		return true
	}

	return false
}

func statusFilterMatchSingle(
	filter *pb.StatusFilter_Filter,
	status *pb.Status,
) bool {
	switch f := filter.Filter.(type) {
	case *pb.StatusFilter_Filter_State:
		return status.State == f.State

	default:
		// unknown filters never match
		return false
	}
}
