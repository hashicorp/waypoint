package cli

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type OndemandRunnerListCommand struct {
	*baseCommand
}

func (c *OndemandRunnerListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListOndemandRunners(c.Ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.OndemandRunners) == 0 {
		return 0
	}

	c.ui.Output("Ondemand Runner Configurations")

	tbl := terminal.NewTable("Id", "Plugin Type", "OCI Url", "Default")

	for _, p := range resp.OndemandRunners {
		def := ""
		if p.Default {
			def = "yes"
		}

		tbl.Rich([]string{
			p.Id,
			p.PluginType,
			p.OciUrl,
			def,
		}, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *OndemandRunnerListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *OndemandRunnerListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OndemandRunnerListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OndemandRunnerListCommand) Synopsis() string {
	return "List all registered ondemand runners."
}

func (c *OndemandRunnerListCommand) Help() string {
	return formatHelp(`
Usage: waypoint ondemand-runner list

  List all registered ondemand runners.

  Ondemand runners are used to dynamically start tasks to execute operations for
  projects such as building, deploying, etc.
`)
}
