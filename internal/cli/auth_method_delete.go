// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type AuthMethodDeleteCommand struct {
	*baseCommand
}

func (c *AuthMethodDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("auth method name required for deletion", terminal.WithErrorStyle())
		return 1
	}
	name := c.args[0]
	ref := &pb.Ref_AuthMethod{Name: name}

	// We special case token here. If it actually exists we let this through
	// (weird edge case, operators shouldn't name it token). If it doesn't
	// exist we notify the user they can't disable token auth.
	if name == "token" {
		_, err := c.project.Client().GetAuthMethod(c.Ctx, &pb.GetAuthMethodRequest{
			AuthMethod: ref,
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.ui.Output(strings.TrimSpace(errDeleteTokenAuth), terminal.WithErrorStyle())
				return 1
			}

			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	_, err := c.project.Client().DeleteAuthMethod(c.Ctx, &pb.DeleteAuthMethodRequest{
		AuthMethod: &pb.Ref_AuthMethod{Name: name},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Auth method deleted.", terminal.WithSuccessStyle())
	return 0
}

func (c *AuthMethodDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *AuthMethodDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthMethodDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMethodDeleteCommand) Synopsis() string {
	return "Delete a previously configured auth method."
}

func (c *AuthMethodDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint auth-method delete NAME

  Delete a previously configured auth method.

  This will not delete any users, although users may no longer be able to
  log in. Already authenticated users will remain logged in even if they
  authenticated using this auth method.

`)
}

const errDeleteTokenAuth = `
The "token" auth method can't be deleted. This auth method is required for
the Waypoint server to function.
`
