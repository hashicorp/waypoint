package state

import (
	"reflect"
	"strings"
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
	Struct interface{}

	// Bucket is the global bucket for all records of this operation.
	Bucket []byte
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
		return op.dbPut(dbTxn, memTxn, update, value)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// Get gets an operation record by ID.
func (op *appOperation) Get(s *State, id string) (interface{}, error) {
	result := op.newStruct()
	err := s.db.View(func(tx *bolt.Tx) error {
		return dbGet(tx.Bucket(op.Bucket), []byte(id), result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
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
		bucket := tx.Bucket(op.Bucket)

		for {
			current := iter.Next()
			if current == nil {
				return nil
			}

			record := current.(*operationIndexRecord)
			if !record.MatchRef(opts.Application) {
				return nil
			}

			value := op.newStruct()
			if err := dbGet(bucket, []byte(record.Id), value); err != nil {
				return err
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
func (op *appOperation) Latest(s *State, ref *pb.Ref_Application) (interface{}, error) {
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
			return nil, nil
		}

		record := raw.(*operationIndexRecord)
		if !record.MatchRef(ref) {
			return nil, nil
		}

		v, err := op.Get(s, record.Id)
		if err != nil {
			return nil, err
		}

		// Shouldn't happen but if it does, return nothing.
		st := op.valueField(v, "Status")
		if st == nil {
			return nil, nil
		}

		// State must be success.
		switch st.(*pb.Status).State {
		case pb.Status_SUCCESS:
			return v, nil
		}
	}
}

// dbPut wites the value to the database and also sets up any index records.
// It expects to hold a write transaction to both bolt and memdb.
func (op *appOperation) dbPut(
	tx *bolt.Tx,
	inmemTxn *memdb.Txn,
	update bool,
	value proto.Message,
) error {
	id := []byte(op.valueField(value, "Id").(string))

	// Get the global bucket and write the value to it.
	b := tx.Bucket(op.Bucket)

	// If we're updating, then this shouldn't already exist
	if update && b.Get(id) == nil {
		return status.Errorf(codes.NotFound, "record with ID %q not found for update", string(id))
	}

	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return op.indexPut(inmemTxn, value)
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
		if err := op.indexPut(memTxn, result); err != nil {
			return err
		}

		return nil
	})
}

// indexPut writes an index record for a single operation record.
func (op *appOperation) indexPut(txn *memdb.Txn, value proto.Message) error {
	var startTime, completeTime time.Time

	statusRaw := op.valueField(value, "Status")
	if statusRaw != nil {
		statusVal := statusRaw.(*pb.Status)
		if statusVal != nil {
			if t := statusVal.StartTime; t != nil {
				st, err := ptypes.Timestamp(t)
				if err != nil {
					return status.Errorf(codes.Internal, "time for build can't be parsed")
				}

				startTime = st
			}

			if t := statusVal.CompleteTime; t != nil {
				ct, err := ptypes.Timestamp(statusVal.CompleteTime)
				if err != nil {
					return status.Errorf(codes.Internal, "time for build can't be parsed")
				}

				completeTime = ct
			}
		}
	}

	ref := op.valueField(value, "Application").(*pb.Ref_Application)
	return txn.Insert(op.memTableName(), &operationIndexRecord{
		Id:           op.valueField(value, "Id").(string),
		Project:      ref.Project,
		App:          ref.Application,
		StartTime:    startTime,
		CompleteTime: completeTime,
	})
}

func (op *appOperation) valueField(value interface{}, field string) interface{} {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	fv := v.FieldByName(field)
	if !fv.IsValid() {
		return nil
	}

	return fv.Interface()
}

// newStruct creates a pointer to a new value of the type of op.Struct.
// The value of op.Struct is usually itself a pointer so the result of this
// is a pointer to a pointer.
func (op *appOperation) newStruct() proto.Message {
	return reflect.New(reflect.TypeOf(op.Struct).Elem()).Interface().(proto.Message)
}

func (op *appOperation) memTableName() string {
	return strings.ToLower(string(op.Bucket))
}

// memSchema is the memdb schema for this operation.
func (op *appOperation) memSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: op.memTableName(),
		Indexes: map[string]*memdb.IndexSchema{
			opIdIndexName: &memdb.IndexSchema{
				Name:         opIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Id",
				},
			},

			opStartTimeIndexName: &memdb.IndexSchema{
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

			opCompleteTimeIndexName: &memdb.IndexSchema{
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
		},
	}
}

type operationIndexRecord struct {
	Id           string
	Project      string
	App          string
	StartTime    time.Time
	CompleteTime time.Time
}

// MatchRef checks if a record matches the ref value. We have to provide
// this because we use LowerBound lookups in memdb and this may return
// a non-matching value at a certain point.
func (rec *operationIndexRecord) MatchRef(ref *pb.Ref_Application) bool {
	return rec.Project == ref.Project && rec.App == ref.Application
}

const (
	opIdIndexName           = "id"
	opStartTimeIndexName    = "start-time"
	opCompleteTimeIndexName = "complete-time"
)

// listOperationsOptions are options that can be set for List calls on
// operations for filtering and limiting the response.
type listOperationsOptions struct {
	Application *pb.Ref_Application
	Status      []*pb.StatusFilter
	Order       *pb.OperationOrder
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
