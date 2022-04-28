package runnerinstall

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type RunnerInstaller interface {
	// Install expects a Waypoint Runner to be installed
	Install(context.Context, *InstallOpts) (bool, error)

	// InstallFlags is called prior to Install and allows the installer to
	// specify flags for the install CLI. The flags should be prefixed with
	// the platform name to avoid conflicts with other flags.
	InstallFlags(*flag.Set)

	// Uninstall should remove the runner(s) installed via Install.
	Uninstall(context.Context, *InstallOpts) (bool, error)

	// UninstallFlags is called prior to Uninstall and allows the Uninstaller to
	// specify flags for the uninstall CLI. The flags should be prefixed with the
	// platform name to avoid conflicts with other flags.
	UninstallFlags(*flag.Set)
}

// InstallOpts are the options sent to RunnerInstaller.Install.
type InstallOpts struct {
	Log hclog.Logger
	UI  terminal.UI
}

type InstallResults struct {
	// Context is the connection context that can be used to connect from
	// the CLI to the server. This will be used to establish an API client.
	Context *clicontext.Config

	// ID is the uuid of the installed runner
	ID string
}

var Platforms = map[string]RunnerInstaller{
	"ecs":        &ECSRunnerInstaller{},
	"kubernetes": &K8sRunnerInstaller{},
	"nomad":      &NomadRunnerInstaller{},
}
