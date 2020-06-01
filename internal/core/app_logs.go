package core

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logviewer"
	"github.com/hashicorp/waypoint/sdk/component"
)

// Logs returns the log viewer for the given deployment.
// TODO(evanphx): test
func (a *App) Logs(ctx context.Context, d *pb.Deployment) (component.LogViewer, error) {
	log := a.logger.Named("logs")

	// First we attempt to query the server for logs for this deployment.
	log.Info("requesting log stream", "deployment_id", d.Id)
	client, err := a.client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		DeploymentId: d.Id,
	})
	if err != nil {
		return nil, err
	}

	// Build our log viewer
	return &logviewer.Viewer{Stream: client}, nil
}
