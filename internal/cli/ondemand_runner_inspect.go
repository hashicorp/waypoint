package cli

import (
	"fmt"

	"github.com/golang/protobuf/jsonpb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type OnDemandRunnerInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *OnDemandRunnerInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoAutoServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("on-demand runner configuration ID required", terminal.WithErrorStyle())
		return 1
	}
	id := c.args[0]

	resp, err := c.project.Client().GetOnDemandRunnerConfig(c.Ctx, &pb.GetOnDemandRunnerConfigRequest{
		Config: &pb.Ref_OnDemandRunnerConfig{Id: id},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		str, err := m.MarshalToString(resp.Config)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(str)
		return 0
	}

	config := resp.Config
	c.ui.Output("On-Demand Runner Configuration:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "ID", Value: config.Id,
		},
		{
			Name: "Default", Value: config.Default,
		},
		{
			Name: "OCI URL", Value: config.OciUrl,
		},
		{
			Name: "Plugin Type", Value: config.PluginType,
		},
		{
			Name: "Environment Variables", Value: config.EnvironmentVariables,
		},
	}, terminal.WithInfoStyle())

	if len(config.PluginConfig) > 0 {
		c.ui.Output("Additional Plugin Configuration:", terminal.WithHeaderStyle())

		// We have to do the %s here in case the plugin config contains
		// formatting chars we don't want to error.
		c.ui.Output("\n%s", string(config.PluginConfig))
	}

	return 0
}

func (c *OnDemandRunnerInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output on-demand runner information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})
	})
}

func (c *OnDemandRunnerInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OnDemandRunnerInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OnDemandRunnerInspectCommand) Synopsis() string {
	return "Show detailed information about a configured auth method"
}

func (c *OnDemandRunnerInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner on-demand inspect NAME

  Show detailed information about an on-demand runner configuration.

`)
}
