// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"encoding/json"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerProfileListCommand struct {
	*baseCommand
}

func (c *RunnerProfileListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListOnDemandRunnerConfigs(c.Ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.Configs) == 0 {
		return 0
	}

	c.ui.Output("Runner profiles")

	tbl := terminal.NewTable("Name", "Plugin Type", "OCI Url", "Target Runner",
		"Default")

	for _, p := range resp.Configs {
		def := ""
		if p.Default {
			def = "yes"
		}

		var targetRunner string
		if p.TargetRunner != nil {
			switch t := p.TargetRunner.Target.(type) {
			case *pb.Ref_Runner_Any:
				targetRunner = "*"
			case *pb.Ref_Runner_Id:
				targetRunner = t.Id.Id
			case *pb.Ref_Runner_Labels:
				s, _ := json.Marshal(t.Labels.Labels)
				targetRunner = "labels: " + string(s)
			}
		}

		tbl.Rich([]string{
			p.Name,
			p.PluginType,
			p.OciUrl,
			targetRunner,
			def,
		}, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *RunnerProfileListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *RunnerProfileListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileListCommand) Synopsis() string {
	return "List all registered runner profiles."
}

func (c *RunnerProfileListCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile list

  List runner profiles.

  Runner profiles are used to dynamically start tasks (i.e. on-demand runners) to execute 
  operations for projects such as building, deploying, etc.
`)
}
