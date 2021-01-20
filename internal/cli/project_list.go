package cli

import (
	"sort"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ProjectListCommand struct {
	*baseCommand
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

	resp, err := c.project.Client().ListProjects(c.Ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
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
	return c.flagSet(0, nil)
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
