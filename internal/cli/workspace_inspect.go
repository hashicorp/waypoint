// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"errors"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/posener/complete"
)

type WorkspaceInspectCommand struct {
	*baseCommand

	flagWorkspaceName string
}

func (c *WorkspaceInspectCommand) Run(args []string) int {
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

	if len(args) > 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	// if a workspace name is given as an argument, use that. This will take
	// precedence over the -workspace flag as it's an argument to the command
	var workspaceName string
	if len(args) > 0 {
		workspaceName = args[0]
	}

	if workspaceName == "" {
		// lookup the default
		wp, err := c.workspace()
		if err != nil {
			c.ui.Output(
				clierrors.Humanize(errors.New("error loading context default workspace")),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		workspaceName = wp
	}

	workspace, err := getWorkspace(c.Ctx, c.project.Client(), workspaceName)
	if err != nil {
		c.ui.Output(
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	var projects []string
	for _, wp := range workspace.Projects {
		projects = append(projects, wp.Project.Project)
	}

	c.ui.Output("Workspace Info:", terminal.WithHeaderStyle())

	var lastActiveTime string
	if workspace.ActiveTime != nil {
		lastActiveTime = humanize.Time(workspace.ActiveTime.AsTime())
	}

	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Name", Value: workspace.Name,
		},
		{
			Name: "Last Updated", Value: lastActiveTime,
		},
		{
			Name: "Projects", Value: strings.Join(projects, ","),
		},
	}, terminal.WithInfoStyle())

	return 0
}

func (c *WorkspaceInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *WorkspaceInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *WorkspaceInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *WorkspaceInspectCommand) Synopsis() string {
	return "Output information for a given Workspace."
}

func (c *WorkspaceInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint workspace inspect [<name>]

  Output information about a Waypoint workspace, including all projects and
  last known activity timestamp

` + c.Flags().Help())
}
