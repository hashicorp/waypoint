package core

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ConfigSync writes all the app configuration in the waypoint.hcl file to
// the server.
func (a *App) ConfigSync(ctx context.Context) error {
	vars, err := a.config.ConfigVars()
	if err != nil {
		return err
	}

	// If we have no vars then we don't want to round trip to the server.
	if len(vars) == 0 {
		return nil
	}

	_, err = a.client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: vars,
	})
	return err
}
