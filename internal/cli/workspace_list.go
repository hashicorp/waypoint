// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/pkg/server/gen"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type WorkspaceListCommand struct {
	*baseCommand
}

func (c *WorkspaceListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListWorkspaces(c.Ctx, &pb.ListWorkspacesRequest{
		Scope: &pb.ListWorkspacesRequest_Global{},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	var result []string
	for _, p := range resp.Workspaces {
		result = append(result, p.Name)
	}

	if len(result) == 0 {
		c.ui.Output("No workspaces found.")
		return 0
	}
	sort.Strings(result)

	table := terminal.NewTable("Name", "Projects")
	for _, workspaceName := range result {
		workspace, err := getWorkspace(c.Ctx, c.project.Client(), workspaceName)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		var projects []string
		for _, wp := range workspace.Projects {
			projects = append(projects, wp.Project.Project)
		}

		table.Rich([]string{
			workspace.Name,
			strings.Join(projects, ","),
		}, nil)
	}
	c.ui.Table(table)

	return 0
}

func (c *WorkspaceListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *WorkspaceListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *WorkspaceListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *WorkspaceListCommand) Synopsis() string {
	return "List workspaces for the current context."
}

func (c *WorkspaceListCommand) Help() string {
	return formatHelp(`
Usage: waypoint workspace list

  Lists all the known workspaces available to the CLI for the current Waypoint server
  context.

` + c.Flags().Help())
}

func getWorkspace(ctx context.Context, client pb.WaypointClient, name string) (*gen.Workspace, error) {
	resp, err := client.GetWorkspace(ctx, &pb.GetWorkspaceRequest{
		Workspace: &pb.Ref_Workspace{
			Workspace: name,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("no workspace found for name %q", name)
		}
		return nil, err
	}

	// this is unlikely to happen without first hitting the codes.NotFound error
	// above, but added here to be defensive.
	if resp.Workspace == nil {
		return nil, fmt.Errorf("no workspace returned for name %q", name)
	}
	return resp.Workspace, nil
}
