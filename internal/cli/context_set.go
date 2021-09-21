package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ContextSetCommand struct {
	*baseCommand
}

func (c *ContextSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	if c.flagWorkspace == "" {
		// Require one argument
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	def, err := c.contextStorage.Default()
	if err != nil {
		// TODO log error
		return 1
	}

	name := def
	// If we still have no name, then we do nothing. We also accept
	// "-" as a valid name that means "do nothing".
	if name == "" || name == "-" {
		// TODO log error
		return 1
	}

	// Load it and set it.
	cfg, err := c.contextStorage.Load(name)
	if err != nil {
		// TODO log error
		return 1
	}

	// set new workspace
	cfg.Workspace = c.flagWorkspace

	// store updated context
	if err := c.contextStorage.Set(name, cfg); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q updated to use %s workspace.", name, cfg.Workspace, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextSetCommand) Synopsis() string {
	return "Set a property of the current context."
}

func (c *ContextSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint context set [options]

  Sets a property of the current context. The only property supported at this
  time is -workspace.

  To use this command, use the global -workspace flag to set the default
  workspace for the current context.

  To restore this CLI context to use the default workspace, use
  -workspace=default

` + c.Flags().Help())
}
