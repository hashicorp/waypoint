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

	flagName string
	flagID   string
}

//ProjectTemplate: &pb.ProjectTemplate{
//	Name:            name,
//	Summary:         "",
//	ExpandedSummary: "",
//	ReadmeTemplate:  "",
//	WaypointProject: &pb.ProjectTemplate_WaypointProject{
//		WaypointHclTemplate: []byte(""),
//	},
//	TerraformNocodeModule: &pb.ProjectTemplate_TerraformNocodeModule{
//		Source:  "",
//		Version: "",
//	},
//	Tags: []string{},
//},

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

	name := ""
	if len(args) == 1 {
		name = args[0]
	}

	var tref pb.Ref_ProjectTemplate
	if name != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Name{
			Name: name,
		}
	} else if c.flagID != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Id{
			Id: c.flagID,
		}
		name = c.flagID
	} else if c.flagName != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Name{
			Name: c.flagName,
		}
		name = c.flagName
	} else {
		c.ui.Output("missing project template name or id", terminal.WithErrorStyle())
		return 1
	}

	_, err := c.project.Client().DeleteProjectTemplate(ctx, &pb.DeleteProjectTemplateRequest{
		ProjectTemplate: &tref,
	})
	if err != nil {
		c.ui.Output("error: ", clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("template %q deleted", name)

	return 0
}

func (c *ProjectTemplateDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagName,
			Default: "",
			Usage:   "Name of project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagID,
			Default: "",
			Usage:   "Id of project template",
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
  projects created from project templates, the terraform workspace will need to
  be cleaned up in addition to deleting the project within Waypoint using
  "waypoint project destroy".

` + c.Flags().Help())
}
