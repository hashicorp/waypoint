package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ContextCreateCommand struct {
	*baseCommand

	flagConfig     clicontext.Config
	flagSetDefault bool
}

func (c *ContextCreateCommand) Run(args []string) int {
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
	if err := c.contextStorage.Set(name, &c.flagConfig); err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q created.", name, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "set-default",
			Target: &c.flagSetDefault,
			Usage:  "Set this context as the new default for the CLI.",
		})
		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Target: &c.flagConfig.Server.Address,
			Usage:  "Address for the server.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:   "server-insecure",
			Target: &c.flagConfig.Server.Insecure,
			Usage:  "If true, will connect to the server over plain TCP.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:   "server-require-auth",
			Target: &c.flagConfig.Server.RequireAuth,
			Usage:  "If true, will send authentication details.",
		})
	})
}

func (c *ContextCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextCreateCommand) Synopsis() string {
	return "Create a context."
}

func (c *ContextCreateCommand) Help() string {
	helpText := `
Usage: waypoint context create [options] NAME

  Creates a context.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
