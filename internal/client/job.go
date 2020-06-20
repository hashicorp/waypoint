package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// job returns the basic job skeleton prepoulated with the correct
// defaults based on how the client is configured. For example, for local
// operations, this will already have the targeting for the local runner.
func (c *Client) job() *pb.Job {
	return &pb.Job{
		Application:  c.application,
		TargetRunner: c.runner,

		DataSource: &pb.Job_Local_{
			Local: &pb.Job_Local{},
		},

		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},
	}
}

// doJob will queue and execute the job. If the client is configured for
// local mode, this will start and target the proper runner.
func (c *Client) doJob(ctx context.Context, job *pb.Job) error {
	log := c.logger

	// In local mode we have to start a runner.
	if c.local {
		log.Info("local mode, starting local runner")
		r, err := c.startRunner()
		if err != nil {
			return err
		}

		log.Info("runner started", "runner_id", r.Id())

		// We defer the close so that we clean up resources. Local mode
		// always blocks and streams the full output so when doJob exits
		// the job is complete.
		defer r.Close()

		// Accept a job. Our local runners execute exactly one job.
		go func() {
			if err := r.Accept(); err != nil {
				log.Error("runner job accept error", "err", err)
			}
		}()

		// Modify the job to target this runner and use the local data source.
		job.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: r.Id(),
				},
			},
		}
	}

	return c.queueAndStreamJob(ctx, job)
}

// queueAndStreamJob will queue the job. If the client is configured to watch the job,
// it'll also stream the output to the configured UI.
func (c *Client) queueAndStreamJob(ctx context.Context, job *pb.Job) error {
	log := c.logger

	// Queue the job
	log.Debug("queueing job", "operation", fmt.Sprintf("%T", job.Operation))
	queueResp, err := c.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: job,
	})
	if err != nil {
		return err
	}
	log = log.With("job_id", queueResp.JobId)

	// Get the stream
	log.Debug("opening job stream")
	stream, err := c.client.GetJobStream(ctx, &pb.GetJobStreamRequest{
		JobId: queueResp.JobId,
	})
	if err != nil {
		return err
	}

	// Wait for open confirmation
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if _, ok := resp.Event.(*pb.GetJobStreamResponse_Open_); !ok {
		return status.Errorf(codes.Aborted,
			"job stream failed to open, got unexpected message %T",
			resp.Event)
	}

	// Process events
	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		if resp == nil {
			// This shouldn't happen, but if it does, just ignore it.
			log.Warn("nil response received, ignoring")
			continue
		}

		switch event := resp.Event.(type) {
		case *pb.GetJobStreamResponse_Complete_:
			if event.Complete.Error == nil {
				log.Info("job completed successfully")
				return nil
			}

			st := status.FromProto(event.Complete.Error)
			log.Warn("job failed", "code", st.Code(), "message", st.Message())
			return st.Err()

		case *pb.GetJobStreamResponse_Error_:
			st := status.FromProto(event.Error.Error)
			log.Warn("job stream failure", "code", st.Code(), "message", st.Message())
			return st.Err()

		case *pb.GetJobStreamResponse_Terminal_:
			for _, line := range event.Terminal.Lines {
				log.Trace("job terminal output", "line", line.Raw)
				c.ui.Output(line.Line)
			}

		default:
			log.Warn("unknown stream event", "event", resp.Event)
		}
	}
}
