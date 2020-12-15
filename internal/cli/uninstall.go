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

	platform      string
	snapshotPath  string
	skipSnapshot  bool
	flagConfirm   bool
	deleteContext bool
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

	sg := c.ui.StepGroup()
	defer sg.Wait()

	// Pre-install work
	// - name the context we'll be uninstalling
	// - generate a snapshot of the current install
	s := sg.Add("")
	contextDefault, err := c.contextStorage.Default()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	s.Update("Default Waypoint server detected as context %q", contextDefault)
	s.Status(terminal.StatusWarn)
	s.Done()
	s = sg.Add("")
	s.Update("Uninstalling Waypoint server using context %q...", contextDefault)
	s.Done()

	s = sg.Add("")
	// Generate a snapshot
	if !c.skipSnapshot {
		s.Update("Generating server snapshot...")
		defer s.Abort()
		// sn := fmt.Sprintf("%s-%d", c.snapshotPath, time.Now().Unix())
		// generate snapshot
		// s.Update("Snapshot %q generated", sn")
	} else {
		s.Update("skip-snapshot set; not generating server snapshot")
		s.Status(terminal.StatusWarn)
	}
	s.Done()

	// Uninstall
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

	// Post-uninstall cleanup of context
	if c.deleteContext {
		if err := c.contextStorage.Delete(contextDefault); err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	c.ui.Output("Waypoint server successfully uninstalled for %s platform", c.platform, terminal.WithSuccessStyle())

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
		f.BoolVar(&flag.BoolVar{
			Name:    "confirm",
			Target:  &c.flagConfirm,
			Default: false,
			Usage:   "Confirm server uninstallation.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "delete-context",
			Target:  &c.deleteContext,
			Default: false,
			Usage:   "Delete the context for the server once it's uninstalled.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "platform",
			Target:  &c.platform,
			Default: "",
			Usage:   "Platform to uninstall the Waypoint server from.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "snapshot-path",
			Target:  &c.snapshotPath,
			Default: "",
			Usage:   "Path of the file to write the snapshot to.",
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
