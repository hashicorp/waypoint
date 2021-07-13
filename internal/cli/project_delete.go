package cli

import (
	"fmt"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type ProjectDeleteCommand struct {
	*baseCommand
	flagForce bool
}

func (c *ProjectDeleteCommand) Run(args []string) int {
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
	// Require one argument
	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}
	name := args[0]

	// TODO: Use the Interactive terminal, when it's implemented
	if !c.ui.Interactive() && !c.flagForce {
		c.ui.Output(
			"operation interrupted, terminal is not interactive and force not provided",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if !c.flagForce {
		// Ask for confirmation before continuing
		c.ui.Output(fmt.Sprintf("project: %s", name), terminal.WithInfoStyle())
		cont, err := c.inputContinue(terminal.WarningBoldStyle)
		if err != nil {
			c.ui.Output(fmt.Sprintf("unable to get input: %v", err), terminal.WithErrorStyle())
			return 1
		}

		if !cont {
			c.ui.Output("operation interrupted", terminal.WithErrorStyle())
			return 2
		}
	}

	resp, err := c.project.Client().DeleteProject(c.Ctx, &gen.DeleteProjectRequest{
		Project: &gen.Ref_Project{Project: name},
	})

	if err != nil {
		c.ui.Output(fmt.Sprintf("unable to delete project: %v", err), terminal.WithErrorStyle())
		return 1
	}

	if !resp.Successful {
		c.ui.Output("unable to delete project: the project doesn't exist", terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output(fmt.Sprintf("project %s deleted successfully", name), terminal.WithSuccessStyle())
	return 0
}

func (c *ProjectDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "force",
			Target:  &c.flagForce,
			Default: false,
			Usage:   "skip the confirmation prompt",
		})
	})
}

func (c *ProjectDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectDeleteCommand) Synopsis() string {
	return "Delete the specified project."
}

func (c *ProjectDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint project delete PROJECT-NAME

  This command deletes a project entirely.

`)
}
