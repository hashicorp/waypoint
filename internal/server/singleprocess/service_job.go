package singleprocess

import (
	"context"
	"math/rand"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/hashicorp/waypoint/internal/serverstate"
)

// TODO: test
func (s *service) GetJob(
	ctx context.Context,
	req *pb.GetJobRequest,
) (*pb.Job, error) {
	job, err := s.state.JobById(req.JobId, nil)
	if err != nil {
		return nil, err
	}
	if job == nil || job.Job == nil {
		return nil, status.Errorf(codes.NotFound, "job not found")
	}

	return job.Job, nil
}

// TODO: test
func (s *service) XListJobs(
	ctx context.Context,
	req *pb.ListJobsRequest,
) (*pb.ListJobsResponse, error) {
	jobs, err := s.state.JobList()
	if err != nil {
		return nil, err
	}

	return &pb.ListJobsResponse{
		Jobs: jobs,
	}, nil
}

func (s *service) CancelJob(
	ctx context.Context,
	req *pb.CancelJobRequest,
) (*empty.Empty, error) {
	if err := s.state.JobCancel(req.JobId, false); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// queueJobMulti queues multiple jobs transactionally. The return value
// are the responses for each in the same order as the requests.
//
// Postcondition: len(result) == len(req) iff err == nil
func (s *service) queueJobMulti(
	ctx context.Context,
	req []*pb.QueueJobRequest,
) ([]*pb.QueueJobResponse, error) {
	jobs := make([]*pb.Job, 0, len(req))
	for _, single := range req {
		job, _, err := s.queueJobReqToJob(ctx, single)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job...)
	}

	// Queue the jobs
	if err := s.state.JobCreate(jobs...); err != nil {
		return nil, err
	}

	// Get the response
	resp := make([]*pb.QueueJobResponse, len(jobs))
	for i, job := range jobs {
		resp[i] = &pb.QueueJobResponse{JobId: job.Id}
	}

	return resp, nil
}

