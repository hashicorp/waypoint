package terminal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
)

// CaptureUI implements terminal.UI and captures the output, calling the
// callback function with the output given. This allows the creator of
// this UI to send the output elsewhere. This is used for example by Waypoint
// runners to stream output to the server.
//
// The caller should call Close when done with the UI to clean up goroutines.
type CaptureUI struct {
	logger hclog.Logger
	real   UI

	// callback is the callback to call when there is captured output.
	// This is guaranteed to not be called concurrently.
	callback CaptureCallback

	// closed when true means this UI is closed. Status updates and
	// outputs no longer work in this case. Setting this can only be done
	// while holding sendLock.
	closed bool

	// sendLock protects sendLines. This must be protected because we may
	// get races through usage of OutputWriters. It is also possible to race
	// just on Output although the docs of terminal.UI state this is unsafe.
	sendLock sync.Mutex

	// pw is the pipe writer for the output writer. pwDoneCh is closed when
	// the output writing goroutine exits. cancelFunc can be called to stop
	// the output writing goroutine.
	pw         *io.PipeWriter
	pwDoneCh   <-chan struct{}
	cancelFunc func()
}

// CaptureCallback is the callback type for CaptureUI.
type CaptureCallback func([]*CaptureLine) error

// CaptureLine is a single line of output that was captured.
type CaptureLine struct {
	Line      string
	Timestamp time.Time
}

// NewCaptureUI creates a CaptureUI for the given job stream. The resulting UI will
// stream any output via the gRPC stream.
func NewCaptureUI(log hclog.Logger, callback CaptureCallback) *CaptureUI {
	// Build a context to cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Build our CaptureUI
	pr, pw := io.Pipe()
	pwDoneCh := make(chan struct{})
	result := &CaptureUI{
		logger:     log,
		real:       &BasicUI{},
		callback:   callback,
		pw:         pw,
		pwDoneCh:   pwDoneCh,
		cancelFunc: cancel,
	}

	// Start a goroutine that handles the outputwriters. This will be cleaned
	// up automatically when jobUi.Close is called.
	go result.streamOutputWriters(ctx, pr, pwDoneCh)

	return result
}

func (ui *CaptureUI) Close() error {
	// Mark we're closed.
	ui.sendLock.Lock()
	if ui.closed {
		ui.sendLock.Unlock()
		return nil
	}
	ui.closed = true
	ui.sendLock.Unlock()

	// Close the writer end
	ui.pw.Close()

	// Cancel the stream
	ui.cancelFunc()

	// Wait for it to cancel
	<-ui.pwDoneCh

	return nil
}

func (ui *CaptureUI) Output(msg string, raw ...interface{}) {
	// Write to our buffer
	var buf bytes.Buffer
	ui.real.Output(msg, append(raw,
		WithWriter(&buf),
	)...)

	ui.parseBuf(&buf)
}

func (ui *CaptureUI) NamedValues(rows []NamedValue, opts ...Option) {
	// Write to our buffer
	var buf bytes.Buffer
	ui.real.NamedValues(rows, append(opts,
		WithWriter(&buf),
	)...)

	ui.parseBuf(&buf)
}

func (ui *CaptureUI) parseBuf(buf *bytes.Buffer) {
	ts := time.Now()

	// Scan and construct lines
	var lines []*CaptureLine
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		lines = append(lines, &CaptureLine{
			Line:      scanner.Text(),
			Timestamp: ts,
		})
	}
	if err := scanner.Err(); err != nil {
		// This really shouldn't happen but just log it.
		ui.logger.Warn("error scanning output lines", "err", err)
	}

	// Send whatever we have
	ui.sendLines(lines)
}

func (ui *CaptureUI) OutputWriters() (io.Writer, io.Writer, error) {
	return ui.pw, ui.pw, nil
}

func (ui *CaptureUI) Status() Status {
	return &captureStatus{UI: ui}
}

func (ui *CaptureUI) streamOutputWriters(ctx context.Context, r *io.PipeReader, doneCh chan<- struct{}) {
	// Signal when we're done
	defer close(doneCh)

	// Close the reader end when we're done
	defer r.Close()

	// We start a goroutine here that just streams lines to us.
	linesCh := make(chan *CaptureLine, 1)
	go func() {
		defer close(linesCh)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			linesCh <- &CaptureLine{
				Line:      scanner.Text(),
				Timestamp: time.Now(),
			}
		}
	}()

	// Accumulate lines
	lines := make([]*CaptureLine, 0, 64)
	for {
		send := false
		select {
		case line := <-linesCh:
			lines = append(lines, line)
			send = len(lines) == cap(lines)

		case <-time.After(100 * time.Millisecond):
			send = len(lines) > 0

		case <-ctx.Done():
			// Drain our lines
			for line := range linesCh {
				lines = append(lines, line)
			}

			send = len(lines) > 0
		}

		// If we have lines to send, send them. Otherwise wait.
		if send {
			ui.sendLines(lines)
			lines = lines[:0]
		}

		// If we're done, then exit
		if ctx.Err() != nil {
			return
		}
	}
}

func (ui *CaptureUI) sendLines(lines []*CaptureLine) {
	if ui.logger.IsTrace() {
		for _, line := range lines {
			ui.logger.Trace("captured output", "line", line.Line)
		}
	}

	if err := ui.callback(lines); err != nil {
		ui.logger.Warn("error sending output line", "err", err)
	}
}

// captureStatus implements terminal.Status to handle status updates over job streams.
//
// This is extremely basic right now and just sends along raw lines via
// the UI. We should make this better.
type captureStatus struct {
	UI *CaptureUI

	lastMsg string
}

func (s *captureStatus) Update(msg string) {
	// Update is often called in a loop since the behavioral expectation
	// is that it clears the line. We don't support that yet so we just
	// make sure we're seeing change.
	if msg == s.lastMsg {
		return
	}
	s.lastMsg = msg

	// Just output to the UI as normal for now
	s.UI.Output(msg)
}

func (s *captureStatus) Close() error {
	return nil
}

func (s *captureStatus) Step(status, msg string) {
	s.UI.Output(fmt.Sprintf("%s: %s", status, msg))
}

var (
	_ UI     = (*CaptureUI)(nil)
	_ Status = (*captureStatus)(nil)
)
