package runner

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// jobUI implements terminal.UI and streams the output to the job stream.
//
// The caller should call Close when done with the UI to clean up goroutines.
type jobUI struct {
	logger hclog.Logger
	client pb.Waypoint_RunnerJobStreamClient
	real   terminal.UI

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

// newJobUI creates a jobUI for the given job stream. The resulting UI will
// stream any output via the gRPC stream.
func newJobUI(log hclog.Logger, client pb.Waypoint_RunnerJobStreamClient) *jobUI {
	// Build a context to cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Build our jobUI
	pr, pw := io.Pipe()
	pwDoneCh := make(chan struct{})
	result := &jobUI{
		logger:     log,
		client:     client,
		real:       &terminal.BasicUI{},
		pw:         pw,
		pwDoneCh:   pwDoneCh,
		cancelFunc: cancel,
	}

	// Start a goroutine that handles the outputwriters. This will be cleaned
	// up automatically when jobUi.Close is called.
	go result.streamOutputWriters(ctx, pr, pwDoneCh)

	return result
}

func (ui *jobUI) Close() error {
	// Mark we're closed.
	ui.sendLock.Lock()
	ui.closed = true
	ui.sendLock.Unlock()

	// Cancel the stream
	ui.cancelFunc()

	// Wait for it to cancel
	<-ui.pwDoneCh

	// Close the writer end
	ui.pw.Close()

	return nil
}

func (ui *jobUI) Output(msg string, raw ...interface{}) {
	// Our timestamp for this is now
	ts := ptypes.TimestampNow()

	// Write to our buffer
	var buf bytes.Buffer
	ui.real.Output(msg, append(raw,
		terminal.WithWriter(&buf),
	)...)

	// Scan and construct lines
	var lines []*pb.GetJobStreamResponse_Terminal_Line
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		lines = append(lines, &pb.GetJobStreamResponse_Terminal_Line{
			Raw:       scanner.Text(),
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

func (ui *jobUI) OutputWriters() (io.Writer, io.Writer, error) {
	return ui.pw, ui.pw, nil
}

func (ui *jobUI) Status() terminal.Status {
	return &jobStatus{UI: ui}
}

func (ui *jobUI) streamOutputWriters(ctx context.Context, r *io.PipeReader, doneCh chan<- struct{}) {
	// Signal when we're done
	defer close(doneCh)

	// Close the reader end when we're done
	defer r.Close()

	// We start a goroutine here that just streams lines to us.
	linesCh := make(chan *pb.GetJobStreamResponse_Terminal_Line, 1)
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			linesCh <- &pb.GetJobStreamResponse_Terminal_Line{
				Raw:       scanner.Text(),
				Line:      scanner.Text(),
				Timestamp: ptypes.TimestampNow(),
			}
		}
	}()

	// Accumulate lines
	lines := make([]*pb.GetJobStreamResponse_Terminal_Line, 0, 64)
	for {
		send := false
		select {
		case line := <-linesCh:
			lines = append(lines, line)
			send = len(lines) == cap(lines)

		case <-time.After(1 * time.Second):
			send = len(lines) > 0

		case <-ctx.Done():
			return
		}

		// If we're not supposed to send lines, wait for more
		if !send {
			continue
		}

		// Send the lines
		ui.sendLines(lines)
	}
}

func (ui *jobUI) sendLines(lines []*pb.GetJobStreamResponse_Terminal_Line) {
	ui.sendLock.Lock()
	defer ui.sendLock.Unlock()

	// If we're closed, we can't send any output because ui.client may be
	// used for other things and its not thread-safe to write concurrently.
	if ui.closed {
		ui.logger.Warn("output after close, dropping")
		return
	}

	if ui.logger.IsTrace() {
		for _, line := range lines {
			ui.logger.Trace("job output", "line", line.Raw)
		}
	}

	if err := ui.client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Lines: lines,
			},
		},
	}); err != nil {
		ui.logger.Warn("error sending output line", "err", err)
	}
}

// jobStatus implements terminal.Status to handle status updates over job streams.
//
// This is extremely basic right now and just sends along raw lines via
// the UI. We should make this better.
type jobStatus struct {
	UI *jobUI

	lastMsg string
}

func (s *jobStatus) Update(msg string) {
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

func (s *jobStatus) Close() error {
	return nil
}

var (
	_ terminal.UI     = (*jobUI)(nil)
	_ terminal.Status = (*jobStatus)(nil)
)
