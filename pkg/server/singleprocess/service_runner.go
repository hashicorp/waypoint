package singleprocess

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/waypoint/pkg/server/hcerr"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint/internal/telemetry/metrics"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logstream"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) ListRunners(
	ctx context.Context,
	req *pb.ListRunnersRequest,
) (*pb.ListRunnersResponse, error) {
	runners, err := s.state(ctx).RunnerList(ctx)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to list runners")
	}
	return &pb.ListRunnersResponse{Runners: runners}, nil
}

// TODO: test
func (s *Service) GetRunner(
	ctx context.Context,
	req *pb.GetRunnerRequest,
) (*pb.Runner, error) {
	result, err := s.state(ctx).RunnerById(ctx, req.RunnerId, nil)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get runner", "id", req.RunnerId)
	}
	return result, err
}

func (s *Service) RunnerGetDeploymentConfig(
	ctx context.Context,
	req *pb.RunnerGetDeploymentConfigRequest,
) (*pb.RunnerGetDeploymentConfigResponse, error) {
	// Get our server config
	serverConfig, err := s.GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get server config to populate runner start job server addr")
	}

	cfg := serverConfig.Config

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

func (s *Service) AdoptRunner(
	ctx context.Context,
	req *pb.AdoptRunnerRequest,
) (*empty.Empty, error) {
	var err error
	if err = serverptypes.ValidateAdoptRunnerRequest(req); err != nil {
		return nil, err
	}
	log := hclog.FromContext(ctx)
	if req.Adopt {
		if err = s.state(ctx).RunnerAdopt(ctx, req.RunnerId, false); err != nil {
			return &empty.Empty{}, hcerr.Externalize(log, err, "failed to adopt runner", "id", req.RunnerId)
		}
	} else {
		if err = s.state(ctx).RunnerReject(ctx, req.RunnerId); err != nil {
			return &empty.Empty{}, hcerr.Externalize(log, err, "failed to reject runner", "id", req.RunnerId)
		}
	}

	return &empty.Empty{}, nil
}

func (s *Service) ForgetRunner(
	ctx context.Context,
	req *pb.ForgetRunnerRequest,
) (*empty.Empty, error) {
	var err error
	if err = serverptypes.ValidateForgetRunnerRequest(req); err != nil {
		return nil, err
	}
	if err = s.state(ctx).RunnerDelete(ctx, req.RunnerId); err != nil {
		return &empty.Empty{}, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to delete runner", "id", req.RunnerId)
	}
	return &empty.Empty{}, nil
}

