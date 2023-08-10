// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// runnerUI Implements terminal.UI and is created by a runner and passed into
// it's operations. The functions send events back to the server to be saved
// and sent to job clients rather than displaying the events directly.
type runnerUI struct {
	ctx    context.Context
	cancel func()
	mu     *sync.Mutex
	evc    pb.Waypoint_RunnerJobStreamClient

	// stepIdx keeps track of the current step "ID" used when talking to
	// the server. Each new stepgroup step gets a new monotonically increasing
	// ID. You must never reuse an old ID.
	stepIdx int32

	stdSetup       sync.Once
	stdout, stderr io.Writer
}

func (u *runnerUI) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.evc = nil
	u.cancel()

	return nil
}

func (u *runnerUI) Input(input *terminal.Input) (string, error) {
	return "", terminal.ErrNonInteractive
}

func (u *runnerUI) Interactive() bool {
	return false
}

// Output outputs a message directly to the terminal. The remaining
// arguments should be interpolations for the format string. After the
// interpolations you may add Options.
func (u *runnerUI) Output(msg string, raw ...interface{}) {
	msg, style, _ := terminal.Interpret(msg, raw...)

	// Extreme java looking code alert!
	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg:   msg,
								Style: style,
							},
						},
					},
				},
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
func (u *runnerUI) NamedValues(tvalues []terminal.NamedValue, _ ...terminal.Option) {
	var values []*pb.GetJobStreamResponse_Terminal_Event_NamedValue

	for _, nv := range tvalues {
		values = append(values, &pb.GetJobStreamResponse_Terminal_Event_NamedValue{
			Name:  nv.Name,
			Value: fmt.Sprintf("%s", nv.Value),
		})
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.evc == nil {
		return
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_NamedValues_{
							NamedValues: &pb.GetJobStreamResponse_Terminal_Event_NamedValues{
								Values: values,
							},
						},
					},
				},
			},
		},
	}

	u.evc.Send(ev)
}

// OutputWriters returns stdout and stderr writers. These are usually
// but not always TTYs. This is useful for subprocesses, network requests,
// etc. Note that writing to these is not thread-safe by default so
// you must take care that there is only ever one writer.
func (u *runnerUI) OutputWriters() (stdout io.Writer, stderr io.Writer, err error) {
	u.stdSetup.Do(func() {
		dr, dw, err := os.Pipe()
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

func (u *runnerUI) sendData(r io.ReadCloser, stderr bool) {
	defer r.Close()

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Raw_{
								Raw: &pb.GetJobStreamResponse_Terminal_Event_Raw{
									Data:   data,
									Stderr: stderr,
								},
							},
						},
					},
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

func (u *runnerUI) Table(tbl *terminal.Table, opts ...terminal.Option) {
	var (
		ptbl *pb.GetJobStreamResponse_Terminal_Event_Table
		rows []*pb.GetJobStreamResponse_Terminal_Event_TableRow
	)

	ptbl.Headers = tbl.Headers

	for _, row := range tbl.Rows {
		var entries []*pb.GetJobStreamResponse_Terminal_Event_TableEntry

		for _, ent := range row {
			entries = append(entries, &pb.GetJobStreamResponse_Terminal_Event_TableEntry{
				Value: ent.Value,
				Color: ent.Color,
			})
		}

		rows = append(rows, &pb.GetJobStreamResponse_Terminal_Event_TableRow{
			Entries: entries,
		})
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.evc == nil {
		return
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Table_{
							Table: ptbl,
						},
					},
				},
			},
		},
	}

	u.evc.Send(ev)
}

// Status returns a live-updating status that can be used for single-line
// status updates that typically have a spinner or some similar style.
func (u *runnerUI) Status() terminal.Status {
	return &runnerUIStatus{u}
}

type runnerUIStatus struct {
	b *runnerUI
}

// Update writes a new status. This should be a single line.
func (u *runnerUIStatus) Update(msg string) {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Status_{
							Status: &pb.GetJobStreamResponse_Terminal_Event_Status{
								Msg: msg,
							},
						},
					},
				},
			},
		},
	}

	u.b.evc.Send(ev)
}

// Indicate that a step has finished, confering an ok, error, or warn upon
// it's finishing state. If the status is not StatusOK, StatusError, or StatusWarn
// then the status text is written directly to the output, allowing for custom
// statuses.
func (u *runnerUIStatus) Step(status string, msg string) {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Status_{
							Status: &pb.GetJobStreamResponse_Terminal_Event_Status{
								Status: status,
								Msg:    msg,
								Step:   true,
							},
						},
					},
				},
			},
		},
	}

	u.b.evc.Send(ev)
}

