package serverinstall

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

// Installer is implemented by the server platforms and is responsible for managing
// the installation of the Waypoint server.
type Installer interface {
	// Install expects the Waypoint server to be installed.
	Install(context.Context, *InstallOpts) (*InstallResults, error)

	// InstallRunner expects a Waypoint runner to be installed.
	InstallRunner(context.Context, *InstallRunnerOpts) error

	// InstallFlags is called prior to Install and allows the installer to
	// specify flags for the install CLI. The flags should be prefixed with
	// the platform name to avoid conflicts with other flags.
	InstallFlags(*flag.Set)
}

// InstallOpts are the options sent to Installer.Install.
type InstallOpts struct {
	Log hclog.Logger
	UI  terminal.UI
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
}

var Platforms = map[string]Installer{
	"kubernetes": &K8sInstaller{},
	"nomad":      &NomadInstaller{},
	"docker":     &DockerInstaller{},
}
