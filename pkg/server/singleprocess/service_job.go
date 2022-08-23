package singleprocess

import (
	"context"
	"math/rand"
	"reflect"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) GetJob(
	ctx context.Context,
	req *pb.GetJobRequest,
) (*pb.Job, error) {
	job, err := s.state(ctx).JobById(req.JobId, nil)
	if err != nil {
		return nil, err
	}
	if job == nil || job.Job == nil {
		return nil, status.Errorf(codes.NotFound, "job not found")
	}

	return job.Job, nil
}

func (s *Service) ListJobs(
	ctx context.Context,
	req *pb.ListJobsRequest,
) (*pb.ListJobsResponse, error) {
	jobs, err := s.state(ctx).JobList(req)
	if err != nil {
		return nil, err
	}

	return &pb.ListJobsResponse{
		Jobs: jobs,
	}, nil
}

func (s *Service) CancelJob(
	ctx context.Context,
	req *pb.CancelJobRequest,
) (*empty.Empty, error) {
	if err := s.state(ctx).JobCancel(req.JobId, req.Force); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// queueJobMulti queues multiple jobs transactionally. The return value
// are the responses for each in the same order as the requests.
//
// Postcondition: len(result) == len(req) iff err == nil
func (s *Service) queueJobMulti(
	ctx context.Context,
	req []*pb.QueueJobRequest,
) ([]*pb.QueueJobResponse, error) {
	jobQueue := make([]*pb.Job, 0, len(req)*4)
	jobIds := make([]string, 0, len(req))
	for _, single := range req {
		jobs, jobId, err := s.queueJobReqToJob(ctx, single)
		if err != nil {
			return nil, err
		}

		jobQueue = append(jobQueue, jobs...)
		jobIds = append(jobIds, jobId)
	}

	// Queue the jobs
	if err := s.state(ctx).JobCreate(jobQueue...); err != nil {
		return nil, err
	}

	// Get the response
	resp := make([]*pb.QueueJobResponse, len(jobIds))
	for i, id := range jobIds {
		resp[i] = &pb.QueueJobResponse{JobId: id}
	}

	return resp, nil
}

// queueJobReqToJob converts a QueueJobRequest to a job to queue, but
// does not queue it. This may return multiple jobs if the queue job
// request requires an on-demand runner. They should all be queued
// atomically with JobCreate.
//
// Precondition: req parameter must be validated
func (s *Service) queueJobReqToJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) ([]*pb.Job, string, error) {
	log := hclog.FromContext(ctx)
	job := req.Job

	// Verify the project exists and use that to set the default data source
	log.Debug("checking job project", "project", job.Application.Project)
	project, err := s.state(ctx).ProjectGet(&pb.Ref_Project{Project: job.Application.Project})
	if status.Code(err) == codes.NotFound {
		return nil, "", status.Errorf(codes.NotFound,
			"Project %q was not found! Please ensure that 'waypoint init' was run with this project.",
			job.Application.Project,
		)
	}

	if job.DataSource == nil {
		if project.DataSource == nil {
			return nil, "", status.Errorf(codes.FailedPrecondition,
				"Project %s does not have a data source configured. Remote jobs "+
					"require a data source such as Git to be configured with the project. "+
					"Data sources can be configured via the CLI or UI. For help, see : "+
					"https://www.waypointproject.io/docs/projects/git#configuring-the-project",
				job.Application.Project,
			)
		}

		job.DataSource = project.DataSource
	}

	// Get the next id
	if job.Id == "" {
		id, err := server.Id()
		if err != nil {
			return nil, "", status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}
		job.Id = id
	}

	// Validate expiry if we have one
	job.ExpireTime = nil
	if req.ExpiresIn != "" {
		dur, err := time.ParseDuration(req.ExpiresIn)
		if err != nil {
			return nil, "", status.Errorf(codes.FailedPrecondition,
				"Invalid expiry duration: %s", err.Error())
		}
		job.ExpireTime = timestamppb.New(time.Now().Add(dur))
	}

	// If the job has any target runner, it is a remote job.
	// Use a default ODR profile if it doesn't already have one assigned.
	if _, ok := job.TargetRunner.Target.(*pb.Ref_Runner_Any); ok {
		if job.OndemandRunner == nil {
			ods, err := s.state(ctx).OnDemandRunnerConfigDefault()
			if err != nil {
				return nil, "", err
			}

			switch len(ods) {
			case 0:
				// ok, no on-demand runners
			case 1:
				job.OndemandRunner = ods[0]
			default:
				job.OndemandRunner = ods[rand.Intn(len(ods))]
				log.Debug("multiple default on-demand runner profiles detected, chose a random one",
					"runner-config-id", job.OndemandRunner.Id)
			}
		}
	}

	// If we have no ODR, our result is just the job
	result := []*pb.Job{job}

	// If we have an ODR profile, then we know that we should be using
	// an on-demand runner for this. Let's wrap the jobs so that we have
	// the full set of start to stop.
	if job.OndemandRunner != nil {
		result, err = s.wrapJobWithRunner(ctx, job)
		if err != nil {
			return nil, "", err
		}

		// If we are skipping, then the job we queued is the watch job.
		if job.OndemandRunnerTask != nil && job.OndemandRunnerTask.SkipOperation {
			for _, j := range result {
				if _, ok := j.Operation.(*pb.Job_WatchTask); ok {
					job = j
					break
				}
			}
		}
	}

	return result, job.Id, nil
}

