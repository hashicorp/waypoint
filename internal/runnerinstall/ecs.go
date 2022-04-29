package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ECSRunnerInstaller struct {
	config ecsConfig
}

func (i *ECSRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *ECSRunnerInstaller) InstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

func (i *ECSRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *ECSRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

type ecsConfig struct {
}
