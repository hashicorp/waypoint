package cli

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ServerBootstrapCommand struct {
	*baseCommand
}

func (c *ServerBootstrapCommand) Run(args []string) int {
	ctx := c.Ctx

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	client := c.project.Client()
	resp, err := client.BootstrapToken(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error bootstrapping the server: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	c.ui.Output(resp.Token)
	return 0
}

func (c *ServerBootstrapCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ServerBootstrapCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerBootstrapCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerBootstrapCommand) Synopsis() string {
	return "Bootstrap the server and retrieve the initial auth token"
}

func (c *ServerBootstrapCommand) Help() string {
	return formatHelp(`
Usage: waypoint server bootstrap [options]

  Bootstrap a new server and retrieve the initial auth token.

  When a server is started for the first time against an empty database,
  it is able to be bootstrapped. The bootstrap process retrieves the initial
  auth token for the server. After the auth token is retrieved, it can never
  be bootstrapped again.

  This command is only required for manually run servers. For servers
  installed with "waypoint install", the bootstrap is done automatically
  during the install process.

` + c.Flags().Help())
}