func (s *Service) QueueJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) (*pb.QueueJobResponse, error) {
	if req.Job == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "job must be set")
	}
	if req.Job.Operation == nil {
		// We special case this check and return "Unimplemented" because
		// the primary case where operation is nil is if a client is sending
		// us an unsupported operation.
		return nil, status.Errorf(codes.Unimplemented, "operation is nil or unknown")
	}
	if err := serverptypes.ValidateJob(req.Job); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	jobs, jobId, err := s.queueJobReqToJob(ctx, req)
	if err != nil {
		return nil, err
	}

	// Queue the job
	if err := s.state(ctx).JobCreate(jobs...); err != nil {
		return nil, err
	}

	return &pb.QueueJobResponse{JobId: jobId}, nil
}

// wrapJobWithRunner takes a job and "wraps" it within an on-demand launched
// runner. This creates a dependency chain that ensures that the runner is
// started and stopped around the given job (hence "wraps").
//
// A diagram of the dependency chain created is shown below. The dashed
// border is the "source" job.
//
//           ┌────────────────┐
//           │   Start Task   │─────────────┐
//           └────────────────┘             │
//                    │                     │
//          ┌─────────┴─────────┐           │
//          ▼                   ▼           │
// ┌────────────────┐   ┌─ ── ── ── ── ──   │
// │   Watch Task   │   │      Job       │  │
// └────────────────┘   └ ── ── ── ── ── ┘  │
//          │                    │          │
//          └─────────┬──────────┘          │
//                    ▼                     │
//           ┌────────────────┐             │
//           │   Stop Task    │◀────────────┘
//           └────────────────┘
//
// Details:
//
//   - Start task launches the on-demand runner.
//   - After it is launched, "job" can run targeting the launched ODR.
//   - Simultaneously, the watch task watches the launched task and records
//     logs, exit code, etc.
//   - Finally, stop task is called to clean up the resources associated
//     with start.
//
func (s *Service) wrapJobWithRunner(
	ctx context.Context,
	source *pb.Job,
) ([]*pb.Job, error) {
	// Get the runner profile we're going to use for this runner.
	od, err := s.state(ctx).OnDemandRunnerConfigGet(source.OndemandRunner)
	if err != nil {
		return nil, err
	}
	if od == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"the on-demand runner config for id %q and job %q was nil",
			source.OndemandRunner.Id, source.Id)
	}

	// Determine if we're skipping this job. This is done for custom tasks.
	skip := source.OndemandRunnerTask != nil && source.OndemandRunnerTask.SkipOperation
	if skip {
		// We only allow noop operations to be skipped out of safety. These
		// make sense to skip, whereas skipping a build or deploy might be
		// a bug in the client.
		_, ok := source.Operation.(*pb.Job_Noop_)
		if !ok {
			return nil, status.Errorf(codes.FailedPrecondition,
				"only noop operations can be skipped with custom tasks, got %T",
				source.Operation)
		}
	}

	// Generate our job to start the ODR
	startJob, runnerId, err := s.onDemandRunnerStartJob(ctx, source, od)
	if err != nil {
		return nil, err
	}

	// Generate our job to watch the ODR
	watchJob, err := s.onDemandRunnerWatchJob(ctx, startJob, source, od)
	if err != nil {
		return nil, err
	}

	// Change our source job to run on the launched ODR.
	source.TargetRunner = &pb.Ref_Runner{
		Target: &pb.Ref_Runner_Id{
			Id: &pb.Ref_RunnerId{
				Id: runnerId,
			},
		},
	}

	// Our source job depends on the starting job.
	source.DependsOn = append(source.DependsOn, startJob.Id)

	// Job to stop the ODR
	stopJob, err := s.onDemandRunnerStopJob(ctx, startJob, watchJob, source, od)
	if err != nil {
		return nil, err
	}

	// For our task tracking, the primary job is usually the job. But
	// if we're skipping, then it is the watch task.
	sourceJob := source
	if skip {
		sourceJob = watchJob
	}

	// Write a Task state with the On-Demand Runner job triple
	task := &pb.Task{
		StartJob: &pb.Ref_Job{Id: startJob.Id},
		TaskJob:  &pb.Ref_Job{Id: sourceJob.Id},
		StopJob:  &pb.Ref_Job{Id: stopJob.Id},
		WatchJob: &pb.Ref_Job{Id: watchJob.Id},
		JobState: pb.Task_PENDING,
	}
	if skip {
		// If we're skipping, the primary task job becomes the watch.
		task.TaskJob = &pb.Ref_Job{Id: watchJob.Id}
	}
	if err := s.state(ctx).TaskPut(task); err != nil {
		return nil, err
	} else {
		task, err := s.state(ctx).TaskGet(&pb.Ref_Task{
			Ref: &pb.Ref_Task_JobId{
				JobId: sourceJob.Id,
			},
		})
		if err != nil {
			// could not find task that was just Put into the db!
			return nil, err
		}

		// assign a task ref to each job for lookup later
		taskRef := &pb.Ref_Task{
			Ref: &pb.Ref_Task_Id{
				Id: task.Id,
			},
		}

		startJob.Task = taskRef
		sourceJob.Task = taskRef
		stopJob.Task = taskRef
		watchJob.Task = taskRef
	}

	jobs := []*pb.Job{startJob, watchJob, stopJob}
	if !skip {
		jobs = append(jobs, sourceJob)
	}

	return jobs, nil
}

