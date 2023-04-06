// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strings"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/posener/complete"
)

type ContextHelpCommand struct {
	*baseCommand

	SynopsisText string
	HelpText     string
}

func (c *ContextHelpCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
		WithNoClient(),
	); err != nil {
		return 1
	}

	return cli.RunResultHelp
}

func (c *ContextHelpCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextHelpCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextHelpCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextHelpCommand) Synopsis() string {
	return strings.TrimSpace(c.SynopsisText)
}

func (c *ContextHelpCommand) Help() string {
	homePath := c.homeConfigPath

	dir, err := homedir.Dir()
	if err == nil {
		if strings.HasPrefix(homePath, dir) {
			homePath = "~" + homePath[len(dir):]
		}
	}
	str := fmt.Sprintf("%s\nContext Info:\n  Config Path: %s\n", c.HelpText, homePath)
	return formatHelp(str)
}
