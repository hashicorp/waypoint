package core

import (
	"context"
	"fmt"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/internal/server/logviewer"
	"github.com/mitchellh/devflow/sdk/component"
)

// Logs returns the log viewer for the given deployment.
// TODO(evanphx): test
func (a *App) Logs(ctx context.Context) (component.LogViewer, error) {
	log := a.logger.Named("logs")

	// First we attempt to query the server for logs for this deployment.
	client, err := a.client.GetLogStream(ctx, &pb.GetLogStreamRequest{})
	if err != nil {
		return nil, err
	}

	// Build our log viewer
	return &logviewer.Viewer{Stream: client}, nil

	ep, ok := a.Platform.(component.LogPlatform)
	if !ok {
		return nil, fmt.Errorf("This platform does not support logs yet")
	}

	lv, err := a.callDynamicFunc(ctx, log, nil, a.Platform, ep.LogsFunc())
	if err != nil {
		return nil, err
	}

	return lv.(component.LogViewer), nil
}
