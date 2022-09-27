package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type ConfigSourceGetCommand struct {
	*baseCommand

	flagType string
}

func (c *ConfigSourceGetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// type is required
	if c.flagType == "" {
		c.ui.Output("A source type must be specified with '-type'.\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	// Get our config source
	client := c.project.Client()
	resp, err := client.GetConfigSource(c.Ctx, &pb.GetConfigSourceRequest{
		Scope: &pb.GetConfigSourceRequest_Global{
			Global: &pb.Ref_Global{},
		},

		Type: c.flagType,
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.ConfigSources) == 0 {
		c.project.UI.Output(
			"Dynamic config source %q is not configured.\n\n"+
				"Note that this doesn't mean that this config source is not usable.\n"+
				"Many config sources work with no explicitly set configurations.",
			c.flagType, terminal.WithErrorStyle())
		return 1
	}

	// we use the first value because this will be the most specific since
	// we do a prefix search.
	cs := resp.ConfigSources[0]
	table := terminal.NewTable("Key", "Value")
	for k, v := range cs.Config {
		table.Rich([]string{
			k,
			v,
		}, []string{
			"",
			"",
		})
	}
	c.ui.Table(table)
	return 0
}

func (c *ConfigSourceGetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "type",
			Target: &c.flagType,
			Usage:  "Dynamic source type to look up, such as 'vault'.",
		})
	})
}

func (c *ConfigSourceGetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSourceGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSourceGetCommand) Synopsis() string {
	return "Get the configuration for a dynamic source plugin"
}

func (c *ConfigSourceGetCommand) Help() string {
	return formatHelp(`
Usage: waypoint config source-get [options]

  Get the configuration for a dynamic configuration source plugin.

  This does not list the dynamic configuration variables for an application.
  This command is for configuring the plugin that is used to fetch dynamic
  configurations globally for the server.

  To use this command, you must specify a "-type" flag.

  Configuration for this command is global. The "-app", "-project", and
  "-workspace" flags are ignored on this command.

` + c.Flags().Help())
}
