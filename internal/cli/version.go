package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/version"
)

type VersionCommand struct {
	*baseCommand

	VersionInfo *version.VersionInfo
}

func (c *VersionCommand) Run(args []string) int {
	flagSet := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	out := c.VersionInfo.FullVersionNumber(true)
	c.ui.Output(out)

	return 0
}

func (c *VersionCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *VersionCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *VersionCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the version of this Waypoint CLI"
}

func (c *VersionCommand) Help() string {
	return formatHelp(`
Usage: waypoint version

  Prints the version of this Waypoint CLI.

  There are no arguments or flags to this command. Any additional arguments or
  flags are ignored.
`)
}
