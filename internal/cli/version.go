// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"errors"
	"time"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/version"
	"github.com/hashicorp/waypoint/pkg/serverclient"
)

type VersionCommand struct {
	*baseCommand

	VersionInfo *version.VersionInfo
}

func (c *VersionCommand) Run(args []string) int {
	flagSet := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
		WithNoClient(),
		WithNoLocalServer(),
	); err != nil {
		return 1
	}

	out := c.VersionInfo.FullVersionNumber(true)
	c.ui.Output("CLI: %s", out)

	// Get our server version. We use a short context here.
	ctx, cancel := context.WithTimeout(c.Ctx, 2*time.Second)
	defer cancel()
	client, err := c.initClient(ctx)
	if err != nil && !errors.Is(err, serverclient.ErrNoServerConfig) {
		c.ui.Output("Error connecting to server to read server version: %s", err.Error())
	}

	if err == nil {
		server := client.ServerVersion()
		c.ui.Output("Server: %s", server.Version)
	}

	return 0
}

func (c *VersionCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *VersionCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *VersionCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the version of this Waypoint CLI"
}

func (c *VersionCommand) Help() string {
	return formatHelp(`
Usage: waypoint version

  Prints the version information for Waypoint.

  This command will show the version of the current Waypoint CLI. If
  the CLI is configured to communicate to a Waypoint server, the server
  version will also be shown.

`)
}