// queueJobReqToJob converts a QueueJobRequest to a job to queue, but
// does not queue it. This may return multiple jobs if the queue job
// request requires an on-demand runner. They should all be queued
// atomically with JobCreate.
func (s *service) queueJobReqToJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) ([]*pb.Job, string, error) {
	log := hclog.FromContext(ctx)
	job := req.Job

	// Validation
	if job == nil {
		return nil, "", status.Errorf(codes.FailedPrecondition, "job must be set")
	}
	if job.Operation == nil {
		// We special case this check and return "Unimplemented" because
		// the primary case where operation is nil is if a client is sending
		// us an unsupported operation.
		return nil, "", status.Errorf(codes.Unimplemented, "operation is nil or unknown")
	}
	if err := serverptypes.ValidateJob(job); err != nil {
		return nil, "", status.Errorf(codes.FailedPrecondition, err.Error())
	}

	// Verify the project exists and use that to set the default data source
	log.Debug("checking job project", "project", job.Application.Project)
	project, err := s.state.ProjectGet(&pb.Ref_Project{Project: job.Application.Project})
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
	id, err := server.Id()
	if err != nil {
		return nil, "", status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	// Validate expiry if we have one
	job.ExpireTime = nil
	if req.ExpiresIn != "" {
		dur, err := time.ParseDuration(req.ExpiresIn)
		if err != nil {
			return nil, "", status.Errorf(codes.FailedPrecondition,
				"Invalid expiry duration: %s", err.Error())
		}
		job.ExpireTime, err = ptypes.TimestampProto(time.Now().Add(dur))
		if err != nil {
			return nil, "", status.Errorf(codes.Aborted, "error configuring expiration: %s", err)
		}
	}

	// If the job can be run by any runner, then we attempt to see if we should spawn
	// an on-demand runner for it. We only consider jobs for any runner because ones
	// that are targeted can not target on-demand runners, because they don't yet exist.
	// If the job has any target runner, it is a remote job.
	if _, anyTarget := job.TargetRunner.Target.(*pb.Ref_Runner_Any); anyTarget {
		if project.OndemandRunner == nil {
			ods, err := s.state.OnDemandRunnerConfigDefault()
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
				log.Debug("multiple default on-demand runners detected, chose a random one",
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
	}

	return result, job.Id, nil
}

func (s *service) QueueJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) (*pb.QueueJobResponse, error) {
	jobs, jobId, err := s.queueJobReqToJob(ctx, req)
	if err != nil {
		return nil, err
	}

	// Queue the job
	if err := s.state.JobCreate(jobs...); err != nil {
		return nil, err
	}

	return &pb.QueueJobResponse{JobId: jobId}, nil
}

// wrapJobWithRunner takes a job and "wraps" it within an on-demand launched
// runner. This creates a dependency chain that ensures that the runner is
// started and stopped around the given job (hence "wraps").
func (s *service) wrapJobWithRunner(
	ctx context.Context,
	source *pb.Job,
) ([]*pb.Job, error) {
	// Get the runner profile we're going to use for this runner.
	od, err := s.state.OnDemandRunnerConfigGet(source.OndemandRunner)
	if err != nil {
		return nil, err
	}

	// Generate our job to start the ODR
	startJob, runnerId, err := s.onDemandRunnerStartJob(ctx, source, od)
	if err != nil {
		return nil, err
	}

	// Change our source job to require being run on the launched ODR.
	source.TargetRunner = &pb.Ref_Runner{
		Target: &pb.Ref_Runner_Id{
			Id: &pb.Ref_RunnerId{
				Id: runnerId,
			},
		},
	}

	// Our source job depends on the starting job.
	source.DependsOn = []string{startJob.Id}

	// Job to stop the ODR
	stopJob, err := s.onDemandRunnerStopJob(ctx, startJob, source, od)
	if err != nil {
		return nil, err
	}

	// These must be in order of dependency currently. This is a limitation
	// of the state.JobCreate API and we should fix it one day. If we get
	// this wrong it'll just error, so we'll know quickly.
	return []*pb.Job{
		startJob,
		source,
		stopJob,
	}, nil
}

func (s *service) onDemandRunnerStartJob(
	ctx context.Context,
	source *pb.Job,
	od *pb.OnDemandRunnerConfig,
) (*pb.Job, string, error) {
	log := hclog.FromContext(ctx)

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
	cfg, err := s.state.ServerConfigGet()
	if err != nil {
		return nil, "", err
	}

	// This should only happen during tests
	if len(cfg.AdvertiseAddrs) == 0 {
		log.Info("server has no advertise addrs, using localhost")
		addr = &pb.ServerConfig_AdvertiseAddr{
			Addr: "localhost:9701",
		}
	} else {
		addr = cfg.AdvertiseAddrs[0]
	}

	// We generate a new login token for each ondemand-runner used. This will inherit
	// the user of the token to be the user that queue'd the original job, which is
	// the correct behavior.
	token, err := s.newToken(60*time.Minute, DefaultKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{Login: &pb.Token_Login{
			UserId: DefaultUserId,
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

	// Arguments for the runner image. Waypoint is ALWAYS assumed to be
	// the entrypoint for ODR images.
	args := []string{"runner", "agent", "-vv", "-id", runnerId, "-odr", "-odr-profile-id", od.Id}

	job := &pb.Job{
		// Inherit the workspace/application of the source job.
		Workspace:   source.Workspace,
		Application: source.Application,

		Operation: &pb.Job_StartTask{
			StartTask: &pb.Job_StartTaskLaunchOp{
				Params: &pb.Job_TaskPluginParams{
					PluginType: od.PluginType,
					HclConfig:  od.PluginConfig,
					HclFormat:  od.ConfigFormat,
				},
				Info: &pb.TaskLaunchInfo{
					OciUrl:               od.OciUrl,
					EnvironmentVariables: envVars,
					Arguments:            args,
				},
			},
		},
	}

	// Get the next id for the job
	id, err := server.Id()
	if err != nil {
		return nil, "", status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	// We're going to wait up to 60s for the job be picked up. No reason it won't be
	// picked up immediately.
	dur, err := time.ParseDuration("60s")
	if err != nil {
		return nil, "", status.Errorf(codes.FailedPrecondition,
			"Invalid expiry duration: %s", err.Error())
	}

	job.ExpireTime, err = ptypes.TimestampProto(time.Now().Add(dur))
	if err != nil {
		return nil, "", status.Errorf(codes.Aborted, "error configuring expiration: %s", err)
	}

	if err != nil {
		return nil, "", status.Errorf(codes.FailedPrecondition,
			"Failed to get on-demand runner config by name %q, id %q: %s",
			job.OndemandRunner.Name, job.OndemandRunner.Id, err)
	}
	// This will be either "Any" or a specific static runner.
	job.TargetRunner = od.TargetRunner

	return job, runnerId, nil
}

func (s *service) onDemandRunnerStopJob(
	ctx context.Context,
	startJob *pb.Job,
	source *pb.Job,
	od *pb.OnDemandRunnerConfig,
) (*pb.Job, error) {
	job := &pb.Job{
		// Inherit the workspace/application of the source job.
		Workspace:   source.Workspace,
		Application: source.Application,

		// We depend on both the start job and the main job. We allow them
		// both to fail, however, because we want to try to stop no matter what.
		DependsOn:             []string{startJob.Id, source.Id},
		DependsOnAllowFailure: []string{startJob.Id, source.Id},

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

func (s *service) ValidateJob(
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
	result.Assignable, err = s.state.JobIsAssignable(ctx, req.Job)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetJobStream(
	req *pb.GetJobStreamRequest,
	server pb.Waypoint_GetJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())
	ctx := server.Context()

	// Get the job
	ws := memdb.NewWatchSet()
	job, err := s.state.JobById(req.JobId, ws)
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
			job, err = s.state.JobById(job.Id, ws)
			if err != nil {
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

			// If we haven't initialized output streaming and the output buffer
			// is now non-nil, initialize that. This will send any buffered
			// data down.
			if eventsCh == nil && job.OutputBuffer != nil {
				eventsCh, err = s.getJobStreamOutputInit(ctx, job, server)
				if err != nil {
					return err
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

func (s *service) readJobLogBatch(r *logbuffer.Reader, block bool) []*pb.GetJobStreamResponse_Terminal_Event {
	entries := r.Read(64, block)
	if entries == nil {
		return nil
	}

	events := make([]*pb.GetJobStreamResponse_Terminal_Event, len(entries))
	for i, entry := range entries {
		events[i] = entry.(*pb.GetJobStreamResponse_Terminal_Event)
	}

	return events
}

func (s *service) getJobStreamOutputInit(
	ctx context.Context,
	job *serverstate.Job,
	server pb.Waypoint_GetJobStreamServer,
) (<-chan []*pb.GetJobStreamResponse_Terminal_Event, error) {
	// Send down all our buffered lines.
	outputR := job.OutputBuffer.Reader(-1)
	go outputR.CloseContext(ctx)
	for {
		events := s.readJobLogBatch(outputR, false)
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

	// Start a goroutine that reads output
	eventsCh := make(chan []*pb.GetJobStreamResponse_Terminal_Event, 1)
	go func() {
		for {
			events := s.readJobLogBatch(outputR, true)
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
