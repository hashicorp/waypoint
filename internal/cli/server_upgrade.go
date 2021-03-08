package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/clisnapshot"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverinstall"
)

type ServerUpgradeCommand struct {
	*baseCommand

	platform     string
	contextName  string
	snapshotName string
	flagSnapshot bool
	confirm      bool
}

func (c *ServerUpgradeCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("upgrade")
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

	// Error handling from input

	if !c.confirm {
		c.ui.Output(confirmReqMsg, terminal.WithErrorStyle())
		return 1
	}

	if c.platform == "" {
		c.ui.Output(
			"A platform is required and must match the server context",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	p, ok := serverinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		c.ui.Output(
			"Error upgrading server on %s: unsupported platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	// Finish error handling

	// Get Server config to preserve existing configurations from context
	var ctxName string
	if c.contextName != "" {
		ctxName = c.contextName
	} else {
		defaultName, err := c.contextStorage.Default()
		if err != nil {
			c.ui.Output(
				"Error getting default context: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		ctxName = defaultName
	}

	originalCfg, err := c.contextStorage.Load(ctxName)
	if err != nil {
		c.ui.Output(
			"Error loading the context %q: %s",
			ctxName,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Upgrade waypoint server
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating server context: %q", ctxName)
	defer func() { s.Abort() }()

	conn, err := serverclient.Connect(ctx, serverclient.FromContextConfig(originalCfg))
	if err != nil {
		c.ui.Output(
			"Error connecting with context %q: %s",
			ctxName,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	s.Update("Verifying connection is valid for context %q...", ctxName)

	client := pb.NewWaypointClient(conn)
	// validate API compat here with new clientpkg
	if _, err := clientpkg.New(ctx,
		clientpkg.WithLogger(c.Log),
		clientpkg.WithClient(client),
	); err != nil {
		c.ui.Output(
			"Error connecting with context %q: %s",
			ctxName,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	resp, err := client.GetVersionInfo(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error retrieving server version info: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle())
		return 1
	}

	initServerVersion := resp.Info.Version

	s.Update("Context %q validated and connected successfully.", ctxName)
	s.Done()

	s = sg.Add("Starting server snapshots")

	// Snapshot server before upgrade
	if c.flagSnapshot {
		s.Update("Taking server snapshot before upgrading")

		snapshotName := c.snapshotName
		if c.snapshotName == defaultSnapshotName {
			// Append timestamps on default snapshot names
			snapshotName = fmt.Sprintf("%s-%d", c.snapshotName, time.Now().Unix())
		}

		s.Update("Taking snapshot of server with name: '%s'", snapshotName)
		writer, err := os.Create(snapshotName)
		if err != nil {
			s.Update("Failed to take server snapshot")
			s.Status(terminal.StatusError)
			s.Done()

			c.ui.Output(fmt.Sprintf("Error opening output: %s", err), terminal.WithErrorStyle())
			os.Remove(snapshotName)
			return 1
		}

		err = clisnapshot.WriteSnapshot(c.Ctx, c.project.Client(), writer)
		writer.Close()

		if err != nil {
			s.Update("Failed to take server snapshot\n")
			s.Status(terminal.StatusError)
			s.Done()

			if status.Code(err) == codes.Unimplemented {
				c.ui.Output(snapshotUnimplementedErr, terminal.WithErrorStyle())
			}

			c.ui.Output(fmt.Sprintf("Error generating Snapshot: %s", err), terminal.WithErrorStyle())
			os.Remove(snapshotName)
			return 1
		}

		s.Update("Snapshot of server written to: '%s'", snapshotName)
		s.Done()
	} else {
		s.Update("Server snapshot disabled on request, this means no snapshot will be taken before upgrades")
		s.Status(terminal.StatusWarn)
		s.Done()
		log.Warn("Server snapshot disabled on request from user, skipping")
	}

	c.ui.Output("Upgrading...", terminal.WithHeaderStyle())

	c.ui.Output("Waypoint server will now upgrade from version %q",
		initServerVersion, terminal.WithInfoStyle())

	installOpts := &serverinstall.InstallOpts{
		Log: log,
		UI:  c.ui,
	}

	// Upgrade in place
	result, err := p.Upgrade(ctx, installOpts, originalCfg.Server)
	if err != nil {
		c.ui.Output(
			"Error upgrading server on %s: %s", c.platform, clierrors.Humanize(err),
			terminal.WithErrorStyle())

		c.ui.Output(upgradeFailHelp)

		return 1
	}

	contextConfig := result.Context
	advertiseAddr := result.AdvertiseAddr
	httpAddr := result.HTTPAddr

	// We update the context config if the server addr has changed between upgrades
	if originalCfg.Server.Address != contextConfig.Server.Address {
		originalCfg.Server.Address = contextConfig.Server.Address

		if err := c.contextStorage.Set(ctxName, originalCfg); err != nil {
			c.ui.Output(
				"Error setting the CLI context: %s\n\n%s",
				clierrors.Humanize(err),
				errInstallRunning,
				terminal.WithErrorStyle(),
			)
			return 1
		}

		c.ui.Output("Server address has changed after upgrade. This client will "+
			"update its context with the new address, however any other clients "+
			"using this server must manually update their server address listed below "+
			"or find the address from `waypoint context list` which lists context %q address",
			ctxName,
			terminal.WithWarningStyle())
	}

	// Connect
	c.ui.Output("Verifying upgrade...", terminal.WithHeaderStyle())

	// New stepgroup to ensure output is after upgrade output
	sg2 := c.ui.StepGroup()
	defer sg2.Wait()

	s2 := sg2.Add("Client attempting to connect to server...")
	defer func() { s2.Abort() }()

	conn, err = serverclient.Connect(ctx,
		serverclient.FromContextConfig(originalCfg),
		serverclient.Timeout(3*time.Minute),
	)
	if err != nil {
		s2.Update("Client failed to connect to server")
		s2.Status(terminal.StatusError)
		s2.Done()

		c.ui.Output(
			"Error connecting to server: %s\n\n%s",
			clierrors.Humanize(err),
			"Check the waypoint server container logs for more information on "+
				"why it could have failed.",
			terminal.WithErrorStyle(),
		)

		return 1
	}
	client = pb.NewWaypointClient(conn)

	resp, err = client.GetVersionInfo(ctx, &empty.Empty{})
	if err != nil {
		s2.Update("Client failed to connect to server")
		s2.Status(terminal.StatusError)
		s2.Done()

		c.ui.Output(
			"Error retrieving server version info: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle())
		return 1
	}

	s2.Update("Server connection verified!")
	s2.Done()

	// Upgrade the runner
	if code := c.upgradeRunner(
		ctx, client, p, installOpts, advertiseAddr,
	); code > 0 {
		return code
	}

	c.ui.Output("\nServer upgrade for platform %q context %q complete!",
		c.platform, ctxName, terminal.WithSuccessStyle())

	c.ui.Output("Waypoint has finished upgrading the server to version %q\n",
		resp.Info.Version, terminal.WithSuccessStyle())

	c.ui.Output(addrSuccess, advertiseAddr.Addr, "https://"+httpAddr,
		terminal.WithSuccessStyle())

	return 0
}

func (c *ServerUpgradeCommand) upgradeRunner(
	ctx context.Context,
	client pb.WaypointClient,
	p serverinstall.Installer,
	installOpts *serverinstall.InstallOpts,
	advertiseAddr *pb.ServerConfig_AdvertiseAddr,
) int {
	// Connect
	c.ui.Output("Upgrading runner if required...", terminal.WithHeaderStyle())

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("")
	defer func() { s.Abort() }()

	// Upgrade the runner
	s.Update("Checking if a runner needs to be upgraded...")
	hasRunner, err := p.HasRunner(ctx, installOpts)
	if err != nil {
		s.Update("Error checking for runner: %s", err)
		s.Status(terminal.StatusError)
		s.Done()
		return 1
	}

	if !hasRunner {
		s.Update("No runners to upgrade.")
		s.Done()
		return 0
	}

	s.Update("Runner found. Uninstalling previous runner...")
	if err := p.UninstallRunner(ctx, installOpts); err != nil {
		c.ui.Output(
			"Error uninstalling runner from %s: %s\n\n"+
				"The runner will not be upgraded.",
			c.platform,
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)

		return 1
	}
	s.Update("Previous runner uninstalled")
	s.Done()

	// TODO(mitchellh): This creates a new auth token for the new runner.
	// In the future, we need to invalidate the old token. We don't have
	// the functionality to do this today.
	return installRunner(ctx, installOpts.Log, client, c.ui, p, advertiseAddr)
}

func (c *ServerUpgradeCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "auto-approve",
			Target:  &c.confirm,
			Default: false,
			Usage:   "Confirm server upgrade.",
		})
		f.StringVar(&flag.StringVar{
			Name:    "context-name",
			Target:  &c.contextName,
			Default: "",
			Usage:   "Waypoint server context to upgrade.",
		})
		f.StringVar(&flag.StringVar{
			Name:    "platform",
			Target:  &c.platform,
			Default: "",
			Usage:   "Platform to upgrade the Waypoint server from.",
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
			Usage:   "Enable or disable taking a snapshot of Waypoint server prior to upgrades.",
		})

		for name, platform := range serverinstall.Platforms {
			platformSet := set.NewSet(name + " Options")
			platform.UpgradeFlags(platformSet)
		}
	})
}

func (c *ServerUpgradeCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerUpgradeCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerUpgradeCommand) Synopsis() string {
	return "Upgrades Waypoint server in the current context to the latest version"
}

func (c *ServerUpgradeCommand) Help() string {
	return formatHelp(`
Usage: waypoint server upgrade [options]

  Upgrade Waypoint server in the current context to the latest version or the
  server image version specified. By default, Waypoint will upgrade to server
  version "hashicorp/waypoint:latest". Before upgrading, a snapshot of the
  server will be taken in case of any upgrade failures.

  If a runner was installed via "waypoint install" then that runner will also
  be upgraded to the latest version after the server is upgraded. Any other
  manually installed runners will not be automatically upgraded.

` + c.Flags().Help())
}

var (
	defaultSnapshotName = "waypoint-server-snapshot"
	confirmReqMsg       = strings.TrimSpace(`
Upgrading Waypoint server requires confirmation.
Rerun the command with '-auto-approve' to continue with the upgrade.
`)
	upgradeFailHelp = strings.TrimSpace(`
Upgrading Waypoint server has failed. To restore from a snapshot, use the command:

waypoint server restore [snapshot-name]

Where 'snapshot-name' is the name of the snapshot taken prior to the upgrade.

More information can be found by runninng 'waypoint server restore -help' or
following the server maintenence guide for backups and restores:
https://www.waypointproject.io/docs/server/run/maintenance#backup-restore
`)
	addrSuccess = strings.TrimSpace(`
Advertise Address: %[1]s
   Web UI Address: %[2]s
`)
)
