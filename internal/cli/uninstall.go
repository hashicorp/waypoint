package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/serverinstall"
)

type UninstallCommand struct {
	*baseCommand

	platform     string
	contextName  string
	snapshotName string
	skipSnapshot bool
	flagConfirm  bool
}

func (c *UninstallCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("install")
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	if !c.flagConfirm {
		c.ui.Output(strings.TrimSpace(confirmReqMsg), terminal.WithErrorStyle())
		return 1
	}

	var err error

	// Generate a snapshot
	if !c.skipSnapshot {
		// sn := fmt.Sprintf("%s-%d", c.snapshotName, time.Now().Unix())
		// generate snapshot
	}

	p, ok := serverinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		c.ui.Output(
			"Error uninstalling server from %s: invalid platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	err = p.Uninstall(ctx, c.ui, log)
	if err != nil {
		// point to current docs on manual server cleanup
		c.ui.Output(
			"Error uninstalling server from %s: %s", c.platform, clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)

		return 1
	}

	// Verify clean state; remove old context

	return 0
}

func (c *UninstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UninstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UninstallCommand) Synopsis() string {
	return "Uninstall the Waypoint server"
}

func (c *UninstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint server uninstall [options]
	Uninstall the Waypoint server.

` + c.Flags().Help())
}

func (c *UninstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "context-name",
			Target: &c.contextName,
			Usage:  "Context of the Waypoint server to uninstall.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "confirm",
			Target:  &c.flagConfirm,
			Default: false,
			Usage:   "Confirm server uninstallation.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "platform",
			Target:  &c.platform,
			Default: "",
			Usage:   "Platform to uninstall the Waypoint server from.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "snapshot-name",
			Target:  &c.snapshotName,
			Default: "",
			Usage:   "Platform to uninstall the Waypoint server from.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "skip-snapshot",
			Target:  &c.skipSnapshot,
			Default: false,
			Usage:   "Skip creating a snapshot of the Waypoint server.",
		})
	})
}

var (
	confirmReqMsg = strings.TrimSpace(`
Uninstalling Waypoint server requires confirmation. 
Rerun the command with ‘-confirm’ to continue with the uninstall.
`)
)
