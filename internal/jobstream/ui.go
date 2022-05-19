package jobstream

import (
	"io"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// UI converts a stream of terminal events from a job stream into
// stateful function calls to a terminal.UI implementation.
type UI struct {
	// The underlying UI that will be written to.
	UI terminal.UI

	// Log for unknown events and other diagnostics.
	Log hclog.Logger

	// Internal state
	status         terminal.Status
	sg             terminal.StepGroup
	stdout, stderr io.Writer
	steps          map[int32]*stepData
}

type stepData struct {
	terminal.Step
	out io.Writer
}

// Write processes the given events and converts them to calls to UI.
// The order of the events matter because events are stateful: some events
// start a table, some end it, etc.
//
// TODO(mitchellh): This was extracted directly from internal/client. It
// had no tests there and it still has no tests here. We should test this
// eventually and it should be pretty easy to do so!
func (s *UI) Write(events []*pb.GetJobStreamResponse_Terminal_Event) error {
	ui := s.UI
	log := s.Log
	if log == nil {
		log = hclog.L()
	}

	for _, ev := range events {
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
			if s.status == nil {
				s.status = ui.Status()
				defer s.status.Close()
			}

			if ev.Status.Msg == "" && !ev.Status.Step {
				s.status.Close()
			} else if ev.Status.Step {
				s.status.Step(ev.Status.Status, ev.Status.Msg)
			} else {
				s.status.Update(ev.Status.Msg)
			}
		case *pb.GetJobStreamResponse_Terminal_Event_Raw_:
			if s.stdout == nil {
				var err error
				s.stdout, s.stderr, err = ui.OutputWriters()
				if err != nil {
					return err
				}
			}

			if ev.Raw.Stderr {
				s.stderr.Write(ev.Raw.Data)
			} else {
				s.stdout.Write(ev.Raw.Data)
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
			if !ev.StepGroup.Close {
				s.sg = ui.StepGroup()
			}
		case *pb.GetJobStreamResponse_Terminal_Event_Step_:
			if s.sg == nil {
				continue
			}

			if s.steps == nil {
				s.steps = map[int32]*stepData{}
			}

			step, ok := s.steps[ev.Step.Id]
			if !ok {
				step = &stepData{
					Step: s.sg.Add(ev.Step.Msg),
				}
				s.steps[ev.Step.Id] = step
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
			// Unknown, ignore.
			log.Warn("unknown UI event", "event", ev)
		}
	}

	return nil
}
