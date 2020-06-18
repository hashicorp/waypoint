package state

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const (
	jobTableName          = "jobs"
	jobIdIndexName        = "id"
	jobStateIndexName     = "state"
	jobQueueTimeIndexName = "queue-time"
	jobTargetIdIndexName  = "target-id"
)

func init() {
	schemas = append(schemas, jobSchema)
}

func jobSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobTableName,
		Indexes: map[string]*memdb.IndexSchema{
			jobIdIndexName: &memdb.IndexSchema{
				Name:         jobIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Id",
				},
			},

			jobStateIndexName: &memdb.IndexSchema{
				Name:         jobStateIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.IntFieldIndex{
					Field: "State",
				},
			},

			jobQueueTimeIndexName: &memdb.IndexSchema{
				Name:         jobQueueTimeIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.IntFieldIndex{
							Field: "State",
						},

						&IndexTime{
							Field: "QueueTime",
							Asc:   true,
						},
					},
				},
			},

			jobTargetIdIndexName: &memdb.IndexSchema{
				Name:         jobTargetIdIndexName,
				AllowMissing: true,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.IntFieldIndex{
							Field: "State",
						},

						&memdb.StringFieldIndex{
							Field:     "TargetRunnerId",
							Lowercase: true,
						},

						&IndexTime{
							Field: "QueueTime",
							Asc:   true,
						},
					},
				},
			},
		},
	}
}

type Job struct {
	Id string

	// QueueTime is the time that the job was queued.
	QueueTime time.Time

	// TargetAny will be true if this job targets anything
	TargetAny bool

	// TargetRunnerId is the ID of the runner to target.
	TargetRunnerId string

	// State is the current state of this job.
	State pb.Job_State

	// TODO(mitchellh): we should persist the job to boltdb rather
	// than storing it in memory so that we can get the value back.
	Job *pb.Job
}

// JobCreate queues the given job.
func (s *State) JobCreate(jobpb *pb.Job) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	if err := s.jobCreate(txn, jobpb); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// JobById looks up a job by ID. The returned Job will be a deep copy
// of the job so it is safe to read/write. If the job can't be found,
// a nil result with no error is returned.
func (s *State) JobById(id string) (*pb.Job, error) {
	// TODO(mitchellh): this isn't actually a deep copy until we persist to bolt

	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	raw, err := memTxn.First(jobTableName, jobIdIndexName, id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	job := raw.(*Job)

	return job.Job, nil
}

// JobAssignForRunner will wait for and assign a job to a specific runner.
// This will automatically evaluate any conditions that the runner and/or
// job may have on assignability.
//
// The assigned job is put into a "waiting" state until the runner
// acks the assignment which can be set with JobAck.
//
// If ctx is provided and assignment has to block waiting for new jobs,
// this will cancel when the context is done.
func (s *State) JobAssignForRunner(ctx context.Context, r *Runner) (*pb.Job, error) {
RETRY_ASSIGN:
	txn := s.inmem.Txn(false)
	defer txn.Abort()

	// candidateQuery finds candidate jobs to assign.
	candidateQuery := []func(*memdb.Txn, *Runner) (*Job, error){
		s.jobCandidateById,
		s.jobCandidateAny,
	}

	// Build the list of candidates
	var candidates []*Job
	for _, f := range candidateQuery {
		job, err := f(txn, r)
		if err != nil {
			return nil, err
		}
		if job != nil {
			candidates = append(candidates, job)
		}
	}

	// If we have no candidates, then we have to wait for a job to show up.
	// We set up a blocking query on the job table for a non-assigned job.
	var watchCh <-chan struct{}
	if len(candidates) == 0 {
		iter, err := txn.Get(jobTableName, jobStateIndexName, pb.Job_QUEUED)
		if err != nil {
			return nil, err
		}
		watchCh = iter.WatchCh()
	}

	// We're done reading so abort the transaction
	txn.Abort()

	// If we have a watch channel set that means we didn't find any
	// results and we need to retry after waiting for changes.
	if watchCh != nil {
		select {
		case <-watchCh:
			goto RETRY_ASSIGN

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// We sort our candidates by queue time so that we can find the earliest
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].QueueTime.Before(candidates[j].QueueTime)
	})

	// Grab a write lock since we're going to delete, modify, add the
	// job that we chose. No need to defer here since the first defer works
	// at the top of the func.
	//
	// Write locks are exclusive so this will ensure we're the only one
	// writing at a time. This lets us be sure we're the only one "assigning"
	// a job candidate.
	txn = s.inmem.Txn(true)
	for _, job := range candidates {
		// Get the job
		raw, err := txn.First(jobTableName, jobIdIndexName, job.Id)
		if err != nil {
			return nil, err
		}
		if raw == nil {
			// The job no longer exists. It may be canceled or something.
			// Invalid candidate, continue to next.
			continue
		}

		// We need to verify that in the time between our candidate search
		// and our write lock acquisition, that this job hasn't been assigned,
		// canceled, etc. If so, this is an invalid candidate.
		job := raw.(*Job)
		if job == nil || job.State != pb.Job_QUEUED {
			continue
		}

		// Delete the job. We have to delete + insert to do an index update.
		if err := txn.Delete(jobTableName, job); err != nil {
			// Some other error. Give up. This should NEVER be an ErrNotFound
			// since we're in a transaction that verified it exists above by
			// doing a Get.
			return nil, err
		}

		job.State = pb.Job_WAITING
		job.Job.State = job.State
		job.Job.AssignTime, err = ptypes.TimestampProto(time.Now())
		if err != nil {
			// This should never happen since encoding a time now should be safe
			panic("time encoding failed: " + err.Error())
		}

		if err := txn.Insert(jobTableName, job); err != nil {
			return nil, err
		}

		txn.Commit()
		return job.Job, nil
	}
	txn.Abort()

	// If we reached here, all of our candidates were invalid, we retry
	goto RETRY_ASSIGN
}

