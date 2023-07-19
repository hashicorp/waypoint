// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/tokenutil"
)

var heartbeatDuration = 5 * time.Second

// AcceptParallel allows up to count jobs to be accepted and executing
// concurrently.
func (r *Runner) AcceptParallel(ctx context.Context, count int) {
	// Create a new cancellable context so we can stop all the goroutines
	// when one exits. We do this because if one exits, its likely that the
	// unrecoverable error exists in all.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start up all the goroutines
	r.logger.Info("accepting jobs concurrently", "count", count)
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer cancel()
			defer wg.Done()
			r.AcceptMany(ctx)
		}()
	}

	// Wait for them to exit
	wg.Wait()
}

// AcceptMany will accept jobs and execute them on after another as they are accepted.
// This is meant to be run in a goroutine and reports its own errors via r's logger.
func (r *Runner) AcceptMany(ctx context.Context) {
	for {
		if err := r.Accept(ctx); err != nil {
			if err == ErrClosed {
				return
			}

			switch status.Code(err) {
			case codes.Canceled:
				// Ideally we'd get ErrClosed, but there are cases where we'll observe
				// the context being closed first, in which case we honor that as a valid
				// reason to stop accepting jobs.
				return
			case codes.PermissionDenied:
				// This means the runner was deregistered and we must exit.
				// This won't be fixed unless the runner is closed and restarted.
				r.logger.Error("runner unexpectedly deregistered, exiting")
				time.Sleep(5 * time.Second)
				return

			case codes.NotFound:
				// This means the runner was deregistered and we must exit.
				// This won't be fixed unless the runner is closed and restarted.
				r.logger.Error("runner unexpectedly deregistered, exiting")
				return
			case codes.Unavailable, codes.Unimplemented:
				// Server became unavailable. Unimplemented likely means that the server
				// is running behind a proxy and is failing health checks.

				// Let's just sleep to give the server time to come back.
				r.logger.Warn("server unavailable, sleeping before retry", "error", err)
				time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
			default:
				r.logger.Error("error running job", "error", err)
			}
		}
	}
}

// Accept will accept and execute a single job. This will block until
// a job is available.
//
// An error is only returned if there was an error internal to the runner.
// Errors during job execution are expected (i.e. a project build is misconfigured)
// and will be reported on the job.
//
// Two specific errors to watch out for are:
//
//   - ErrClosed (in this package) which means that the runner is closed
//     and Accept can no longer be called.
//   - code = NotFound which means that the runner was deregistered. This
//     means the runner has to be fully recycled: Close called, a new runner
//     started.
//
// This is safe to be called concurrently which can be used to execute
// multiple jobs in parallel as a runner.
func (r *Runner) Accept(ctx context.Context) error {
	return r.accept(ctx, "")
}

// AcceptExact is the same as Accept except that it accepts only
// a job with exactly the given ID. This is used by Waypoint only in
// local execution mode as an extra security measure to prevent other
// jobs from being assigned to the runner.
func (r *Runner) AcceptExact(ctx context.Context, id string) error {
	return r.accept(ctx, id)
}

var testRecvDelay time.Duration

