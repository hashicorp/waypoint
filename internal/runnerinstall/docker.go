package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type DockerRunnerInstaller struct {
}

func (d DockerRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (d DockerRunnerInstaller) InstallFlags(set *flag.Set) {
	//TODO implement me
}

func (d DockerRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (d DockerRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}
