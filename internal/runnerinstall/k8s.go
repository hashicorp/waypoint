package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type K8sRunnerInstaller struct {
	config k8sConfig
}

func (i *K8sRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (i *K8sRunnerInstaller) InstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

func (i *K8sRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (i *K8sRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

type k8sConfig struct {
}

// Use as base - suffix will be ID
const (
	serviceName                  = "waypoint"
	runnerRoleBindingName        = "waypoint-runner-rolebinding"
	runnerClusterRoleName        = "waypoint-runner"
	runnerClusterRoleBindingName = "waypoint-runner"
)
