package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type PipelineRunCommand struct {
	*baseCommand
}

func (c *PipelineRunCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	return 0
}

func (c *PipelineRunCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *PipelineRunCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineRunCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineRunCommand) Synopsis() string {
	return "Manually execute a pipeline"
}

func (c *PipelineRunCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline run [options] <pipeline-id>

  Run a pipeline by id.

` + c.Flags().Help())
}
