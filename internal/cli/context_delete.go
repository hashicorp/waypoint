package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ContextDeleteCommand struct {
	*baseCommand

	flagConfig     clicontext.Config
	flagSetDefault bool
}

func (c *ContextDeleteCommand) Run(args []string) int {
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
	if err := c.contextStorage.Delete(name); err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q deleted.", name, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextDeleteCommand) Synopsis() string {
	return "Delete a context."
}

func (c *ContextDeleteCommand) Help() string {
	helpText := `
Usage: waypoint context delete [options] NAME

  Deletes a context. This will succeed if the context is already deleted.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
