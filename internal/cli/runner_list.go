// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"
	"google.golang.org/protobuf/encoding/protojson"

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
		c.ui.Output("No runners found")
		return 0
	}

	if c.flagJson {
		m := protojson.MarshalOptions{
			Indent: "\t",
		}
		for _, t := range resp.Runners {
			data, err := m.Marshal(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(string(data))
		}
		return 0
	}

	tblHeaders := []string{"ID", "State", "Kind", "Labels", "Last Registered"}
	tbl := terminal.NewTable(tblHeaders...)

	var kindStr string
	var lastSeenStr string
	var stateStr string

	for _, r := range resp.Runners {
		switch r.Kind.(type) {
		case *pb.Runner_Odr:
			kindStr = "on-demand"
		case *pb.Runner_Local_:
			kindStr = "local"
		case *pb.Runner_Remote_:
			kindStr = "remote"
		default:
			kindStr = "unknown"
		}

		if r.LastSeen != nil {
			lastSeenStr = humanize.Time(r.LastSeen.AsTime())
		}

		stateStr = strings.ToLower(r.AdoptionState.String())
		if stateStr == "" {
			stateStr = "unknown"
		}
		// Omit label that the user didn't set from the output
		delete(r.Labels, "waypoint.hashicorp.com/runner-hash")

		var labelStr string
		for k, v := range r.Labels {
			labelStr += k + ":" + v + " "
		}

		tblColumn := []string{
			r.Id,
			stateStr,
			kindStr,
			labelStr,
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
