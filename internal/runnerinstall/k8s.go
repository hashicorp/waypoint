package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint/internal/installutil/helm"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"time"
)

type K8sRunnerInstaller struct {
	config k8sConfig
}

func (i *K8sRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	// Initialize Helm settings
	opts.Log.Debug("Getting settings for Helm.")
	settings, err := helm.SettingsInit()
	if err != nil {
		return err
	}

	opts.Log.Debug("Getting path options for Helm chart.")
	actionConfig, err := helm.ActionInit(opts.Log, i.config.KubeconfigPath, i.config.K8sContext)
	if err != nil {
		return err
	}

	chartNS := ""
	if v := i.config.Namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

	opts.Log.Debug("Creating new install action client.")
	client := action.NewInstall(actionConfig)
	client.ClientOnly = false
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.WaitForJobs = false
	client.Devel = true
	client.DependencyUpdate = false
	client.Timeout = 300 * time.Second
	client.Namespace = chartNS
	client.ReleaseName = "waypoint"
	client.GenerateName = false
	client.NameTemplate = ""
	client.OutputDir = ""
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = true
	client.DisableOpenAPIValidation = false
	client.Replace = false
	client.Description = ""
	client.CreateNamespace = true

	opts.Log.Debug("Locating chart.")
	// TODO: Add support for targeting specific versions of the chart, but default to latest
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/v0.1.8.tar.gz", settings)
	if err != nil {
		return err
	}

	opts.Log.Debug("Locating chart.")
	c, err := loader.Load(path)
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": false,
		},
		"runner": map[string]interface{}{
			"id": opts.Id,
			"server": map[string]interface{}{
				"addr":        opts.ServerAddr,
				"cookie":      opts.Cookie,
				"tokenSecret": "",
			},
			"image": map[string]interface{}{
				"repository": i.config.RunnerImage,
				"tag":        i.config.RunnerImageTag,
			},
			"pullPolicy": "always",
		},
	}
	opts.Log.Debug("Installing Waypoint Helm chart: " + c.Name())
	_, err = client.RunWithContext(ctx, c, values)
	if err != nil {
		return err
	}

	return nil
}

func (i *K8sRunnerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-config-path",
		Usage:  "Path to the kubeconfig file to use,",
		Target: &i.config.KubeconfigPath,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.config.K8sContext,
		Usage: "The Kubernetes context to install the Waypoint runner to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-helm-version",
		Target: &i.config.Version,
		Usage:  "The version of the Helm chart to use for the Waypoint runner install.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-namespace",
		Target: &i.config.Namespace,
		Usage: "The namespace in the Kubernetes cluster into which the Waypoint " +
			"runner will be installed.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-runner-image",
		Target:  &i.config.RunnerImage,
		Default: "hashicorp/waypoint",
		Usage:   "Docker image for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-runner-image-tag",
		Target:  &i.config.RunnerImageTag,
		Default: "latest",
		Usage:   "Tag of the Docker image for the Waypoint runner.",
	})
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
	K8sContext     string
	Version        string
	Namespace      string
	RunnerImage    string
	RunnerImageTag string
}

// Use as base - suffix will be ID
const (
	serviceName                  = "waypoint"
	runnerRoleBindingName        = "waypoint-runner-rolebinding"
	runnerClusterRoleName        = "waypoint-runner"
	runnerClusterRoleBindingName = "waypoint-runner"
)
