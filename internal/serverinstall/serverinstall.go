package serverinstall

import (
	"context"

	"github.com/hashicorp/waypoint/internal/runnerinstall"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

// Installer is implemented by the server platforms and is responsible for managing
// the installation of the Waypoint server.
type Installer interface {
	// HasRunner returns true if a runner is installed.
	HasRunner(context.Context, *InstallOpts) (bool, error)

	// Install expects the Waypoint server to be installed.
	// Returns InstallResults, a bootstrap token (if platform sets one up), or an error
	Install(context.Context, *InstallOpts) (*InstallResults, string, error)

	// InstallRunner expects a Waypoint runner to be installed.
	InstallRunner(context.Context, *runnerinstall.InstallOpts) error

	// InstallFlags is called prior to Install and allows the installer to
	// specify flags for the install CLI. The flags should be prefixed with
	// the platform name to avoid conflicts with other flags.
	InstallFlags(*flag.Set)

	// Upgrade expects the Waypoint server to be upgraded from a previous install.
	// After upgrading the server, this should also upgrade the primary
	// runner that was installed with InstallRunner, if it exists.
	Upgrade(ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (*InstallResults, error)

	// UpgradeFlags is called prior to Upgrade and allows the upgrader to
	// specify flags for the upgrade CLI. The flags should be prefixed with
	// the platform name to avoid conflicts with other flags.
	UpgradeFlags(*flag.Set)

	// Uninstall expects the Waypoint server to be uninstalled. This should
	// also look up to see if any runners exist (installed via InstallRunner)
	// and remove those as well. Runners manually installed outside of this
	// interface should not be touched.
	Uninstall(context.Context, *InstallOpts) error

	// UninstallRunner should remove the runner(s) installed via InstallRunner.
	//
	// No runners may exist. Runners installed manually by the user should be
	// ignored (i.e. InstallRunner should set some identifiers that can be used
	// to distinguish between automatically installed vs. manually installed).
	UninstallRunner(context.Context, *runnerinstall.InstallOpts) error

	// UninstallFlags is called prior to Uninstall and allows the Uninstaller to
	// specify flags for the uninstall CLI. The flags should be prefixed with the
	// platform name to avoid conflicts with other flags.
	UninstallFlags(*flag.Set)
}

// InstallOpts are the options sent to Installer.Install.
type InstallOpts struct {
	Log            hclog.Logger
	UI             terminal.UI
	ServerRunFlags []string
}

// InstallResults are the results expected for a successful Installer.Install.
type InstallResults struct {
	// Context is the connection context that can be used to connect from
	// the CLI to the server. This will be used to establish an API client.
	Context *clicontext.Config

	// AdvertiseAddr is the configuration for the advertised address
	// that entrypoints (deployed workloads) will use to communicate back
	// to the server. This may be different from the context info because this
	// may be a private address.
	AdvertiseAddr *pb.ServerConfig_AdvertiseAddr

	// HTTPAddr is the address to the HTTP listener on the server. This generally
	// is reachable from the CLI immediately and not a private address.
	HTTPAddr string
}

// InstallRunnerOpts are the options sent to Installer.InstallRunner.
type InstallRunnerOpts struct {
	Log hclog.Logger
	UI  terminal.UI

	// AuthToken is an auth token that can be used for this runner.
	AuthToken string

	// AdvertiseAddr is the advertised address configuration currently set
	// for the server. This is likely the same information you want to use
	// for the runner to connect to the server, but doesn't have to be.
	AdvertiseAddr *pb.ServerConfig_AdvertiseAddr

	// AdvertiseClient is the serverconfig.Client information for connecting
	// to the server via the AdvertiseAddr information. This also has the auth
	// token already set. This is provided as a convenience since it is common
	// to build this immediately.
	AdvertiseClient *serverconfig.Client

	Cookie string

	Id string
}

var Platforms = map[string]Installer{
	"ecs":        &ECSInstaller{},
	"kubernetes": &K8sInstaller{},
	"nomad":      &NomadInstaller{},
	"docker":     &DockerInstaller{},
}

const (
	serverName = "waypoint-server"
	runnerName = "waypoint-runner"
)

// Default server ports to use
var (
	// todo: remove these and just use the serverconfig constants.
	defaultGrpcPort = serverconfig.DefaultGRPCPort
	defaultHttpPort = serverconfig.DefaultHTTPPort
)
