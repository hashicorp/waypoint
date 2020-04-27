package cli

import (
	"fmt"
	"os"

	"github.com/hashicorp/waypoint/sdk/terminal"
)

type UpCommand struct {
	*baseCommand
}

func (c *UpCommand) Run([]string) int {
	ctx := c.Ctx
	log := c.Log.Named("up")

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

	// Build
	fmt.Fprintf(os.Stdout, "==> Building\n")
	build, err := app.Build(ctx)
	if err != nil {
		log.Error("error running builder", "error", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "==> Pushing artifact\n")
	pushedArtifact, err := app.PushBuild(ctx, build)
	if err != nil {
		log.Error("error pushing artifact to registry", "error", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "==> Deploying\n")
	_, err = app.Deploy(ctx, pushedArtifact)
	if err != nil {
		log.Error("error deploying", "error", err)
		return 1
	}

	return 0
}

func (c *UpCommand) Synopsis() string {
	return ""
}

func (c *UpCommand) Help() string {
	return ""
}
