package cli

import (
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ContextListCommand struct {
	*baseCommand
}

func (c *ContextListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	// Get our direct stdout handle cause we're going to be writing colors
	// and want color detection to work.
	out, _, err := c.ui.OutputWriters()
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	// Get our contexts
	names, err := c.contextStorage.List()
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	if len(names) == 0 {
		c.ui.Output("No contexts. Create one with `waypoint context create`.")
		return 0
	}

	// Get our default
	def, err := c.contextStorage.Default()
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"", "Name", "Server Address"})
	table.SetBorder(false)
	for _, name := range names {
		ctxConfig, err := c.contextStorage.Load(name)
		if err != nil {
			c.ui.Output("Error loading context %q: %s", name, err.Error(), terminal.WithErrorStyle())
			return 1
		}

		// Determine our bullet
		defStatus := ""
		if name == def {
			defStatus = "*"
		}

		table.Rich([]string{
			defStatus,
			name,
			ctxConfig.Server.Address,
		}, []tablewriter.Colors{
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
		})
	}
	table.Render()

	return 0
}

func (c *ContextListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextListCommand) Synopsis() string {
	return "List contexts."
}

func (c *ContextListCommand) Help() string {
	helpText := `
Usage: waypoint context list [options]

  Lists the contexts available to the CLI.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