func (s *Service) RunnerToken(
	ctx context.Context,
	req *pb.RunnerTokenRequest,
) (*pb.RunnerTokenResponse, error) {
	log := hclog.FromContext(ctx)
	record := req.Runner

	// Get our token because our behavior changes a bit with different tokens.
	// Token may be nil because this is an unauthenticated endpoint.
	if tok := s.decodedTokenFromContext(ctx); tok != nil {
		switch k := tok.Kind.(type) {
		case *pb.Token_Login_:
			// Legacy (pre WP 0.8) token. We accept these as preadopted. We just
			// return an empty token here meaning to not change.
			// NOTE(mitchellh): One day, we should reject these because modern
			// preadoption should be via runner tokens.
			log.Debug("valid login token provided, adoption will be skipped")
			return &pb.RunnerTokenResponse{}, nil

		case *pb.Token_Runner_:
			if k.Runner.Id != "" {
				runnerId, err := s.decodeId(k.Runner.Id)
				if err != nil {
					log.Error("Failed to parse hcp id", "id", k.Runner.Id, "err", err)
					return nil, status.Errorf(codes.InvalidArgument, "invalid runner id format")
				}

				// If the runner token has an ID set and it doesn't match this one,
				// then the token is invalid and we should kick off the adoption process.
				if runnerId != record.Id {
					break
				}
			}

			// If the token has a label hash, then we need to validate it.
			// If the label hash does not match what we know about the runner,
			// we need to trigger adoption.
			if expected := k.Runner.LabelHash; expected > 0 {
				actual, err := serverptypes.RunnerLabelHash(record.Labels)
				if err != nil {
					return nil, err
				}

				if expected != actual {
					log.Info("runner token has invalid label hash, restarting adoption")
					break
				}
			}

			// Seemingly valid runner token. If our logic is wrong its okay
			// because RunnerConfig will reject them.
			log.Debug("valid runner token provided, adoption will be skipped")
			return &pb.RunnerTokenResponse{}, nil
		}

		// Any other token type we just continue with the adoption process.
	}

	// We require a cookie. We only need to check emptiness cause if its
	// set it will be validated in auth.go. We do NOT require the cookie if
	// we receive a valid token so its important to have this check after the
	// above token check.
	if CookieFromRequest(ctx) == "" {
		return nil, status.Errorf(codes.PermissionDenied,
			"RunnerToken requires the 'cookie' metadata value to be set")
	}

	// Create our record
	log = log.With("runner_id", record.Id)
	log.Trace("registering runner")
	if err := s.state(ctx).RunnerCreate(ctx, record); err != nil {
		return nil, hcerr.Externalize(log, err, "failed to create runner", "id", record.Id)
	}

	// When we exit, mark the runner as offline. This will delete the record
	// if we're never adopted.
	defer func() {
		log.Trace("marking runner as offline")
		if err := s.state(ctx).RunnerOffline(ctx, record.Id); err != nil {
			log.Error("failed to mark runner as offline. This should not happen.", "err", err)
		}
	}()

	// Get the runner
	r, err := s.state(ctx).RunnerById(ctx, record.Id, nil)
	if status.Code(err) == codes.NotFound {
		err = nil
		r = nil
	}
	if err != nil {
		return nil, hcerr.Externalize(log, err, "unknown runner connected", "id", record.Id)
	}
	prevAdopted := r != nil && r.AdoptionState == pb.Runner_ADOPTED

	// If we reached this point and we're previously adopted, then it is an
	// error. If we're previously adopted, we expect that runners will have
	// the token from that adoption. If we allowed this through, then any
	// guest with the runner ID could get a token -- a big security issue.
	if prevAdopted {
		return nil, status.Errorf(codes.PermissionDenied,
			"runner is already adopted, use the previously issued runner token")
	}

	log.Debug("token provided is not a runner token, waiting for adoption")
	for {
		// Get the runner
		ws := memdb.NewWatchSet()
		r, err := s.state(ctx).RunnerById(ctx, record.Id, ws)
		if err != nil {
			return nil, hcerr.Externalize(log, err, "failed to get runner while waiting for adoption state to change", "id", record.Id)
		}

		switch r.AdoptionState {
		case pb.Runner_REJECTED:
			// Runner is explicitly rejected. Return and error.
			return nil, status.Errorf(codes.PermissionDenied,
				"runner adoption is explicitly rejected")

		case pb.Runner_ADOPTED:
			// Runner explicitly adopted, create token and return!

			hash, err := serverptypes.RunnerLabelHash(record.Labels)
			if err != nil {
				return nil, hcerr.Externalize(log, err, "failed to calculate runner hash", "id", record.Id)
			}

			encodedId, err := s.encodeId(ctx, record.Id)
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					fmt.Errorf("failed to encode runner id %s", record.Id),
					"failed to encode runner id",
					"id",
					record.Id,
				)
			}

			tok, err := s.newToken(ctx,
				// Doesn't expire because we can expire it by unadopting.
				// NOTE(mitchellh): At some point, we should make these
				// expire and introduce rotation as a feature of adoption.
				0,

				s.activeAuthKeyId,
				nil,
				&pb.Token{
					Kind: &pb.Token_Runner_{
						Runner: &pb.Token_Runner{
							Id:        encodedId,
							LabelHash: hash,
						},
					},
				},
			)
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to generate new runner token",
				)
			}

			return &pb.RunnerTokenResponse{Token: tok}, nil
		}

		// Wait for changes
		log.Trace("runner is not adopted, waiting for state change")
		if err := ws.WatchCtx(ctx); err != nil {
			return nil, hcerr.Externalize(
				log,
				err,
				"error waiting for runner state changes",
			)
		}
	}
}

