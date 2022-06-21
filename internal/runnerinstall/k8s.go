package runnerinstall

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	helminstallutil "github.com/hashicorp/waypoint/internal/installutil/helm"
	k8sinstallutil "github.com/hashicorp/waypoint/internal/installutil/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/mitchellh/mapstructure"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sRunnerInstaller struct {
	k8sinstallutil.K8sInstaller
	Config K8sConfig
}

type K8sConfig struct {
	KubeconfigPath       string `hcl:"kubeconfig,optional"`
	K8sContext           string `hcl:"context,optional"`
	Version              string `hcl:"version,optional"`
	Namespace            string `hcl:"namespace,optional"`
	RunnerImage          string `hcl:"runner_image,optional"`
	RunnerImageTag       string `hcl:"runner_image_tag,optional"`
	CpuRequest           string `hcl:"runner_cpu_request,optional"`
	MemRequest           string `hcl:"runner_mem_request,optional"`
	CreateServiceAccount bool   `hcl:"odr_service_account_init,optional"`
}

const (
	defaultRunnerMemory   = "256Mi"
	defaultRunnerCPU      = "250m"
	defaultRunnerImageTag = "latest"
)

type InstalledRunnerConfig struct {
	Id string `mapstructure:"id"`
}

func (i *K8sRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	// Initialize Helm settings
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Getting Helm configs...")
	defer func() { s.Abort() }()
	settings, err := helminstallutil.SettingsInit()
	if err != nil {
		return err
	}
	s.Update("Helm settings retrieved")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Getting Helm action configuration...")
	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeconfigPath, i.Config.K8sContext)
	if err != nil {
		return err
	}
	s.Update("Helm action initialized")
	s.Status(terminal.StatusOK)
	s.Done()

	chartNS := ""
	if v := i.Config.Namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

	// This setup for Helm install matches the setup for the Helm platform plugin
	s = sg.Add("Creating new Helm install object...")
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
	client.ReleaseName = "waypoint-" + strings.ToLower(opts.Id)
	client.GenerateName = false
	client.NameTemplate = ""
	client.OutputDir = ""
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = true
	client.DisableOpenAPIValidation = false
	client.Replace = false
	client.Description = "Static runner for executing remote operations for Hashicorp Waypoint."
	client.CreateNamespace = true
	s.Update("Helm install created")
	s.Status(terminal.StatusOK)
	s.Done()

	version := i.Config.Version
	if version == "" {
		version = helminstallutil.DefaultHelmChartVersion
	}

	s = sg.Add("Locating chart...")
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/v"+version+".tar.gz", settings)
	if err != nil {
		return err
	}
	s.Update("Helm chart located")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Loading Helm chart...")
	c, err := loader.Load(path)
	if err != nil {
		return err
	}
	s.Update("Helm chart loaded")
	s.Status(terminal.StatusOK)
	s.Done()

	var memory, cpu, image, tag string
	if i.Config.MemRequest == "" {
		memory = defaultRunnerMemory
	} else {
		memory = i.Config.MemRequest
	}
	if i.Config.CpuRequest == "" {
		cpu = defaultRunnerCPU
	} else {
		cpu = i.Config.CpuRequest
	}
	if i.Config.RunnerImage == "" {
		image = defaultRunnerImage
	} else {
		image = i.Config.RunnerImage
	}
	if i.Config.RunnerImageTag == "" {
		tag = defaultRunnerImageTag
	} else {
		tag = i.Config.RunnerImageTag
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": false,
		},
		"runner": map[string]interface{}{
			"id": opts.Id,
			"image": map[string]interface{}{
				"repository": image,
				"tag":        tag,
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"memory": memory,
					"cpu":    cpu,
				},
			},
			"server": map[string]interface{}{
				"addr":          opts.ServerAddr,
				"tls":           opts.AdvertiseClient.Tls,
				"tlsSkipVerify": opts.AdvertiseClient.TlsSkipVerify,
				"cookie":        opts.Cookie,
			},
			"serviceAccount": map[string]interface{}{
				"create": i.Config.CreateServiceAccount,
			},

			"pullPolicy": "always",
		},
	}
	s = sg.Add("Installing Waypoint Helm chart with runner options: " + c.Name())
	_, err = client.RunWithContext(ctx, c, values)
	if err != nil {
		return err
	}
	s.Update("Waypoint runner installed with Helm!")
	s.Status(terminal.StatusOK)
	s.Done()

	return nil
}

func (i *K8sRunnerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-Config-path",
		Usage:  "Path to the kubeconfig file to use,",
		Target: &i.Config.KubeconfigPath,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.Config.K8sContext,
		Usage: "The Kubernetes context to install the Waypoint runner to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-helm-version",
		Target: &i.Config.Version,
		Usage:  "The version of the Helm chart to use for the Waypoint runner install.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Default: "default",
		Usage: "The namespace in the Kubernetes cluster into which the Waypoint " +
			"runner will be installed.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-runner-image",
		Target:  &i.Config.RunnerImage,
		Default: defaultRunnerImage,
		Usage:   "Docker image for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-runner-image-tag",
		Target:  &i.Config.RunnerImageTag,
		Default: defaultRunnerImageTag,
		Usage:   "Tag of the Docker image for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-request",
		Target:  &i.Config.CpuRequest,
		Default: defaultRunnerCPU,
		Usage:   "Requested amount of CPU for Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-request",
		Target:  &i.Config.MemRequest,
		Default: defaultRunnerMemory,
		Usage:   "Requested amount of memory for Waypoint runner.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.Config.CreateServiceAccount,
		Default: true,
		Usage:   "Create the service account if it does not exist. The default is true.",
	})
}

func (i *K8sRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Preparing Helm...")
	defer func() { s.Abort() }()

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeconfigPath, i.Config.K8sContext)
	if err != nil {
		return err
	}
	s.Update("Helm settings retrieved")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Uninstallation Pre-check...")
	helmRunnerId := "waypoint-" + strings.ToLower(opts.Id)
	verifyClient := action.NewGetValues(actionConfig)
	cfg, err := verifyClient.Run(helmRunnerId)
	if err != nil {
		return err
	}

	var runnerCfg InstalledRunnerConfig
	err = mapstructure.Decode(cfg["runner"], &runnerCfg)
	if err != nil {
		return err
	}

	// Check if the runner we are uninstalling matches the helm chart
	// This should always be true and is a sanity check to make sure this is a
	// proper runner installation and that we are uninstalling what we think we
	// should be uninstalling.
	if strings.ToLower(runnerCfg.Id) != strings.ToLower(opts.Id) {
		return errors.New("Runner not found")
	}
	s.Update("Runner %q found", opts.Id)
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Uninstalling Runner...")
	client := action.NewUninstall(actionConfig)
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.Timeout = 300 * time.Second
	client.Description = ""

	_, err = client.Run(helmRunnerId)
	if err != nil {
		return err
	}
	s.Update("Runner Uninstalled")
	s.Status(terminal.StatusOK)
	s.Done()

	// Delete left over runner persistent volume claim
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", helmRunnerId),
	}
	err = i.CleanPVC(ctx, opts.UI, opts.Log, listOptions)
	if err != nil {
		return err
	}

	return nil
}

func (i *K8sRunnerInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-config-path",
		Usage:  "Path to the kubeconfig file to use,",
		Target: &i.Config.KubeconfigPath,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.Config.K8sContext,
		Usage: "The Kubernetes context to install the Waypoint runner to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Default: "default",
		Usage: "The namespace in the Kubernetes cluster into which the Waypoint " +
			"runner will be installed.",
	})
}
