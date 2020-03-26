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

type DeploymentCreateCommand struct {
	*baseCommand
}

func (c *DeploymentCreateCommand) Run(args []string) int {
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
		// Get the most recent pushed artifact
		push, err := client.GetLatestPushedArtifact(ctx, &empty.Empty{})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Push it
		_, err = app.Deploy(ctx, push)
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})

	return 0
}

func (c *DeploymentCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *DeploymentCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *DeploymentCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DeploymentCreateCommand) Synopsis() string {
	return "Deploy a pushed artifact."
}

func (c *DeploymentCreateCommand) Help() string {
	helpText := `
Usage: devflow deployment deploy [options]

  Deploy an application. This will deploy the most recent successful
  pushed artifact by default. You can view a list of recent artifacts
  using the "artifact list" command.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
