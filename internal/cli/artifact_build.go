package cli

import (
	"context"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ArtifactBuildCommand struct {
	*baseCommand

	flagPush bool
}

func (c *ArtifactBuildCommand) Run(args []string) int {
	//ctx := c.Ctx
	//log := c.Log.Named("artifact").Named("build")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	c.DoApp(c.Ctx, func(ctx context.Context, app *core.App) error {
		_, err := app.Build(ctx)
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})

	return 0
}

func (c *ArtifactBuildCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetLabel, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "push",
			Target: &c.flagPush,
			Usage:  "Push the artifact to the configured registry.",
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
	return "Build a new versioned artifact from source."
}

func (c *ArtifactBuildCommand) Help() string {
	helpText := `
Usage: waypoint artifact build [options]
Alias: waypoint build

  Build a new versioned artifact from source.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
