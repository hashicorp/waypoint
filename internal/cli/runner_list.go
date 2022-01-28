package cli

import (
	"fmt"
	"reflect"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerListCommand struct {
	*baseCommand

	flagJson    bool
	flagPending bool
}

func (c *RunnerListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	resp, err := c.project.Client().ListRunners(ctx, &pb.ListRunnersRequest{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.Runners) == 0 {
		return 0
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, t := range resp.Runners {
			str, err := m.MarshalToString(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	tblHeaders := []string{"State", "ID", "Kind", "Last Registered"}
	tbl := terminal.NewTable(tblHeaders...)

	kindMap := map[reflect.Type]string{
		reflect.TypeOf((*pb.Runner_Odr)(nil)):     "on-demand",
		reflect.TypeOf((*pb.Runner_Local_)(nil)):  "local",
		reflect.TypeOf((*pb.Runner_Remote_)(nil)): "remote",
	}

	stateMap := map[pb.Runner_AdoptionState]string{
		pb.Runner_NEW:        "pending",
		pb.Runner_PREADOPTED: "pre-adopted",
		pb.Runner_ADOPTED:    "adopted",
		pb.Runner_REJECTED:   "rejected",
	}

	for _, r := range resp.Runners {
		var kindStr = "unknown"
		if v, ok := kindMap[reflect.TypeOf(r.Kind)]; ok {
			kindStr = v
		}

		var lastSeenStr string
		if v, err := ptypes.Timestamp(r.LastSeen); err == nil {
			lastSeenStr = humanize.Time(v)
		}

		var stateStr = "unknown"
		if v, ok := stateMap[r.AdoptionState]; ok {
			stateStr = v
		}

		tblColumn := []string{
			stateStr,
			r.Id,
			kindStr,
			lastSeenStr,
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *RunnerListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "pending",
			Target: &c.flagPending,
			Usage:  "List only runners pending adoption.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output runner configuration list information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})

	})
}

func (c *RunnerListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerListCommand) Synopsis() string {
	return "List registered runners"
}

func (c *RunnerListCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner list [options]

  List runners that are registered with the Waypoint server.

  This will list all runners that the Waypoint server currently knows
  about. This list does not guarantee each runner is online; Waypoint currently
  does not expose online/offline status. This lists runners that have registered
  at least once with the server.

  This can be used to find pending runners that need to be adopted. Runners
  that are pending (not adopted) will not be sent any jobs or configuration.
  Runners that are accepted (adopted) are sent jobs. Runners that are
  "preadopted" are sent jobs but have avoided the manual adoption process by
  being preconfigured with a valid runner token (see "waypoint runner token").
  Runners that are "rejected" are never given jobs, and error immediately if
  they try to register.

  If a runner is pending, you can adopt it using "waypoint runner adopt ID"
  where "ID" comes from the output from this command. You can explicitly
  reject a runner using "waypoint runner reject ID". A runner can be rejected
  at any time, even after it is already adopted.

` + c.Flags().Help())
}
