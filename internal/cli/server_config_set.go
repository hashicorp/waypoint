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

type ServerConfigSetCommand struct {
	*baseCommand

	flagAdvertiseAddr pb.ServerConfig_AdvertiseAddr
}

func (c *ServerConfigSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	c.ui.Output(
		"Modifying server configuration with the following settings:",
		terminal.WithHeaderStyle(),
	)

	addr := c.flagAdvertiseAddr.Addr
	if addr == "" {
		addr = "<empty>"
	}

	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name:  "advertise-addr",
			Value: addr,
		},
		{
			Name:  "advertise-tls",
			Value: c.flagAdvertiseAddr.Tls,
		},
		{
			Name:  "advertise-tls-skip-verify",
			Value: c.flagAdvertiseAddr.TlsSkipVerify,
		},
	})

	cfg := &pb.ServerConfig{
		AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
			&c.flagAdvertiseAddr,
		},
	}

	client := c.project.Client()
	_, err := client.SetServerConfig(c.Ctx, &pb.SetServerConfigRequest{
		Config: cfg,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Server configuration set!", terminal.WithSuccessStyle())
	return 0
}

func (c *ServerConfigSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "advertise-addr",
			Target: &c.flagAdvertiseAddr.Addr,
			Usage: "Address to advertise for the server. This is used by the entrypoints\n" +
				"binaries to communicate back to the server. If this is blank, then\n" +
				"the entrypoints will not communicate to the server. Features such as\n" +
				"logs, exec, etc. will not work.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls",
			Target:  &c.flagAdvertiseAddr.Tls,
			Usage:   "If true, the advertised address should be connected to with TLS.",
			Default: true,
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls-skip-verify",
			Target:  &c.flagAdvertiseAddr.TlsSkipVerify,
			Usage:   "Do not verify the TLS certificate presented by the server.",
			Default: false,
		})
	})
}

func (c *ServerConfigSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerConfigSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerConfigSetCommand) Synopsis() string {
	return "Set the server online configuration"
}

func (c *ServerConfigSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint server config-set [options]

  Set the online configuration for a running Waypoint server.

  The configuration that can be set here is different from the configuration
  given via the startup file. This configuration is persisted in the server
  database.

  Each flag represents a setting and all settings are transmitted to the server
  on submission. To correctly set the configuration, provide all flags
  together in one call.

` + c.Flags().Help())
}
