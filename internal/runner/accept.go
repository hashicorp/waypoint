package runner

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Accept will accept and execute a single job. This will block until
// a job is available.
//
// This is safe to be called concurrently which can be used to execute
// multiple jobs in parallel as a runner.
func (r *Runner) Accept() error {
	if r.closed() {
		return ErrClosed
	}

	log := r.logger

	// Open a new job stream
	log.Debug("opening job stream")
	client, err := r.client.RunnerJobStream(r.ctx)
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

	ctx, cancel := context.WithCancel(r.ctx)

	// For our UI, we will use a manually set UI if available. Otherwise,
	// we setup the runner UI which streams the output to the server.
	ui := r.ui
	if ui == nil {
		ui = &runnerUI{
			ctx:    ctx,
			cancel: cancel,
			evc:    client,
		}
	}

	// Execute the job. We have to close the UI right afterwards to
	// ensure that no more output is writting to the client.
	log.Info("starting job execution")
	result, err := r.executeJob(r.ctx, log, ui, assignment.Assignment.Job)
	if ui, ok := ui.(*runnerUI); ok {
		ui.Close()
	}

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

		return err
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
	_, err = client.Recv()
	if err == io.EOF {
		return nil
	}

	return err
}
