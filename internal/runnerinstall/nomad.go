package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type NomadRunnerInstaller struct {
	config NomadConfig
}

func (i *NomadRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *NomadRunnerInstaller) InstallFlags(set *flag.Set) {
	//TODO implement me
}

func (i *NomadRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *NomadRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

type NomadConfig struct {
}
