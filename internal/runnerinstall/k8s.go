package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"helm.sh/helm/v3/pkg/cli"
)

type K8sRunnerInstaller struct {
	config k8sConfig
}

func (i *K8sRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	// Initialize Helm settings
	settings := cli.New()

	// Get our K8S API
	// TODO: Get kubeconfig path & context
	_, ns, rc, err := k8s.Clientset(i.config.KubeconfigPath, i.config.Context)
	if err != nil {
		return err
	}

	panic("implement me")
}

func (i *K8sRunnerInstaller) InstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

func (i *K8sRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *K8sRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

type k8sConfig struct {
	KubeconfigPath string
	Context        string
}

// Use as base - suffix will be ID
const (
	serviceName                  = "waypoint"
	runnerRoleBindingName        = "waypoint-runner-rolebinding"
	runnerClusterRoleName        = "waypoint-runner"
	runnerClusterRoleBindingName = "waypoint-runner"
)
