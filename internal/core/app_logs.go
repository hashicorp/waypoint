// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
	"time"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/ceb/virtualceb"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Logs launches a logs plugin. Logs plugins are only used if the plugin's
// platforms plugin wishes to implement the LogsFunc protocol.
// Under traditional platform scenarios, we don't need to run a logs plugin, instead
// the logs command returns data buffered on the server sent via the entrypoint binary.
// The result of running this task is that the platform plugin is called
// and made available as a virtual instance with the given id.
// startTime inidcates the time horizon a log entry must be beyond before it is returned.
// limit controls how many log entries to emit.
func (a *App) Logs(ctx context.Context, id string, d *pb.Deployment, startTime time.Time, limit int) error {
	// Add our build to our config
	var evalCtx hcl.EvalContext

	// Start the plugin
	c, err := componentCreatorMap[component.PlatformType].Create(ctx, a, &evalCtx)
	if err != nil {
		a.logger.Error("error creating component in platform", "error", err)
		return err
	}
	defer c.Close()

	a.logger.Debug("spooling logs operation")

	logs, ok := c.Value.(component.LogPlatform)
	if !ok || logs.LogsFunc() == nil {
		a.logger.Debug("component is not an Logger or has no LogsFunc()")
		return nil
	}

	a.logger.Debug("spawn virtual ceb to handle logs", "instance-id", id)

	virt, err := virtualceb.New(a.logger, virtualceb.Config{
		InstanceId: id,
		Client:     a.client,
	})

	if err != nil {
		return err
	}

	runFn := func(ctx context.Context, lv *component.LogViewer) error {
		_, err := a.callDynamicFunc(ctx,
			a.logger,
			nil,
			c,
			logs.LogsFunc(),
			plugin.ArgNamedAny("deployment", d.Deployment),
			argmapper.Typed(lv),
		)
		if err != nil {
			a.logger.Error("error executing plugin function", "error", err)
			return err
		}

		a.logger.Info("plugin logs function finished")

		return nil
	}

	return virt.RunLogs(ctx, startTime, limit, runFn)
}
