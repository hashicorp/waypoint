package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ContextUseCommand struct {
	*baseCommand
}

func (c *ContextUseCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()

	// Require one argument
	if len(args) != 1 {
		c.ui.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	// Get our contexts
	if err := c.contextStorage.SetDefault(name); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("Context %q doesn't exist.", name)
		}

		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q is now the default.", name, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextUseCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextUseCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextUseCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextUseCommand) Synopsis() string {
	return "Set the default context."
}

func (c *ContextUseCommand) Help() string {
	helpText := `
Usage: waypoint context use [options] NAME

  Set the default context for the CLI.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
