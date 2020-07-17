package terminal

import (
	"context"
	"io"
	"sync"
)

const (
	TermRows    = 10
	TermColumns = 100
)

// fancyStepGroup implements StepGroup with live updating and a display
// "window" for live terminal output (when using TermOutput).
type fancyStepGroup struct {
	ctx    context.Context
	cancel func()

	display *Display

	wg sync.WaitGroup
}

// Start a step in the output
func (f *fancyStepGroup) Add(str string, args ...interface{}) Step {
	f.wg.Add(1)

	ent := f.display.NewStatus(0)

	ent.StartSpinner()
	ent.Update(str, args...)

	return &fancyStep{
		sg:  f,
		ent: ent,
	}
}

func (f *fancyStepGroup) Wait() {
	f.wg.Wait()
	f.cancel()

	f.display.Close()
}

type fancyStep struct {
	sg  *fancyStepGroup
	ent *DisplayEntry

	done bool

	term *Term
}

func (f *fancyStep) TermOutput() io.Writer {
	if f.term == nil {
		t, err := NewTerm(f.sg.ctx, f.ent, TermRows, TermColumns)
		if err != nil {
			panic(err)
		}

		f.term = t
	}

	return f.term
}

func (f *fancyStep) Update(str string, args ...interface{}) {
	f.ent.Update(str, args...)
}

func (f *fancyStep) Status(status string) {
	f.ent.SetStatus(status)
}

func (f *fancyStep) Done() {
	if f.done {
		return
	}

	f.ent.StopSpinner()
	f.Status(StatusOK)
	f.done = true
	f.sg.wg.Done()
}

func (f *fancyStep) Abort() {
	if f.done {
		return
	}

	f.ent.StopSpinner()
	f.Status(StatusError)

	f.done = true
	f.sg.wg.Done()
}
