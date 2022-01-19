package singleprocess

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	"github.com/hashicorp/waypoint/internal/serverstate"
)

func (s *service) ListRunners(
	ctx context.Context,
	req *pb.ListRunnersRequest,
) (*pb.ListRunnersResponse, error) {
	runners, err := s.state.RunnerList()
	if err != nil {
		return nil, err
	}
	return &pb.ListRunnersResponse{Runners: runners}, nil
}

// TODO: test
func (s *service) GetRunner(
	ctx context.Context,
	req *pb.GetRunnerRequest,
) (*pb.Runner, error) {
	return s.state.RunnerById(req.RunnerId, nil)
}

func (s *service) RunnerGetDeploymentConfig(
	ctx context.Context,
	req *pb.RunnerGetDeploymentConfigRequest,
) (*pb.RunnerGetDeploymentConfigResponse, error) {
	// Get our server config
	cfg, err := s.state.ServerConfigGet()
	if err != nil {
		return nil, err
	}

	// If we have no config set yet, this is an error.
	if cfg == nil {
		return nil, status.Errorf(codes.Aborted,
			"server configuration for deployment information not yet set.")
	}

	// If we have no advertise addresses, then we just send back empty values.
	// This disables any entrypoint settings.
	if len(cfg.AdvertiseAddrs) == 0 {
		return &pb.RunnerGetDeploymentConfigResponse{}, nil
	}

	// Our addr for now is just the first one since we don't support
	// multiple addresses yet. In the future we will want to support more
	// advanced choicing.
	addr := cfg.AdvertiseAddrs[0]

	return &pb.RunnerGetDeploymentConfigResponse{
		ServerAddr:          addr.Addr,
		ServerTls:           addr.Tls,
		ServerTlsSkipVerify: addr.TlsSkipVerify,
	}, nil
}

func (s *service) RunnerToken(
	ctx context.Context,
	req *pb.RunnerTokenRequest,
) (*pb.RunnerTokenResponse, error) {
	log := hclog.FromContext(ctx)
	record := req.Runner

	// Create our record
	log = log.With("runner_id", record.Id)
	log.Trace("registering runner")
	if err := s.state.RunnerCreate(record); err != nil {
		return nil, err
	}

	// When we exit, mark the runner as offline. This will delete the record
	// if we're never adopted.
	defer func() {
		log.Trace("marking runner as offline")
		if err := s.state.RunnerOffline(record.Id); err != nil {
			log.Error("failed to mark runner as offline. This should not happen.", "err", err)
		}
	}()

	for {
		// Get the runner
		ws := memdb.NewWatchSet()
		r, err := s.state.RunnerById(record.Id, ws)
		if err != nil {
			return nil, err
		}

		switch r.AdoptionState {
		case pb.Runner_REJECTED:
			// Runner is explicitly rejected. Return and error.
			return nil, status.Errorf(codes.PermissionDenied,
				"runner adoption is explicitly rejected")

		case pb.Runner_ADOPTED:
			// Runner explicitly adopted, create token and return!
			tok, err := s.newToken(
				// Doesn't expire because we can expire it by unadopting.
				// NOTE(mitchellh): At some point, we should make these
				// expire and introduce rotation as a feature of adoption.
				0,

				DefaultKeyId,
				nil,
				&pb.Token{
					Kind: &pb.Token_Runner_{
						Runner: &pb.Token_Runner{
							Id: record.Id,
						},
					},
				},
			)
			if err != nil {
				return nil, err
			}

			return &pb.RunnerTokenResponse{Token: tok}, nil
		}

		// Wait for changes
		if err := ws.WatchCtx(ctx); err != nil {
			return nil, err
		}
	}
}

