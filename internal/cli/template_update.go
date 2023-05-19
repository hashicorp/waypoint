package cli

import (
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateUpdateCommand struct {
	*baseCommand

	flagID   string
	flagName string

	flagSummary                    string
	flagExpandedSummary            string
	flagReadmeMarkdownTemplatePath string
	flagWaypointHCLTemplatePath    string
	flagTFCNoCodeModuleSource      string
	flagTFCNoCodeModuleVersion     string
	flagTags                       []string
}

func (c *ProjectTemplateUpdateCommand) Run(args []string) int {
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
		return 1
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
		c.ui.Output("Missing project template name or id.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	checkResp, err := c.project.Client().GetProjectTemplate(ctx, &pb.GetProjectTemplateRequest{
		ProjectTemplate: &pb.Ref_ProjectTemplate{
			Ref: &pb.Ref_ProjectTemplate_Name{
				Name: name,
			},
		},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	if checkResp.ProjectTemplate == nil {
		c.ui.Output(
			"Project template %q does not exist", checkResp.ProjectTemplate.Name,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	template := checkResp.ProjectTemplate

	template.Name = name
	if c.flagSummary != "" {
		template.Summary = c.flagSummary
	}
	if c.flagExpandedSummary != "" {
		template.ExpandedSummary = c.flagExpandedSummary
	}

	if c.flagReadmeMarkdownTemplatePath != "" {
		rmt, err := os.ReadFile(c.flagReadmeMarkdownTemplatePath)
		if err != nil {
			c.ui.Output("Unable to read readme.md template file: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		template.ReadmeMarkdownTemplate = rmt
	}

	if c.flagWaypointHCLTemplatePath != "" {
		wpt, err := os.ReadFile(c.flagWaypointHCLTemplatePath)
		if err != nil {
			c.ui.Output("Unable to read waypoint.hcl template file: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		template.WaypointProject = &pb.ProjectTemplate_WaypointProject{
			WaypointHclTemplate: wpt,
		}
	}

	if c.flagTFCNoCodeModuleSource != "" && c.flagTFCNoCodeModuleVersion == "" {
		c.ui.Output("Terraform no code module version required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagTFCNoCodeModuleSource == "" && c.flagTFCNoCodeModuleVersion != "" {
		c.ui.Output("Terraform no code module source required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagTFCNoCodeModuleSource != "" && c.flagTFCNoCodeModuleVersion != "" {
		template.TerraformNocodeModule = &pb.ProjectTemplate_TerraformNocodeModule{
			Source:  c.flagTFCNoCodeModuleSource,
			Version: c.flagTFCNoCodeModuleVersion,
		}
	}

	template.Tags = c.flagTags

	_, err = c.project.Client().UpdateProjectTemplate(ctx, &pb.UpdateProjectTemplateRequest{
		ProjectTemplate: template,
	})
	if err != nil {
		c.ui.Output("Error updating project template: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("template updated!")

	return 0
}

func (c *ProjectTemplateUpdateCommand) Flags() *flag.Sets {
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

		f.StringVar(&flag.StringVar{
			Name:    "summary",
			Target:  &c.flagSummary,
			Default: "",
			Usage:   "Summary for the project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "expanded-summary",
			Target:  &c.flagExpandedSummary,
			Default: "",
			Usage:   "Expanded Summary for the project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "readme-markdown-template-path",
			Target:  &c.flagReadmeMarkdownTemplatePath,
			Default: "",
			Usage:   "Path to a markdown readme template for projects created from a project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "waypoint-hcl-template-path",
			Target:  &c.flagWaypointHCLTemplatePath,
			Default: "",
			Usage:   "Path to a templated waypoint.hcl file for projects created from a project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "tfc-nocode-module-source",
			Target:  &c.flagTFCNoCodeModuleSource,
			Default: "",
			Usage:   "The name of the Terraform no-code module from a Terraform registry that the template should use to provision infrastructure for Waypoint projects created from the template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "tfc-nocode-module-version",
			Target:  &c.flagTFCNoCodeModuleVersion,
			Default: "",
			Usage:   "The version of the Terraform no-code module from a Terraform registry that the template should use to provision infrastructure for Waypoint projects created from the template",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "tag",
			Target: &c.flagTags,
			Usage:  "A tag to add to the project template",
		})
	})
}

func (c *ProjectTemplateUpdateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateUpdateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateUpdateCommand) Synopsis() string {
	return "Update a project template."
}

func (c *ProjectTemplateUpdateCommand) Help() string {
	return formatHelp(`
Usage: waypoint template create [options] [NAME]

  Update a project template.

  This will update an existing project template with the given options.

` + c.Flags().Help())
}
