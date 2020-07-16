package terminal

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/creack/pty"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/hashicorp/waypoint/sdk/proto"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// UIPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the terminal.UI interface.
type UIPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    terminal.UI       // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *UIPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterTerminalUIServiceServer(s, &uiServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *UIPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := pb.NewTerminalUIServiceClient(c)
	evstream, err := client.Events(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &uiBridge{
		ctx:    ctx,
		cancel: cancel,
		evc:    evstream,
	}, nil
}

// uiServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type uiServer struct {
	Impl    terminal.UI
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (s *uiServer) Output(
	ctx context.Context,
	req *pb.TerminalUI_OutputRequest,
) (*empty.Empty, error) {
	for _, line := range req.Lines {
		s.Impl.Output(line)
	}

	return &empty.Empty{}, nil
}

func (s *uiServer) Events(stream pb.TerminalUIService_EventsServer) error {
	type stepData struct {
		terminal.Step
		out io.Writer
	}

	var (
		status terminal.Status
		stdout io.Writer
		stderr io.Writer

		sg    terminal.StepGroup
		steps = map[int32]*stepData{}
	)

	for {
		ev, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		switch ev := ev.Event.(type) {
		case *pb.TerminalUI_Event_Line_:
			s.Impl.Output(ev.Line.Msg, terminal.WithStyle(ev.Line.Style))
		case *pb.TerminalUI_Event_NamedValues_:
			var values []terminal.NamedValue

			for _, nv := range ev.NamedValues.Values {
				values = append(values, terminal.NamedValue{
					Name:  nv.Name,
					Value: nv.Value,
				})
			}

			s.Impl.NamedValues(values)
		case *pb.TerminalUI_Event_Status_:
			if ev.Status.Msg == "" && !ev.Status.Step {
				if status != nil {
					status.Close()
				}
			} else {
				if status == nil {
					status = s.Impl.Status()
					defer status.Close()
				}

				if ev.Status.Step {
					status.Step(ev.Status.Status, ev.Status.Msg)
				} else {
					status.Update(ev.Status.Msg)
				}
			}
		case *pb.TerminalUI_Event_Raw_:
			if stdout == nil {
				stdout, stderr, err = s.Impl.OutputWriters()
				if err != nil {
					return err
				}
			}

			if ev.Raw.Stderr {
				stderr.Write(ev.Raw.Data)
			} else {
				stdout.Write(ev.Raw.Data)
			}
		case *pb.TerminalUI_Event_Table_:
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

			s.Impl.Table(tbl)
		case *pb.TerminalUI_Event_StepGroup_:
			if sg != nil {
				sg.Wait()
			}

			if !ev.StepGroup.Close {
				sg = s.Impl.StepGroup()
			}
		case *pb.TerminalUI_Event_Step_:
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
				step.Status(ev.Step.Status)
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
			s.Logger.Error("Unknown terminal event seen", "type", hclog.Fmt("%T", ev))
		}
	}
}

type uiBridge struct {
	ctx    context.Context
	cancel func()
	mu     sync.Mutex
	evc    pb.TerminalUIService_EventsClient

	stdSetup       sync.Once
	stdout, stderr io.Writer
}

func (u *uiBridge) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, err := u.evc.CloseAndRecv()
	u.evc = nil
	u.cancel()

	return err
}

// Output outputs a message directly to the terminal. The remaining
// arguments should be interpolations for the format string. After the
// interpolations you may add Options.
func (u *uiBridge) Output(msg string, raw ...interface{}) {
	msg, style, _ := terminal.Interpret(msg, raw...)

	ev := &pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Line_{
			Line: &pb.TerminalUI_Event_Line{
				Msg:   msg,
				Style: style,
			},
		},
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.evc == nil {
		return
	}

	u.evc.Send(ev)
}

