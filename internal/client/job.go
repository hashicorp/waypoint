package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/finalcontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// job returns the basic job skeleton prepoulated with the correct
// defaults based on how the client is configured. For example, for local
// operations, this will already have the targeting for the local runner.
func (c *Project) job() *pb.Job {
	job := &pb.Job{
		TargetRunner: c.runner,
		Labels:       c.labels,
		Variables:    c.variables,
		Workspace:    c.workspace,
		Application: &pb.Ref_Application{
			Project: c.project.Project,
		},

		DataSource: &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Local{
				Local: &pb.Job_Local{},
			},
		},
		DataSourceOverrides: c.dataSourceOverrides,

		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},
	}

	// If we're not local, we set a nil data source so it defaults to
	// wahtever the project has remotely.
	if !c.local {
		job.DataSource = nil
	}

	return job
}

// doJob will queue and execute the job. If the client is configured for
// local mode, this will start and target the proper runner.
func (c *Project) doJob(ctx context.Context, job *pb.Job, ui terminal.UI) (*pb.Job_Result, error) {
	return c.doJobMonitored(ctx, job, ui, nil)
}

// Same as doJob, but with the addition of a mon channel that can be used
// to monitor the job status as it changes.
// The receiver must be careful to not block sending to mon as it will block
// the job state processing loop.
func (c *Project) doJobMonitored(ctx context.Context, job *pb.Job, ui terminal.UI, monCh chan pb.Job_State) (*pb.Job_Result, error) {
	// Be sure that the monitor is closed so the reciever knows for sure the job isn't going
	// anymore.
	if monCh != nil {
		defer close(monCh)
	}

	// In local mode we have to start a runner.
	if c.local {
		// Modify the job to target this runner and use the local data source.
		// The runner will have been started when we created the Project value and be
		// used for all local jobs.
		job.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: c.activeRunner.Id(),
				},
			},
		}
	}

	return c.queueAndStreamJob(ctx, job, ui, monCh)
}

// queueAndStreamJob will queue the job. If the client is configured to watch the job,
// it'll also stream the output to the configured UI.
func (c *Project) queueAndStreamJob(
	ctx context.Context,
	job *pb.Job,
	ui terminal.UI,
	monCh chan pb.Job_State,
) (*pb.Job_Result, error) {
	log := c.logger

	// When local, we set an expiration here in case we can't gracefully
	// cancel in the event of an error. This will ensure that the jobs don't
	// remain queued forever. This is only for local ops.
	expiration := ""
	if c.local {
		expiration = "30s"
	}

	// Queue the job
	log.Debug("queueing job", "operation", fmt.Sprintf("%T", job.Operation))
	queueResp, err := c.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job:       job,
		ExpiresIn: expiration,
	})
	if err != nil {
		return nil, err
	}
	log = log.With("job_id", queueResp.JobId)

	// Get the stream
	log.Debug("opening job stream")
	stream, err := c.client.GetJobStream(ctx, &pb.GetJobStreamRequest{
		JobId: queueResp.JobId,
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

	type stepData struct {
		terminal.Step

		out io.Writer
	}

	// Process events
	var (
		completed bool

		stateEventTimer *time.Timer
		tstatus         terminal.Status

		stdout, stderr io.Writer

		sg    terminal.StepGroup
		steps = map[int32]*stepData{}
	)

	if c.local {
		defer func() {
			// If we completed then do nothing, or if the context is still
			// active since this means that we're not cancelled.
			if completed || ctx.Err() == nil {
				return
			}

			ctx, cancel := finalcontext.Context(log)
			defer cancel()

			log.Warn("canceling job")
			_, err := c.client.CancelJob(ctx, &pb.CancelJobRequest{
				JobId: queueResp.JobId,
			})
			if err != nil {
				log.Warn("error canceling job", "err", err)
			} else {
				log.Info("job cancelled successfully")
			}
		}()
	}

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

		case *pb.GetJobStreamResponse_Terminal_:
			// Ignore this for local jobs since we're using our UI directly.
			if c.local {
				continue
			}

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
							return nil, err
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
					c.logger.Error("Unknown terminal event seen", "type", hclog.Fmt("%T", ev))
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

			if monCh != nil {
				select {
				case <-ctx.Done():
					break
				case monCh <- event.State.Current:
					// ok
				}
			}

		default:
			log.Warn("unknown stream event", "event", resp.Event)
		}
	}
}

const stateEventPause = 1500 * time.Millisecond
