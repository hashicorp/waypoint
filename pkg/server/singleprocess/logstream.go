package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
	"github.com/hashicorp/waypoint/pkg/server/logstream"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// singleProcessLogStreamProvider implements logstream.Provider
// Prefixed with "singleprocess" to indicate that it makes
// singleprocess server assumptions and is unsafe to use
// outside of that context.
type singleProcessLogStreamProvider struct{}

func (s *singleProcessLogStreamProvider) StartWriter(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) (logstream.Tracker, error) {
	return &singleProcessLogStreamTracker{
		log: log,
		job: job,
	}, nil
}

// singleProcessLogStreamTracker implements logstream.Tracker
type singleProcessLogStreamTracker struct {
	log hclog.Logger
	job *serverstate.Job
}

func (l *singleProcessLogStreamTracker) Flush(ctx context.Context) {
	// No work required
	return
}

func (l *singleProcessLogStreamTracker) NewEvent(ctx context.Context, event *pb.RunnerJobStreamRequest_Terminal) error {
	log := l.log

	// NOTE(izaak): Someday we should probably refactor job.OutputBuffer
	// to be the abstraction around writing new events, rather than
	// this singleProcessLogStreamTracker.
	if l.job.OutputBuffer == nil {
		log.Warn("got terminal event but internal output buffer is nil, dropping lines")
		return nil
	}

	// Write the entries to the output buffer
	entries := make([]logbuffer.Entry, len(event.Terminal.Events))
	for i, ev := range event.Terminal.Events {
		entries[i] = ev
	}

	l.job.OutputBuffer.Write(entries...)
	return nil
}
