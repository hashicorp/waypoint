package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type PipelineListCommand struct {
	*baseCommand
}

func (c *PipelineListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	return 0
}

func (c *PipelineListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *PipelineListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineListCommand) Synopsis() string {
	return "List all pipelines for a project."
}

func (c *PipelineListCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline list

  List all pipelines for a project.

`)
}
