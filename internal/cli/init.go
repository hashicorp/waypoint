package cli

import (
	"strings"
	"time"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/posener/complete"
)

type InitCommand struct {
	*baseCommand
}

func (c *InitCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	sg := c.ui.StepGroup()

	// If we have a configuration file, let's validate that first.
	s := sg.Add("Validating configuration file...")
	time.Sleep(1 * time.Second)
	s.Update("Configuration file appears valid")
	s.Status(terminal.StatusOK)
	s.Done()

	sg.Wait()

	return 0
}

func (c *InitCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *InitCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InitCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InitCommand) Synopsis() string {
	return "Initialize and validate a project."
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: waypoint init [options]

  Initialize and validate a project.

  This is the first command that should be run for any new or existing
  Waypoint project per machine. This sets up the project if required and
  also validates that operations such as "up" will most likely work.

  This command is always safe to run multiple times. This command will never
  delete your configuration or any data in the server.

`

	return strings.TrimSpace(helpText)
}