//nolint:govet,lostcancel
func (r *Runner) accept(ctx context.Context, id string) error {
	if r.readState(&r.stateExit) > 0 {
		return ErrClosed
	}

	log := r.logger

	// The runningCtx has the token that is set during runner adoption.
	// This is required for API calls to succeed. Put the token into ctx
	// as well so that this can be used for API calls.
	if tok := tokenutil.TokenFromContext(r.runningCtx); tok != "" {
		ctx = tokenutil.TokenWithContext(ctx, tok)
	}

	// Retry tracks whether we're trying a job stream connection or not.
	// We use this so that the first attempt fails fast so we can log it.
	// Subsequent attempts block.
	retry := false

	// State we want to initialize outside the retry label.
	var client pb.Waypoint_RunnerJobStreamClient
	var streamCtx context.Context
	var streamCancel context.CancelFunc
	var streamCtxLock sync.Mutex
	var stateGen uint64
	var err error

	// We wrap this in a func() so that we use the latest client value
	// and don't stack defers on retry.
	defer func() {
		if client != nil {
			client.CloseSend()
		}

		streamCtxLock.Lock()
		defer streamCtxLock.Unlock()
		if streamCancel != nil {
			streamCancel()
		}
	}()

	// If we have a timeout, then we setup a timer for accepting.
	var acceptTimer *time.Timer
	var canceled int32
	if r.acceptTimeout > 0 {
		acceptTimer = time.AfterFunc(r.acceptTimeout, func() {
			log.Error("runner timed out waiting for a job",
				"timeout", r.acceptTimeout.String())

			// Grab this lock before updating canceled. You don't
			// need to have this lock to touch canceled (we use atomic ops)
			// but it is used when the streamCancel is being reset so that
			// we don't set it up and race with cancellation.
			streamCtxLock.Lock()
			defer streamCtxLock.Unlock()

			// Mark that we canceled
			atomic.StoreInt32(&canceled, 1)

			// Cancel the context
			if streamCancel != nil {
				streamCancel()
			}
		})
	}

RESTART_JOB_STREAM:
	// If we're retrying, these might be non-nil and we want to do some clean-up
	if retry {
		log.Warn("server down before accepting a job, will reconnect")

		if client != nil {
			log.Debug("closing client send")
			client.CloseSend()
		}
	}

	// Setup a new context that we can cancel at any time to close the stream.
	// We use this for timeouts.
	//
	// Note: we disable the lostcancel linter for streamCancel because
	// golangci-lint is not detecting that we have the defer above the
	// label as well as the retry block above.
	log.Debug("acquiring stream ctx lock")
	streamCtxLock.Lock()
	if streamCancel != nil {
		log.Debug("canceling stream context")
		streamCancel()
	}
	if atomic.LoadInt32(&canceled) > 0 {
		log.Debug("stream context is canceled - unlocking and returning")
		streamCtxLock.Unlock()
		return ErrTimeout
	}
	streamCtx, streamCancel = context.WithCancel(r.runningCtx)
	streamCtxLock.Unlock()

	// Since this is a disconnect, we have to wait for our
	// RunnerConfig stream to re-establish. We wait for the config
	// generation to increment.
	if retry {
		log.Debug("waiting for state greater than", "r.stateConfig", r.stateConfig, "stateGen", stateGen)
		if r.waitStateGreater(&r.stateConfig, stateGen) {
			log.Debug("early exit while waiting for reconnect")
			return status.Error(codes.Internal, "early exit while waiting for reconnect")
		}
		log.Debug("I bet we don't see this log line! If you see this, it means we're not getting stuck in waitStateGreater")
	}

	log.Debug("Stream (re)-established")

	// Get our configuration state value. We use this so that we can detect
	// when we've reconnected during failures.
	stateGen = r.readState(&r.stateConfig)

	// Open a new job stream. This retries on connection errors. Note that
	// this retry loop doesn't respect the accept timeout because gRPC has no way
	// to time out of a "WaitForReady" RPC call (it ignores context cancellation,
	// too). TODO: do a manual backoff with WaitForReady(false) so we can
	// weave in accept timeout.
	log.Debug("opening job stream", "retry", retry)
	client, err = r.client.RunnerJobStream(streamCtx, grpc.WaitForReady(retry))
	retry = true
	if err != nil {
		if atomic.LoadInt32(&canceled) > 0 ||
			status.Code(err) == codes.Unavailable ||
			status.Code(err) == codes.NotFound {
			// Throttle ourselves so that we don't hammer the server in the case that
			// we've been deleted and are not likely to return.
			time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)

			log.Trace("Restarting the accept loop due to a cancellation and we got an error establishing the runner job stream. I don't think we'll see this.", "err", err)
			goto RESTART_JOB_STREAM
		}

		return err
	}

	// Send our request
	log.Trace("sending job request")
	if err := client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: r.id,
			},
		},
	}); err != nil {
		if atomic.LoadInt32(&canceled) > 0 ||
			status.Code(err) == codes.Unavailable ||
			status.Code(err) == codes.NotFound {
			log.Trace("Restarting the accept loop due to a cancellation and we got an error sending on the job stream. I don't think we'll see this.", "err", err)
			goto RESTART_JOB_STREAM
		}

		return err
	}

	// Wait for an assignment
	log.Info("waiting for job assignment")

	// NOTE: if r.runningCtx is canceled, because the runner has finished closing,
	// any job sent won't be acked, but the server will see an error on waiting
	// for us to ack the job, and auto-nack it.
	resp, err := client.Recv()
	if err != nil {
		if atomic.LoadInt32(&canceled) > 0 ||
			status.Code(err) == codes.Unavailable ||
			status.Code(err) == codes.NotFound {
			log.Trace("Restarting the accept loop due to a cancellation and we got an error receiving the runner job stream. I don't think we'll see this.", "err", err)
			goto RESTART_JOB_STREAM
		}

		return err
	}

	// Be sure to stop the timer so that we don't cancel the context after this point.
	if acceptTimer != nil {
		acceptTimer.Stop()
	}

	// We received an assignment!
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	if !ok {
		return status.Errorf(codes.Aborted,
			"expected job assignment, server sent %T",
			resp.Event)
	}
	jobId := assignment.Assignment.Job.Id
	log = log.With(
		"job_id", jobId,
		"job_op", fmt.Sprintf("%T", assignment.Assignment.Job.Operation),
	)
	log.Info("job assignment received")

	// Used to test the behavior of accepting a job while
	// the runner is shutting down.
	if testRecvDelay != 0 {
		time.Sleep(testRecvDelay)
	}

	// We need to register ourselves as being worthy to be waited on
	// since prior to this, if Close() were called, we could go ahead
	// and return.

	r.runningCond.L.Lock()
	shutdown := r.readState(&r.stateExit) > 0
	if !shutdown {
		r.runningJobs++
	}
	r.runningCond.L.Unlock()

	defer func() {
		r.runningCond.L.Lock()
		defer r.runningCond.L.Unlock()

		r.runningJobs--
		r.runningCond.Broadcast()
	}()

	if shutdown {
		return errors.Wrapf(ErrClosed, "runner shutdown, dropped job: %s", jobId)
	}

	// If this isn't the job we expected then we nack and error.
	if id != "" {
		if jobId != id {
			log.Warn("unexpected job id for exact match, nacking")
			if err := client.Send(&pb.RunnerJobStreamRequest{
				Event: &pb.RunnerJobStreamRequest_Error_{
					Error: &pb.RunnerJobStreamRequest_Error{},
				},
			}); err != nil {
				// We don't restart the accept here on disconnect because
				// we already know we're in an error state that was truly
				// unexpected and erroneous: the server gave us a job that
				// wasn't assignd to us! Let's return.

				return err
			}

			return status.Errorf(codes.Aborted, "server sent us an invalid job")
		}

		log.Trace("assigned job matches expected ID for local mode")
	}

	// Ack the assignment
	log.Trace("acking job assignment")
	if err := client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}); err != nil {
		// This is sort of sketchy situation, but this comment is here to tell
		// you why its safe. At this point, the error should only be if the ack
		// failed to send so the server shouldn't have received the ack. However,
		// if they did, this goto will abandon the job. That's okay, it'll be
		// stuck for the heartbeat period and then the job manager will kill
		// it. That's unfortunate but unlikely to happen in practice and not
		// a bad outcome since no logic is ever executed for the job.
		if status.Code(err) == codes.Unavailable ||
			status.Code(err) == codes.NotFound {
			log.Trace("Congratulations, you are in an unlikely sketchy situation.", "err", err)
			goto RESTART_JOB_STREAM
		}

		return err
	}

	// Now that we've acked the job, we can create the re-attachable client.
	// Note: we use this context and not the new one below so that we
	// continue to reconnect even if our job is done since we need to still
	// send job complete messages.
	client = &reattachClient{
		ctx:    ctx,
		client: client,
		log:    log.Named("job_stream").With("job_id", jobId),
		runner: r,
		jobId:  jobId,
	}

	// Create a cancelable context so we can stop if job is canceled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// We need a mutex to protect against simultaneous sends to the client.
	var sendMutex sync.Mutex

	// For our UI, we always send output to the server. If we have a local UI
	// set, we mirror to that as well.
	var ui terminal.UI = &runnerUI{
		ctx:    ctx,
		cancel: cancel,
		evc:    client,
		mu:     &sendMutex,
	}
	if r.ui != nil {
		ui = &multiUI{
			UIs: []terminal.UI{r.ui, ui},
		}
	}

	// Start up a goroutine to listen for any other events
	errCh := make(chan error, 1)
	go func() {
		for {
			// Wait for the connection to close. We do this because this ensures
			// that the server received our completion and updated the database.
			resp, err = client.Recv()
			if err != nil {
				errCh <- err
				return
			}

			// Determine the event
			switch resp.Event.(type) {
			case *pb.RunnerJobStreamResponse_Cancel:
				log.Info("job cancellation request received, canceling")
				cancel()

			default:
				log.Info("unknown job event", "event", resp.Event)
			}
		}
	}()

	// Heartbeat
	go func() {
		tick := time.NewTicker(heartbeatDuration)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
			}

			sendMutex.Lock()
			err := client.Send(&pb.RunnerJobStreamRequest{
				Event: &pb.RunnerJobStreamRequest_Heartbeat_{
					Heartbeat: &pb.RunnerJobStreamRequest_Heartbeat{},
				},
			})
			sendMutex.Unlock()
			if err != nil && err != io.EOF {
				log.Warn("error during heartbeat", "err", err)
			}
		}
	}()

	// The job stream setup is done. Actually run the job, download any
	// data necessary, setup the core, etc
	log.Info("starting job execution")
	result, err := r.prepareAndExecuteJob(ctx, log, ui, &sendMutex, client, assignment.Assignment)
	log.Debug("job finished", "error", err)

	// We won't output anything else to the UI anymore.
	if ui, ok := ui.(*runnerUI); ok {
		ui.Close()
	}

	// Check if we were force canceled. If so, then just exit now. Realistically
	// we could also be force cancelled at any point below but this is the
	// most likely spot to catch it and the error scenario below is not bad.
	if ctx.Err() != nil {
		select {
		case err := <-errCh:
			// If we got an EOF then we were force cancelled.
			if err == io.EOF {
				log.Info("job force canceled")
				return nil
			}
		default:
		}
	}

	// If we have a nil result, then use an empty result so its set to
	// SOMETHING but just with no values.
	if result == nil {
		result = &pb.Job_Result{}
	}

	// For the remainder of the job, we're going to hold the mutex. We are
	// just sending quick status updates so this should not block anything
	// for very long.
	sendMutex.Lock()
	defer sendMutex.Unlock()

	// Handle job execution errors
	if err != nil {
		st, _ := status.FromError(err)

		log.Warn("error during job execution", "err", err)
		if rpcerr := client.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Error_{
				Error: &pb.RunnerJobStreamRequest_Error{
					Error: st.Proto(),
				},
			},
		}); rpcerr != nil {
			log.Warn("error sending error event, job may be dangling", "err", rpcerr)
			return rpcerr
		}
	} else {
		// Complete the job
		log.Debug("sending job completion")
		if err := client.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Complete_{
				Complete: &pb.RunnerJobStreamRequest_Complete{
					Result: result,
				},
			},
		}); err != nil {
			log.Error("error sending job complete message", "error", err)
			return err
		}
	}

	// Wait for the connection to close. We do this because this ensures
	// that the server received our completion and updated the database.
	err = <-errCh
	if err == io.EOF {
		return nil
	}

	return err
}

