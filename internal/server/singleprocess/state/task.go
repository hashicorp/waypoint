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

var taskBucket = []byte("task")

func init() {
	dbBuckets = append(dbBuckets, taskBucket)
	dbIndexers = append(dbIndexers, (*State).taskIndexInit)
	schemas = append(schemas, taskIndexSchema)
}

// TaskPut creates or updates the given Task.
func (s *State) TaskPut(t *pb.Task) error {
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

		return s.taskPut(dbTxn, memTxn, t)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TaskGet gets a task by reference.
func (s *State) TaskGet(ref *pb.Ref_Task) (*pb.Task, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Task
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.taskGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// TaskDelete deletes a task by reference.
func (s *State) TaskDelete(ref *pb.Ref_Task) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.taskDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TaskList returns the list of tasks.
func (s *State) TaskList() ([]*pb.Task, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.taskList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.Task

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.taskGet(dbTxn, memTxn, ref)
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

func (s *State) taskPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Task,
) error {
	id := s.taskId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(taskBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.taskIndexSet(memTxn, id, value)
}

func (s *State) taskGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Task,
) (*pb.Task, error) {
	var result pb.Task
	b := dbTxn.Bucket(taskBucket)

	var taskId string
	switch r := ref.Ref.(type) {
	case *pb.Ref_Task_Id:
		s.log.Info("looking up task by id", "id", r.Id)
		taskId = r.Id
	case *pb.Ref_Task_JobId:
		s.log.Info("looking up task by job id", "job_id", r.JobId)
		// Look up Task by jobid
		task, err := s.taskByJobId(r.JobId)
		if err != nil {
			return nil, err
		}

		s.log.Info("found task id", "id", task.Id)
		taskId = task.Id
	default:
		return nil, status.Error(codes.FailedPrecondition, "No valid ref id provided in Task ref to taskGet")
	}

	return &result, dbGet(b, []byte(strings.ToLower(taskId)), &result)
}

func (s *State) taskList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_Task, error) {
	iter, err := memTxn.Get(taskIndexTableName, taskIndexId+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_Task
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*taskIndexRecord)

		result = append(result, &pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: idx.Id,
			},
		})
	}

	return result, nil
}

func (s *State) taskDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Task,
) error {
	// Get the task. If it doesn't exist then we're successful.
	_, err := s.taskGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id, err := s.taskIdByRef(ref)
	if err != nil {
		return err
	}

	if err := dbTxn.Bucket(taskBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(taskIndexTableName, &taskIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// taskIndexSet writes an index record for a single task.
func (s *State) taskIndexSet(txn *memdb.Txn, id []byte, value *pb.Task) error {
	record := &taskIndexRecord{
		Id:    string(id),
		JobId: value.TaskJob.Id,
	}

	// Insert the index
	return txn.Insert(taskIndexTableName, record)
}

// taskIndexInit initializes the task index from persisted data.
func (s *State) taskIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(taskBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Task
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		if err := s.taskIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) taskId(t *pb.Task) []byte {
	return []byte(strings.ToLower(t.Id))
}

func (s *State) taskIdByRef(ref *pb.Ref_Task) ([]byte, error) {
	var taskId string
	switch t := ref.Ref.(type) {
	case *pb.Ref_Task_Id:
		taskId = t.Id
	case *pb.Ref_Task_JobId:
		// Look up Task by jobid
		task, err := s.taskByJobId(t.JobId)
		if err != nil {
			return nil, err
		}

		taskId = task.Id
	default:
		return nil, status.Error(codes.FailedPrecondition, "No valid ref id provided in Task ref to taskIdByRef")
	}

	return []byte(strings.ToLower(taskId)), nil
}

func (s *State) taskByJobId(jobId string) (*pb.Task, error) {
	trackedTasks, err := s.TaskList()
	if err != nil {
		return nil, err
	}

	for _, t := range trackedTasks {
		if t.TaskJob.Id == jobId {
			return t, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "A Task with job id %q was not found", jobId)
}

func taskIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: taskIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			taskIndexId: {
				Name:         taskIndexId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
			taskIndexJobId: {
				Name:         taskIndexJobId,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "JobId",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	taskIndexTableName = "task-index"
	taskIndexId        = "id"
	taskIndexJobId     = "jobid"
)

type taskIndexRecord struct {
	Id    string
	JobId string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *taskIndexRecord) Copy() *taskIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
