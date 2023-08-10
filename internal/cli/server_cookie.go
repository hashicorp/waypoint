// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ServerCookieCommand struct {
	*baseCommand

	flagJson    bool
	flagPending bool
}

func (c *ServerCookieCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	resp, err := c.project.Client().GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// We use fmt here and not the UI helpers because UI helpers will
	// trim output horizontally on terminals that are narrow.
	fmt.Println(resp.Config.Cookie)
	return 0
}

func (c *ServerCookieCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ServerCookieCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerCookieCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerCookieCommand) Synopsis() string {
	return "Output server cookie value"
}

func (c *ServerCookieCommand) Help() string {
	return formatHelp(`
Usage: waypoint server cookie [options]

  Output the server cookie value.

  The server cookie is used in API calls to superficially ensure that
  you're communicating with the proper cluster. This isn't mean to be a
  security mechanism. This is an optional way to prevent errant API calls
  to the incorrect Waypoint cluster (if you're running multiple).

  Some unauthenticated API endpoints require the cookie value be set to
  protect against random noise, such as the runner registration endpoint.

  While the cookie isn't a security mechanism, it should be kept secret
  to prevent unnecessary API noise to a cluster.

` + c.Flags().Help())
}
