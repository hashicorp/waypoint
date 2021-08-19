package clijob

import (
	"context"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const stateEventPause = 1500 * time.Millisecond

// WatchJobStream will take a given job id that has already been sent off
// to be queued by the server and attempt to stream the output directly
// into the ui terminal passed in through from a CLI caller. Note that the
// original implementation of this can be found in 'client.queueAndStreamJob'.
func WatchJobStream(
	ctx context.Context,
	log hclog.Logger,
	client pb.WaypointClient,
	ui terminal.UI,
	jobId string,
) error {
	log = log.Named("clijob_stream")
	log = log.With("job_id", jobId)

	// Get the stream
	log.Debug("opening job stream")
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{
		JobId: jobId,
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

	type stepData struct {
		terminal.Step

		out io.Writer
	}

	// Process events
	var (
		stateEventTimer *time.Timer
		tstatus         terminal.Status

		stdout, stderr io.Writer

		sg    terminal.StepGroup
		steps = map[int32]*stepData{}
	)

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
			for _, ev := range event.Terminal.Events {
				log.Trace("job terminal output", "event", ev)

				switch ev := ev.Event.(type) {
				case *pb.GetJobStreamResponse_Terminal_Event_Line_:
					ui.Output(ev.Line.Msg, terminal.WithStyle(ev.Line.Style))
				case *pb.GetJobStreamResponse_Terminal_Event_NamedValues_:
					var values []terminal.NamedValue

					for _, tnv := range ev.NamedValues.Values {
						values = append(values, terminal.NamedValue{
							Name:  tnv.Name,
							Value: tnv.Value,
						})
					}

					ui.NamedValues(values)
				case *pb.GetJobStreamResponse_Terminal_Event_Status_:
					if tstatus == nil {
						tstatus = ui.Status()
						defer tstatus.Close()
					}

					if ev.Status.Msg == "" && !ev.Status.Step {
						tstatus.Close()
					} else if ev.Status.Step {
						tstatus.Step(ev.Status.Status, ev.Status.Msg)
					} else {
						tstatus.Update(ev.Status.Msg)
					}
				case *pb.GetJobStreamResponse_Terminal_Event_Raw_:
					if stdout == nil {
						stdout, stderr, err = ui.OutputWriters()
						if err != nil {
							return err
						}
					}

					if ev.Raw.Stderr {
						stderr.Write(ev.Raw.Data)
					} else {
						stdout.Write(ev.Raw.Data)
					}
				case *pb.GetJobStreamResponse_Terminal_Event_Table_:
					tbl := terminal.NewTable(ev.Table.Headers...)

					for _, row := range ev.Table.Rows {
						var trow []terminal.TableEntry

						for _, ent := range row.Entries {
							trow = append(trow, terminal.TableEntry{
								Value: ent.Value,
								Color: ent.Color,
							})
						}
					}

					ui.Table(tbl)
				case *pb.GetJobStreamResponse_Terminal_Event_StepGroup_:
					if sg != nil {
						sg.Wait()
					}

					if !ev.StepGroup.Close {
						sg = ui.StepGroup()
					}
				case *pb.GetJobStreamResponse_Terminal_Event_Step_:
					if sg == nil {
						continue
					}

					step, ok := steps[ev.Step.Id]
					if !ok {
						step = &stepData{
							Step: sg.Add(ev.Step.Msg),
						}
						steps[ev.Step.Id] = step
					} else {
						if ev.Step.Msg != "" {
							step.Update(ev.Step.Msg)
						}
					}

					if ev.Step.Status != "" {
						if ev.Step.Status == terminal.StatusAbort {
							step.Abort()
						} else {
							step.Status(ev.Step.Status)
						}
					}

					if len(ev.Step.Output) > 0 {
						if step.out == nil {
							step.out = step.TermOutput()
						}

						step.out.Write(ev.Step.Output)
					}

					if ev.Step.Close {
						step.Done()
					}
				default:
					log.Error("Unknown terminal event seen", "type", hclog.Fmt("%T", ev))
				}
			}
		case *pb.GetJobStreamResponse_State_:
			// Stop any state event timers if we have any since the state
			// has changed and we don't want to output that information anymore.
			if stateEventTimer != nil {
				stateEventTimer.Stop()
				stateEventTimer = nil
			}

			// For certain states, we do a quality of life UI message if
			// the wait time ends up being long.
			switch event.State.Current {
			case pb.Job_QUEUED:
				stateEventTimer = time.AfterFunc(stateEventPause, func() {
					ui.Output("Operation is queued. Waiting for runner assignment...",
						terminal.WithHeaderStyle())
					ui.Output("If you interrupt this command, the job will still run in the background.",
						terminal.WithInfoStyle())
				})

			case pb.Job_WAITING:
				stateEventTimer = time.AfterFunc(stateEventPause, func() {
					ui.Output("Operation is assigned to a runner. Waiting for start...",
						terminal.WithHeaderStyle())
					ui.Output("If you interrupt this command, the job will still run in the background.",
						terminal.WithInfoStyle())
				})
			}

		default:
			log.Warn("unknown stream event", "event", resp.Event)
		}
	}
}
