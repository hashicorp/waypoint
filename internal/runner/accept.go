package runner

import (
	"context"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var heartbeatDuration = 5 * time.Second

// Accept will accept and execute a single job. This will block until
// a job is available.
//
// An error is only returned if there was an error internal to the runner.
// Errors during job execution are expected (i.e. a project build is misconfigured)
// and will be reported on the job.
//
// This is safe to be called concurrently which can be used to execute
// multiple jobs in parallel as a runner.
func (r *Runner) Accept(ctx context.Context) error {
	if r.closed() {
		return ErrClosed
	}

	log := r.logger

	// Open a new job stream. NOTE: we purposely do NOT use ctx above
	// since if the context is cancelled we want to continue reporting
	// errors.
	log.Debug("opening job stream")
	client, err := r.client.RunnerJobStream(context.Background())
	if err != nil {
		return err
	}
	defer client.CloseSend()

	// Send our request
	log.Trace("sending job request")
	if err := client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: r.id,
			},
		},
	}); err != nil {
		return err
	}

	// Wait for an assignment
	log.Info("waiting for job assignment")
	resp, err := client.Recv()
	if err != nil {
		return err
	}

	// We received an assignment!
	assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	if !ok {
		return status.Errorf(codes.Aborted,
			"expected job assignment, server sent %T",
			resp.Event)
	}
	log = log.With("job_id", assignment.Assignment.Job.Id)
	log.Info("job assignment received")

	// We increment the waitgroup at this point since prior to this if we're
	// forcefully quit, we shouldn't have acked. This is somewhat brittle so
	// a todo here is to build a better notification mechanism that we've quit
	// and exit here.
	r.acceptWg.Add(1)
	defer r.acceptWg.Done()

	// Ack the assignment
	log.Trace("acking job assignment")
	if err := client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}); err != nil {
		return err
	}

	// Create a cancelable context so we can stop if job is canceled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// We need a mutex to protect against simultaneous sends to the client.
	var sendMutex sync.Mutex

	// For our UI, we will use a manually set UI if available. Otherwise,
	// we setup the runner UI which streams the output to the server.
	ui := r.ui
	if ui == nil {
		ui = &runnerUI{
			ctx:    ctx,
			cancel: cancel,
			evc:    client,
			mu:     &sendMutex,
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

	// Execute the job. We have to close the UI right afterwards to
	// ensure that no more output is writting to the client.
	log.Info("starting job execution")
	result, err := r.executeJob(ctx, log, ui, assignment.Assignment.Job)
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
		}

		return nil
	}

	// Complete the job
	if err := client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{
				Result: result,
			},
		},
	}); err != nil {
		return err
	}

	// Wait for the connection to close. We do this because this ensures
	// that the server received our completion and updated the database.
	err = <-errCh
	if err == io.EOF {
		return nil
	}

	return err
}
