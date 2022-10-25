package boltdbstate

import (
	"context"
	"strings"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var taskBucket = []byte("task")

func init() {
	dbBuckets = append(dbBuckets, taskBucket)
	dbIndexers = append(dbIndexers, (*State).taskIndexInit)
	schemas = append(schemas, taskIndexSchema)
}

// TaskPut creates or updates the given Task.
func (s *State) TaskPut(ctx context.Context, t *pb.Task) error {
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
func (s *State) TaskGet(ctx context.Context, ref *pb.Ref_Task) (*pb.Task, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Task
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.taskGet(ctx, dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// TaskDelete deletes a task by reference.
func (s *State) TaskDelete(ctx context.Context, ref *pb.Ref_Task) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.taskDelete(ctx, dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// TaskCancel cancels a tasks jobs by task id or run job id reference.
// NOTE(briancain): this way means each cancel is its own transaction and commit, not the greatest.
// Previously I attempted to implement this with taskCancel where the entire job
// triple was canceled inside a single transaction. This ended up deadlocking because
// when we cancel a job, we also call a `db.Update` on each job when we update its
// state in the database. This caused a deadlock. For now we cancel each job
// separately.
func (s *State) TaskCancel(ctx context.Context, ref *pb.Ref_Task) error {
	task, err := s.TaskGet(ctx, ref)
	if err != nil {
		return err
	}

	s.log.Trace("canceling start job for task", "task id", task.Id, "start job id", task.StartJob.Id)
	err = s.JobCancel(task.StartJob.Id, false)
	if err != nil {
		return err
	}

	s.log.Trace("canceling task job for task", "task id", task.Id, "task job id", task.TaskJob.Id)
	err = s.JobCancel(task.TaskJob.Id, false)
	if err != nil {
		return err
	}

	if task.WatchJob != nil {
		s.log.Trace("canceling watch job for task", "task id", task.Id, "watch job id", task.WatchJob.Id)
		err = s.JobCancel(task.WatchJob.Id, false)
		if err != nil {
			return err
		}
	}

	s.log.Trace("canceling stop job for task", "task id", task.Id, "stop job id", task.StopJob.Id)
	err = s.JobCancel(task.StopJob.Id, false)
	if err != nil {
		return err
	}

	return nil
}

// GetJobsByTaskRef will look up every job triple by Task ref in a single
// memdb transaction. This is often used via the API for building out
// a complete picture of a task beyond the job ID refs.
func (s *State) JobsByTaskRef(
	ctx context.Context,
	task *pb.Task,
) (startJob *pb.Job, taskJob *pb.Job, stopJob *pb.Job, watchJob *pb.Job, err error) {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		job, err := s.jobById(dbTxn, task.StartJob.Id)
		if err != nil {
			return err
		} else if job == nil {
			return status.Errorf(codes.NotFound, "start job %q not found", task.StartJob.Id)
		} else {
			startJob = job
		}

		job, err = s.jobById(dbTxn, task.TaskJob.Id)
		if err != nil {
			return err
		} else if job == nil {
			return status.Errorf(codes.NotFound, "task job %q not found", task.TaskJob.Id)
		} else {
			taskJob = job
		}

		job, err = s.jobById(dbTxn, task.StopJob.Id)
		if err != nil {
			return err
		} else if job == nil {
			return status.Errorf(codes.NotFound, "stop job %q not found", task.StopJob.Id)
		} else {
			stopJob = job
		}

		job, err = s.jobById(dbTxn, task.WatchJob.Id)
		if err != nil {
			return err
		} else if job == nil {
			return status.Errorf(codes.NotFound, "watch job %q not found", task.WatchJob.Id)
		} else {
			watchJob = job
		}

		return nil
	})

	return startJob, taskJob, stopJob, watchJob, err
}

// TaskList returns the list of tasks.
func (s *State) TaskList(ctx context.Context, req *pb.ListTaskRequest) ([]*pb.Task, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	refs, err := s.taskList(memTxn)
	if err != nil {
		return nil, err
	}

	var out []*pb.Task

	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, ref := range refs {
			val, err := s.taskGet(ctx, dbTxn, memTxn, ref)
			if err != nil {
				return err
			}

			// filter any tasks by request
			if len(req.TaskState) > 0 {
				found := false
				for _, state := range req.TaskState {
					if val.JobState == state {
						found = true
						break
					}
				}
				if !found {
					continue
				}
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
	ctx context.Context,
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Task,
) (*pb.Task, error) {
	var result pb.Task
	b := dbTxn.Bucket(taskBucket)

	var taskId string
	switch r := ref.Ref.(type) {
	case *pb.Ref_Task_Id:
		s.log.Trace("looking up task", "id", r.Id)
		taskId = r.Id
	case *pb.Ref_Task_JobId:
		s.log.Trace("looking up task by job id", "job_id", r.JobId)
		// Look up Task by jobid
		task, err := s.taskByJobId(ctx, r.JobId)
		if err != nil {
			return nil, err
		}

		s.log.Trace("found task id", "id", task.Id)
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
	ctx context.Context,
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Task,
) error {
	// Get the task. If it doesn't exist then we're successful.
	_, err := s.taskGet(ctx, dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id, err := s.taskIdByRef(ctx, ref)
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

/*
NOTE(briancain): This was intentionally left commented out. In the future, if
we ever decide to cancel Task jobs in a single transaction (see the note above `TaskCancel`)
we can use this function to do the actual cancelation in that transaction.
func (s *State) taskCancel(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	task *pb.Task,
) error {
	s.log.Info("canceling task", "id", task.Id)
	// call jobCancel on the job triple

	// Get the start job
	raw, err := memTxn.First(jobTableName, jobIdIndexName, task.StartJob.Id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "StartJob for task %q not found: %s", task.Id, task.StartJob.Id)
	}
	startJob := raw.(*jobIndex)

	// Cancel the job
	s.log.Info("canceling start job", "id", startJob.Id)
	if err = s.jobCancel(memTxn, startJob, false); err != nil {
		return err
	}

	// Get the task job
	raw, err = memTxn.First(jobTableName, jobIdIndexName, task.TaskJob.Id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "TaskJob for task %q not found: %s", task.Id, task.TaskJob.Id)
	}
	taskJob := raw.(*jobIndex)

	// Cancel the job
	s.log.Info("canceling task job", "id", taskJob.Id)
	if err = s.jobCancel(memTxn, taskJob, false); err != nil {
		return err
	}

	// Get the stop job
	raw, err = memTxn.First(jobTableName, jobIdIndexName, task.StopJob.Id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "StopJob for task %q not found: %s", task.Id, task.StopJob.Id)
	}
	stopJob := raw.(*jobIndex)

	// Cancel the job
	s.log.Info("canceling stop job", "id", stopJob.Id)
	if err = s.jobCancel(memTxn, stopJob, false); err != nil {
		return err
	}

	return nil
}
*/

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

func (s *State) taskIdByRef(ctx context.Context, ref *pb.Ref_Task) ([]byte, error) {
	var taskId string
	switch t := ref.Ref.(type) {
	case *pb.Ref_Task_Id:
		taskId = t.Id
	case *pb.Ref_Task_JobId:
		// Look up Task by jobid
		task, err := s.taskByJobId(ctx, t.JobId)
		if err != nil {
			return nil, err
		}

		taskId = task.Id
	default:
		return nil, status.Error(codes.FailedPrecondition, "No valid ref id provided in Task ref to taskIdByRef")
	}

	return []byte(strings.ToLower(taskId)), nil
}

func (s *State) taskByJobId(ctx context.Context, jobId string) (*pb.Task, error) {
	trackedTasks, err := s.TaskList(ctx, &pb.ListTaskRequest{})
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
