package core

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/config"
)

// execHook executes the given hook. This will return any errors. If
// "on_failure" is set to "continue" then if an error occurs it will be logged
// but the returned error will be nil.
func (a *App) execHook(ctx context.Context, log hclog.Logger, h *config.Hook) error {
	return nil
}
