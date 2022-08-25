package cli

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type RunnerInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *RunnerInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("Runner name required.", terminal.WithErrorStyle())
		return 1
	}
	id := c.args[0]

	resp, err := c.project.Client().GetRunner(c.Ctx, &pb.GetRunnerRequest{RunnerId: id})
	if err != nil {
		if status.Code(err) != codes.NotFound {
			c.ui.Output("runner not found", terminal.WithErrorStyle())
			return 1
		}
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		data, err := protojson.MarshalOptions{
			Indent: "\t",
		}.Marshal(resp)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(string(data))
		return 0
	}

	var kindStr string
	switch resp.Kind.(type) {
	case *pb.Runner_Odr:
		kindStr = "on-demand"
	case *pb.Runner_Local_:
		kindStr = "local"
	case *pb.Runner_Remote_:
		kindStr = "remote"
	default:
		kindStr = "unknown"
	}

	var lastSeenStr string
	if resp.LastSeen != nil {
		lastSeenStr = humanize.Time(resp.LastSeen.AsTime())
	}

	stateStr := strings.ToLower(resp.AdoptionState.String())
	if stateStr == "" {
		stateStr = "unknown"
	}

	// Omit label that the user didn't set from the output
	if _, ok := resp.Labels["waypoint.hashicorp.com/runner-hash"]; ok {
		delete(resp.Labels, "waypoint.hashicorp.com/runner-hash")
	}

	c.ui.Output("Runner:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "ID", Value: resp.Id,
		},
		{
			Name: "Adoption State", Value: stateStr,
		},
		{
			Name: "Kind", Value: kindStr,
		},
		{
			Name: "Last Registered", Value: lastSeenStr,
		},
		{
			Name: "Labels", Value: resp.Labels,
		},
		{
			Name: "Online", Value: resp.Online,
		},
	}, terminal.WithInfoStyle())

	return 0
}

func (c *RunnerInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output runner information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})
	})
}

func (c *RunnerInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerInspectCommand) Synopsis() string {
	return "Show detailed information about a runner."
}

func (c *RunnerInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner inspect <id>

  Show detailed information about a runner.

`)
}
