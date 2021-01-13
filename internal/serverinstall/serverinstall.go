package serverinstall

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Installer is implemented by the server platforms and is responsible for managing
// the installation of the Waypoint server.
type Installer interface {
	Install(ctx context.Context, ui terminal.UI, log hclog.Logger) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, error)
	InstallFlags(*flag.Set)

	// Uninstall expects the Waypoint server to be uninstalled.
	Uninstall(context.Context, *InstallOpts) error

	// UninstallFlags is called prior to Uninstall and allows the Uninstaller to
	// specify flags for the uninstall CLI. The flags should be prefixed with the
	// platform name to avoid conflicts with other flags.
	UninstallFlags(*flag.Set)
}

var Platforms = map[string]Installer{
	"kubernetes": &K8sInstaller{},
	"nomad":      &NomadInstaller{},
	"docker":     &DockerInstaller{},
}

const (
	serverName = "waypoint-server"
)
