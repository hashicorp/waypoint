package jobstream

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/finalcontext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Stream a single job and return the result from that job. This function
// blocks until the job reaches a terminal state (success or fail).
func Stream(ctx context.Context, jobId string, opts ...Option) (*pb.Job_Result, error) {
	s := &stream{
		jobId: jobId,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	if s.log == nil {
		s.log = hclog.L()
	}
	if s.client == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"client must be set")
	}

	return s.Run(ctx)
}

type stream struct {
	jobId          string
	log            hclog.Logger
	client         pb.WaypointClient
	ui             terminal.UI
	cancelOnErr    bool
	ignoreTerminal bool
	stateCh        chan<- pb.Job_State
}

// Get the job stream for a single job, handle all the events, and
// return the final result of the job execution.
func (s *stream) Run(ctx context.Context) (*pb.Job_Result, error) {
	log := s.log

	// Get the stream
	log.Debug("opening job stream")
	stream, err := s.client.GetJobStream(ctx, &pb.GetJobStreamRequest{
		JobId: s.jobId,
	})
	if err != nil {
		return nil, err
	}

	// Wait for open confirmation
	resp, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	if _, ok := resp.Event.(*pb.GetJobStreamResponse_Open_); !ok {
		return nil, status.Errorf(codes.Aborted,
			"job stream failed to open, got unexpected message %T",
			resp.Event)
	}

	// This timer is used to track whether we're stuck in certain states for
	// too long and show a UI message. For example, if we're queued for a long
	// time we notify the user we're queued.
	var stateEventTimer *time.Timer

	// The UI that will translate terminal events into UI calls.
	ui := s.ui
	streamUI := &UI{UI: ui}

	// If we're canceling the job on non-successful exit, then setup
	// the defer so that we do that.
	var completed bool
	if s.cancelOnErr {
		defer func() {
			// If we completed then do nothing, or if the context is still
			// active since this means that we're not cancelled.
			if completed || ctx.Err() == nil {
				return
			}

			ctx, cancel := finalcontext.Context(log)
			defer cancel()

			log.Warn("canceling job")
			_, err := s.client.CancelJob(ctx, &pb.CancelJobRequest{
				JobId: s.jobId,
			})
			if err != nil {
				log.Warn("error canceling job", "err", err)
			} else {
				log.Info("job cancelled successfully")
			}
		}()
	}

	var assignedRunner *pb.Ref_RunnerId
	for {
		resp, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		if resp == nil {
			// This shouldn't happen, but if it does, just ignore it.
			log.Warn("nil response received, ignoring")
			continue
		}

		switch event := resp.Event.(type) {
		case *pb.GetJobStreamResponse_Complete_:
			completed = true

			if event.Complete.Error == nil {
				log.Info("job completed successfully")
				return event.Complete.Result, nil
			}

			st := status.FromProto(event.Complete.Error)
			log.Warn("job failed", "code", st.Code(), "message", st.Message())
			return nil, st.Err()

		case *pb.GetJobStreamResponse_Error_:
			completed = true

			st := status.FromProto(event.Error.Error)
			log.Warn("job stream failure", "code", st.Code(), "message", st.Message())
			return nil, st.Err()

		case *pb.GetJobStreamResponse_Download_:
			if ui != nil {
				ui.Output("Downloading from Git", terminal.WithHeaderStyle())

				// Assume git type for now
				git := event.Download.DataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git)

				ui.Output("Git Commit: %s", git.Git.Commit, terminal.WithInfoStyle())
				ui.Output(" Timestamp: %s", git.Git.Timestamp.AsTime(), terminal.WithInfoStyle())
				ui.Output("   Message: %s", git.Git.CommitMessage, terminal.WithInfoStyle())
			}

		case *pb.GetJobStreamResponse_Terminal_:
			if s.ignoreTerminal {
				continue
			}

			if ui != nil {
				if err := streamUI.Write(event.Terminal.Events); err != nil {
					log.Warn("job stream UI failure", "err", err)
				}
			}

		case *pb.GetJobStreamResponse_Job:
			// Job changed, we don't use this information

		case *pb.GetJobStreamResponse_State_:
			// Stop any state event timers if we have any since the state
			// has changed and we don't want to output that information anymore.
			if stateEventTimer != nil {
				stateEventTimer.Stop()
				stateEventTimer = nil
			}

			// Check if this job has been assigned a runner for the first time
			if event.State != nil &&
				event.State.Job != nil &&
				event.State.Job.AssignedRunner != nil &&
				assignedRunner == nil {

				assignedRunner = event.State.Job.AssignedRunner

				runner, err := s.client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: assignedRunner.Id})
				if err != nil {
					if ui != nil {
						ui.Output("Failed to inspect the runner (id %q) assigned for this operation: %s", assignedRunner.Id, err, terminal.WithErrorStyle())
					}

					break
				}
				switch runnerType := runner.Kind.(type) {
				case *pb.Runner_Local_:
					if ui != nil {
						ui.Output("Performing operation locally", terminal.WithHeaderStyle())
					}
				case *pb.Runner_Remote_:
					if ui != nil {
						ui.Output("Performing this operation on a remote runner with id %q", runner.Id, terminal.WithInfoStyle())
					}
				case *pb.Runner_Odr:
					log.Debug("Executing operation with on-demand runner from profile", "runner_profile_id", runnerType.Odr.ProfileId)
					profile, err := s.client.GetOnDemandRunnerConfig(
						ctx, &pb.GetOnDemandRunnerConfigRequest{
							Config: &pb.Ref_OnDemandRunnerConfig{
								Id: runnerType.Odr.ProfileId,
							},
						})
					if ui != nil {
						if err != nil {
							ui.Output("Performing operation on an on-demand runner from profile with ID %q", runnerType.Odr.ProfileId, terminal.WithInfoStyle())
							ui.Output("Failed inspecting runner profile with id %q: %s", runnerType.Odr.GetProfileId(), err, terminal.WithErrorStyle())
						} else {
							ui.Output("Performing operation on %q with runner profile %q", profile.Config.PluginType, profile.Config.Name, terminal.WithInfoStyle())
						}
					}
				}
			}

			// For certain states, we do a quality of life UI message if
			// the wait time ends up being long.
			switch event.State.Current {
			case pb.Job_QUEUED:
				if ui != nil {
					stateEventTimer = time.AfterFunc(stateEventPause, func() {
						ui.Output(
							"Operation is queued waiting for job %q. Waiting for runner assignment...",
							s.jobId,
							terminal.WithHeaderStyle())
						ui.Output(
							"If you interrupt this command, the job will still run in the background.",
							terminal.WithInfoStyle())
					})
				}

			case pb.Job_WAITING:
				if ui != nil {
					stateEventTimer = time.AfterFunc(stateEventPause, func() {
						ui.Output("Operation is assigned to a runner. Waiting for start...",
							terminal.WithHeaderStyle())
						ui.Output("If you interrupt this command, the job will still run in the background.",
							terminal.WithInfoStyle())
					})
				}
			}

			if s.stateCh != nil {
				select {
				case <-ctx.Done():
					break
				case s.stateCh <- event.State.Current:
					// ok
				}
			}

		default:
			log.Warn("unknown stream event",
				"type", fmt.Sprintf("%T", resp.Event),
				"event", resp.Event,
			)
		}
	}
}

// The time here is meant to encompass the typical case for an operation to begin.
// With the introduction of ondemand runners, we bumped it up from 1500 to 3000
// to accomidate the additional time before the job was picked up when testing in
// local Docker.
const stateEventPause = 3000 * time.Millisecond
