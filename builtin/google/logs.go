package google

import (
	"context"
	"time"

	"github.com/mitchellh/devflow/sdk/component"
)

type LogViewer struct{}

func (v *LogViewer) NextLogBatch(ctx context.Context) ([]component.LogEvent, error) {
	time.Sleep(3 * time.Second)

	return []component.LogEvent{
		component.LogEvent{
			Timestamp: time.Now(),
			Message:   "hello",
		},
	}, nil
}
