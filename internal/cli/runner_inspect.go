package cli

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"reflect"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	if err != nil && status.Code(err) != codes.NotFound {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		str, err := m.MarshalToString(resp)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(str)
		return 0
	}

	kindMap := map[reflect.Type]string{
		reflect.TypeOf((*pb.Runner_Odr)(nil)):     "on-demand",
		reflect.TypeOf((*pb.Runner_Local_)(nil)):  "local",
		reflect.TypeOf((*pb.Runner_Remote_)(nil)): "remote",
	}
	var kindStr = "unknown"
	if v, ok := kindMap[reflect.TypeOf(resp.Kind)]; ok {
		kindStr = v
	}
	var lastSeenStr string
	if v, err := ptypes.Timestamp(resp.LastSeen); err == nil {
		lastSeenStr = humanize.Time(v)
	}
	var stateStr = "unknown"
	stateStr = strings.ToLower(resp.AdoptionState.String())

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
Usage: waypoint runner inspect ID

  Show detailed information about a runner.

`)
}
