// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type WorkspaceCreateCommand struct {
	*baseCommand

	flagWorkspaceName string
}

func (c *WorkspaceCreateCommand) Run(args []string) int {
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

	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	workspaceName := args[0]

	client := c.project.Client()
	resp, err := client.UpsertWorkspace(c.Ctx, &pb.UpsertWorkspaceRequest{
		Workspace: &pb.Workspace{
			Name: workspaceName,
		},
	})
	if err != nil {
		c.ui.Output(
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// this is unlikely to happen with a nil error above, but added here to be
	// defensive.
	if resp.Workspace == nil {
		c.ui.Output(
			"no workspace returned for create command with name %q", workspaceName,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// the UpsertWorkspace call is idempotent, and does not return any
	// indication if the workspace was created or if it already existed, so we
	// report a generic response
	c.ui.Output("Workspace %q registered with the server", workspaceName, terminal.WithSuccessStyle())

	return 0
}

func (c *WorkspaceCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *WorkspaceCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *WorkspaceCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *WorkspaceCreateCommand) Synopsis() string {
	return "Create a workspace with a given name."
}

func (c *WorkspaceCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint workspace create [options] <name>

  Create a workspace in Waypoint with the given value. If a workspace with the
  given name already exists, no error will be returned. This command ignores
  the -workspace flag.

` + c.Flags().Help())
}