func (s *Service) RunnerConfig(
	srv pb.Waypoint_RunnerConfigServer,
) error {
	log := hclog.FromContext(srv.Context())
	ctx, cancel := context.WithCancel(srv.Context())
	defer cancel()

	// Get the request
	event, err := srv.Recv()
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"error receiving runner config request",
		)
	}
	req, ok := event.Event.(*pb.RunnerConfigRequest_Open_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"expected open event, got %T", event)
	}
	record := req.Open.Runner

	// Get our token and reverify that we are adopted.
	if err := s.runnerVerifyToken(log, ctx, record.Id, record.Labels); err != nil {
		return hcerr.Externalize(
			log,
			err,
			"error verifying runner token",
		)
	}

	// Create our record
	log = log.With("runner_id", record.Id)
	log.Trace("registering runner")

	if err := s.state(ctx).RunnerCreate(ctx, record); err != nil {
		return hcerr.Externalize(log, err, "failed to create runner", "id", record.Id)
	}

	state := s.state(ctx)
	// Mark the runner as offline if they disconnect from the config stream loop.
	defer func() {
		log.Trace("marking runner as offline")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := state.RunnerOffline(ctx, record.Id); err != nil {
			log.Error("failed to mark runner as offline. This should not happen.", "err", err)
		}
	}()

	// If the runner we just registered is explicitly rejected then we
	// do not allow it to continue, even with a preadoption token.
	r, err := s.state(ctx).RunnerById(ctx, record.Id, nil)
	if err != nil {
		return hcerr.Externalize(log, err, "failed to get newly-registered runner", "id", record.Id)
	}
	if r.AdoptionState == pb.Runner_REJECTED {
		return status.Errorf(codes.PermissionDenied,
			"runner is explicitly rejected (unadopted)")
	}
	if r.AdoptionState != pb.Runner_ADOPTED {
		if err := s.state(ctx).RunnerAdopt(ctx, record.Id, true); err != nil {
			return hcerr.Externalize(log, err, "failed to adopt runner", "id", record.Id)
		}
	}

	// Start a goroutine that listens on the recvmsg so we can detect
	// when the client exited.
	go func() {
		defer cancel()

		for {
			_, err := srv.Recv()
			if err != nil {
				if err != io.EOF && err != context.Canceled {
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
		sjob, err := s.state(ctx).JobPeekForRunner(ctx, record)
		if err != nil {
			return hcerr.Externalize(log, err, "failed to get job for runner", "id", record.Id)
		}
		if sjob == nil {
			return status.Errorf(codes.FailedPrecondition,
				"no pending job for this on-demand runner. A pending job "+
					"must be registered prior to registering the runner.")
		}

		// Set our job
		job = sjob.Job

		// We know a job was accepted, so it shouldn't be hanging around because this runner
		// is available.
		log.Trace("updating expiry time for job to be 60 seconds now that runner has been assigned job")
		dur, err := time.ParseDuration("60s")
		if err != nil {
			return status.Errorf(codes.FailedPrecondition,
				"Invalid expiry duration: %s", err.Error())
		}

		newExpireTime := timestamppb.New(time.Now().Add(dur))
		if err := s.state(ctx).JobUpdateExpiry(ctx, job.Id, newExpireTime); err != nil {
			return hcerr.Externalize(log, err, "failed to update job expiry time after runner accepted job", "id", job.Id)
		}

		log.Debug("runner is scoped for config",
			"project/application", job.Application,
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

		vars, err := s.state(ctx).ConfigGetWatch(ctx, configReq, ws)
		if err != nil {
			return hcerr.Externalize(log, err, "failed to get configuration variables")
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
			getConfigSourceRequest := &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},
			}
			if job != nil {
				getConfigSourceRequest.Workspace = job.Workspace
				getConfigSourceRequest.Scope = &pb.GetConfigSourceRequest_Application{
					Application: job.Application,
				}
			}
			sources, err := s.state(ctx).ConfigSourceGetWatch(ctx, getConfigSourceRequest, ws)
			if err != nil {
				return hcerr.Externalize(log, err, "failed to get the configuration for a dynamic source plugin")
			}

			config.ConfigSources = sources
		}

		// Send new config
		if err := srv.Send(&pb.RunnerConfigResponse{
			Config: config,
		}); err != nil {
			return hcerr.Externalize(
				log,
				err,
				"error sending runner config response",
			)
		}

		// Nil out the stuff we used so that if we're waiting awhile we can GC
		config = nil

		// Wait for any changes
		if err := ws.WatchCtx(ctx); err != nil {
			return hcerr.Externalize(
				log,
				err,
				"error waiting for runner config state changes",
			)
		}
	}
}

func runnerKind(runner *pb.Runner) string {
	switch runner.Kind.(type) {
	case *pb.Runner_Remote_:
		return "remote"
	case *pb.Runner_Local_:
		return "local"
	case *pb.Runner_Odr:
		return "odr"
	case *pb.Runner_DeprecatedIsOdr:
		return "deprecated-is-odr"
	default:
		return "unknown"
	}
}

func (s *Service) RunnerJobStream(
	server pb.Waypoint_RunnerJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	// Receive our opening message so we can determine the runner ID.
	req, err := server.Recv()
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to receive first message for RunnerJobStrem",
		)
	}
	reqEvent, ok := req.Event.(*pb.RunnerJobStreamRequest_Request_)
	if !ok {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.FailedPrecondition,
				"first message must be a Request event"),
			"first message to RunnerJobStream must be a Request event",
		)
	}

	// Get the runner to validate it is registered
	runnerId := reqEvent.Request.RunnerId

	runner, err := s.state(ctx).RunnerById(ctx, runnerId, nil)
	if err != nil {
		return hcerr.Externalize(log, err, "failed to get this runner", "id", runnerId)
	}
	log = log.With("runner_id", reqEvent.Request.RunnerId, "runner_type", runnerKind(runner))

	// The runner must be adopted to get a job.
	if runner.AdoptionState != pb.Runner_ADOPTED &&
		runner.AdoptionState != pb.Runner_PREADOPTED {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.FailedPrecondition,
				"runner must be adopted prior to requesting jobs"),
			"runner must be adopted prior to requesting jobs",
		)
	}

	// Verify our token matches the request
	if err := s.runnerVerifyToken(log, ctx, runner.Id, runner.Labels); err != nil {
		return hcerr.Externalize(
			log,
			err,
			"runner token verification failed",
		)
	}

	// Get the job for this runner. If this is a reattach, we lookup
	// the preexisting job. Otherwise, we assign a new job.
	var job *serverstate.Job
	reattach := false
	var jobId string
	if jobId = reqEvent.Request.ReattachJobId; jobId != "" {
		reattach = true

		log.Info("runner reattaching to an existing job", "job_id", jobId)
		job, err = s.state(ctx).JobById(ctx, jobId, nil)
		if err != nil {
			return hcerr.Externalize(log, err, "failed to get job", "id", jobId)
		}

		// If the job is not found, that is an error.
		if job == nil {
			return hcerr.Externalize(
				log,
				status.Errorf(codes.InvalidArgument,
					"reattach job ID does not exist"),
				"reattach job ID does not exist",
				"id",
				jobId,
			)
		}

		// The runner reattaching must be the assigned runner.
		assigned := job.Job.AssignedRunner
		if assigned == nil || assigned.Id != runner.Id {
			return hcerr.Externalize(
				log,
				status.Errorf(codes.InvalidArgument,
					"reattach job is not assigned to this runner"),
				"reattach job is not assigned to this runner",
				"id",
				jobId,
			)
		}

		// NOTE(mitchellh): things we should check in the future:
		// * job stream already open for this job ID
		// * job already in a terminal state
	} else {
		// Get a job assignment for this runner. This will block until
		// a job is available for the runner.
		log.Info("waiting for job assignment")
		job, err = s.state(ctx).JobAssignForRunner(ctx, runner)
		if err != nil {
			// TODO don't report this as an error in the logs if ctx is canceled, because
			// we stopped looking for a job when the context is canceled.
			return hcerr.Externalize(log, err, "failed to get job assignment for runner")
		}
		jobId = job.Id
	}
	if job == nil || job.Job == nil {
		panic("job is nil, should never be nil at this point")
	}
	log = log.With("job_id", jobId)

	operation := operationString(job.Job)
	defer func(start time.Time) {
		metrics.MeasureOperation(ctx, start, operation)
	}(time.Now())
	metrics.CountOperation(ctx, operation)

	log.Info("received job", "operation", operation)

	// Load config sourcers to send along with the job assignment
	cfgSrcs, err := s.state(ctx).ConfigSourceGetWatch(ctx, &pb.GetConfigSourceRequest{
		Scope: &pb.GetConfigSourceRequest_Application{
			Application: job.Application,
		},
		Workspace: job.Workspace,
	}, nil)
	if err != nil {
		return hcerr.Externalize(log, err, "failed to get the configuration for a dynamic source plugin to send with job assignment")
	}
	log.Trace("loaded config sources for job", "total_sourcers", len(cfgSrcs))

	log.Debug("sending job assignment to runner")
	// Send the job assignment.
	//
	// If this has an error, we continue to accumulate the error until
	// we set the ack status in the DB. We do this because if we fail to
	// send the job assignment we want to nack the job so it is queued again.
	err = server.Send(&pb.RunnerJobStreamResponse{
		Event: &pb.RunnerJobStreamResponse_Assignment{
			Assignment: &pb.RunnerJobStreamResponse_JobAssignment{
				Job:           job.Job,
				ConfigSources: cfgSrcs,
			},
		},
	})
	if err != nil {
		log.Warn("error sending job assignment to runner, will wait for ack", "err", err)
	}

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

	// We only ack if we're not reattached. If we reattached, then we can
	// only reattach to an already-acked job.
	if !reattach {
		// Send the ack OR nack, based on the value of +ack+.
		var ackerr error
		job, ackerr = s.state(ctx).JobAck(ctx, jobId, ack)
		if ackerr != nil {
			// If this fails, we just log, there is nothing more we can do.
			log.Warn("job ack failed", "outer_error", err, "error", ackerr)

			// Check if job is nil, so not to panic later on
			if job == nil {
				return hcerr.Externalize(log, ackerr, "job is nil, db might not be open")
			}
			// If we had no outer error, set the ackerr so that we exit. If
			// we do have an outer error, then the ack error only shows up in
			// the log.
			if err == nil {
				err = ackerr
			}
		}
	} else {
		// If we acked, we do nothing, cause reattachment only works
		// with already-acked job. We still require the ack from the client
		// to sync progress, but it has no state impact. If we nack, however,
		// we cancel the job.
		if !ack {
			log.Warn("reattach job was nacked, force cancelling")
			err = s.state(ctx).JobCancel(ctx, jobId, true)
		}
	}

	// If we have an error, return that. We also return if we didn't ack for
	// any reason. This error can be set at any point since job assignment.
	if err != nil {
		return hcerr.Externalize(log, err, "failed to ack the job or the job was cancelled", "id", jobId)
	}

	if ack {
		log.Info("runner as acknowledge the job")
	} else {
		log.Info("runner rejected the job")
		// If runners don't ack the job, this means close the stream
		return nil
	}

	var logStreamWriter logstream.Writer
	if s.logStreamProvider != nil {
		logStreamWriter, err = s.logStreamProvider.StartWriter(ctx, log, s.state(ctx), job)
		if err != nil {
			return hcerr.Externalize(log, err, "failed to start a log writer to handle jog logs")
		}
	}

	// We don't want the log stream writer to use the request context, because we want to
	// ensure that flushing occurs even if it needs to happen after the request context
	// is closed.
	logStreamCtx := context.Background()
	defer logStreamWriter.Flush(logStreamCtx)

	// Start a goroutine that watches for job changes
	jobCh := make(chan *serverstate.Job, 1)
	errCh := make(chan error, 1)
	go func() {
		for {
			ws := memdb.NewWatchSet()
			checkJob, err := s.state(ctx).JobById(ctx, jobId, ws)
			if err != nil {
				// if the context was canceled, we're all good, we're supposed to just exit.
				// It's critically important we don't put something on errCh here because
				// it can cause a job that finished fine to error out mistakenly.
				if err == context.Canceled {
					return
				}
				errCh <- errors.Wrapf(err, "error getting job by id from state %q", jobId)
				return
			}
			if checkJob == nil {
				errCh <- status.Errorf(codes.Internal, "failed to find job for id %q", jobId)
				return
			}

			// Send the job
			select {
			case jobCh <- checkJob:
			case <-ctx.Done():
				return
			}

			// Wait for the job to update
			if err := ws.WatchCtx(ctx); err != nil {
				if ctx.Err() == nil {
					errCh <- errors.Wrapf(err, "error on context while waiting for job %q", jobId)
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
				errCh <- errors.Wrapf(err, "failed receiving from runner for job %q", jobId)
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

	// This is wrapped in a closure so that we can easily capture the error return value
	err = func() error {

		// Recv events in a loop
		var lastJob *pb.Job

		for {
			select {
			case <-ctx.Done():
				// We need to drain the event channel
				for {
					select {
					case req := <-eventCh:
						if err := s.handleJobStreamRequest(log, job, server, req, logStreamWriter); err != nil {
							return hcerr.Externalize(log, err, "error handling job stream request during drain", "req", req)
						}
					default:
						return nil
					}
				}

			case err := <-errCh:
				return hcerr.Externalize(log, err, "err from err channel")

			case req := <-eventCh:
				if err := s.handleJobStreamRequest(log, job, server, req, logStreamWriter); err != nil {
					return hcerr.Externalize(log, err, "error handling job stream request", "req", req)
				}

			case updatedJob := <-jobCh:
				if lastJob == updatedJob.Job {
					continue
				}

				// If the job is canceled, send that event. We send this each time
				// the cancel time changes. The cancel time only changes if multiple
				// cancel requests are made.
				if updatedJob.CancelTime != nil &&
					(lastJob == nil || !lastJob.CancelTime.AsTime().Equal(updatedJob.CancelTime.AsTime())) {
					log.Trace("job cancellation request received")

					// The job is forced if we're in an error state. This must be true
					// because we would've already exited the loop if we naturally
					// got a terminal event.
					force := updatedJob.State == pb.Job_ERROR

					err := server.Send(&pb.RunnerJobStreamResponse{
						Event: &pb.RunnerJobStreamResponse_Cancel{
							Cancel: &pb.RunnerJobStreamResponse_JobCancel{
								Force: force,
							},
						},
					})
					if err != nil {
						return hcerr.Externalize(log, err, "error sending job cancel event to runner")
					}

					// On force we exit immediately.
					if force {
						return nil
					}
				}

				log.Trace("updating job from state store", "last_job", lastJob, "job", updatedJob.Job)
				lastJob = updatedJob.Job
			}
		}
	}()

	if err != nil {
		log.Info("job processing complete with error")
	} else {
		log.Info("job processing complete")
	}

	return err
}

func (s *Service) handleJobStreamRequest(
	log hclog.Logger,
	job *serverstate.Job,
	srv pb.Waypoint_RunnerJobStreamServer,
	req *pb.RunnerJobStreamRequest,
	logStreamWriter logstream.Writer,
) error {
	ctx := srv.Context()
	log.Trace("event received", "event", req.Event)
	switch event := req.Event.(type) {
	case *pb.RunnerJobStreamRequest_Complete_:
		if event.Complete == nil {
			log.Error("event.Complete was not filled when we expected it to be", "id", job.Id, "event", req.Event)
			return hcerr.Externalize(log, fmt.Errorf("failed to complete job"), "failed to complete job", "id", job.Id)
		}
		if err := s.state(ctx).JobComplete(ctx, job.Id, event.Complete.Result, nil); err != nil {
			return hcerr.Externalize(log, err, "failed to complete job", "id", job.Id)
		}
		log.Info("job completed properly")
	case *pb.RunnerJobStreamRequest_Error_:
		var remoteError error

		if event.Error != nil {
			remoteError = status.FromProto(event.Error.Error).Err()
		} else {
			log.Warn("job error did not include error description")
			remoteError = status.New(codes.Unknown, "remote error not correctly reported").Err()
		}

		if err := s.state(ctx).JobComplete(ctx, job.Id, nil, remoteError); err != nil {
			return hcerr.Externalize(log, err, "failed to complete job", "id", job.Id)
		}
		log.Info("job completed with error", "error", remoteError)
	case *pb.RunnerJobStreamRequest_Heartbeat_:
		if err := s.state(ctx).JobHeartbeat(ctx, job.Id); err != nil {
			return hcerr.Externalize(log, err, "job heartbeat failed", "id", job.Id)
		}
	case *pb.RunnerJobStreamRequest_Download:
		if event.Download == nil {
			log.Warn("invalid download event detected")
			return status.Error(codes.InvalidArgument, "download information missing")
		}

		if err := s.state(ctx).JobUpdateRef(ctx, job.Id, event.Download.DataSourceRef); err != nil {
			return hcerr.Externalize(log, err, "failed to update the job reference", "id", job.Id)
		}

		if err := s.state(ctx).ProjectUpdateDataRef(ctx, &pb.Ref_Project{
			Project: job.Application.Project,
		}, job.Workspace, event.Download.DataSourceRef); err != nil {
			return hcerr.Externalize(log, err, "failed to update the project", "project", job.Application.Project)
		}

	case *pb.RunnerJobStreamRequest_ConfigLoad_:
		if event.ConfigLoad == nil {
			log.Warn("invalid configload event detected")
			return status.Error(codes.InvalidArgument, "config information missing")
		}

		if err := s.state(ctx).JobUpdate(ctx, job.Id, func(jobpb *pb.Job) error {
			jobpb.Config = event.ConfigLoad.Config
			return nil
		}); err != nil {
			return hcerr.Externalize(log, err, "failed to update the job with config", "id", job.Id)
		}

	case *pb.RunnerJobStreamRequest_VariableValuesSet_:
		if event.VariableValuesSet == nil {
			log.Warn("invalid variableset event detected")
			return status.Error(codes.InvalidArgument, "variableset information missing")
		}

		if err := s.state(ctx).JobUpdate(ctx, job.Id, func(jobpb *pb.Job) error {
			jobpb.VariableFinalValues = event.VariableValuesSet.FinalValues
			return nil
		}); err != nil {
			return hcerr.Externalize(log, err, "failed to update the job with variables", "id", job.Id)
		}

	case *pb.RunnerJobStreamRequest_Terminal:
		// Write the events
		logStreamWriter.NewEvent(ctx, event)

		return nil

	default:
		log.Warn("unexpected event received", "event", hclog.Fmt("%T", req.Event))
	}

	return nil
}

// This verifies that the token present in the context "ctx" is valid
// for the runner ID specified in the function arguments.
func (s *Service) runnerVerifyToken(
	log hclog.Logger,
	ctx context.Context,
	realRunnerId string, // real runner ID
	runnerLabels map[string]string, // real runner labels
) error {
	// Get our token and reverify that we are adopted.
	tok := s.decodedTokenFromContext(ctx)
	if tok == nil {
		log.Error("no token, should not be possible")
		return status.Errorf(codes.Unauthenticated, "no token")
	}

	switch k := tok.Kind.(type) {
	case *pb.Token_Login_:
		// Legacy (pre WP 0.8) token. We accept these as preadopted.
		// NOTE(mitchellh): One day, we should reject these because modern
		// preadoption should be via runner tokens.

	case *pb.Token_Runner_:
		runnerId, err := s.decodeId(k.Runner.Id)
		if err != nil {
			log.Error("Failed to decode runner id while verifying runner token", "id", k.Runner.Id, "error", err)
			return status.Errorf(codes.Internal, "invalid runner id within token")
		}

		// A runner token. We validate here that we're not explicitly rejected.
		// We have to check again here because runner tokens can be created
		// for ANY runner, but we can reject a SPECIFIC runner.
		if runnerId != "" && !strings.EqualFold(runnerId, realRunnerId) {
			return status.Errorf(codes.PermissionDenied,
				"provided runner token is for a different runner")
		}

		// If the token has a label hash and it does not match our record,
		// then it is an error.
		if expected := k.Runner.LabelHash; expected > 0 {
			actual, err := serverptypes.RunnerLabelHash(runnerLabels)
			if err != nil {
				return err
			}

			if expected != actual {
				return status.Errorf(codes.PermissionDenied,
					"provided runner token is for a different set of runner labels")
			}
		}

	default:
		return status.Errorf(codes.PermissionDenied, "not a valid runner token")
	}

	return nil
}

func operationString(job *pb.Job) string {
	// Types that are assignable to Operation:
	switch job.Operation.(type) {
	case *pb.Job_Noop_:
		return "noop"
	case *pb.Job_Build:
		return "build"
	case *pb.Job_Push:
		return "push"
	case *pb.Job_Deploy:
		return "deploy"
	case *pb.Job_Destroy:
		return "destroy"
	case *pb.Job_Release:
		return "release"
	case *pb.Job_Validate:
		return "validate"
	case *pb.Job_Auth:
		return "auth"
	case *pb.Job_Docs:
		return "docs"
	case *pb.Job_ConfigSync:
		return "config_sync"
	case *pb.Job_Exec:
		return "exec"
	case *pb.Job_Up:
		return "up"
	case *pb.Job_Logs:
		return "logs"
	case *pb.Job_QueueProject:
		return "queue_project"
	case *pb.Job_Poll:
		return "poll"
	case *pb.Job_StatusReport:
		return "status_report"
	case *pb.Job_StartTask:
		return "start_task"
	case *pb.Job_StopTask:
		return "stop_task"
	case *pb.Job_WatchTask:
		return "watch_task"
	case *pb.Job_Init:
		return "init"
	}
	return "unknown"
}
