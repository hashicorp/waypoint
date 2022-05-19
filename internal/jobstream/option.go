package jobstream

import (
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Option specifies an option for streaming.
type Option func(s *stream) error

// Set a logger for the stream watcher.
func WithLogger(log hclog.Logger) Option {
	return func(s *stream) error {
		s.log = log
		return nil
	}
}

// Set the client for the stream watcher. This is required.
func WithClient(client pb.WaypointClient) Option {
	return func(s *stream) error {
		s.client = client
		return nil
	}
}

// Set the UI for all user-facing output to be sent to.
func WithUI(ui terminal.UI) Option {
	return func(s *stream) error {
		s.ui = ui
		return nil
	}
}

// WithCancelOnError causes the job being watched to be canceled (with
// CancelJob) if the streamer exits unsuccessfully. Defaults to false.
func WithCancelOnError(v bool) Option {
	return func(s *stream) error {
		s.cancelOnErr = v
		return nil
	}
}

// WithIgnoreTerminal ignores the terminal events from the job. Other UI
// output may happen such as queue delays but the actual terminal events
// are hidden.
func WithIgnoreTerminal(v bool) Option {
	return func(s *stream) error {
		s.ignoreTerminal = v
		return nil
	}
}

// WithStateCh sets a channel that is sent all the job state changes.
// This must not block. If the channel blocks, the entire stream watcher may
// be blocked.
func WithStateCh(v chan<- pb.Job_State) Option {
	return func(s *stream) error {
		s.stateCh = v
		return nil
	}
}
