package runnerinstall

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

const (
	defaultRunnerImage = "hashicorp/waypoint"
	runnerName         = "waypoint-runner"
)

type RunnerInstaller interface {
	// Install expects a Waypoint Runner to be installed
	Install(context.Context, *InstallOpts) error

	// InstallFlags is called prior to Install and allows the installer to
	// specify flags for the install CLI. The flags should be prefixed with
	// the platform name to avoid conflicts with other flags.
	InstallFlags(*flag.Set)

	// Uninstall should remove the runner(s) installed via Install.
	Uninstall(context.Context, *InstallOpts) error

	// UninstallFlags is called prior to Uninstall and allows the Uninstaller to
	// specify flags for the uninstall CLI. The flags should be prefixed with the
	// platform name to avoid conflicts with other flags.
	UninstallFlags(*flag.Set)
}

// InstallOpts are the options sent to RunnerInstaller.Install.
type InstallOpts struct {
	Log hclog.Logger
	UI  terminal.UI

	// Cookie is the server cookie that can be used for this runner
	Cookie string

	// ServerAddr is the address of the server to which the runner
	// connects
	ServerAddr string

	// AdvertiseClient is the serverconfig.Client information for connecting
	// to the server via the AdvertiseAddr information. This also has the auth
	// token already set. This is provided as a convenience since it is common
	// to build this immediately.
	AdvertiseClient *serverconfig.Client

	// Unique ID for the runner.
	Id string

	// TODO: Description
	RunnerRunFlags []string
}

var Platforms = map[string]RunnerInstaller{
	"ecs":        &ECSRunnerInstaller{},
	"kubernetes": &K8sRunnerInstaller{},
	"nomad":      &NomadRunnerInstaller{},
	"docker":     &DockerRunnerInstaller{},
}

func defaultRunnerName(id string) string {
	return "waypoint-" + id + "-runner"
}

const (
	DefaultRunnerTagName = "waypoint-runner"
)