// Output data as a table of data. Each entry is a row which will be output
// with the columns lined up nicely.
func (u *uiBridge) NamedValues(tvalues []terminal.NamedValue, _ ...terminal.Option) {
	var values []*pb.TerminalUI_Event_NamedValue

	for _, nv := range tvalues {
		values = append(values, &pb.TerminalUI_Event_NamedValue{
			Name:  nv.Name,
			Value: fmt.Sprintf("%s", nv.Value),
		})
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.evc == nil {
		return
	}

	u.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_NamedValues_{
			NamedValues: &pb.TerminalUI_Event_NamedValues{
				Values: values,
			},
		},
	})
}

// OutputWriters returns stdout and stderr writers. These are usually
// but not always TTYs. This is useful for subprocesses, network requests,
// etc. Note that writing to these is not thread-safe by default so
// you must take care that there is only ever one writer.
func (u *uiBridge) OutputWriters() (stdout io.Writer, stderr io.Writer, err error) {
	u.stdSetup.Do(func() {
		dr, dw, err := pty.Open()
		if err != nil {
			panic(err)
		}

		err = pty.Setsize(dw, &pty.Winsize{
			Rows: uint16(terminal.TermRows),
			Cols: uint16(terminal.TermColumns),
		})

		if err != nil {
			panic(err)
		}

		go u.sendData(dr, false)

		er, ew, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		go u.sendData(er, true)

		go func() {
			<-u.ctx.Done()
			dr.Close()
			dw.Close()
			er.Close()
			ew.Close()
		}()

		u.stdout = dw
		u.stderr = ew
	})

	return u.stdout, u.stderr, nil
}

func (u *uiBridge) sendData(r io.ReadCloser, stderr bool) {
	defer r.Close()

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		ev := &pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Raw_{
				Raw: &pb.TerminalUI_Event_Raw{
					Data:   data,
					Stderr: stderr,
				},
			},
		}

		u.mu.Lock()
		if u.evc == nil {
			u.mu.Unlock()
			return
		}

		u.evc.Send(ev)
		u.mu.Unlock()
	}
}

func (u *uiBridge) Table(tbl *terminal.Table, opts ...terminal.Option) {
	var (
		ptbl *pb.TerminalUI_Event_Table
		rows []*pb.TerminalUI_Event_TableRow
	)

	ptbl.Headers = tbl.Headers

	for _, row := range tbl.Rows {
		var entries []*pb.TerminalUI_Event_TableEntry

		for _, ent := range row {
			entries = append(entries, &pb.TerminalUI_Event_TableEntry{
				Value: ent.Value,
				Color: ent.Color,
			})
		}

		rows = append(rows, &pb.TerminalUI_Event_TableRow{
			Entries: entries,
		})
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.evc == nil {
		return
	}

	u.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Table_{
			Table: ptbl,
		},
	})
}

type uiBridgeSGStep struct {
	sg   *uiBridgeSG
	id   int32
	done bool

	stdSetup sync.Once
	stdout   io.Writer
}

func (u *uiBridgeSGStep) TermOutput() io.Writer {
	u.stdSetup.Do(func() {
		dr, dw, err := pty.Open()
		if err != nil {
			panic(err)
		}

		go u.sendData(dr, false)

		go func() {
			<-u.sg.ctx.Done()
			dr.Close()
			dw.Close()
		}()

		u.stdout = dw
	})

	return u.stdout
}

func (u *uiBridgeSGStep) sendData(r io.ReadCloser, stderr bool) {
	defer r.Close()

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		ev := &pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Step_{
				Step: &pb.TerminalUI_Event_Step{
					Id:     u.id,
					Output: data,
				},
			},
		}

		u.sg.ui.mu.Lock()
		if u.sg.ui.evc == nil {
			u.sg.ui.mu.Unlock()
			return
		}

		u.sg.ui.evc.Send(ev)
		u.sg.ui.mu.Unlock()
	}
}

func (u *uiBridgeSGStep) Update(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)

	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.sg.ui.evc != nil {
		u.sg.ui.evc.Send(&pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Step_{
				Step: &pb.TerminalUI_Event_Step{
					Id:  u.id,
					Msg: msg,
				},
			},
		})
	}
}

func (u *uiBridgeSGStep) Status(status string) {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.sg.ui.evc != nil {
		u.sg.ui.evc.Send(&pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Step_{
				Step: &pb.TerminalUI_Event_Step{
					Id:     u.id,
					Status: status,
				},
			},
		})
	}
}

