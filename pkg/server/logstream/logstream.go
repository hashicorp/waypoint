// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logstream

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// Provider provides a log stream tracker
type Provider interface {

	// StartWriter starts a new log writer. Requires state to persist logs.
	StartWriter(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) (Writer, error)

	// StartReader starts a new log reader.
	StartReader(ctx context.Context, log hclog.Logger, job *serverstate.Job) (Reader, error)

	// ReadCompleted returns all the log entries for a job that has been completed,
	// by reading them out of persistent storage.
	ReadCompleted(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) ([]*pb.GetJobStreamResponse_Terminal_Event, error)
}

// Writer collects and tracks loggable events
type Writer interface {
	Flush(ctx context.Context)
	NewEvent(ctx context.Context, event *pb.RunnerJobStreamRequest_Terminal) error
}

// Reader reads terminal events for a given
type Reader interface {

	// ReadStream returns a batch of log entries for a job that's currenly active.
	// If zero exist and block is true, this will block waiting for
	// available entries. If block is false and no more log entries exist,
	// this will return nil.
	ReadStream(ctx context.Context, block bool) ([]*pb.GetJobStreamResponse_Terminal_Event, error)
}