func (s *service) RunnerConfig(
	srv pb.Waypoint_RunnerConfigServer,
) error {
	log := hclog.FromContext(srv.Context())
	ctx, cancel := context.WithCancel(srv.Context())
	defer cancel()

	// Get the request
	event, err := srv.Recv()
	if err != nil {
		return err
	}
	req, ok := event.Event.(*pb.RunnerConfigRequest_Open_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"expected open event, got %T", event)
	}
	record := req.Open.Runner

	// Create our record
	log = log.With("runner_id", record.Id)
	log.Trace("registering runner")
	if err := s.state.RunnerCreate(record); err != nil {
		return err
	}

	// Mark the runner as offline if they disconnect from the config stream loop.
	defer func() {
		log.Trace("marking runner as offline")
		if err := s.state.RunnerOffline(record.Id); err != nil {
			log.Error("failed to mark runner as offline. This should not happen.", "err", err)
		}
	}()

	// Get our token and reverify that we are adopted.
	tok := s.tokenFromContext(ctx)
	if tok == nil {
		log.Error("no token, should not be possible")
		return status.Errorf(codes.Unauthenticated, "no token")
	}
	switch tok.Kind.(type) {
	case *pb.Token_Login_:
		// Legacy (pre WP 0.8) token. We accept these as preadopted.
		// NOTE(mitchellh): One day, we should reject these because modern
		// preadoption should be via runner tokens.
		if err := s.state.RunnerAdopt(record.Id, true); err != nil {
			return err
		}

	case *pb.Token_Runner_:
		// A runner token. We validate here that we're not explicitly rejected.
		// We have to check again here because runner tokens can be created
		// for ANY runner, but we can reject a SPECIFIC runner.
		r, err := s.state.RunnerById(record.Id, nil)
		if err != nil {
			return err
		}
		if r.AdoptionState == pb.Runner_REJECTED {
			return status.Errorf(codes.Unauthenticated,
				"runner is explicitly rejected (unadopted)")
		}

		// Mark it as preadopted if it isn't adopted
		if r.AdoptionState != pb.Runner_ADOPTED {
			if err := s.state.RunnerAdopt(record.Id, true); err != nil {
				return err
			}
		}

	default:
		return status.Errorf(codes.Unauthenticated, "not a valid runner token")
	}

	// Start a goroutine that listens on the recvmsg so we can detect
	// when the client exited.
	go func() {
		defer cancel()

		for {
			_, err := srv.Recv()
			if err != nil {
				if err != io.EOF {
					log.Warn("unknown error from recvmsg", "err", err)
				}

				return
			}
		}
	}()

	// If this is an ODR runner, then we query the job it is waiting for
	// in order to build up other information about this runner such as the
	// project/app scope, workspace, etc.
	//
	// It is REQUIRED that an ODR has its target job queued BEFORE the
	// ODR is launched. If we can't find a job, we error and exit which
	// will also exit the runner.
	var job *pb.Job

	if _, ok := record.Kind.(*pb.Runner_Odr); ok {
		// Get a job assignment for this runner, non-blocking
		sjob, err := s.state.JobPeekForRunner(ctx, record)
		if err != nil {
			return err
		}
		if sjob == nil {
			return status.Errorf(codes.FailedPrecondition,
				"no pending job for this on-demand runner. A pending job "+
					"must be registered prior to registering the runner.")
		}

		// Set our job
		job = sjob.Job

		log.Debug("runner is scoped for config",
			"application", job.Application,
			"workspace", job.Workspace,
			"labels", job.Labels)
	}

	// Build our config in a loop.
	for {
		ws := memdb.NewWatchSet()

		// Build our config
		config := &pb.RunnerConfig{}

		// Build our config var request. This is always runner-scoped, but
		// if we're ODR then job should be non-nil and we set the proper
		// project/app, workspace, labels, etc.
		configReq := &pb.ConfigGetRequest{
			Runner: &pb.Ref_RunnerId{
				Id: record.Id,
			},
		}
		if job != nil {
			configReq.Scope = &pb.ConfigGetRequest_Application{
				Application: job.Application,
			}
			configReq.Workspace = job.Workspace
			configReq.Labels = job.Labels
		}

		vars, err := s.state.ConfigGetWatch(configReq, ws)
		if err != nil {
			return err
		}
		config.ConfigVars = vars

		// Get the config sources we need for our vars. We only do this if
		// at least one var has a dynamic value.
		//
		// We also do this if the runner is NOT local, in case it is processing
		// jobs that are using config variables with dynamic defaults.
		_, isLocal := record.Kind.(*pb.Runner_Local_)
		if varContainsDynamic(vars) || !isLocal {
			// Important: we've discussed optimizing this to send down only the
			// config sourcers that are needed by vars. We cannot do that because
			// waypoint.hcl config can now source dynamic config too and we can't
			// know those in advance perfectly. Always send down all config sources.
			sources, err := s.state.ConfigSourceGetWatch(&pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},
			}, ws)
			if err != nil {
				return err
			}

			config.ConfigSources = sources
		}

		// Send new config
		if err := srv.Send(&pb.RunnerConfigResponse{
			Config: config,
		}); err != nil {
			return err
		}

		// Nil out the stuff we used so that if we're waiting awhile we can GC
		config = nil

		// Wait for any changes
		if err := ws.WatchCtx(ctx); err != nil {
			return err
		}
	}
}

