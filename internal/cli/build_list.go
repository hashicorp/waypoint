// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serversort "github.com/hashicorp/waypoint/pkg/server/sort"
)

type BuildListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
	flagId           idFormat
}

func (c *BuildListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	// Get our API client
	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		var wsRef *pb.Ref_Workspace
		if !c.flagWorkspaceAll {
			wsRef = c.project.WorkspaceRef()
		}

		// List builds
		resp, err := client.ListBuilds(c.Ctx, &pb.ListBuildsRequest{
			Application: app.Ref(),
			Workspace:   wsRef,
		})
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		sort.Sort(serversort.BuildStartDesc(resp.Builds))

		const bullet = "●"

		table := terminal.NewTable("", "ID", "Workspace", "Builder", "Started", "Completed", "Pipeline")
		for _, b := range resp.Builds {
			// Determine our bullet
			status := ""
			statusColor := ""
			switch b.Status.State {
			case pb.Status_RUNNING:
				status = bullet
				statusColor = terminal.Yellow

			case pb.Status_SUCCESS:
				status = "✔"
				statusColor = terminal.Green

			case pb.Status_ERROR:
				status = "✖"
				statusColor = terminal.Red
			}

			// Parse our times
			var startTime, completeTime string
			if b.Status.StartTime != nil {
				startTime = humanize.Time(b.Status.StartTime.AsTime())
			}
			if b.Status.CompleteTime != nil {
				completeTime = humanize.Time(b.Status.CompleteTime.AsTime())
			}

			var pipeline string
			j, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: b.JobId,
			})
			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return err
			}
			if j.Pipeline != nil {
				pipeline = "name: " + j.Pipeline.PipelineName + ", run: " + strconv.FormatUint(j.Pipeline.RunSequence, 10) + ", step: " + j.Pipeline.Step
			}

			table.Rich([]string{
				status,
				c.flagId.FormatId(b.Sequence, b.Id),
				b.Workspace.Workspace,
				b.Component.Name,
				startTime,
				completeTime,
				pipeline,
			}, []string{
				statusColor,
			})
		}

		app.UI.Output("%s", app.Ref().Application, terminal.WithHeaderStyle())
		c.ui.Table(table)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *BuildListCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "workspace-all",
			Target: &c.flagWorkspaceAll,
			Usage:  "List builds in all workspaces for this project and application.",
		})

		initIdFormat(f, &c.flagId)
	})
}

func (c *BuildListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *BuildListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *BuildListCommand) Synopsis() string {
	return "List builds."
}

func (c *BuildListCommand) Help() string {
	return formatHelp(`
Usage: waypoint artifact list-builds [options]

List artifacts created from a build. An artifact is the result of a build or
registry. This is the metadata only. The binary contents of an artifact are
expected to be stored in a registry.

` + c.Flags().Help())
}