// onDemandRunnerStartJob generates a StartJob template for a Task.
func (s *Service) onDemandRunnerStartJob(
	ctx context.Context,
	source *pb.Job,
	od *pb.OnDemandRunnerConfig,
) (*pb.Job, string, error) {
	log := hclog.FromContext(ctx)

	if od == nil {
		return nil, "", status.Errorf(codes.FailedPrecondition,
			"the on-demand runner config for id %q and job %q was nil",
			source.OndemandRunner.Id, source.Id)
	}

	// Generate a unique ID for the runner
	runnerId, err := server.Id()
	if err != nil {
		return nil, "", err
	}
	log.Info("requesting ondemand runner via task start", "runner-id", runnerId)

	// Follow the same logic as RunnerGetDeploymentConfig to get our advertise
	// address for the runner to connect to the Waypoint server. Note:
	// Our addr for now is just the first one since we don't support
	// multiple addresses yet. In the future we will want to support more
	// advanced choicing.
	var addr *pb.ServerConfig_AdvertiseAddr
	serverConfig, err := s.GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to get server config to populate runner start job server addr")
	}

	cfg := serverConfig.Config

	// This should only happen during tests
	if len(cfg.AdvertiseAddrs) == 0 {
		log.Info("server has no advertise addrs, using localhost")
		addr = &pb.ServerConfig_AdvertiseAddr{
			Addr: "localhost:9701",
		}
	} else {
		addr = cfg.AdvertiseAddrs[0]
	}

	encodedDefaultUserId, err := s.encodeId(ctx, DefaultUserId)
	if err != nil {
		msg := "failed to encode the default user id when starting and ODR job"
		log.Error(msg, "id", DefaultUserId, "err", err)
		return nil, "", status.Error(codes.Internal, msg)
	}

	// We generate a new login token for each ondemand-runner used. This will inherit
	// the user of the token to be the user that queued the original job, which is
	// the correct behavior.
	token, err := s.newToken(ctx, 60*time.Minute, s.activeAuthKeyId, nil, &pb.Token{
		// TODO(emp) should this be a Token_Runner_?
		Kind: &pb.Token_Login_{Login: &pb.Token_Login{
			UserId: encodedDefaultUserId,
		}},
	})
	if err != nil {
		return nil, "", err
	}

	// Build up our env vars to connect to the server, and add in our runner ID.
	scfg := serverconfig.Client{
		Address:       addr.Addr,
		Tls:           addr.Tls,
		TlsSkipVerify: addr.TlsSkipVerify,
		RequireAuth:   true,
		AuthToken:     token,
	}
	envVars := scfg.EnvMap()
	envVars["WAYPOINT_RUNNER_ID"] = runnerId

	// Add any env vars that our profile overrides
	for k, v := range od.EnvironmentVariables {
		envVars[k] = v
	}

	// Build our task launch info.
	launchInfo := &pb.TaskLaunchInfo{}
	if override := source.OndemandRunnerTask; override != nil {
		if info := override.LaunchInfo; info != nil {
			launchInfo = info
		}
	}
	if launchInfo.OciUrl == "" {
		launchInfo.OciUrl = od.OciUrl

		// Arguments for the runner image. Waypoint is ALWAYS assumed to be
		// the entrypoint for ODR images if no custom one is specified.
		launchInfo.Arguments = []string{
			"runner", "agent", "-vv", "-id", runnerId, "-odr", "-odr-profile-id", od.Id,
		}
	}

	// We always default our env vars so that a custom image can still
	// behave like a runner and has access to the token, runner ID, etc.
	for k, v := range launchInfo.EnvironmentVariables {
		envVars[k] = v
	}
	launchInfo.EnvironmentVariables = envVars

	job := &pb.Job{
		// Inherit the workspace/application of the source job.
		Workspace:   source.Workspace,
		Application: source.Application,

		// Depend on the same dependencies as the source job. This way,
		// we don't start up the ODR very early when the job is not ready
		// to execute.
		DependsOn:             source.DependsOn,
		DependsOnAllowFailure: source.DependsOnAllowFailure,

		Operation: &pb.Job_StartTask{
			StartTask: &pb.Job_StartTaskLaunchOp{
				Params: &pb.Job_TaskPluginParams{
					PluginType: od.PluginType,
					HclConfig:  od.PluginConfig,
					HclFormat:  od.ConfigFormat,
				},
				Info: launchInfo,
			},
		},
	}

	// Get the next id for the job
	id, err := server.Id()
	if err != nil {
		return nil, "", status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	// This will be either "Any" or a specific static runner.
	job.TargetRunner = od.TargetRunner

	return job, runnerId, nil
}

