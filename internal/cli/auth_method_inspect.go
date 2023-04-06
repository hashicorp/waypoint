// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type AuthMethodInspectCommand struct {
	*baseCommand
}

func (c *AuthMethodInspectCommand) Run(args []string) int {
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
		c.ui.Output("auth method name required", terminal.WithErrorStyle())
		return 1
	}
	name := c.args[0]

	// Special case token cause its not real.
	if name == "token" {
		c.ui.Output(outInspectToken)
		return 0
	}

	resp, err := c.project.Client().GetAuthMethod(c.Ctx, &pb.GetAuthMethodRequest{
		AuthMethod: &pb.Ref_AuthMethod{Name: name},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	am := resp.AuthMethod

	// Some redaction
	switch method := am.Method.(type) {
	case *pb.AuthMethod_Oidc:
		method.Oidc.ClientSecret = "[REDACTED: client secret]"
	}

	fmt.Println(am.String())
	return 0
}

func (c *AuthMethodInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *AuthMethodInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthMethodInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMethodInspectCommand) Synopsis() string {
	return "Show detailed information about a configured auth method"
}

func (c *AuthMethodInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint auth-method inspect NAME

  Show detailed information about a configured auth method.

`)
}

const outInspectToken = `
The token auth method has no additional configuration. The token auth
method is the most fundamental auth method in Waypoint and can't be disabled.
This auth method accepts an API token for authentication.
`