// Close should be called when the live updating is complete. The
// status will be cleared from the line.
func (u *runnerUIStatus) Close() error {
	u.b.mu.Lock()
	defer u.b.mu.Unlock()

	if u.b.evc == nil {
		return nil
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Status_{
							Status: &pb.GetJobStreamResponse_Terminal_Event_Status{},
						},
					},
				},
			},
		},
	}

	u.b.evc.Send(ev)

	return nil
}

type runnerUISGStep struct {
	sg   *runnerUISG
	id   int32
	done bool

	stdSetup sync.Once
	stdout   io.Writer
}

func (u *runnerUISGStep) TermOutput() io.Writer {
	u.stdSetup.Do(func() {
		dr, dw, err := os.Pipe()
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

func (u *runnerUISGStep) sendData(r io.ReadCloser, stderr bool) {
	defer r.Close()

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
								Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
									Id:     u.id,
									Output: data,
								},
							},
						},
					},
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

func (u *runnerUISGStep) Update(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)

	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.sg.ui.evc != nil {
		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
								Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
									Id:  u.id,
									Msg: msg,
								},
							},
						},
					},
				},
			},
		}

		u.sg.ui.evc.Send(ev)
	}
}

func (u *runnerUISGStep) Status(status string) {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.sg.ui.evc != nil {
		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
								Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
									Id:     u.id,
									Status: status,
								},
							},
						},
					},
				},
			},
		}

		u.sg.ui.evc.Send(ev)
	}
}

func (u *runnerUISGStep) Done() {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.done {
		return
	}
	u.done = true

	if u.sg.ui.evc != nil {
		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
								Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
									Id:    u.id,
									Close: true,
								},
							},
						},
					},
				},
			},
		}

		u.sg.ui.evc.Send(ev)
	}

	u.sg.wg.Done()
}

func (u *runnerUISGStep) Abort() {
	u.sg.ui.mu.Lock()
	defer u.sg.ui.mu.Unlock()

	if u.done {
		return
	}
	u.done = true

	if u.sg.ui.evc != nil {
		ev := &pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Terminal{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events: []*pb.GetJobStreamResponse_Terminal_Event{
						{
							Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
								Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
									Id:     u.id,
									Close:  true,
									Status: terminal.StatusAbort,
								},
							},
						},
					},
				},
			},
		}

		u.sg.ui.evc.Send(ev)
	}

	u.sg.wg.Done()
}

type runnerUISG struct {
	ctx    context.Context
	cancel func()

	ui *runnerUI
	wg sync.WaitGroup

	steps []*runnerUISGStep
}

// Start a step in the output
func (u *runnerUISG) Add(str string, args ...interface{}) terminal.Step {
	msg := fmt.Sprintf(str, args...)

	u.ui.mu.Lock()
	defer u.ui.mu.Unlock()

	u.wg.Add(1)

	stepIdx := u.ui.stepIdx
	u.ui.stepIdx++

	step := &runnerUISGStep{
		sg: u,
		id: stepIdx,
	}

	u.steps = append(u.steps, step)

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
							Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
								Id:  step.id,
								Msg: msg,
							},
						},
					},
				},
			},
		},
	}

	u.ui.evc.Send(ev)

	return step
}

func (u *runnerUISG) Wait() {
	u.wg.Wait()
	u.cancel()

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_StepGroup_{
							StepGroup: &pb.GetJobStreamResponse_Terminal_Event_StepGroup{
								Close: true,
							},
						},
					},
				},
			},
		},
	}

	u.ui.evc.Send(ev)
}

func (u *runnerUI) StepGroup() terminal.StepGroup {
	u.mu.Lock()
	defer u.mu.Unlock()

	ctx, cancel := context.WithCancel(u.ctx)

	sg := &runnerUISG{
		ui:     u,
		ctx:    ctx,
		cancel: cancel,
	}

	ev := &pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_StepGroup_{
							StepGroup: &pb.GetJobStreamResponse_Terminal_Event_StepGroup{},
						},
					},
				},
			},
		},
	}

	u.evc.Send(ev)

	return sg
}

var (
	_ terminal.UI     = (*runnerUI)(nil)
	_ terminal.Status = (*runnerUIStatus)(nil)
)