func (s *service) RunnerJobStream(
	server pb.Waypoint_RunnerJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	// Receive our opening message so we can determine the runner ID.
	req, err := server.Recv()
	if err != nil {
		return err
	}
	reqEvent, ok := req.Event.(*pb.RunnerJobStreamRequest_Request_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"first message must be a Request event")
	}
	log = log.With("runner_id", reqEvent.Request.RunnerId)

	// Get the runner to validate it is registered
	runner, err := s.state.RunnerById(reqEvent.Request.RunnerId, nil)
	if err != nil {
		log.Error("unknown runner connected", "id", reqEvent.Request.RunnerId)
		return err
	}

	// The runner must be adopted to get a job.
	if runner.AdoptionState != pb.Runner_ADOPTED &&
		runner.AdoptionState != pb.Runner_PREADOPTED {
		return status.Errorf(codes.FailedPrecondition,
			"runner must be adopted prior to requesting jobs")
	}

	// Get a job assignment for this runner
	job, err := s.state.JobAssignForRunner(ctx, runner)
	if err != nil {
		return err
	}

	// Send the job assignment.
	//
	// If this has an error, we continue to accumulate the error until
	// we set the ack status in the DB. We do this because if we fail to
	// send the job assignment we want to nack the job so it is queued again.
	err = server.Send(&pb.RunnerJobStreamResponse{
		Event: &pb.RunnerJobStreamResponse_Assignment{
			Assignment: &pb.RunnerJobStreamResponse_JobAssignment{
				Job: job.Job,
			},
		},
	})

	// Wait for an ack. We only do this if the job assignment above
	// succeeded. If it didn't succeed, the client will never send us
	// an ack.
	ack := false
	if err == nil { // if sending the job assignment was a success
		req, err = server.Recv()

		// If we received a message we inspect it. If we failed to
		// receive a message, we've set the `err` value and we keep
		// ack to false so that we nack the job later.
		if err == nil {
			switch req.Event.(type) {
			case *pb.RunnerJobStreamRequest_Ack_:
				ack = true

			case *pb.RunnerJobStreamRequest_Error_:
				ack = false

			default:
				ack = false
				err = status.Errorf(codes.FailedPrecondition,
					"ack expected, got: %T", req.Event)
			}
		} else {
			ack = false
		}
	}

	// Send the ack OR nack, based on the value of +ack+.
	job, ackerr := s.state.JobAck(job.Id, ack)
	if ackerr != nil {
		// If this fails, we just log, there is nothing more we can do.
		log.Warn("job ack failed", "outer_error", err, "error", ackerr)

		// If we had no outer error, set the ackerr so that we exit. If
		// we do have an outer error, then the ack error only shows up in
		// the log.
		if err == nil {
			err = ackerr
		}
	}

	// If we have an error, return that. We also return if we didn't ack for
	// any reason. This error can be set at any point since job assignment.
	if err != nil || !ack {
		return err
	}

	// Start a goroutine that watches for job changes
	jobCh := make(chan *serverstate.Job, 1)
	errCh := make(chan error, 1)
	go func() {
		for {
			ws := memdb.NewWatchSet()
			job, err = s.state.JobById(job.Id, ws)
			if err != nil {
				errCh <- err
				return
			}
			if job == nil {
				errCh <- status.Errorf(codes.Internal, "job disappeared")
				return
			}

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
		}
	}()

	// Create a goroutine that just waits for events. We have to do this
	// so we can exit properly on client side close.
	eventCh := make(chan *pb.RunnerJobStreamRequest, 1)
	go func() {
		defer cancel()

		for {
			log.Trace("waiting for job stream event")
			req, err := server.Recv()
			if err == io.EOF {
				// On EOF, this means the client closed their write side.
				// In this case, we assume we have exited and exit accordingly.
				return
			}

			if err != nil {
				// For any other error, we send the error along and exit the
				// read loop. The sent error will be picked up and sent back
				// as a result to the client.
				errCh <- err
				return
			}
			log.Trace("event received", "event", req.Event)

			// Send the event down
			select {
			case eventCh <- req:
			case <-ctx.Done():
				return
			}

			// If this is a terminating event, we exit this loop
			switch event := req.Event.(type) {
			case *pb.RunnerJobStreamRequest_Complete_:
				log.Debug("job stream recv loop exiting due to completion")
				return
			case *pb.RunnerJobStreamRequest_Error_:
				log.Debug("job stream recv loop exiting due to error",
					"error", event.Error.Error.Message)
				return
			}
		}
	}()

	// Recv events in a loop
	var lastJob *pb.Job
	for {
		select {
		case <-ctx.Done():
			// We need to drain the event channel
			for {
				select {
				case req := <-eventCh:
					if err := s.handleJobStreamRequest(log, job, server, req); err != nil {
						return err
					}
				default:
					return nil
				}
			}

		case err := <-errCh:
			return err

		case req := <-eventCh:
			if err := s.handleJobStreamRequest(log, job, server, req); err != nil {
				return err
			}

		case job := <-jobCh:
			if lastJob == job.Job {
				continue
			}

			// If the job is canceled, send that event. We send this each time
			// the cancel time changes. The cancel time only changes if multiple
			// cancel requests are made.
			if job.CancelTime != nil &&
				(lastJob == nil || !lastJob.CancelTime.AsTime().Equal(job.CancelTime.AsTime())) {
				// The job is forced if we're in an error state. This must be true
				// because we would've already exited the loop if we naturally
				// got a terminal event.
				force := job.State == pb.Job_ERROR

				err := server.Send(&pb.RunnerJobStreamResponse{
					Event: &pb.RunnerJobStreamResponse_Cancel{
						Cancel: &pb.RunnerJobStreamResponse_JobCancel{
							Force: force,
						},
					},
				})
				if err != nil {
					return err
				}

				// On force we exit immediately.
				if force {
					return nil
				}
			}

			lastJob = job.Job
		}
	}
}

