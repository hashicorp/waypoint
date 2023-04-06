// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectListCommand struct {
	*baseCommand

	flagJson bool
}

func (c *ProjectListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListProjects(c.Ctx, &pb.ListProjectsRequest{Pagination: &pb.PaginationRequest{}})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, p := range resp.Projects {
			str, err := m.MarshalToString(p)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	var result []string
	for _, p := range resp.Projects {
		result = append(result, p.Project)
	}

	if len(result) == 0 {
		c.ui.Output("No projects found.")
		return 0
	}
	sort.Strings(result)
	for _, p := range result {
		c.ui.Output(p)
	}

	return 0
}

func (c *ProjectListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the Project names as json.",
		})
	})
}

func (c *ProjectListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectListCommand) Synopsis() string {
	return "List all registered projects."
}

func (c *ProjectListCommand) Help() string {
	return formatHelp(`
Usage: waypoint project list

  List all registered projects.

  Projects usually map to a single version control repository and contain
  exactly one "waypoint.hcl" configuration. A project may contain multiple
  applications.

  A project is registered via the web UI, "waypoint project apply",
  or "waypoint init".

`)
}
