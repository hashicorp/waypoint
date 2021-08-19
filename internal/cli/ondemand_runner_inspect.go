package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type OnDemandRunnerInspectCommand struct {
	*baseCommand
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
	fmt.Println(proto.MarshalTextString(resp.Config))
	return 0
}

func (c *OnDemandRunnerInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
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
