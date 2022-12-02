package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateSetCommand struct {
	*baseCommand

	flagDescription string

	flagFromProject string

	flagGithubTemplateRepo  string
	flagGithubTemplateOwner string

	flagGithubRepoOwner   string
	flagGithubRepoPrivate bool
}

func (c *ProjectTemplateSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
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

	if len(args) != 1 {
		c.ui.Output("Single argument required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Checking for an existing project named: %s", c.flagFromProject)
	defer func() { s.Abort() }()

	// TODO: validate flagFromProject set
	projectResp, err := c.project.Client().GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: c.flagFromProject,
		},
	})
	if err != nil {
		c.ui.Output(
			"Error getting settings source project: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// TODO: support updating
	template := &pb.ProjectTemplate{
		Name:        name,
		Description: c.flagDescription,
		SourceCodePlatform: &pb.ProjectTemplate_Github{
			Github: &pb.ProjectTemplate_SourceCodePlatformGithub{
				Source: &pb.ProjectTemplate_SourceCodePlatformGithub_Source{
					Owner: c.flagGithubTemplateOwner,
					Repo:  c.flagGithubTemplateRepo,
				},
				Destination: &pb.ProjectTemplate_SourceCodePlatformGithub_Destination{
					Private:            c.flagGithubRepoPrivate,
					IncludeAllBranches: true,
				},
			},
		},
		Tokens:          nil,
		ProjectSettings: projectResp.Project,
	}

	s.Update("Setting project template")
	_, err = c.project.Client().UpsertProjectTemplate(ctx, &pb.UpsertProjectTemplateRequest{
		ProjectTemplate: template,
	})
	if err != nil {
		c.ui.Output(
			"Error upserting project template", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	s.Update("Project template %q set", name)

	s.Done()

	return 0
}

func (c *ProjectTemplateSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "description",
			Target:  &c.flagDescription,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "from-project",
			Target:  &c.flagFromProject,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "github-template-repo",
			Target:  &c.flagGithubTemplateRepo,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "github-template-owner",
			Target:  &c.flagGithubTemplateOwner,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "github-repo-owner",
			Target:  &c.flagGithubRepoOwner,
			Default: "",
			Usage:   "todo",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "github-repo-private",
			Target:  &c.flagGithubRepoPrivate,
			Default: true,
			Usage:   "todo",
		})
	})
}

func (c *ProjectTemplateSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateSetCommand) Synopsis() string {
	return "Create or update a project template."
}

func (c *ProjectTemplateSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint project template set [options] NAME

  Create or update a project template.

  This will create a new project template with the given options. If a 
  project template with the same name already exists, this will update 
  the existing project template using the fields that are set.

` + c.Flags().Help())
}