func (s *service) handleJobStreamRequest(
	log hclog.Logger,
	job *serverstate.Job,
	srv pb.Waypoint_RunnerJobStreamServer,
	req *pb.RunnerJobStreamRequest,
) error {
	log.Trace("event received", "event", req.Event)
	switch event := req.Event.(type) {
	case *pb.RunnerJobStreamRequest_Complete_:
		return s.state.JobComplete(job.Id, event.Complete.Result, nil)

	case *pb.RunnerJobStreamRequest_Error_:
		return s.state.JobComplete(job.Id, nil, status.FromProto(event.Error.Error).Err())

	case *pb.RunnerJobStreamRequest_Heartbeat_:
		return s.state.JobHeartbeat(job.Id)

	case *pb.RunnerJobStreamRequest_Download:
		if err := s.state.JobUpdateRef(job.Id, event.Download.DataSourceRef); err != nil {
			return err
		}

		return s.state.ProjectUpdateDataRef(&pb.Ref_Project{
			Project: job.Application.Project,
		}, job.Workspace, event.Download.DataSourceRef)

	case *pb.RunnerJobStreamRequest_ConfigLoad_:
		return s.state.JobUpdate(job.Id, func(jobpb *pb.Job) error {
			jobpb.Config = event.ConfigLoad.Config
			return nil
		})

	case *pb.RunnerJobStreamRequest_Terminal:
		// This shouldn't happen but we want to protect against it to prevent
		// a panic.
		if job.OutputBuffer == nil {
			log.Warn("got terminal event but internal output buffer is nil, dropping lines")
			return nil
		}

		// Write the entries to the output buffer
		entries := make([]logbuffer.Entry, len(event.Terminal.Events))
		for i, ev := range event.Terminal.Events {
			entries[i] = ev
		}

		// Write the events
		job.OutputBuffer.Write(entries...)

		return nil

	default:
		log.Warn("unexpected event received", "event", hclog.Fmt("%T", req.Event))
	}

	return nil
}