func (u *uiBridgeSGStep) Done() {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.done {
		return
	}

	u.done = true

	if u.sg.ui.evc != nil {
		u.sg.ui.evc.Send(&pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Step_{
				Step: &pb.TerminalUI_Event_Step{
					Id:    u.id,
					Close: true,
				},
			},
		})
	}

	u.sg.wg.Done()
}

func (u *uiBridgeSGStep) Abort() {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.done {
		return
	}

	u.done = true

	if u.sg.ui.evc != nil {
		u.sg.ui.evc.Send(&pb.TerminalUI_Event{
			Event: &pb.TerminalUI_Event_Step_{
				Step: &pb.TerminalUI_Event_Step{
					Id:     u.id,
					Close:  true,
					Status: terminal.ErrorStyle,
				},
			},
		})
	}

	u.sg.wg.Done()
}

type uiBridgeSG struct {
	ctx    context.Context
	cancel func()

	ui *uiBridge
	wg sync.WaitGroup

	steps []*uiBridgeSGStep
}

// Start a step in the output
func (u *uiBridgeSG) Add(str string, args ...interface{}) terminal.Step {
	msg := fmt.Sprintf(str, args...)

	u.ui.mu.Lock()
	defer u.ui.mu.Unlock()

	u.wg.Add(1)

	step := &uiBridgeSGStep{
		sg: u,
		id: int32(len(u.steps)),
	}

	u.steps = append(u.steps, step)

	u.ui.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Step_{
			Step: &pb.TerminalUI_Event_Step{
				Id:  step.id,
				Msg: msg,
			},
		},
	})

	return step
}

func (u *uiBridgeSG) Wait() {
	u.wg.Wait()
	u.cancel()

	u.ui.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_StepGroup_{
			StepGroup: &pb.TerminalUI_Event_StepGroup{
				Close: true,
			},
		},
	})

}

func (u *uiBridge) StepGroup() terminal.StepGroup {
	ctx, cancel := context.WithCancel(u.ctx)

	sg := &uiBridgeSG{
		ui:     u,
		ctx:    ctx,
		cancel: cancel,
	}

	u.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_StepGroup_{
			StepGroup: &pb.TerminalUI_Event_StepGroup{},
		},
	})

	return sg
}

// Status returns a live-updating status that can be used for single-line
// status updates that typically have a spinner or some similar style.
func (u *uiBridge) Status() terminal.Status {
	return &uiBridgeStatus{u}
}

type uiBridgeStatus struct {
	b *uiBridge
}

// Update writes a new status. This should be a single line.
func (u *uiBridgeStatus) Update(msg string) {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return
	}

	u.b.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Status_{
			Status: &pb.TerminalUI_Event_Status{
				Msg: msg,
			},
		},
	})
}

// Indicate that a step has finished, confering an ok, error, or warn upon
// it's finishing state. If the status is not StatusOK, StatusError, or StatusWarn
// then the status text is written directly to the output, allowing for custom
// statuses.
func (u *uiBridgeStatus) Step(status string, msg string) {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return
	}

	u.b.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Status_{
			Status: &pb.TerminalUI_Event_Status{
				Status: status,
				Msg:    msg,
				Step:   true,
			},
		},
	})
}

// Close should be called when the live updating is complete. The
// status will be cleared from the line.
func (u *uiBridgeStatus) Close() error {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return nil
	}

	u.b.evc.Send(&pb.TerminalUI_Event{
		Event: &pb.TerminalUI_Event_Status_{
			Status: &pb.TerminalUI_Event_Status{},
		},
	})

	return nil
}

var (
	_ plugin.Plugin              = (*UIPlugin)(nil)
	_ plugin.GRPCPlugin          = (*UIPlugin)(nil)
	_ pb.TerminalUIServiceServer = (*uiServer)(nil)
	_ terminal.UI                = (*uiBridge)(nil)
	_ terminal.Status            = (*uiBridgeStatus)(nil)
)
