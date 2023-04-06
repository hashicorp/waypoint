// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ArtifactBuildCommand struct {
	*baseCommand

	flagPush bool
}

func (c *ArtifactBuildCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		app.UI.Output("Building %s...", app.Ref().Application, terminal.WithHeaderStyle())
		buildResult, err := app.Build(ctx, &pb.Job_BuildOp{
			DisablePush: !c.flagPush,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if buildResult.Push != nil {
			app.UI.Output("\nCreated artifact v%d", buildResult.Push.Sequence)
		}

		// Show input variable values used in build
		app.UI.Output("Variables used:", terminal.WithHeaderStyle())
		resp, err := c.project.Client().GetJob(ctx, &pb.GetJobRequest{
			JobId: buildResult.Build.JobId,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		tbl := fmtVariablesOutput(resp.VariableFinalValues)
		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ArtifactBuildCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "push",
			Target:  &c.flagPush,
			Default: true,
			Usage:   "Push the artifact to the configured registry.",
		})
	})
}

func (c *ArtifactBuildCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ArtifactBuildCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ArtifactBuildCommand) Synopsis() string {
	return "Build a new versioned artifact from source"
}

func (c *ArtifactBuildCommand) Help() string {
	return formatHelp(`
Usage: waypoint artifact build [options]
Alias: waypoint build [options]

  Build a new versioned artifact from source.

` + c.Flags().Help())
}
