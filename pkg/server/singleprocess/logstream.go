// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

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

func (s *singleProcessLogStreamProvider) StartWriter(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) (logstream.Writer, error) {
	return &singleProcessLogStreamWriter{
		log: log,
		job: job,
	}, nil
}

func (s *singleProcessLogStreamProvider) StartReader(ctx context.Context, log hclog.Logger, job *serverstate.Job) (logstream.Reader, error) {
	var outputR *logbuffer.Reader

	// NOTE(izaak): Not having an output buffer on the job isn't an error condition.
	// Not sure when it would happen though.
	if job.OutputBuffer != nil {
		outputR = job.OutputBuffer.Reader(-1)
		go outputR.CloseContext(ctx)
	}

	return &singleProcessLogStreamReader{
		log:     log,
		outputR: outputR,
	}, nil
}

// ReadCompleted reads all the buffered logs for the specified job.
// NOTE: It doesn't verify that the job has, in fact, completed -
// If called on a non-completed job, it may still return logs.
// ALSO NOTE: The writer implementation does not persist job logs
// permanently, so this reader will not give logs for a long-completed job.
func (l *singleProcessLogStreamProvider) ReadCompleted(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) ([]*pb.GetJobStreamResponse_Terminal_Event, error) {
	r, err := l.StartReader(ctx, log, job)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start reader to get completed logs")
	}

	var events []*pb.GetJobStreamResponse_Terminal_Event

	// Read events for this job until
	for {
		eventsBatch, err := r.ReadStream(ctx, false)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read log batch for completed job")
		}
		if len(eventsBatch) == 0 {
			break
		}
		events = append(events, eventsBatch...)
	}
	return events, nil
}

// singleProcessLogStreamWriter implements logstream.Writer
type singleProcessLogStreamWriter struct {
	log hclog.Logger
	job *serverstate.Job
}

func (l *singleProcessLogStreamWriter) Flush(ctx context.Context) {
	// No work required
}

func (l *singleProcessLogStreamWriter) NewEvent(ctx context.Context, event *pb.RunnerJobStreamRequest_Terminal) error {
	log := l.log

	// NOTE(izaak): Someday we should probably refactor job.OutputBuffer
	// to be the abstraction around writing new events, rather than
	// this singleProcessLogStreamWriter.
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

// singleProcessLogStreamReader implements logstream.Reader
type singleProcessLogStreamReader struct {
	log     hclog.Logger
	outputR *logbuffer.Reader
}

func (l *singleProcessLogStreamReader) ReadStream(ctx context.Context, block bool) ([]*pb.GetJobStreamResponse_Terminal_Event, error) {

	// NOTE(izaak): Not having an output buffer on the job isn't an error condition.
	// Not sure when it would happen though.
	if l.outputR == nil {
		return []*pb.GetJobStreamResponse_Terminal_Event{}, nil
	}

	entries := l.outputR.Read(64, block)
	if entries == nil {
		return nil, nil
	}

	events := make([]*pb.GetJobStreamResponse_Terminal_Event, len(entries))
	for i, entry := range entries {
		events[i] = entry.(*pb.GetJobStreamResponse_Terminal_Event)
	}

	return events, nil
}
