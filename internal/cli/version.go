package cli

import (
	"errors"
	"time"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/serverclient"
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
		WithNoAutoServer(),
	); err != nil {
		return 1
	}

	out := c.VersionInfo.FullVersionNumber(true)
	c.ui.Output("CLI: %s", out)

	// Get our server version. We use a short context here.
	// TODO(izaak): validate this times out after 2 seconds
	_, err := c.initClient(serverclient.Timeout(2 * time.Second))
	if err != nil && !errors.Is(err, serverclient.ErrNoServerConfig) {
		c.ui.Output("Error connecting to server to read server version: %s", err.Error())
	}

	// version is saved on the base command when we initialize the server
	if err == nil {
		c.ui.Output("Server: %s", c.serverVersion)
	}

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

  Prints the version information for Waypoint.

  This command will show the version of the current Waypoint CLI. If
  the CLI is configured to communicate to a Waypoint server, the server
  version will also be shown.

`)
}