// onDemandRunnerWatchJob generates a WatchJob template for a Task.
func (s *Service) onDemandRunnerWatchJob(
	ctx context.Context,
	startJob *pb.Job,
	source *pb.Job,
	od *pb.OnDemandRunnerConfig,
) (*pb.Job, error) {
	job := &pb.Job{
		// Inherit the workspace/application of the source job.
		Workspace:   source.Workspace,
		Application: source.Application,

		// We depend on the starting job. We don't run if the start job fails.
		DependsOn: []string{startJob.Id},

		// Use the same targeting as the start job. We assume the start job
		// had proper access to stop, too, so we just copy it.
		TargetRunner: startJob.TargetRunner,

		// Watch
		Operation: &pb.Job_WatchTask{
			WatchTask: &pb.Job_WatchTaskOp{
				StartJob: &pb.Ref_Job{Id: startJob.Id},
			},
		},
	}

	// Get the next id for the job
	id, err := server.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	return job, nil
}

// onDemandRunnerStopJob generates a StopJob template for a Task.
func (s *Service) onDemandRunnerStopJob(
	ctx context.Context,
	startJob *pb.Job,
	watchJob *pb.Job,
	source *pb.Job,
	od *pb.OnDemandRunnerConfig,
) (*pb.Job, error) {
	depends := []string{startJob.Id, watchJob.Id}

	// Only add the source job if we're not skipping it.
	if over := source.OndemandRunnerTask; over == nil || !over.SkipOperation {
		depends = append(depends, source.Id)
	}

	job := &pb.Job{
		// Inherit the workspace/application of the source job.
		Workspace:   source.Workspace,
		Application: source.Application,

		// We depend on both the start job and the main job. We allow them
		// both to fail, however, because we want to try to stop no matter what.
		DependsOn:             depends,
		DependsOnAllowFailure: depends,

		// Use the same targeting as the start job. We assume the start job
		// had proper access to stop, too, so we just copy it.
		TargetRunner: startJob.TargetRunner,

		// Stop
		Operation: &pb.Job_StopTask{
			StopTask: &pb.Job_StopTaskLaunchOp{
				Params: &pb.Job_TaskPluginParams{
					PluginType: od.PluginType,
					HclConfig:  od.PluginConfig,
					HclFormat:  od.ConfigFormat,
				},

				// Get our state from the start job.
				State: &pb.Job_StopTaskLaunchOp_StartJobId{
					StartJobId: startJob.Id,
				},
			},
		},
	}

	// Get the next id for the job
	id, err := server.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	return job, nil
}