func (r *Runner) prepareAndExecuteJob(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	sendMutex *sync.Mutex,
	client pb.Waypoint_RunnerJobStreamClient,
	assignment *pb.RunnerJobStreamResponse_JobAssignment,
) (*pb.Job_Result, error) {
	job := assignment.Job
	log.Trace("preparing to execute job operation", "type", hclog.Fmt("%T", job.Operation))

	// Some operation types don't need to download data, execute those here.
	switch job.Operation.(type) {
	case *pb.Job_Poll:
		return r.executePollOp(ctx, log, ui, job)
	case *pb.Job_StartTask:
		return r.executeStartTaskOp(ctx, log, ui, job)
	case *pb.Job_StopTask:
		return r.executeStopTaskOp(ctx, log, ui, job)
	case *pb.Job_WatchTask:
		return r.executeWatchTaskOp(ctx, log, ui, job)
	}

	// We need to get our data source next prior to executing.
	var result *pb.Job_Result
	wd, ref, closer, err := r.downloadJobData(
		ctx,
		log,
		ui,
		job.DataSource,
		job.DataSourceOverrides,
	)
	if err == nil {
		log.Debug("job data downloaded (or local)",
			"pwd", wd,
			"ref", fmt.Sprintf("%#v", ref),
		)

		if closer != nil {
			defer func() {
				log.Debug("cleaning up downloaded data")
				if err := closer(); err != nil {
					log.Warn("error cleaning up data", "err", err)
				}
			}()
		}

		// Send our download info
		if ref != nil {
			log.Debug("sending download event")

			sendMutex.Lock()
			err = client.Send(&pb.RunnerJobStreamRequest{
				Event: &pb.RunnerJobStreamRequest_Download{
					Download: &pb.GetJobStreamResponse_Download{
						DataSourceRef: ref,
					},
				},
			})
			sendMutex.Unlock()
		}

		// We want the working directory to always be absolute.
		if err == nil && !filepath.IsAbs(wd) {
			err = status.Errorf(codes.Internal,
				"data working directory should be absolute. This is a bug, please report it.")
		}

		if err == nil {
			// Execute the job. We have to close the UI right afterwards to
			// ensure that no more output is written to the client.
			result, err = r.executeJob(ctx, log, ui, assignment, wd, sendMutex, client)
		}
	}

	return result, err
}
