package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/serverstop"
)

type ServerDockerStopCommand struct {
	*baseCommand
}

func (c *ServerDockerStopCommand) Run(args []string) int {
	ctx := c.Ctx

	if err := c.Init(
		WithArgs(args),
		WithNoConfig(),
		WithFlags(c.Flags()),
		WithClient(false),
	); err != nil {
		return 1
	}

	status := c.ui.Status()
	defer func() {
		_ = status.Close()
	}()

	if err := serverstop.StopDocker(ctx, status); err != nil {
		c.ui.Output(
			"Error stopping server installed in Docker: %s",
			clierrors.Humanize(err), terminal.WithErrorStyle(),
		)
		return 1
	}

	c.ui.Output(
		"Successfully stopped and removed server installed in Docker",
		terminal.WithSuccessStyle(),
	)

	return 0
}

func (c *ServerDockerStopCommand) Synopsis() string {
	return "Stop and remove the Waypoint Docker container and its volume"
}

func (c *ServerDockerStopCommand) Help() string {
	return formatHelp(`
Usage: waypoint server docker stop

  Stops and removes a Waypoint server installed in Docker along with its volume.
  ...`)
}

func (c *ServerDockerStopCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {})
}