func (s *Service) ValidateJob(
	ctx context.Context,
	req *pb.ValidateJobRequest,
) (*pb.ValidateJobResponse, error) {
	var err error
	result := &pb.ValidateJobResponse{Valid: true}

	// Struct validation
	if err := serverptypes.ValidateJob(req.Job); err != nil {
		result.Valid = false
		result.ValidationError = status.New(codes.FailedPrecondition, err.Error()).Proto()
		return result, nil
	}

	// Check assignability
	result.Assignable, err = s.state(ctx).JobIsAssignable(ctx, req.Job)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetJobStream(
	req *pb.GetJobStreamRequest,
	server pb.Waypoint_GetJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())
	ctx := server.Context()

	// Get the job
	ws := memdb.NewWatchSet()
	job, err := s.state(ctx).JobById(req.JobId, ws)
	if err != nil {
		return err
	}
	if job == nil {
		return status.Errorf(codes.NotFound, "job not found for ID: %s", req.JobId)
	}
	log = log.With("job_id", job.Id)

	// We always send the open message as confirmation the job was found.
	if err := server.Send(&pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Open_{
			Open: &pb.GetJobStreamResponse_Open{},
		},
	}); err != nil {
		return err
	}

	// Start a goroutine that watches for job changes
	jobCh := make(chan *serverstate.Job, 1)
	errCh := make(chan error, 1)
	go func() {
		for {
			// Send the job
			select {
			case jobCh <- job:
			case <-ctx.Done():
				return
			}

			// Wait for the job to update
			if err := ws.WatchCtx(ctx); err != nil {
				if ctx.Err() == nil {
					errCh <- err
				}

				return
			}

			// Updated job, requery it
			ws = memdb.NewWatchSet()
			job, err = s.state(ctx).JobById(job.Id, ws)
			if err != nil {
				log.Error("error acquiring job by id", "error", err, "id", req.JobId)
				errCh <- err
				return
			}
			if job == nil {
				errCh <- status.Errorf(codes.Internal, "job disappeared for ID: %s", req.JobId)
				return
			}
		}
	}()

	// Track that we only send these events once. We could use a bitmask for
	// this if we cared about that level of optimization but it hurts readability
	// and we don't need the performance yet.
	var cancelSent bool
	var downloadSent bool

	// Enter the event loop
	var lastState pb.Job_State
	var lastJob *pb.Job
	var eventsCh <-chan []*pb.GetJobStreamResponse_Terminal_Event
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errCh:
			return err

		case job := <-jobCh:
			log.Debug("job state change", "state", job.State)

			// If we have a state change, send that event down. We also send
			// down a state change if we enter a "cancelled" scenario.
			canceling := job.CancelTime != nil
			if lastState != job.State || cancelSent != canceling {
				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_State_{
						State: &pb.GetJobStreamResponse_State{
							Previous:  lastState,
							Current:   job.State,
							Job:       job.Job,
							Canceling: canceling,
						},
					},
				}); err != nil {
					return err
				}

				lastState = job.State
				cancelSent = canceling
			}

			// If we have a data source ref set, then we need to send the download event.
			if !downloadSent && job.DataSourceRef != nil {
				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Download_{
						Download: &pb.GetJobStreamResponse_Download{
							DataSourceRef: job.DataSourceRef,
						},
					},
				}); err != nil {
					return err
				}

				downloadSent = true
			}

			// If our job changed then we send down a job change notification.
			// We use reflect.DeepEqual here which isn't super exact but errors
			// on the side of false positives rather than false negatives so
			// at worst it'll send down a few more noisy job updates rather than
			// miss any. Because of this, we use it for simplicity.
			if lastJob == nil || !reflect.DeepEqual(lastJob, job.Job) {
				lastJob = job.Job

				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Job{
						Job: &pb.GetJobStreamResponse_JobChange{
							Job: job.Job,
						},
					},
				}); err != nil {
					return err
				}
			}

			if eventsCh == nil {
				switch job.State {
				case pb.Job_RUNNING:

					// We're seeing the job start up live, so we'll initialize the
					// event channel and use the job streamer to stream the logs
					// in as they arrive via the runner interface.
					// If the job OutputBuffer is nil, the streamer won't send
					// any events.
					eventsCh, err = s.getJobStreamOutputInit(ctx, log, job, server)
					if err != nil {
						msg := "failed to init job output stream"
						log.Error(msg, "job.id", job.Id, "error", err)
						return status.Error(codes.Internal, msg)
					}

				// NOTE: at present (2022-04-20) the CLI does not utilize this
				// code path but the UI does for historical job viewing.
				case pb.Job_SUCCESS, pb.Job_ERROR:
					// This means that the requested stream finished before GetJobStream was
					// called. As such, we'll replay the output events.
					events, err := s.logStreamProvider.ReadCompleted(ctx, log, s.state(ctx), job)
					if err != nil {
						msg := "failed to stream logs for completed job"
						log.Error(msg, "job.Id", job.Id, "error", err)
						return status.Error(codes.Internal, msg)
					}

					// We're doing this synchronously so that the client receives the events
					// before we send down the completion event.
					if err := server.Send(&pb.GetJobStreamResponse{
						Event: &pb.GetJobStreamResponse_Terminal_{
							Terminal: &pb.GetJobStreamResponse_Terminal{
								Events:   events,
								Buffered: true,
							},
						},
					}); err != nil {
						log.Error("failed to send logs for completed job", "job.Id", job.Id, "error", err)
						return err
					}
				}
			}

			switch job.State {

			case pb.Job_SUCCESS, pb.Job_ERROR:
				// TODO(mitchellh): we should drain the output buffer

				// Job is done. For success, error will be nil, so this
				// populates the event with the proper values.
				return server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Complete_{
						Complete: &pb.GetJobStreamResponse_Complete{
							Error:  job.Error,
							Result: job.Result,
						},
					},
				})
			}

		case events := <-eventsCh:
			if err := server.Send(&pb.GetJobStreamResponse{
				Event: &pb.GetJobStreamResponse_Terminal_{
					Terminal: &pb.GetJobStreamResponse_Terminal{
						Events: events,
					},
				},
			}); err != nil {
				return err
			}
		}
	}
}

