package logviewer

import (
	"context"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
)

// Viewer implements component.LogViewer over the server-side log stream endpoint.
//
// TODO(mitchellh): we should support some form of reconnection in the event of
// network errors.
type Viewer struct {
	// Stream is the log stream client to use.
	Stream pb.Devflow_GetLogStreamClient
}

// NextLogBatch implements component.LogViewer
func (v *Viewer) NextLogBatch(ctx context.Context) ([]component.LogEvent, error) {
	// Get the next entry. Note that we specifically do NOT buffer here because
	// we want to provide the proper amount of backpressure and we expect our
	// downstream caller to be calling these as quickly as possible.
	entry, err := v.Stream.Recv()
	if err != nil {
		return nil, err
	}

	events := make([]component.LogEvent, len(entry.Lines))
	for i, entry := range entry.Lines {
		events[i] = component.LogEvent{
			Message: entry.Line,
		}
	}

	return events, nil
}

var _ component.LogViewer = (*Viewer)(nil)
