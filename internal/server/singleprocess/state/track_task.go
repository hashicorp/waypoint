package state

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var tracktaskBucket = []byte("tracktask")

func init() {
	dbBuckets = append(dbBuckets, tracktaskBucket)
	dbIndexers = append(dbIndexers, (*State).tracktaskIndexInit)
	schemas = append(schemas, tracktaskIndexSchema)
}

// TrackTaskPut creates or updates the given TrackTask.
func (s *State) TrackTaskPut(t *pb.TrackTask) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		if t.TaskJob == nil {
			return status.Error(codes.FailedPrecondition,
				"a Job ref for the TaskJob is required")
		}

		if t.Id == "" {
			id, err := ulid()
			if err != nil {
				return err
			}

			t.Id = id
		}

		return s.tracktaskPut(dbTxn, memTxn, t)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TrackTaskGet gets a tracktask by reference.
func (s *State) TrackTaskGet(ref *pb.Ref_TrackTask) (*pb.TrackTask, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.TrackTask
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.tracktaskGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// TrackTaskDelete deletes a tracktask by reference.
func (s *State) TrackTaskDelete(ref *pb.Ref_TrackTask) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.tracktaskDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TrackTaskList returns the list of tracktasks.
func (s *State) TrackTaskList() ([]*pb.TrackTask, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.tracktaskList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.TrackTask

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.tracktaskGet(dbTxn, memTxn, ref)
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

func (s *State) tracktaskPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.TrackTask,
) error {
	id := s.tracktaskId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(tracktaskBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.tracktaskIndexSet(memTxn, id, value)
}

func (s *State) tracktaskGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_TrackTask,
) (*pb.TrackTask, error) {
	var result pb.TrackTask
	b := dbTxn.Bucket(tracktaskBucket)

	var taskId string
	switch r := ref.Ref.(type) {
	case *pb.Ref_TrackTask_Id:
		s.log.Info("looking up tracktask by id", "id", r.Id)
		taskId = r.Id
	case *pb.Ref_TrackTask_JobId:
		s.log.Info("looking up tracktask by job id", "job_id", r.JobId)
		// Look up TrackTask by jobid
		task, err := s.tracktaskByJobId(r.JobId)
		if err != nil {
			return nil, err
		}

		s.log.Info("found tracktask id", "id", task.Id)
		taskId = task.Id
	default:
		return nil, status.Error(codes.FailedPrecondition, "No valid ref id provided in TrackTask ref to tracktaskGet")
	}

	return &result, dbGet(b, []byte(strings.ToLower(taskId)), &result)
}

func (s *State) tracktaskList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_TrackTask, error) {
	iter, err := memTxn.Get(tracktaskIndexTableName, tracktaskIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_TrackTask
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*tracktaskIndexRecord)

		result = append(result, &pb.Ref_TrackTask{
			Ref: &pb.Ref_TrackTask_Id{
				Id: idx.Id,
			},
		})
	}

	return result, nil
}

func (s *State) tracktaskDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_TrackTask,
) error {
	// Get the tracktask. If it doesn't exist then we're successful.
	_, err := s.tracktaskGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id, err := s.tracktaskIdByRef(ref)
	if err != nil {
		return err
	}

	if err := dbTxn.Bucket(tracktaskBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(tracktaskIndexTableName, &tracktaskIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// tracktaskIndexSet writes an index record for a single tracktask.
func (s *State) tracktaskIndexSet(txn *memdb.Txn, id []byte, value *pb.TrackTask) error {
	record := &tracktaskIndexRecord{
		Id: string(id),
	}

	// Insert the index
	return txn.Insert(tracktaskIndexTableName, record)
}

// tracktaskIndexInit initializes the tracktask index from persisted data.
func (s *State) tracktaskIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(tracktaskBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.TrackTask
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		if err := s.tracktaskIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) tracktaskId(t *pb.TrackTask) []byte {
	return []byte(strings.ToLower(t.Id))
}

func (s *State) tracktaskIdByRef(ref *pb.Ref_TrackTask) ([]byte, error) {
	var taskId string
	switch t := ref.Ref.(type) {
	case *pb.Ref_TrackTask_Id:
		taskId = t.Id
	case *pb.Ref_TrackTask_JobId:
		// Look up TrackTask by jobid
		task, err := s.tracktaskByJobId(t.JobId)
		if err != nil {
			return nil, err
		}

		taskId = task.Id
	default:
		return nil, status.Error(codes.FailedPrecondition, "No valid ref id provided in TrackTask ref to trackTaskIdByRef")
	}

	return []byte(strings.ToLower(taskId)), nil
}

func (s *State) tracktaskByJobId(jobId string) (*pb.TrackTask, error) {
	trackedTasks, err := s.TrackTaskList()
	if err != nil {
		return nil, err
	}

	for _, t := range trackedTasks {
		if t.TaskJob.Id == jobId {
			return t, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "A TrackTask with job id %q was not found", jobId)
}

func tracktaskIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: tracktaskIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			tracktaskIndexId: {
				Name:         tracktaskIndexId,
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
	tracktaskIndexTableName = "tracktask-index"
	tracktaskIndexId        = "id"
)

type tracktaskIndexRecord struct {
	Id string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *tracktaskIndexRecord) Copy() *tracktaskIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
