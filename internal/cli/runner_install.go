package cli

import (
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/posener/complete"
)

type RunnerInstallCommand struct {
	*baseCommand

	platform     string
	mode         string
	serverUrl    string
	serverCookie string
}

func (c *RunnerInstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerInstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerInstallCommand) Flags() *flag.Sets {
	return &flag.Sets{}
}

func (c *RunnerInstallCommand) Help() string {
	return ""
}

func (c *RunnerInstallCommand) Run(args []string) int {
	return 0
}

func (c *RunnerInstallCommand) Synopsis() string {
	return ""
}