// JobAck acknowledges that a job has been accepted or rejected by the runner.
// If ack is false, then this will move the job back to the queued state
// and be eligible for assignment.
func (s *State) JobAck(id string, ack bool) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Get the job
	raw, err := txn.First(jobTableName, jobIdIndexName, id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "job not found: %s", id)
	}
	job := raw.(*Job)

	// If the job is not in the assigned state, then this is an error.
	if job.State != pb.Job_WAITING {
		return status.Errorf(codes.FailedPrecondition,
			"job can't be acked from state: %s",
			job.State.String())
	}

	if ack {
		// Set to accepted
		job.State = pb.Job_RUNNING
		job.Job.State = job.State
		job.Job.AckTime, err = ptypes.TimestampProto(time.Now())
		if err != nil {
			// This should never happen since encoding a time now should be safe
			panic("time encoding failed: " + err.Error())
		}
	} else {
		// Set to queued
		job.State = pb.Job_QUEUED
		job.Job.State = job.State
		job.Job.AssignTime = nil
	}

	// Insert to update
	if err := txn.Insert(jobTableName, job); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// JobComplete marks a running job as complete. If an error is given,
// the job is marked as failed (a completed state). If no error is given,
// the job is marked as successful.
func (s *State) JobComplete(id string, cerr error) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Get the job
	raw, err := txn.First(jobTableName, jobIdIndexName, id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "job not found: %s", id)
	}
	job := raw.(*Job)

	// If the job is not in the assigned state, then this is an error.
	if job.State != pb.Job_RUNNING {
		return status.Errorf(codes.FailedPrecondition,
			"job can't be completed from state: %s",
			job.State.String())
	}

	// Set to complete, assume success for now
	job.State = pb.Job_SUCCESS
	job.Job.State = job.State
	job.Job.CompleteTime, err = ptypes.TimestampProto(time.Now())
	if err != nil {
		// This should never happen since encoding a time now should be safe
		panic("time encoding failed: " + err.Error())
	}

	if cerr != nil {
		job.State = pb.Job_ERROR
		job.Job.State = job.State

		st, _ := status.FromError(cerr)
		job.Job.Error = st.Proto()
	}

	// Insert to update
	if err := txn.Insert(jobTableName, job); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *State) jobCreate(memTxn *memdb.Txn, jobpb *pb.Job) error {
	rec := &Job{
		Id:        jobpb.Id,
		State:     pb.Job_QUEUED,
		QueueTime: time.Now(),
		Job:       jobpb,
	}

	switch v := jobpb.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Any:
		rec.TargetAny = true

	case *pb.Ref_Runner_Id:
		rec.TargetRunnerId = v.Id.Id

	default:
		return fmt.Errorf("unknown runner target value: %#v", jobpb.TargetRunner.Target)
	}

	// Setup our initial job state
	var err error
	jobpb.State = rec.State
	jobpb.QueueTime, err = ptypes.TimestampProto(rec.QueueTime)
	if err != nil {
		return err
	}

	// Insert into the DB
	if err := memTxn.Insert(jobTableName, rec); err != nil {
		return err
	}

	return nil
}

// jobCandidateById returns the most promising candidate job to assign
// that is targeting a specific runner by ID.
func (s *State) jobCandidateById(memTxn *memdb.Txn, r *Runner) (*Job, error) {
	iter, err := memTxn.LowerBound(
		jobTableName,
		jobTargetIdIndexName,
		pb.Job_QUEUED,
		r.Id,
		time.Unix(0, 0),
	)
	if err != nil {
		return nil, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}

		job := raw.(*Job)
		if job.State != pb.Job_QUEUED || job.TargetRunnerId == "" {
			continue
		}

		return job, nil
	}

	return nil, nil
}

// jobCandidateAny returns the first candidate job that targets any runner.
func (s *State) jobCandidateAny(memTxn *memdb.Txn, r *Runner) (*Job, error) {
	iter, err := memTxn.LowerBound(
		jobTableName,
		jobQueueTimeIndexName,
		pb.Job_QUEUED,
		time.Unix(0, 0),
	)
	if err != nil {
		return nil, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}

		job := raw.(*Job)
		if job.State != pb.Job_QUEUED || !job.TargetAny {
			continue
		}

		return job, nil
	}

	return nil, nil
}
