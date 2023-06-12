// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateDeleteCommand struct {
	*baseCommand

	flagID string
}

func (c *ProjectTemplateDeleteCommand) Run(args []string) int {
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()
	ctx := c.Ctx

	if len(args) > 1 {
		c.ui.Output("Only one project template may be specified at a time.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := ""
	if len(args) == 1 {
		name = args[0]
	}

	if name != "" && c.flagID != "" {
		c.ui.Output("Name argument and id flag may not be specified together.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}
	if name == "" && c.flagID == "" {
		c.ui.Output("Missing project template name or id.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	var tref pb.Ref_ProjectTemplate
	if name != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Name{
			Name: name,
		}
	}
	if c.flagID != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Id{
			Id: c.flagID,
		}
		name = c.flagID
	}

	_, err := c.project.Client().DeleteProjectTemplate(ctx, &pb.DeleteProjectTemplateRequest{
		ProjectTemplate: &tref,
	})
	if err != nil {
		c.ui.Output("Encountered an error while deleting the project template: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Template %q deleted", name, terminal.WithSuccessStyle())

	return 0
}

func (c *ProjectTemplateDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagID,
			Default: "",
			Usage:   "Id of the project template to delete. Mutually exclusive with Name argument.",
		})
	})
}

func (c *ProjectTemplateDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateDeleteCommand) Synopsis() string {
	return "Delete a project template."
}

func (c *ProjectTemplateDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint template delete [options] [NAME]

  Delete a project template.

  This will delete a project template with a given name or id.

  Deleting a project template only deletes the template and does not delete the
  projects which have been created from the project template. To delete
  projects created from project templates, the Terraform workspace will need to
  be cleaned up in addition to deleting the project within Waypoint using
  "waypoint project destroy".

` + c.Flags().Help())
}
