package logstream

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// Provider provides a log stream tracker
type Provider interface {
	StartWriter(ctx context.Context, log hclog.Logger, state serverstate.Interface, job *serverstate.Job) (Tracker, error)
}

// Tracker collects and tracks loggable events
type Tracker interface {
	Flush(ctx context.Context)
	NewEvent(ctx context.Context, event *pb.RunnerJobStreamRequest_Terminal) error
}
