package state

import (
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

		return op.dbGet(tx, []byte(id), result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
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
			"not found for sequence number %d", ref.Number)
	}

	idx := raw.(*operationIndexRecord)
	return idx.Id, nil
}

// List lists all the records.
func (op *appOperation) List(s *State, opts *listOperationsOptions) ([]interface{}, error) {
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
			if err := op.dbGet(tx, []byte(record.Id), value); err != nil {
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

	return result, nil
}

// Latest gets the latest operation that was completed successfully.
func (op *appOperation) Latest(
	s *State,
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (interface{}, error) {
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
		switch st.(*pb.Status).State {
		case pb.Status_SUCCESS:
			return v, nil
		}
	}

	return nil, status.Error(codes.NotFound, "none available")
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
		matchFields := []string{"Sequence"}
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

	// If we're not updating, then set the sequence number up if we have one.
	if !update {
		if f := op.valueFieldReflect(value, "Sequence"); f.IsValid() {
			seq := atomic.AddUint64(op.appSeq(appRef), 1)
			f.Set(reflect.ValueOf(seq))
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
	return bucket.ForEach(func(k, v []byte) error {
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

		return nil
	})
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
				st, err := ptypes.Timestamp(t)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "time for build can't be parsed")
				}

				startTime = st
			}

			if t := statusVal.CompleteTime; t != nil {
				ct, err := ptypes.Timestamp(statusVal.CompleteTime)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "time for build can't be parsed")
				}

				completeTime = ct
			}
		}
	}

	var sequence uint64
	if v := op.valueField(value, "Sequence"); v != nil {
		sequence = v.(uint64)
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
		StartTime:    startTime,
		CompleteTime: completeTime,
	}
	return rec, txn.Insert(op.memTableName(), rec)
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
)

// listOperationsOptions are options that can be set for List calls on
// operations for filtering and limiting the response.
type listOperationsOptions struct {
	Application   *pb.Ref_Application
	Workspace     *pb.Ref_Workspace
	Status        []*pb.StatusFilter
	Order         *pb.OperationOrder
	PhysicalState pb.Operation_PhysicalState
}

func buildListOperationsOptions(ref *pb.Ref_Application, opts ...ListOperationOption) *listOperationsOptions {
	var result listOperationsOptions
	result.Application = ref
	for _, opt := range opts {
		opt(&result)
	}

	return &result
}

// ListOperationOption is an exported type to set configuration for listing operations.
type ListOperationOption func(opts *listOperationsOptions)

// ListWithStatusFilter sets a status filter.
func ListWithStatusFilter(f ...*pb.StatusFilter) ListOperationOption {
	return func(opts *listOperationsOptions) {
		opts.Status = f
	}
}

// ListWithOrder sets ordering on the list operation.
func ListWithOrder(f *pb.OperationOrder) ListOperationOption {
	return func(opts *listOperationsOptions) {
		opts.Order = f
	}
}

// ListWithPhysicalState sets ordering on the list operation.
func ListWithPhysicalState(f pb.Operation_PhysicalState) ListOperationOption {
	return func(opts *listOperationsOptions) {
		opts.PhysicalState = f
	}
}

// ListWithWorkspace sets ordering on the list operation.
func ListWithWorkspace(f *pb.Ref_Workspace) ListOperationOption {
	return func(opts *listOperationsOptions) {
		opts.Workspace = f
	}
}

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
