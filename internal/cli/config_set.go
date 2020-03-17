package cli

import (
	"fmt"
	"os"

	"github.com/mitchellh/devflow/sdk/terminal"
)

type ConfigSetCommand struct {
	*baseCommand
}

func (c *ConfigSetCommand) Run(args []string) int {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "config-set requires 2 arguments: a variable name and it's value")
		return 1
	}

	ctx := c.Ctx
	log := c.Log.Named("exec")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(); err != nil {
		return 1
	}

	cfg := c.cfg
	proj := c.project

	// NOTE(mitchellh): temporary restriction
	if len(cfg.Apps) != 1 {
		proj.UI.Output("only one app is supported at this time", terminal.WithErrorStyle())
		return 1
	}

	// Get our app
	app, err := proj.App(cfg.Apps[0].Name)
	if err != nil {
		c.logError(c.Log, "failed to initialize app", err)
		return 1
	}

	err = app.ConfigSet(ctx, args[0], args[1])
	if err != nil {
		log.Error("error exec", "error", err)
		return 1
	}

	return 0
}

func (c *ConfigSetCommand) Synopsis() string {
	return ""
}

func (c *ConfigSetCommand) Help() string {
	return ""
}
