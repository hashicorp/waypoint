// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
	"github.com/skratchdot/open-golang/open"
)

type UICommand struct {
	*baseCommand

	flagAuthenticate bool
}

func (c *UICommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if c.project.Local() {
		c.project.UI.Output("Waypoint must be configured in server mode to access the UI", terminal.WithWarningStyle())
	}

	// Get our API client
	client := c.project.Client()

	// Get our default context (used context)
	name, err := c.contextStorage.Default()
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if name == "" {
		// No default context found, do they have any at all?
		if allContexts, err := c.contextStorage.List(); len(allContexts) == 0 || err != nil {
			if err != nil {
				c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			c.ui.Output("\n"+noContextFoundError, terminal.WithWarningStyle())
			return 1
		}

		// They have some context, but no default set
		c.ui.Output("\n"+wpNoServerContext, terminal.WithWarningStyle())
		return 1
	}

	ctxConfig, err := c.contextStorage.Load(name)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	var inviteToken string
	if c.flagAuthenticate {
		c.ui.Output("Creating invite token", terminal.WithStyle(terminal.HeaderStyle))
		c.ui.Output("This invite token will be exchanged for an authentication \ntoken that your browser stores.")

		resp, err := client.GenerateInviteToken(c.Ctx, &pb.InviteTokenRequest{
			Duration: (2 * time.Minute).String(),
		})
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		inviteToken = resp.Token
	}

	// todo(mitchellh: current default port is hardcoded, cannot configure http address)
	addr := strings.Split(ctxConfig.Server.Address, ":")[0]
	// Default Docker platform HTTP port, for now
	port := 9702 // TODO(briancain): properly look this up from server cfg
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Opening browser", terminal.WithStyle(terminal.HeaderStyle))

	uiAddr := fmt.Sprintf("https://%s:%d", addr, port)
	if c.flagAuthenticate {
		uiAddr = fmt.Sprintf("%s/auth/invite?token=%s&cli=true", uiAddr, inviteToken)
	}

	open.Run(uiAddr)

	return 0
}

func (c *UICommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "authenticate",
			Target:  &c.flagAuthenticate,
			Default: false,
			Usage:   "Creates a new invite token and passes it to the UI for authorization",
		})

	})
}

func (c *UICommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UICommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UICommand) Synopsis() string {
	return "Open the web UI"
}

func (c *UICommand) Help() string {
	return formatHelp(`
Usage: waypoint ui [options]

  Opens the new UI. When provided a flag, will automatically open the
  token invite page with an invite token for authentication.

` + c.Flags().Help())
}

var (
	noContextFoundError = strings.TrimSpace(`
Attempted to open the ui but found no Waypoint contexts. Please either create a new
context that uses an existing Waypoint server with 'waypoint context create'
or install a server using 'waypoint server install' which will set up a context for you.
`)
)
