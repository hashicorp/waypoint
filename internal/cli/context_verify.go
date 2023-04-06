// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
)

type ContextVerifyCommand struct {
	*baseCommand

	flagContext        string
	flagContextDefault bool
}

func (c *ContextVerifyCommand) Run(args []string) int {
	ctx := c.Ctx

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoClient(),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(args) > 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	var name string
	if len(args) == 1 {
		name = args[0]
	}
	if name == "" {
		def, err := c.contextStorage.Default()
		if err != nil {
			c.ui.Output(
				"Error getting default context: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		name = def
	}

	config, err := c.contextStorage.Load(name)
	if err != nil {
		c.ui.Output(
			"Error loading the context %q: %s",
			name,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	sg := c.ui.StepGroup()
	step := sg.Add("Connecting with context %q...", name)
	defer step.Abort()

	conn, err := serverclient.Connect(ctx, serverclient.FromContextConfig(config))
	if err != nil {
		c.ui.Output(
			"Error connecting with context %q: %s",
			name,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	step.Update("Verifying connection is valid for context %q...", name)

	client := pb.NewWaypointClient(conn)
	if _, err := clientpkg.New(ctx,
		clientpkg.WithLogger(c.Log),
		clientpkg.WithClient(client),
	); err != nil {
		c.ui.Output(
			"Error connecting with context %q: %s",
			name,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	step.Update("Context %q connected successfully.", name)
	step.Done()

	return 0
}

func (c *ContextVerifyCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextVerifyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextVerifyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextVerifyCommand) Synopsis() string {
	return "Verify server access with a context"
}

func (c *ContextVerifyCommand) Help() string {
	return formatHelp(`
Usage: waypoint context verify [options] [NAME]

  Verify the connection information for a context is valid.

  This will use the provided context (or default) configuration,
  connect to the server, and perform test API calls to ensure the
  connection information is valid.

` + c.Flags().Help())
}
