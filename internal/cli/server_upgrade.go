// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/builtin/k8s"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/clisnapshot"
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	"github.com/hashicorp/waypoint/internal/serverinstall"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
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
		WithNoLocalServer(),
	); err != nil {
		return 1
	}

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

	// Default the platform to the platform from the context
	if c.platform == "" {
		c.platform = originalCfg.Server.Platform
	}

	// Error handling from input

	if !c.confirm {
		proceed, err := c.ui.Input(&terminal.Input{
			Prompt: "Would you like to proceed with the Waypoint server upgrade? Only 'yes' will be accepted to approve: ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error upgrading server: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
		} else if strings.ToLower(proceed) != "yes" {
			c.ui.Output(upgradeConfirmMsg, terminal.WithErrorStyle())
			return 1
		}
		if c.platform == "" {
			c.ui.Output(platformReqMsg, terminal.WithErrorStyle())
			c.ui.Output(c.Help(), terminal.WithErrorStyle())
			return 1
		}
	} else if c.platform == "" {
		c.ui.Output(
			platformReqMsg,
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
		if err == nil {
			err = writer.Close()
		}
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

	// TODO(demophoon): Remove when we can handle automatic snapshot backup and
	// restore for kubernetes servers.
	if strings.ToLower(c.platform) == "kubernetes" {
		upgradeFrom, _ := version.NewVersion(initServerVersion)
		postHelmVersion, _ := version.NewVersion("v0.9.0")
		if upgradeFrom.LessThan(postHelmVersion) {
			c.ui.Output(upgradeToHelmRefused, terminal.WithErrorStyle())
			return 1
		}
	}

	c.ui.Output("Waypoint server will now upgrade from version %q",
		initServerVersion, terminal.WithInfoStyle())

	installOpts := &serverinstall.InstallOpts{
		Log:            log,
		UI:             c.ui,
		ServerRunFlags: c.args,
	}
	runnerOpts := &runnerinstall.InstallOpts{
		Log: log,
		UI:  c.ui,
		Id:  "static",
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
	if originalCfg.Server.Address != contextConfig.Server.Address ||
		originalCfg.Server.Platform != c.platform {
		// Update the platform here, basically to upgrade an older context that didn't
		// have platform set.
		originalCfg.Server.Platform = c.platform

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
		ctx, client, p, installOpts, runnerOpts, advertiseAddr,
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
	runnerOpts *runnerinstall.InstallOpts,
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

	s.Update("Runner found on Waypoint server. Uninstalling previous runner...")
	if err := p.UninstallRunner(ctx, runnerOpts); err != nil {
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

	if odc, ok := p.(installutil.OnDemandRunnerConfigProvider); ok {
		odr := odc.OnDemandRunnerConfig()

		runnerConfigName := odr.PluginType + "-bootstrap-profile"
		// We attempt to look up the default runner profile from the previous
		// installation. If we find it, we get the ID, so we can delete it after
		// the new runner profile is set up
		oldRunnerConfig, err := client.GetOnDemandRunnerConfig(ctx, &pb.GetOnDemandRunnerConfigRequest{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Name: runnerConfigName,
			}})

		if err != nil && status.Code(err) != codes.NotFound {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		} else if err != nil && status.Code(err) == codes.NotFound {
			c.ui.Output("Waypoint runner profile %q not found, creating new profile", runnerConfigName, terminal.WithWarningStyle())
		} else {
			ociUrl := odr.OciUrl
			if ociUrl == "" {
				ociUrl = installutil.DefaultODRImage
			}
			odr = &pb.OnDemandRunnerConfig{
				Id:                   oldRunnerConfig.Config.Id,
				Name:                 oldRunnerConfig.Config.Name,
				TargetRunner:         oldRunnerConfig.Config.TargetRunner,
				OciUrl:               ociUrl,
				EnvironmentVariables: oldRunnerConfig.Config.EnvironmentVariables,
				PluginType:           oldRunnerConfig.Config.PluginType,
				PluginConfig:         oldRunnerConfig.Config.PluginConfig,
				ConfigFormat:         oldRunnerConfig.Config.ConfigFormat,
				Default:              true,
			}
		}

		// Look at the existing on-demand runner configs and let the user know
		// they should not have multiple defaults
		// NOTE(briancain): A better way to handle this going forward is to enforce
		// a single default at the Waypoint state level. That should likely happen
		// soon.
		resp, err := client.ListOnDemandRunnerConfigs(ctx, &empty.Empty{})
		if err != nil {
			c.ui.Output("Failed to determine if default runner profiles are already set: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		if len(resp.Configs) > 0 {
			var (
				runnerUnsetStr     []string
				runnerDefaultNames []string
			)

			for _, cfg := range resp.Configs {
				if cfg.Name != runnerConfigName {
					runnerDefaultNames = append(runnerDefaultNames, fmt.Sprintf(runnerDefaultName, cfg.Name))
					runnerUnsetStr = append(runnerUnsetStr, fmt.Sprintf(runnerUnsetDefault, cfg.PluginType, cfg.Name))
				}

				// Checking if the default runner profile for a platform, pre-0.9, has the correct task launcher configs
				if c.platform == cfg.Name {
					switch c.platform {
					case "kubernetes":
						// attempt to parse the runner profile config into the correct task launcher config struct
						var result *k8s.TaskLauncherConfig
						// NOTE(briancain): This is here due to a k8s task plugin bug. When
						// we attempt to upgrade if we detect the previous mistake we warn
						// users that certain key values in their plugin config are wrong.
						if cfg.ConfigFormat == pb.Hcl_JSON {
							err = json.Unmarshal(cfg.PluginConfig, result)
							if err != nil {
								c.ui.Output(runnerProfileUpgradeConfigError, cfg.Name,
									cfg.Name, clierrors.Humanize(err), terminal.WithWarningStyle())
								var content map[string]interface{}
								err = json.Unmarshal(cfg.PluginConfig, &content)
								if err != nil {
									c.ui.Output("Error parsing plugin content: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
									return 1
								}
								if content["cpu"] != nil {
									cpuBody := content["cpu"].(map[string]interface{})
									if cpuBody["Requested"] != nil {
										c.ui.Output("The 'Requested' key specified for the CPU resources should instead be 'request'",
											terminal.WithWarningStyle())
									}
								}
								if content["memory"] != nil {
									memBody := content["memory"].(map[string]interface{})
									if memBody["Requested"] != nil {
										c.ui.Output("The 'Requested' key specified for the Memory resources should instead be 'request'",
											terminal.WithWarningStyle())
									}
								}
							}
						}
					default:
					}

				}
			}

			if len(runnerDefaultNames) > 0 {
				c.ui.Output("")
				c.ui.Output(runnerMultiDefault, strings.Join(runnerDefaultNames[:], "\n"), strings.Join(runnerUnsetStr[:], "\n"), terminal.WithWarningStyle())
				c.ui.Output("")
			}
		}

		// TODO(mitchellh): This creates a new auth token for the new runner.
		// In the future, we need to invalidate the old token. We don't have
		// the functionality to do this today.
		return installRunner(ctx, installOpts.Log, client, c.ui, p, advertiseAddr, odr, true, true)
	}
	return 0
}

func (c *ServerUpgradeCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "auto-approve",
			Target:  &c.confirm,
			Default: false,
			Usage:   "Auto-approve server upgrade. If unset, confirmation will be requested.",
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
			Usage: "Platform to upgrade the Waypoint server from, " +
				"defaults to the platform stored in the context.",
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

		// Add platforms in alphabetical order. A consistent order is important for repeatable doc generation.
		i := 0
		sortedPlatformNames := make([]string, len(serverinstall.Platforms))
		for name := range serverinstall.Platforms {
			sortedPlatformNames[i] = name
			i++
		}
		sort.Strings(sortedPlatformNames)

		for _, name := range sortedPlatformNames {
			platform := serverinstall.Platforms[name]
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
	upgradeConfirmMsg   = strings.TrimSpace(`
Upgrading Waypoint server requires confirmation.
`)
	platformReqMsg = strings.TrimSpace(`
A platform is required and must match the server context.
Rerun the command with '-platform=' and include the platform of the context to
upgrade.
`)
	upgradeFailHelp = strings.TrimSpace(`
Upgrading Waypoint server has failed. To restore from a snapshot, use the command:

waypoint server restore [snapshot-name]

Where 'snapshot-name' is the name of the snapshot taken prior to the upgrade.

More information can be found by running 'waypoint server restore -help' or
following the server maintenance guide for backups and restores:
https://www.waypointproject.io/docs/server/run/maintenance#backup-restore
`)
	addrSuccess = strings.TrimSpace(`
Advertise Address: %[1]s
   Web UI Address: %[2]s
`)

	runnerUnsetDefault = strings.TrimSpace(`
waypoint runner profile set -default=false -plugin-type=%[1]s -name=%[2]s
`)

	runnerDefaultName = "=> %[1]s"

	runnerMultiDefault = strings.TrimSpace(`
Waypoint expects only one runner profile to be a default. During the upgrade,
we have detected that there are multiple default runner profiles. This can
cause issues with launching on-demand runner tasks. The following profile names
have been set to be a default runner profile:

%[1]s

Please run the following commands if you wish to unset these runner profiles
from being the default:

%[2]s
`)

	upgradeToHelmRefused = strings.TrimSpace(`
Upgrading directly to 0.9.0 is not currently supported via this method on Kubernetes.

You can manually perform the upgrade by taking a snapshot of your Waypoint
server and restoring the snapshot to a fresh install of Waypoint Server using
the commands:

waypoint server snapshot
waypoint install -platform=kubernetes
waypoint server restore [snapshot-name]
`)

	runnerProfileUpgradeConfigError = strings.TrimSpace(`
The plugin config for runner profile %[1]s failed to correctly be parsed. 
Plugin Config Error: %[3]s

Please run the following with the corrected plugin configuration to fix this.

waypoint runner profile set -name=%[2]s -plugin-config=<path_to_config_file>

Starting in v0.9.0, ODR plugin configurations are validated for correctness.
The previous configuration for this profile is invalid, and ODRs will not launch
unless it is updated.
`)
)
