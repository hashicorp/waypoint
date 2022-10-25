package boltdbstate

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/copystructure"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

var (
	jobBucket = []byte("jobs")
)

const (
	jobTableName            = "jobs"
	jobIdIndexName          = "id"
	jobStateIndexName       = "state"
	jobQueueTimeIndexName   = "queue-time"
	jobTargetIdIndexName    = "target-id"
	jobSingletonIdIndexName = "singleton-id"
	jobDependsOnIndexName   = "depends-on"

	maximumJobsIndexed = 10000
)

func init() {
	dbBuckets = append(dbBuckets, jobBucket)
	dbIndexers = append(dbIndexers, (*State).jobIndexInit)
	schemas = append(schemas, jobSchema)
}

func jobSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobTableName,
		Indexes: map[string]*memdb.IndexSchema{
			jobIdIndexName: {
				Name:         jobIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Id",
				},
			},

			jobStateIndexName: {
				Name:         jobStateIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.IntFieldIndex{
					Field: "State",
				},
			},

			jobQueueTimeIndexName: {
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

			jobTargetIdIndexName: {
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

			jobSingletonIdIndexName: {
				Name:         jobSingletonIdIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "SingletonId",
							Lowercase: true,
						},

						&memdb.IntFieldIndex{
							Field: "State",
						},
					},
				},
			},

			jobDependsOnIndexName: {
				Name:         jobDependsOnIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.StringSliceFieldIndex{
					Field:     "DependsOn",
					Lowercase: true,
				},
			},
		},
	}
}

type jobIndex struct {
	Id string

	// SingletonId matches singleton_id if set on the job.
	SingletonId string

	// DependsOn is the list of jobs that this job depends on. If these
	// don't exist, they're assumed to have COMPLETED due to pruning,
	// since we don't allow job creation unless they existed in the past,
	// and if they errored then they would've immediately set this job
	// state to error.
	DependsOn []string

	// OpType is the operation type for the job.
	OpType reflect.Type

	// The project/workspace that this job is part of. This is used
	// to determine if the job is blocked. See job_assigned.go for more details.
	Application *pb.Ref_Application
	Workspace   *pb.Ref_Workspace

	// QueueTime is the time that the job was queued.
	QueueTime time.Time

	// TargetAny will be true if this job targets anything
	TargetAny bool

	// TargetRunnerId is the ID of the runner to target.
	TargetRunnerId string

	// TargetRunnerLabels are the labels of the runner to target.
	TargetRunnerLabels map[string]string

	// State is the current state of this job.
	State pb.Job_State

	// StateTimer holds a timer that is usually acting as a timeout mechanism
	// on the current state. When the state changes, the timer should be cancelled.
	StateTimer *time.Timer

	// OutputBuffer stores the terminal output
	OutputBuffer *logbuffer.Buffer
}

// A helper, pulled out rather than on a value to allow it to be used against
// pb.Job,s and jobIndex's alike.
func jobIsCompleted(state pb.Job_State) bool {
	switch state {
	case pb.Job_ERROR, pb.Job_SUCCESS:
		return true
	default:
		return false
	}
}

