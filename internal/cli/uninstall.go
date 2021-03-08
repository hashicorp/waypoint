package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/clisnapshot"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/serverinstall"
)

type UninstallCommand struct {
	*baseCommand

	platform          string
	snapshotName      string
	flagSnapshot      bool
	autoApprove       bool
	deleteContext     bool
	ignoreRunnerError bool
}

func (c *UninstallCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("uninstall")
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoAutoServer(),
	); err != nil {
		return 1
	}

	if !c.autoApprove {
		c.ui.Output(strings.TrimSpace(autoApproveMsg), terminal.WithErrorStyle())
		return 1
	}

	// output the context we'll be uninstalling
	contextDefault, err := c.contextStorage.Default()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	c.ui.Output(
		"Uninstalling Waypoint server with context %q",
		contextDefault,
		terminal.WithSuccessStyle(),
	)

	// Get the platform early so we can validate it.
	p, ok := serverinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		if c.platform == "" {
			c.ui.Output(
				"The -platform flag is required.",
				terminal.WithErrorStyle(),
			)

			return 1
		}

		c.ui.Output(
			"Error uninstalling server from %s: invalid platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	// Pre-uninstall work
	// - generate a snapshot of the current install
	s := sg.Add("")
	defer func() { s.Abort() }()

	// Generate a snapshot
	if c.flagSnapshot {
		s.Update("Generating server snapshot...")

		// set config snapshot name with default + timestamp or flag value
		snapshotName := c.snapshotName
		if c.snapshotName == defaultSnapshotName {
			// Append timestamps on default snapshot names
			snapshotName = fmt.Sprintf("%s-%d", c.snapshotName, time.Now().Unix())
		}

		// take the snapshot
		s.Update("Taking snapshot of server with name: '%s'", snapshotName)
		w, err := os.Create(snapshotName)
		if err != nil {
			s.Update("Failed to take server snapshot\n")
			s.Status(terminal.StatusError)
			s.Done()

			c.ui.Output("Error creating snapshot file: %s", err, terminal.WithErrorStyle())
			os.Remove(snapshotName)
			return 1
		}
		if err = clisnapshot.WriteSnapshot(ctx, c.project.Client(), w); err != nil {
			s.Update("Failed to take server snapshot\n")
			s.Status(terminal.StatusError)
			s.Done()

			if status.Code(err) == codes.Unimplemented {
				c.ui.Output(snapshotUnimplementedErr, terminal.WithErrorStyle())
			}

			c.ui.Output("Error generating snapshot: %s", err, terminal.WithErrorStyle())
			os.Remove(snapshotName)
			return 1
		}
		s.Update("Snapshot %q generated", snapshotName)
	} else {
		s.Update("skip-snapshot set; not generating server snapshot")
		s.Status(terminal.StatusWarn)
	}
	s.Done()

	installOpts := &serverinstall.InstallOpts{
		Log: log,
		UI:  c.ui,
	}

	// We first uninstall any runners.
	log.Trace("calling UninstallRunner")
	if err := p.UninstallRunner(ctx, installOpts); err != nil {
		if !c.ignoreRunnerError {
			c.ui.Output(
				"Error uninstalling runners from %s: %s\n\n"+
					"Server uninstallation aborted. You can force server uninstallation "+
					"with runner uninstallation failures by using the -ignore-runner-error. "+
					"Note that this may leave your runners dangling.\n\n"+
					"See Troubleshooting docs "+
					"for guidance on manual uninstall: https://www.waypointproject.io/docs/troubleshooting",
				c.platform,
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		c.ui.Output(
			"Error uninstalling runners from %s: %s\n\n"+
				"-ignore-runner-error is specified so this will be ignored.",
			c.platform,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
	}

	log.Trace("calling Uninstall")
	if err := p.Uninstall(ctx, installOpts); err != nil {
		c.ui.Output(
			"Error uninstalling server from %s: %s\nSee Troubleshooting docs "+
				"for guidance on manual uninstall: https://www.waypointproject.io/docs/troubleshooting",
			c.platform,
			clierrors.Humanize(err),
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

  Uninstall the Waypoint server. The platform should be specified as kubernetes,
  nomad, or docker. '-auto-approve' is required.

  By default, this command deletes the default server's context.

  This command does not destroy Waypoint resources, such as deployments and
  releases. Clear all workspaces prior to uninstall to prevent hanging resources.

  If a runner was installed via "waypoint install", the runner will also be
  uninstalled. Manually installed runners (outside of the "waypoint install"
  command) will not be affected.

` + c.Flags().Help())
}

func (c *UninstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "auto-approve",
			Target:  &c.autoApprove,
			Default: false,
			Usage:   "Auto-approve server uninstallation.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "delete-context",
			Target:  &c.deleteContext,
			Default: true,
			Usage:   "Delete the context for the server once it's uninstalled.",
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
			Default: defaultSnapshotName,
			Usage: "Filename to write the snapshot to. If no name is specified, by" +
				" default a timestamp will be appended to the default snapshot name.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "snapshot",
			Target:  &c.flagSnapshot,
			Default: true,
			Usage:   "Enable or disable taking a snapshot of Waypoint server prior to uninstall.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "ignore-runner-error",
			Target:  &c.ignoreRunnerError,
			Default: false,
			Usage: "Ignore any errors encountered while uninstalling runners. This allows " +
				"the server to be uninstalled even if runner uninstallation fails. Note that " +
				"this may leave runners dangling since future 'uninstall' runs will do nothing if " +
				"the server is uninstalled.",
		})

		for name, platform := range serverinstall.Platforms {
			platformSet := set.NewSet(name + " Options")
			platform.UninstallFlags(platformSet)
		}
	})
}

var (
	autoApproveMsg = strings.TrimSpace(`
Uninstalling Waypoint server requires approval.
Rerun the command with -auto-approve to continue with the uninstall.
`)
)
