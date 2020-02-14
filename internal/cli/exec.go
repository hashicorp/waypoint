package cli

type ExecCommand struct {
	*baseCommand
}

func (c *ExecCommand) Run([]string) int {
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
		c.ui.Error("only one app is supported at this time")
		return 1
	}

	// Get our app
	app, err := proj.App(cfg.Apps[0].Name)
	if err != nil {
		c.logError(c.Log, "failed to initialize app", err)
		return 1
	}

	err = app.Exec(ctx)
	if err != nil {
		log.Error("error exec", "error", err)
		return 1
	}

	return 0
}

func (c *ExecCommand) Synopsis() string {
	return ""
}

func (c *ExecCommand) Help() string {
	return ""
}
