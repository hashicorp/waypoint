package cli

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"

	"github.com/mitchellh/devflow/internal/core"
	"github.com/mitchellh/devflow/internal/pkg/flag"
	"github.com/mitchellh/devflow/sdk/terminal"
)

type ArtifactPushCommand struct {
	*baseCommand
}

func (c *ArtifactPushCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	client := c.project.Client()

	c.DoApp(c.Ctx, func(ctx context.Context, app *core.App) error {
		// Get the most recent build
		build, err := client.GetLatestBuild(ctx, &empty.Empty{})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Push it
		_, err = app.PushBuild(ctx, build)
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})

	return 0
}

func (c *ArtifactPushCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
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
	helpText := `
Usage: devflow artifact push [options]
Alias: devflow push

  Push a local build to a registry. This will push the most recent
  successful local build by default. You can view a list of the recent
  local builds using "artifact list-builds" command.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
