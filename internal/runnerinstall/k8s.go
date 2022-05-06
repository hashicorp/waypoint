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
	client.ReleaseName = "waypoint-" + opts.Id
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

	var version string
	if i.config.Version == "" {
		version = defaultHelmChartVersion
	} else {
		version = i.config.Version
	}

	opts.Log.Debug("Locating chart.")
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/v"+version+".tar.gz", settings)
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
			"image": map[string]interface{}{
				"repository": i.config.RunnerImage,
				"tag":        i.config.RunnerImageTag,
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"memory": i.config.MemRequest,
					"cpu":    i.config.CpuRequest,
				},
			},
			"server": map[string]interface{}{
				"addr":   opts.ServerAddr,
				"cookie": opts.Cookie,
			},
			"serviceAccount": map[string]interface{}{
				"create": i.config.CreateServiceAccount,
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

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-request",
		Target:  &i.config.CpuRequest,
		Default: "250m",
		Usage:   "Requested amount of CPU for Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-request",
		Target:  &i.config.MemRequest,
		Default: "256Mi",
		Usage:   "Requested amount of memory for Waypoint runner.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.config.CreateServiceAccount,
		Default: true,
		Usage:   "Create the service account if it does not exist. The default is true.",
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
	KubeconfigPath       string
	K8sContext           string
	Version              string
	Namespace            string
	RunnerImage          string
	RunnerImageTag       string
	CpuRequest           string
	MemRequest           string
	CreateServiceAccount bool
}

const (
	defaultHelmChartVersion = "0.1.8"
)
