package cli

import (
	"fmt"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type HostnameListCommand struct {
	*baseCommand
}

func (c *HostnameListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListHostnames(c.Ctx, &pb.ListHostnamesRequest{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	table := terminal.NewTable("Hostname", "FQDN", "Labels")
	for _, hostname := range resp.Hostnames {
		table.Rich([]string{
			hostname.Hostname,
			hostname.Fqdn,
			fmt.Sprintf("%v", hostname.TargetLabels),
		}, nil)
	}

	c.ui.Table(table, terminal.WithStyle("Simple"))
	return 0
}

func (c *HostnameListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *HostnameListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *HostnameListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *HostnameListCommand) Synopsis() string {
	return "List all registered hostnames."
}

func (c *HostnameListCommand) Help() string {
	return formatHelp(`
Usage: waypoint hostname list

  List all registered hostnames.

  This will list all the registered hostnames for all applications.

` + c.Flags().Help())
}
