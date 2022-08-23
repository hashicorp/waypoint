package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"strings"
)

type ProjectDestroyCommand struct {
	*baseCommand

	skipDestroyResources bool
	confirm              bool
}

func (c *ProjectDestroyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// Get the project we're destroying
	project, err := c.project.Client().GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: c.project.Ref(),
	})
	if err != nil {
		return 1
	}

	// Confirmation is required for destroying a project &/or its resources
	if !c.confirm {
		proceed, err := c.ui.Input(&terminal.Input{
			Prompt: "Do you really want to destroy project \"" + project.Project.Name + "\" and its resources? Only 'yes' will be accepted to approve: ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		} else if strings.ToLower(proceed) != "yes" {
			c.ui.Output("Destroying project %q and resources requires confirmation.", project.Project.Name, terminal.WithWarningStyle())
			return 1
		}
	}
	_, err = c.project.DestroyProject(c.Ctx, &pb.Job_DestroyProjectOp{
		Project:              project.Project,
		SkipDestroyResources: c.skipDestroyResources,
	},
	)
	if err != nil {
		c.ui.Output("Error destroying project: %s", err.Error(), terminal.WithErrorStyle())
		return 1
	}
	c.ui.Output("Project %q destroyed!", project.Project.Name, terminal.WithSuccessStyle())

	return 0
}

func (c *ProjectDestroyCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "skip-destroy-resources",
			Usage:   "Skips destroying resources created for the Waypoint project.",
			Default: false,
			Target:  &c.skipDestroyResources,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "auto-approve",
			Usage:   "Destroy the project and all of its resources without confirmation.",
			Default: false,
			Target:  &c.confirm,
		})
	})
}

func (c *ProjectDestroyCommand) Synopsis() string {
	return "Delete the specified project and optionally destroy its resources."
}

func (c *ProjectDestroyCommand) Help() string {
	return formatHelp(`
Usage: waypoint project destroy [options]

  Delete the project and all resources created for all apps in the project, within
  the platform each app was deployed to.

  You can optionally skip destroying the resources by setting
  -skip-destroy-resources to true.
` + c.Flags().Help())
}
