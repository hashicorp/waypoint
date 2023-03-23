// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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

	if c.flagProject == "" {
		c.ui.Output("Must explicitly set -project (-p) flag to destroy project.\n %s", c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	// Verify the project we're destroying exists
	project, err := c.project.Client().GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: c.project.Ref(),
	})
	if err != nil {
		c.ui.Output("Project %q not found.", c.project.Ref().Project, terminal.WithErrorStyle())
		return 1
	}

	// Confirmation required without `-auto-approve` flag
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

	// If project has a remote data source, queue destroy operation.
	// Otherwise, directly call server API to delete from the database.
	if project.Project.DataSource == nil {
		_, err = c.project.Client().DestroyProject(c.Ctx, &pb.DestroyProjectRequest{
			Project: c.project.Ref(),
		})
	} else {
		_, err = c.project.DestroyProject(c.Ctx, &pb.Job_DestroyProjectOp{
			Project:              &pb.Ref_Project{Project: project.Project.Name},
			SkipDestroyResources: c.skipDestroyResources,
		})
	}
	if err != nil {
		c.ui.Output("Error destroying project %q: %s", project.Project.Name, err.Error(), terminal.WithErrorStyle())
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
Usage: waypoint project destroy [options] -p <project>

  Delete the project and all resources created for all apps in the project, within
  the platform each app was deployed to.

  You must explicitly specify the project to destroy with the -project or -p flag.

  You can skip destroying app resources with the -skip-destroy-resources flag.

  You can skip the manual confirmation prompt with the -auto-approve flag.
` + c.Flags().Help())
}
