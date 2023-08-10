// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

type ArtifactPushCommand struct {
	*baseCommand
}

func (c *ArtifactPushCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Get the most recent build
		app.UI.Output("Pushing artifact for %s...", app.Ref().Application, terminal.WithHeaderStyle())
		build, err := client.GetLatestBuild(ctx, &pb.GetLatestBuildRequest{
			Application: app.Ref(),
			Workspace:   c.project.WorkspaceRef(),
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Push it
		err = app.PushBuild(ctx, &pb.Job_PushOp{
			Build: build,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ArtifactPushCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *ArtifactPushCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ArtifactPushCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ArtifactPushCommand) Synopsis() string {
	return "Push a build's artifact to a registry"
}

func (c *ArtifactPushCommand) Help() string {
	return formatHelp(`
Usage: waypoint artifact push [options]
Alias: waypoint push [options]

  Push a local build to a registry. This will push the most recent
  successful local build by default. You can view a list of the recent
  local builds using "artifact list-builds" command.

` + c.Flags().Help())
}