// JobCreate queues the given jobs. If any job fails to queue, no jobs
// are queued. If partial failures are acceptable, call this multiple times
// with a single job.
func (s *State) JobCreate(jobs ...*pb.Job) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Before we do any job creation, we go through and verify that
	// any dependencies being created do not result in a cycle. We know
	// that jobs already created couldn't have depended on these new jobs
	// because a dependency must exist at time of creation, so we only
	// need to check our new dependencies for cycles between each other.
	jobMap := map[string]*pb.Job{}
	var depGraph graph.Graph
	for _, job := range jobs {
		// Add our job
		jobId := strings.ToLower(job.Id)
		jobMap[jobId] = job
		depGraph.Add(jobId)

		// Add any dependencies
		for _, depId := range job.DependsOn {
			depId = strings.ToLower(depId)
			depGraph.Add(depId)
			depGraph.AddEdge(depId, jobId)
		}
	}
	if cycles := depGraph.Cycles(); len(cycles) > 0 {
		return status.Errorf(codes.FailedPrecondition,
			"Job dependencies contain one or more cycles: %#v", cycles)
	}

	// Get our order that we'll create the jobs in. We do a topological
	// sort so that we get we create dependencies first. This is so that
	// if someone submits jobs A, B, C that depend on each other in that order
	// in argumen torder of C, A, B, we can still create, since jobCreate
	// requires that all dependencies exist.
	order := depGraph.KahnSort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		// Go through the jobs in DFS-order. This ensures that we create
		// the jobs that are dependencies first.
		for _, id := range order {
			// Create the job. Its okay if it doesn't exist in the jobMap,
			// that means its just a root node or some other dep already exists.
			job, ok := jobMap[id.(string)]
			if ok {
				if err := s.jobCreate(dbTxn, txn, job); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err == nil {
		txn.Commit()
	}

	return err
}

// Given a Project Ref and a Job Template, this function will generate a slice
// of QueueJobRequests for every application inside the requested Project
func (s *State) JobProjectScopedRequest(
	ctx context.Context,
	pRef *pb.Ref_Project,
	jobTemplate *pb.Job,
) ([]*pb.QueueJobRequest, error) {
	var result []*pb.QueueJobRequest
	project, err := s.ProjectGet(ctx, pRef)
	if err != nil {
		return nil, err
	}

	for _, app := range project.Applications {
		copyJob, err := copystructure.Copy(jobTemplate)
		if err != nil {
			return nil, status.Errorf(codes.Internal,
				"failed to copy job template for project scoped request: %s", err)
		}
		tempJob, ok := copyJob.(*pb.Job)
		if !ok {
			return nil, status.Errorf(codes.Internal,
				"failed to convert copied job template into a Job message: %s", err)
		}

		tempJob.Application = &pb.Ref_Application{
			Project:     project.Name,
			Application: app.Name,
		}

		jobReq := &pb.QueueJobRequest{Job: tempJob}
		result = append(result, jobReq)
	}

	return result, nil
}

// JobList returns the list of jobs.
func (s *State) JobList(
	req *pb.ListJobsRequest,
) ([]*pb.Job, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(jobTableName, jobIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Job
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*jobIndex)

		var job *pb.Job
		err = s.db.View(func(dbTxn *bolt.Tx) error {
			job, err = s.jobById(dbTxn, idx.Id)
			return err
		})

		// filter job list by request
		if req.Workspace != nil {
			if job.Workspace.Workspace != req.Workspace.Workspace {
				continue
			}
		}

		if req.Project != nil {
			if job.Application.Project != req.Project.Project {
				continue
			}
		}

		if req.Application != nil {
			if job.Application.Application != req.Application.Application {
				continue
			}
			if job.Application.Project != "" && job.Application.Project != req.Application.Project {
				continue
			}
		}

		if len(req.JobState) > 0 {
			found := false
			for _, state := range req.JobState {
				if job.State == state {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if req.TargetRunner != nil {
			switch tr := job.TargetRunner.Target.(type) {
			case *pb.Ref_Runner_Any:
				// job was set to Any runner
				_, ok := req.TargetRunner.Target.(*pb.Ref_Runner_Any)
				if !ok {
					// request is not targeted to Any runner, so don't include in list
					continue
				}
			case *pb.Ref_Runner_Id:
				// the job is targed to a specific runner id

				id, ok := req.TargetRunner.Target.(*pb.Ref_Runner_Id)
				if !ok {
					// request was not for a target runner by id
					continue
				} else if id.Id.Id != tr.Id.Id {
					// the requested id doesn't match the target runner id on the job
					continue
				}
			case *pb.Ref_Runner_Labels:
				// the job is targeted by runner labels

				// if _any_ label matches, include it
				reqLabels, ok := req.TargetRunner.Target.(*pb.Ref_Runner_Labels)
				if !ok {
					// Request was not for target runner by labels
					continue
				}

				// look for any matching label from the request on the job
				match := false
				for key, value := range reqLabels.Labels.Labels {
					v, ok := tr.Labels.Labels[key]
					if !ok {
						// requested key not found in job, continue searching through label loop
						continue
					}
					if v == value {
						// a key was found, and its value matches
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
		}

		if req.Pipeline != nil {
			if job.Pipeline == nil {
				continue
			}
			// check whether the job is a match for pipeline name or id, and run sequence (if specified).
			if req.Pipeline.RunSequence != 0 && (req.Pipeline.RunSequence != job.Pipeline.RunSequence) {
				continue
			}
			if req.Pipeline.PipelineName != "" && (req.Pipeline.PipelineName != job.Pipeline.PipelineName) {
				continue
			}
			if req.Pipeline.PipelineId != "" && (req.Pipeline.PipelineId != job.Pipeline.PipelineId) {
				continue
			}
			if (req.Pipeline.PipelineId == job.Pipeline.PipelineId || req.Pipeline.PipelineName == job.Pipeline.PipelineName) && req.Pipeline.RunSequence != job.Pipeline.RunSequence {
				continue
			}
		}

		result = append(result, job)
	}

	return result, nil
}

// JobById looks up a job by ID. The returned Job will be a deep copy
// of the job so it is safe to read/write. If the job can't be found,
// a nil result with no error is returned.
func (s *State) JobById(id string, ws memdb.WatchSet) (*serverstate.Job, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	watchCh, raw, err := memTxn.FirstWatch(jobTableName, jobIdIndexName, id)
	if err != nil {
		return nil, err
	}

	ws.Add(watchCh)

	if raw == nil {
		return nil, nil
	}
	jobIdx := raw.(*jobIndex)

	// Get blocked status if it is queued.
	var blocked bool
	if jobIdx.State == pb.Job_QUEUED {
		blocked, err = s.jobIsBlocked(memTxn, jobIdx, ws)
		if err != nil {
			return nil, err
		}
	}

	var job *pb.Job
	err = s.db.View(func(dbTxn *bolt.Tx) error {
		job, err = s.jobById(dbTxn, jobIdx.Id)
		return err
	})

	result := jobIdx.Job(job)
	result.Blocked = blocked

	return result, err
}

// JobPeekForRunner effectively simulates JobAssignForRunner with two changes:
// (1) jobs are not actually assigned (they remain queued) and (2) this will
// not block if a job isn't available. If a job isn't available, this will
// return (nil, nil).
func (s *State) JobPeekForRunner(ctx context.Context, r *pb.Runner) (*serverstate.Job, error) {
	// The false,false here will (1) not block and (2) not assign
	return s.jobAssignForRunner(ctx, r, false, false)
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
func (s *State) JobAssignForRunner(ctx context.Context, r *pb.Runner) (*serverstate.Job, error) {
	return s.jobAssignForRunner(ctx, r, true, true)
}

func (s *State) jobAssignForRunner(ctx context.Context, r *pb.Runner, block, assign bool) (*serverstate.Job, error) {
	var txn *memdb.Txn

RETRY_ASSIGN:
	// If our transaction is not nil that means this is a repeated time around.
	// If we aren't blocking, return now.
	if txn != nil && !block {
		return nil, nil
	}

	// WatchSet we'll trigger to retry assignment
	ws := memdb.NewWatchSet()

	// If our runner exists, it must be adopted. We're lax right about allowing
	// runners that don't exist for JobPeek. If it doesn't exist, callers should
	// handle this. For example, RunnerJobStream handles this itself.
	rCheck, err := s.RunnerById(r.Id, ws)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if rCheck != nil {
		if rCheck.AdoptionState != pb.Runner_ADOPTED && rCheck.AdoptionState != pb.Runner_PREADOPTED {
			return nil, status.Errorf(codes.FailedPrecondition,
				"cannot assign jobs to a runner that isn't adopted")
		}
	}

	txn = s.inmem.Txn(false)
	defer txn.Abort()

	// Turn our runner into a runner record so we can more efficiently assign
	runnerIdx := newRunnerIndex(r)

	// candidateQuery finds candidate jobs to assign.
	type candidateFunc func(*memdb.Txn, memdb.WatchSet, *runnerIndex, bool) (*jobIndex, error)
	candidateQuery := []candidateFunc{
		s.jobCandidateById,
		s.jobCandidateByLabels,
		s.jobCandidateAny,
	}

	// If the runner is by id only, then explicitly set it to by id only.
	// We explicitly set the full list so that if we add more candidate
	// searches in the future, we're unlikely to break this.
	if r.ByIdOnly {
		candidateQuery = []candidateFunc{s.jobCandidateById}
	}

	// Build the list of candidates
	var candidates []*jobIndex
	for _, f := range candidateQuery {
		job, err := f(txn, ws, runnerIdx, assign)
		if err != nil {
			return nil, err
		}
		if job == nil {
			continue
		}

		candidates = append(candidates, job)
	}

	// If we have no candidates, then we have to wait for a job to show up.
	// We set up a blocking query on the job table for a non-assigned job.
	if len(candidates) == 0 {
		iter, err := txn.Get(jobTableName, jobStateIndexName, pb.Job_QUEUED)
		if err != nil {
			return nil, err
		}

		ws.Add(iter.WatchCh())
	}

	// We're done reading so abort the transaction
	txn.Abort()

	// If we have a watch channel set that means we didn't find any
	// results and we need to retry after waiting for changes.
	if len(candidates) == 0 {
		if block {
			ws.WatchCtx(ctx)
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}

		goto RETRY_ASSIGN
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
	//
	// Note: we only grab a write lock if we're assigning. If we're not
	// assigning then we grab a read lock.
	txn = s.inmem.Txn(assign)
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
		job := raw.(*jobIndex)
		if job == nil || job.State != pb.Job_QUEUED {
			continue
		}

		// If we've been requested to not assign, then we found our result.
		//
		// Importantly we do this above the blocked check below. This is tested.
		// We want to check this before checking if the job is blocked because
		// it is possible for all jobs to be blocked and therefore peek to return
		// nil. We want to return any job that could possibly be next.
		if !assign {
			// We're no longer going to use the memdb txn
			txn.Abort()

			var pbjob *pb.Job
			err = s.db.View(func(dbTxn *bolt.Tx) error {
				pbjob, err = s.jobById(dbTxn, job.Id)
				return err
			})
			if err != nil {
				return nil, err
			}

			return job.Job(pbjob), nil
		}

		// We also need to recheck that we aren't blocked. If we're blocked
		// now then we need to skip this job.
		if blocked, err := s.jobIsBlocked(txn, job, nil); blocked {
			continue
		} else if err != nil {
			return nil, err
		}

		// We're now modifying this job, so perform a copy
		job = job.Copy()

		// Update our state and update our on-disk job
		job.State = pb.Job_WAITING
		result, err := s.jobReadAndUpdate(txn, job.Id, func(jobpb *pb.Job) error {
			jobpb.State = job.State
			jobpb.AssignTime = timestamppb.Now()
			jobpb.AssignedRunner = &pb.Ref_RunnerId{Id: r.Id}

			return nil
		})
		if err != nil {
			return nil, err
		}

		// Create our timer to requeue this if it isn't acked
		job.StateTimer = time.AfterFunc(serverstate.JobWaitingTimeout, func() {
			s.log.Info("job ack timer expired", "job", job.Id, "timeout", serverstate.JobWaitingTimeout)
			s.JobAck(job.Id, false)
		})

		if err := txn.Insert(jobTableName, job); err != nil {
			return nil, err
		}

		// Update our assignment state
		if err := s.jobAssignedSet(txn, job, true); err != nil {
			s.JobAck(job.Id, false)
			return nil, err
		}

		txn.Commit()
		return job.Job(result), nil
	}
	txn.Abort()

	// If we reached here, all of our candidates were invalid, we retry
	goto RETRY_ASSIGN
}

// JobAck acknowledges that a job has been accepted or rejected by the runner.
// If ack is false, then this will move the job back to the queued state
// and be eligible for assignment.
// Additionally, if a job is associated with an on-demand runner task and/or pipeline,
// this func will progress the Tasks and Pipeline state machines depending on which job
// has currently been acked.
func (s *State) JobAck(id string, ack bool) (*serverstate.Job, error) {
	ctx := context.Background()
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Get the job
	raw, err := txn.First(jobTableName, jobIdIndexName, id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "job not found: %s", id)
	}
	job := raw.(*jobIndex)

	// If the job is not in the assigned state, then this is an error.
	if job.State != pb.Job_WAITING {
		return nil, status.Errorf(codes.FailedPrecondition,
			"job can't be acked from state: %s",
			job.State.String())
	}

	// We're now modifying this job, so perform a copy
	job = job.Copy()

	result, err := s.jobReadAndUpdate(txn, job.Id, func(jobpb *pb.Job) error {
		if ack {
			// Set to accepted
			job.State = pb.Job_RUNNING
			jobpb.State = job.State
			jobpb.AckTime = timestamppb.Now()

			// We also initialize the output buffer here because we can
			// expect output to begin streaming in.
			job.OutputBuffer = logbuffer.New()
		} else {
			// Set to queued
			job.State = pb.Job_QUEUED
			jobpb.State = job.State
			jobpb.AssignTime = nil
			jobpb.AssignedRunner = nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Cancel our timer
	if job.StateTimer != nil {
		job.StateTimer.Stop()
		job.StateTimer = nil
	}

	// Create a new timer that we'll use for our heartbeat. After this
	// timer expires, the job will immediately move to an error state.
	job.StateTimer = time.AfterFunc(serverstate.JobHeartbeatTimeout, func() {
		s.log.Info("canceling job due to heartbeat timeout", "job", job.Id)
		// Force cancel
		if err := s.JobCancel(job.Id, true); err != nil {
			s.log.Error("error canceling job due to heartbeat failure", "error", err, "job", job.Id)
		}
	})

	s.log.Debug("heartbeat timer set", "job", job.Id, "timeout", serverstate.JobHeartbeatTimeout)

	// Insert to update
	if err := txn.Insert(jobTableName, job); err != nil {
		return nil, err
	}

	// Update our assigned state if we nacked
	if !ack {
		if err := s.jobAssignedSet(txn, job, false); err != nil {
			return nil, err
		}
	}

	txn.Commit()

	// Update the task state machine if acked job is part of an on-demand runner task.
	if err := s.taskAck(ctx, job.Id); err != nil {
		s.log.Error("error updating task state", "error", err, "job", job.Id)
		return nil, err
	}
	if err := s.pipelineAck(ctx, job.Id); err != nil {
		s.log.Error("error updating pipeline state", "error", err, "job", job.Id)
		return nil, err
	}
	return job.Job(result), nil
}

// JobUpdateRef sets the data_source_ref field for a job. This job can be
// in any state.
func (s *State) JobUpdateRef(id string, ref *pb.Job_DataSource_Ref) error {
	return s.JobUpdate(id, func(jobpb *pb.Job) error {
		jobpb.DataSourceRef = ref
		return nil
	})
}

// JobUpdateExpiry will update the jobs expiration time with the new value. This
// method is used when a runner has accepted a job and peeks at its runner
// config. By this point, we know the runner has accepted a job to work on,
// so the runner should handle the job soon.
func (s *State) JobUpdateExpiry(id string, newExpire *timestamppb.Timestamp) error {
	return s.JobUpdate(id, func(jobpb *pb.Job) error {
		jobpb.ExpireTime = newExpire
		return nil
	})
}

// JobUpdate calls the given callback to update fields on the job data.
// The callback is called in the context of a database write lock so it
// should NOT compute anything and should be fast. The callback can return
// an error to abort the transaction.
//
// Job states should NOT be modified using this, only metadata.
func (s *State) JobUpdate(id string, cb func(jobpb *pb.Job) error) error {
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
	job := raw.(*jobIndex)

	_, err = s.jobReadAndUpdate(txn, job.Id, func(jobpb *pb.Job) error {
		return cb(jobpb)
	})
	if err != nil {
		return err
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
func (s *State) JobComplete(id string, result *pb.Job_Result, cerr error) error {
	ctx := context.Background()
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
	job := raw.(*jobIndex)

	// Update our assigned state
	if err := s.jobAssignedSet(txn, job, false); err != nil {
		return err
	}

	// If the job is not in the assigned state, then this is an error.
	if job.State != pb.Job_RUNNING {
		return status.Errorf(codes.FailedPrecondition,
			"job can't be completed from state: %s",
			job.State.String())
	}

	// We're now modifying this job, so perform a copy
	job = job.Copy()

	_, err = s.jobReadAndUpdate(txn, job.Id, func(jobpb *pb.Job) error {
		// Set to complete, assume success for now
		job.State = pb.Job_SUCCESS
		jobpb.State = job.State
		jobpb.Result = result
		jobpb.CompleteTime = timestamppb.Now()

		if cerr != nil {
			job.State = pb.Job_ERROR
			jobpb.State = job.State

			st, _ := status.FromError(cerr)
			jobpb.Error = st.Proto()
		}

		return nil
	})
	if err != nil {
		return err
	}

	// End the job
	job.End()

	// Insert to update
	if err := txn.Insert(jobTableName, job); err != nil {
		return err
	}

	txn.Commit()

	// If the job is part of an on-demand runner task, update the task state machine
	// to mark the job in the task as complete
	if err := s.taskComplete(ctx, job.Id); err != nil {
		s.log.Error("error updating task state for complete", "error", err, "job", job.Id)
		return err
	}

	// If the job is part of a pipeline, update the pipeline state machine
	// to mark the pipeline as complete
	if err := s.pipelineComplete(ctx, job.Id); err != nil {
		s.log.Error("error updating pipeline state for complete", "error", err, "job", job.Id)
		return err
	}

	return nil
}

// JobCancel marks a job as cancelled. This will set the internal state
// and request the cancel but if the job is running then it is up to downstream
// to listen for and react to Job changes for cancellation.
func (s *State) JobCancel(id string, force bool) error {
	ctx := context.Background()
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
	job := raw.(*jobIndex)

	if err := s.jobCancel(txn, job, force); err != nil {
		return err
	}

	txn.Commit()

	if err := s.pipelineCancel(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *State) jobCancel(txn *memdb.Txn, job *jobIndex, force bool) error {
	job = job.Copy()
	oldState := job.State

	s.log.Debug("attempting to cancel job", "job", job.Id)

	// How we handle cancel depends on the state
	switch job.State {
	case pb.Job_ERROR, pb.Job_SUCCESS:
		s.log.Debug("attempted to cancel completed job", "state", job.State.String(), "job", job.Id)
		// Jobs that are already completed do nothing for cancellation.
		// We do not mark that they were requested as cancelled since they
		// completed fine.
		return nil

	case pb.Job_QUEUED:
		// For queued jobs, we immediately transition them to an error state.
		job.State = pb.Job_ERROR

	case pb.Job_WAITING, pb.Job_RUNNING:
		// For these states, we just need to mark it as cancelled and have
		// downstream listeners complete the job. However, if we are forcing
		// then we immediately transition to error.
		if force {
			job.State = pb.Job_ERROR
			job.End()
		}
	}

	s.log.Debug("changing job state for cancel", "old-state", oldState.String(), "new-state", job.State.String(), "job", job.Id, "force", force)

	if force && job.State == pb.Job_ERROR {
		// Update our assigned state to unblock future jobs
		if err := s.jobAssignedSet(txn, job, false); err != nil {
			return err
		}
	}

	// Persist the on-disk data
	_, err := s.jobReadAndUpdate(txn, job.Id, func(jobpb *pb.Job) error {
		jobpb.State = job.State
		jobpb.CancelTime = timestamppb.Now()

		// If we transitioned to the error state we note that we were force
		// cancelled. We can only be in the error state under that scenario
		// since otherwise we would've returned early.
		if jobpb.State == pb.Job_ERROR {
			jobpb.Error = status.New(codes.Canceled, "canceled").Proto()
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Store the inmem data
	// This will be seen by a currently running RunnerJobStream goroutine, which
	// will then see that the job has been canceled and send the request to cancel
	// down to the runner.
	if err := txn.Insert(jobTableName, job); err != nil {
		return err
	}

	return nil
}

// JobHeartbeat resets the heartbeat timer for a running job. If the job
// is not currently running this does nothing, it will not return an error.
// If the job doesn't exist then this will return an error.
func (s *State) JobHeartbeat(id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	if err := s.jobHeartbeat(txn, id); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *State) jobHeartbeat(txn *memdb.Txn, id string) error {
	// Get the job
	raw, err := txn.First(jobTableName, jobIdIndexName, id)
	if err != nil {
		return err
	}
	if raw == nil {
		return status.Errorf(codes.NotFound, "job not found: %s", id)
	}
	job := raw.(*jobIndex)

	// If the job is not in the running state, we do nothing.
	if job.State != pb.Job_RUNNING {
		return nil
	}

	// If the state timer is nil... that is weird but we ignore it here.
	// It is up to other parts of the job system to ensure a running
	// job has a heartbeat timer.
	if job.StateTimer == nil {
		s.log.Info("job with no start timer detected", "job", id)
		return nil
	}

	// Reset the timer
	job.StateTimer.Reset(serverstate.JobHeartbeatTimeout)

	return nil
}

// JobExpire expires a job. This will cancel the job if it is still queued.
func (s *State) JobExpire(id string) error {
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
	job := raw.(*jobIndex)

	// How we handle depends on the state
	switch job.State {
	case pb.Job_QUEUED, pb.Job_WAITING:
		if err := s.jobCancel(txn, job, false); err != nil {
			return err
		}

	default:
	}

	txn.Commit()
	return nil
}

// JobIsAssignable returns whether there is a registered runner that
// meets the requirements to run this job.
//
// If this returns true, the job if queued should eventually be assigned
// successfully to a runner. An assignable result does NOT mean that it will be
// in queue a short amount of time.
//
// Note the result is a point-in-time result. If the only candidate runners
// deregister between this returning true and queueing, the job may still
// sit in a queue indefinitely.
func (s *State) JobIsAssignable(ctx context.Context, jobpb *pb.Job) (bool, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	// If we have no runners, we cannot be assigned
	empty, err := s.runnerEmpty(memTxn)
	if err != nil {
		return false, err
	}
	if empty {
		return false, nil
	}

	// If we have a special targeting constraint, that has to be met
	var iter memdb.ResultIterator
	var targetCheck func(*pb.Runner) (bool, error)
	switch v := jobpb.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Any:
		// We need a special target check that disallows by ID only
		targetCheck = func(r *pb.Runner) (bool, error) {
			return !r.ByIdOnly, nil
		}

		iter, err = memTxn.LowerBound(runnerTableName, runnerIdIndexName, "")

	case *pb.Ref_Runner_Id:
		iter, err = memTxn.Get(runnerTableName, runnerIdIndexName, v.Id.Id)

	case *pb.Ref_Runner_Labels:
		targetCheck = func(r *pb.Runner) (bool, error) {
			for k, v := range v.Labels.Labels {
				if val, ok := r.Labels[k]; ok && v != val {
					return false, nil
				}
			}
			return true, nil
		}
		iter, err = memTxn.LowerBound(runnerTableName, runnerIdIndexName, "")

	default:
		return false, fmt.Errorf("unknown runner target value: %#v", jobpb.TargetRunner.Target)
	}
	if err != nil {
		return false, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			// We're out of candidates and we found none.
			return false, nil
		}
		runner := raw.(*runnerIndex)

		// Check our target-specific check
		if targetCheck != nil {
			check, err := targetCheck(runner.Runner)
			if err != nil {
				return false, err
			}
			if !check {
				continue
			}
		}

		// This works!
		return true, nil
	}
}

// jobIndexInit initializes the config index from persisted data.
func (s *State) jobIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(jobBucket)
	c := bucket.Cursor()

	var cnt int

	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		var value pb.Job
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		// if we still have headroom for more indexed jobs OR the job hasn't finished yet,
		// index it.
		if cnt < maximumJobsIndexed || !jobIsCompleted(value.State) {
			cnt++
			idx, err := s.jobIndexSet(memTxn, k, &value)
			if err != nil {
				return err
			}

			// If the job was running or waiting, set it as assigned.
			if value.State == pb.Job_RUNNING || value.State == pb.Job_WAITING {
				if err := s.jobAssignedSet(memTxn, idx, true); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// jobIndexSet writes an index record for a single job.
func (s *State) jobIndexSet(txn *memdb.Txn, id []byte, jobpb *pb.Job) (*jobIndex, error) {
	rec := &jobIndex{
		Id:          jobpb.Id,
		SingletonId: jobpb.SingletonId,
		DependsOn:   jobpb.DependsOn,
		State:       jobpb.State,
		Application: jobpb.Application,
		Workspace:   jobpb.Workspace,
		OpType:      reflect.TypeOf(jobpb.Operation),
	}

	// Target
	if jobpb.TargetRunner == nil {
		return nil, fmt.Errorf("job target runner must be set")
	}
	switch v := jobpb.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Any:
		rec.TargetAny = true

	case *pb.Ref_Runner_Id:
		rec.TargetRunnerId = v.Id.Id

	case *pb.Ref_Runner_Labels:
		rec.TargetRunnerLabels = v.Labels.Labels

	default:
		return nil, fmt.Errorf("unknown runner target value: %#v", jobpb.TargetRunner.Target)
	}

	// Timestamps
	timestamps := []struct {
		Field *time.Time
		Src   *timestamppb.Timestamp
	}{
		{&rec.QueueTime, jobpb.QueueTime},
	}
	for _, ts := range timestamps {
		*ts.Field = ts.Src.AsTime()
	}

	// If this job is assigned. Then we have to start a nacking timer.
	// We reset the nack timer so it gives runners time to reconnect.
	if rec.State == pb.Job_WAITING {
		// Create our timer to requeue this if it isn't acked
		rec.StateTimer = time.AfterFunc(serverstate.JobWaitingTimeout, func() {
			s.JobAck(rec.Id, false)
		})
	}

	// If this job is running, we need to restart a heartbeat timeout.
	// This should only happen on reinit. This is tested.
	if rec.State == pb.Job_RUNNING {
		rec.StateTimer = time.AfterFunc(serverstate.JobHeartbeatTimeout, func() {
			// Force cancel
			s.JobCancel(rec.Id, true)
		})
	}

	// If we have an expiry, we need to set a timer to expire this job.
	if jobpb.ExpireTime != nil {
		now := time.Now()

		t := jobpb.ExpireTime.AsTime()

		dur := t.Sub(now)
		if dur < 0 {
			dur = 1
		}

		time.AfterFunc(dur, func() { s.JobExpire(jobpb.Id) })
	}

	// Insert the index
	return rec, txn.Insert(jobTableName, rec)
}

func (s *State) jobCreate(dbTxn *bolt.Tx, memTxn *memdb.Txn, jobpb *pb.Job) error {
	// Setup our initial job state
	var err error
	jobpb.State = pb.Job_QUEUED
	jobpb.QueueTime = timestamppb.Now()

	id := []byte(jobpb.Id)
	bucket := dbTxn.Bucket(jobBucket)

	// If singleton ID is set, we need to delete (cancel) any previous job
	// with the same singleton ID if it is still queued.
	if jobpb.SingletonId != "" {
		result, err := memTxn.First(
			jobTableName,
			jobSingletonIdIndexName,
			jobpb.SingletonId,
			pb.Job_QUEUED,
		)
		if err != nil {
			return err
		}

		if result != nil {
			// Note we don't have to worry about jobAssignedSet here because
			// we only run this block of code if the job is in the QUEUED state.

			// Note we don't need to Copy here like other places because
			// we never modify this old jobIndex in-place.
			old := result.(*jobIndex)

			oldpb, err := s.jobById(dbTxn, old.Id)
			if err != nil {
				return err
			}
			oldpb.State = pb.Job_ERROR
			oldpb.Error = status.Newf(codes.Canceled,
				"replaced by job %s", id).Proto()

			// Update the index and data
			if err := dbPut(bucket, []byte(old.Id), oldpb); err != nil {
				return err
			}
			if _, err = s.jobIndexSet(memTxn, []byte(old.Id), oldpb); err != nil {
				return err
			}

			// Copy the queue time from the old one so we retain our position
			jobpb.QueueTime = oldpb.QueueTime
		}
	}

	// If we have dependencies, they all must exist.
	dependsMap := map[string]struct{}{}
	if len(jobpb.DependsOn) > 0 {
		// Let's remove any duplicates
		jobpb.DependsOn = uniqueStr(jobpb.DependsOn)

		// Go through and ensure that each exists.
		for _, id := range jobpb.DependsOn {
			if id == "" {
				return status.Errorf(codes.FailedPrecondition,
					"job %q has an empty string id for depends_on", jobpb.Id)
			}

			dependsMap[id] = struct{}{}

			_, err := s.jobById(dbTxn, id)
			if err != nil {
				return err
			}
		}
	}

	// Allowed failures must be in the DependsOn list
	for _, id := range jobpb.DependsOnAllowFailure {
		if _, ok := dependsMap[id]; !ok {
			return status.Errorf(codes.FailedPrecondition,
				"job %q in depends_on_allow_failure but not in depends_on",
				id,
			)
		}
	}

	// Insert into bolt
	if err := dbPut(dbTxn.Bucket(jobBucket), id, jobpb); err != nil {
		return err
	}

	// Insert into the DB
	_, err = s.jobIndexSet(memTxn, id, jobpb)
	if err != nil {
		return err
	}

	s.pruneMu.Lock()
	defer s.pruneMu.Unlock()

	s.indexedJobs++

	return nil
}

func (s *State) jobsPruneOld(memTxn *memdb.Txn, max int) (int, error) {
	return pruneOld(memTxn, pruneOp{
		lock:      &s.pruneMu,
		table:     jobTableName,
		index:     jobQueueTimeIndexName,
		indexArgs: []interface{}{pb.Job_QUEUED, time.Unix(0, 0)},
		max:       max,
		cur:       &s.indexedJobs,
		check: func(raw interface{}) bool {
			job := raw.(*jobIndex)
			return !jobIsCompleted(job.State)
		},
	})
}

func (s *State) jobById(dbTxn *bolt.Tx, id string) (*pb.Job, error) {
	var result pb.Job
	b := dbTxn.Bucket(jobBucket)
	return &result, dbGet(b, []byte(id), &result)
}

func (s *State) jobReadAndUpdate(
	memTxn *memdb.Txn, id string, f func(*pb.Job) error) (*pb.Job, error) {
	var result *pb.Job
	var err error
	return result, s.db.Update(func(dbTxn *bolt.Tx) error {
		result, err = s.jobById(dbTxn, id)
		if err != nil {
			return err
		}

		// Modify
		if err := f(result); err != nil {
			return err
		}

		// Cascade our state if we have to. We do this within the database
		// transaction so that error updates the full dependent chain or
		// nothing.
		if err := s.jobCascadeDependentState(memTxn, dbTxn, result); err != nil {
			return err
		}

		// Commit
		return dbPut(dbTxn.Bucket(jobBucket), []byte(id), result)
	})
}

// jobCascadeDependentState cascades certain state changes to all dependent
// jobs. For example, if a job transitions to an error state, then all
// dependents must also error.
func (s *State) jobCascadeDependentState(
	memTxn *memdb.Txn,
	dbTxn *bolt.Tx,
	parent *pb.Job,
) error {
	// If the state isn't error, we don't cascade anything.
	if parent.State != pb.Job_ERROR {
		return nil
	}

	// Look for any dependents
	iter, err := memTxn.Get(jobTableName, jobDependsOnIndexName, parent.Id)
	if err != nil {
		return err
	}

	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*jobIndex)

		// The state SHOULD be queued if the job is still waiting on us.
		// We only handle this scenario. If it isn't queued, we just ignore
		// it because we're in a state we don't understand.
		if idx.State != pb.Job_QUEUED {
			continue
		}

		// Read this dependent job from disk so we can update the proto struct.
		job, err := s.jobById(dbTxn, idx.Id)
		if status.Code(err) == codes.NotFound {
			continue
		}
		if err != nil {
			return err
		}

		// If this dependency is allowed to fail, do not cascade.
		failAllowed := false
		for _, id := range job.DependsOnAllowFailure {
			if strings.EqualFold(id, parent.Id) {
				failAllowed = true
				break
			}
		}
		if failAllowed {
			continue
		}

		// Cascade the error state
		job.State = pb.Job_ERROR
		job.Error = status.New(codes.Canceled, fmt.Sprintf(
			"Job dependency %q errored: %s",
			parent.Id,
			parent.Error.Message,
		)).Proto()
		job.CancelTime = parent.CancelTime

		// Write to disk
		if err := dbPut(dbTxn.Bucket(jobBucket), []byte(job.Id), job); err != nil {
			return err
		}

		// Write to memory index. Error state is terminal so we also call End.
		idx = idx.Copy()
		idx.State = job.State
		idx.End()
		if err := memTxn.Insert(jobTableName, idx); err != nil {
			return err
		}

		// Cascade further if necessary
		if err := s.jobCascadeDependentState(memTxn, dbTxn, job); err != nil {
			return err
		}
	}

	return nil
}

// jobCandidateById returns the most promising candidate job to assign
// that is targeting a specific runner by ID.
func (s *State) jobCandidateById(
	memTxn *memdb.Txn, ws memdb.WatchSet, r *runnerIndex, assign bool,
) (*jobIndex, error) {
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

		job := raw.(*jobIndex)
		if job.State != pb.Job_QUEUED || job.TargetRunnerId != r.Id {
			continue
		}

		// If this job is blocked, it is not a candidate.
		if blocked, err := s.jobIsBlocked(memTxn, job, ws); err != nil {
			return nil, err
		} else if blocked && assign {
			continue
		}

		return job, nil
	}

	return nil, nil
}

// jobCandidateByLabels returns the most promising candidate job to assign
// that is targeting a specific runner by labels.
func (s *State) jobCandidateByLabels(
	memTxn *memdb.Txn, ws memdb.WatchSet, r *runnerIndex, assign bool,
) (*jobIndex, error) {
	// NOTE(xx): This query forces us to search all queued jobs for matching labels.
	// A more efficient query for label searching would be preferable in the future.
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

		job := raw.(*jobIndex)
		if job.State != pb.Job_QUEUED || job.TargetRunnerLabels == nil {
			continue
		}

		// If this job is blocked, it is not a candidate.
		if blocked, err := s.jobIsBlocked(memTxn, job, ws); err != nil {
			return nil, err
		} else if blocked && assign {
			continue
		}

		// Check whether job target labels match with runner labels
		match := true
		for k, v := range job.TargetRunnerLabels {
			if val, ok := r.Runner.Labels[k]; !ok || v != val {
				match = false
				break
			}
		}
		if !match {
			continue
		}

		return job, nil
	}

	return nil, nil
}

// jobCandidateAny returns the first candidate job that targets any runner.
func (s *State) jobCandidateAny(
	memTxn *memdb.Txn, ws memdb.WatchSet, r *runnerIndex, assign bool,
) (*jobIndex, error) {
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

		job := raw.(*jobIndex)
		if job.State != pb.Job_QUEUED || !job.TargetAny {
			continue
		}

		// If this job is blocked, it is not a candidate.
		if blocked, err := s.jobIsBlocked(memTxn, job, ws); err != nil {
			return nil, err
		} else if blocked && assign {
			continue
		}

		return job, nil
	}

	return nil, nil
}

// taskAck checks if the referenced job id has a Task ref associated with it,
// and if so, Ack the specific job inside the Task job triple to progress the
// Task state machine.
func (s *State) taskAck(ctx context.Context, jobId string) error {
	job, err := s.JobById(jobId, nil)
	if err != nil {
		s.log.Error("error getting job by id", "job", jobId, "err", err)
		return err
	} else if job.Task == nil {
		s.log.Trace("job is not an on-demand runner task", "job", jobId)
		return nil
	}

	// grab the full task message based on the Task Ref on the job
	task, err := s.TaskGet(ctx, job.Task)
	if err != nil {
		s.log.Error("failed to get task to ack job", "job", job.Id, "task", job.Task.Ref)
		return err
	}

	// now figure out which job has been acked for the task, and update
	// the task state
	switch job.Id {
	case task.StartJob.Id:
		// StartJob has been Acked
		task.JobState = pb.Task_STARTING
		s.log.Trace("start job is starting", "job", job.Id, "task", task.Id)
	case task.TaskJob.Id:
		// TaskJob has been Acked
		task.JobState = pb.Task_RUNNING
		s.log.Trace("task job is running", "job", job.Id, "task", task.Id)
	case task.StopJob.Id:
		// StopJob has been Acked
		task.JobState = pb.Task_STOPPING
		s.log.Trace("stop job is running", "job", job.Id, "task", task.Id)
	case task.WatchJob.Id:
		s.log.Trace("ignoring watch task job")
	default:
		return status.Errorf(codes.Internal, "no task job id matches the requested job id %q", job.Id)
	}

	// TaskPut the new state
	if err := s.TaskPut(ctx, task); err != nil {
		s.log.Error("failed to ack task state", "job", job.Id, "task", task.Id)
		return err
	}

	return nil
}

// taskComplete will look up the referenced job to see if it has a Task ref
// associated with it, and if so, mark the specific job inside the Task job triple
// as complete to progress the task state machine.
func (s *State) taskComplete(ctx context.Context, jobId string) error {
	job, err := s.JobById(jobId, nil)
	if err != nil {
		s.log.Error("error getting job by id", "job", jobId, "err", err)
		return err
	} else if job.Task == nil {
		s.log.Trace("job is not an on-demand runner task", "job", jobId)
		return nil
	}

	// grab the full task message based on the Task Ref on the job
	task, err := s.TaskGet(ctx, job.Task)
	if err != nil {
		s.log.Error("failed to get task to mark complete", "job", job.Id, "task", job.Task.Ref)
		return err
	}

	// now figure out which job has been completed this is for the task, and update
	// the task state
	switch job.Id {
	case task.StartJob.Id:
		// StartJob has completed
		task.JobState = pb.Task_STARTED
		s.log.Trace("start job has completed", "job", job.Id, "task", task.Id)
	case task.TaskJob.Id:
		// TaskJob has completed
		task.JobState = pb.Task_COMPLETED
		s.log.Trace("task job has completed", "job", job.Id, "task", task.Id)
	case task.StopJob.Id:
		// StopJob has completed, the whole task is finished
		task.JobState = pb.Task_STOPPED
		s.log.Trace("stop job has completed", "job", job.Id, "task", task.Id)
	case task.WatchJob.Id:
		s.log.Trace("ignoring watch task job")
	default:
		return status.Errorf(codes.Internal, "no task job id matches the requested job id %q", job.Id)
	}

	// TaskPut the new state
	if err := s.TaskPut(ctx, task); err != nil {
		s.log.Error("failed to complete task state", "job", job.Id, "task", task.Id)
		return err
	}

	return nil
}

// pipelineComplete will look up the referenced job to see if it has a PipelineStep ref
// associated with it, and if so, progress the pipelineRun state machine.
func (s *State) pipelineComplete(ctx context.Context, jobId string) error {
	// grab the job
	job, err := s.JobById(jobId, nil)
	if err != nil {
		s.log.Error("error getting job by id", "job", jobId, "err", err)
		return err
	} else if job.Pipeline == nil {
		s.log.Trace("job is not part of a pipeline", "job", jobId)
		return nil
	}
	// grab the pipeline run
	run, err := s.PipelineRunGetByJobId(ctx, jobId)
	if run == nil {
		return nil
	}
	if err != nil {
		s.log.Error("failed to retrieve pipeline to complete", "job", job.Id)
		return err
	}

	if job.State == pb.Job_ERROR {
		run.State = pb.PipelineRun_ERROR
	} else if job.State == pb.Job_SUCCESS {
		// Look at all job ids in a run and check if any are not SUCCESS
		runComplete := true
		for _, j := range run.Jobs {
			if j.Id == job.Id {
				continue
			}

			rj, err := s.JobById(j.Id, nil)
			if err != nil {
				return err
			}

			if rj.State != pb.Job_SUCCESS {
				runComplete = false
				break
			}
		}

		if runComplete {
			run.State = pb.PipelineRun_SUCCESS
			s.log.Trace("pipeline run is complete", "job", job.Id, "pipeline", job.Pipeline.PipelineId, "run", run.Sequence)
		}
	}

	// PipelineRunPut the new state
	if err = s.PipelineRunPut(ctx, run); err != nil {
		s.log.Error("failed to complete pipeline run", "job", job.Id, "pipeline", job.Pipeline.PipelineId, "run", job.Pipeline.RunSequence)
		return err
	}

	return nil
}

// pipelineAck checks if the referenced job id has a PipelineStep ref associated with it,
// and if so, progress the PipelineRun state machine.
func (s *State) pipelineAck(ctx context.Context, jobId string) error {
	// grab pipeline run that triggered the job based on the PipelineTask Ref
	run, err := s.PipelineRunGetByJobId(ctx, jobId)
	if run == nil {
		return nil
	}
	if err != nil {
		s.log.Error("failed to retrieve pipeline to complete", "job", jobId)
		return err
	}

	// Update the new pipeline run state if it's not already running
	if run.State != pb.PipelineRun_RUNNING {
		run.State = pb.PipelineRun_RUNNING
	}
	s.log.Trace("pipeline is running", "job", jobId, "pipeline", run.Pipeline, "run", run.Sequence)
	if err := s.PipelineRunPut(ctx, run); err != nil {
		s.log.Error("failed to ack pipeline run state", "job", jobId, "pipeline", run.Pipeline, "run", run.Sequence)
		return err
	}

	return nil
}

// pipelineCancel will look up the referenced job to see if it has a PipelineStep ref
// associated with it, and if so, progress the pipelineRun state machine.
func (s *State) pipelineCancel(ctx context.Context, jobId string) error {
	job, err := s.JobById(jobId, nil)
	if err != nil {
		s.log.Error("error getting job by id", "job", jobId, "err", err)
		return err
	} else if job.Pipeline == nil {
		s.log.Trace("job is not part of a pipeline", "job", jobId)
		return nil
	}

	run, err := s.PipelineRunGetByJobId(ctx, jobId)
	if run == nil {
		return nil
	}
	if err != nil {
		s.log.Error("failed to retrieve pipeline to complete", "job", job.Id)
		return err
	}

	if job.State == pb.Job_SUCCESS {
		return nil
	} else {
		run.State = pb.PipelineRun_CANCELLED
		s.log.Trace("pipeline run cancelled", "job", job.Id, "pipeline", job.Pipeline.PipelineId, "run", run.Sequence)
	}
	// PipelineRunPut the new state
	if err := s.PipelineRunPut(ctx, run); err != nil {
		s.log.Error("failed to cancel pipeline run", "job", job.Id, "pipeline", job.Pipeline.PipelineId, "run", job.Pipeline.RunSequence)
		return err
	}

	return nil
}

// Copy should be called prior to any modifications to an existing jobIndex.
func (idx *jobIndex) Copy() *jobIndex {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}

// Job returns the Job for an index.
func (idx *jobIndex) Job(jobpb *pb.Job) *serverstate.Job {
	return &serverstate.Job{
		Job:          jobpb,
		OutputBuffer: idx.OutputBuffer,
	}
}

// End notes this job is complete and performs any cleanup on the index.
func (idx *jobIndex) End() {
	if idx.StateTimer != nil {
		idx.StateTimer.Stop()
		idx.StateTimer = nil
	}
}

// uniqueStr is a little helper to ensure a string slice only has unique values.
func uniqueStr(s []string) []string {
	keys := map[string]struct{}{}
	result := make([]string, 0, len(s))
	for _, value := range s {
		if _, ok := keys[value]; !ok {
			keys[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}