func (s *Service) getJobStreamOutputInit(
	ctx context.Context,
	log hclog.Logger,
	job *serverstate.Job,
	server pb.Waypoint_GetJobStreamServer,
) (<-chan []*pb.GetJobStreamResponse_Terminal_Event, error) {
	// Start a log stream reader for this job
	lsReader, err := s.logStreamProvider.StartReader(ctx, log, job)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start log reader")
	}

	// Send down all our buffered lines.
	for {
		events, err := lsReader.ReadStream(ctx, false)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read buffered log batch")
		}
		if events == nil {
			break
		}

		if err := server.Send(&pb.GetJobStreamResponse{
			Event: &pb.GetJobStreamResponse_Terminal_{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events:   events,
					Buffered: true,
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	// Start a goroutine that reads output. If things go wrong in here,
	// we cancel the whole job stream request.
	ctx, cancel := context.WithCancel(ctx)
	eventsCh := make(chan []*pb.GetJobStreamResponse_Terminal_Event, 1)
	go func() {
		for {
			events, err := lsReader.ReadStream(ctx, true)
			if err != nil {
				// In the event of a reader error, we shut the stream down.
				// It's up to the reader to retry if ephemeral errors are common.

				msg := "failed to read streaming log batch"
				log.Error(msg, "error", err)

				// Let the client know we're terminating due to an error.
				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Terminal_{
						Terminal: &pb.GetJobStreamResponse_Terminal{
							Events: []*pb.GetJobStreamResponse_Terminal_Event{{
								Timestamp: timestamppb.Now(),

								// NOTE(izaak): really not sure if this is the right kind of terminal event
								Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
									Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
										Style: terminal.ErrorStyle,
										Msg:   msg,
									},
								},
							}},
							Buffered: false,
						},
					},
				}); err != nil {
					// This waypoint server must be experiencing some kind of catastrophic failure - it can't
					// stream logs or send messages to the client.
					log.Error("failed to inform client of read streaming error", "error", err)
				}

				cancel()
				return
			}
			if events == nil {
				return
			}

			select {
			case eventsCh <- events:
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventsCh, nil
}
