package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/mitchellh/devflow/internal/pkg/flag"
)

type ArtifactBuildCommand struct {
	*baseCommand

	flagPush bool
}

func (c *ArtifactBuildCommand) Run([]string) int {
	//ctx := c.Ctx
	//log := c.Log.Named("artifact").Named("build")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(); err != nil {
		return 1
	}

	c.project.UI.Output("Coming soon")
	return 0
}

func (c *ArtifactBuildCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "push",
			Target: &c.flagPush,
			Usage:  "Push the artifact to the configured registry.",
		})
	})
}

func (c *ArtifactBuildCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ArtifactBuildCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ArtifactBuildCommand) Synopsis() string {
	return "Build a new versioned artifact from source."
}

func (c *ArtifactBuildCommand) Help() string {
	helpText := `
Usage: devflow artifact build [options]
Alias: devflow build

  Build a new versioned artifact from source.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
