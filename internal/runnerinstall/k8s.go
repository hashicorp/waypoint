package runnerinstall

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/installutil/helm"
	"github.com/hashicorp/waypoint/internal/installutil/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/mitchellh/mapstructure"
)

type K8sRunnerInstaller struct {
	k8s.K8sInstaller
	config k8sConfig
}

type InstalledRunnerConfig struct {
	Id string `mapstructure:"id"`
	Server map[string]string `mapstructure:"server"`
}

func (i *K8sRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	// Initialize Helm settings
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Getting Helm configs...")
	defer func() { s.Abort() }()
	settings, err := helm.SettingsInit()
	if err != nil {
		return err
	}
	s.Update("Helm settings retrieved")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Getting Helm action configuration...")
	actionConfig, err := helm.ActionInit(opts.Log, i.config.KubeconfigPath, i.config.K8sContext)
	if err != nil {
		return err
	}
	s.Update("Helm action initialized")
	s.Status(terminal.StatusOK)
	s.Done()

	chartNS := ""
	if v := i.config.Namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

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
	client.Description = ""
	client.CreateNamespace = true
	s.Update("Helm install created")
	s.Status(terminal.StatusOK)
	s.Done()

	var version string
	if i.config.Version == "" {
		version = defaultHelmChartVersion
	} else {
		version = i.config.Version
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
	if i.config.MemRequest == "" {
		memory = defaultRunnerMemory
	} else {
		memory = i.config.MemRequest
	}
	if i.config.CpuRequest == "" {
		cpu = defaultRunnerCPU
	} else {
		cpu = i.config.CpuRequest
	}
	if i.config.RunnerImage == "" {
		image = defaultRunnerImage
	} else {
		image = i.config.RunnerImage
	}
	if i.config.RunnerImageTag == "" {
		tag = defaultRunnerImageTag
	} else {
		tag = i.config.RunnerImageTag
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
				"addr":   opts.ServerAddr,
				"cookie": opts.Cookie,
			},
			"serviceAccount": map[string]interface{}{
				"create": i.config.CreateServiceAccount,
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
		Default: defaultRunnerImage,
		Usage:   "Docker image for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-runner-image-tag",
		Target:  &i.config.RunnerImageTag,
		Default: defaultRunnerImageTag,
		Usage:   "Tag of the Docker image for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-request",
		Target:  &i.config.CpuRequest,
		Default: defaultRunnerCPU,
		Usage:   "Requested amount of CPU for Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-request",
		Target:  &i.config.MemRequest,
		Default: defaultRunnerMemory,
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
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Preparing Helm...")
	defer func() { s.Abort() }()

	actionConfig, err := helm.ActionInit(opts.Log, i.config.KubeconfigPath, i.config.K8sContext)
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
		opts.Log.Debug("Could not find runner for id %q", opts.Id)
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

	if ! i.config.KeepPVC {
		// Delete left over runner persistent volume claim
		s = sg.Add("Gathering Persistent Volume Claim...")
		clientset, err := i.NewClient()
		if err != nil {
			return err
		}
		pvClient := clientset.CoreV1().PersistentVolumeClaims(i.config.Namespace)
		listOptions := metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", helmRunnerId),
		}
		if list, err := pvClient.List(ctx, listOptions); err != nil {
			return err
		} else if len(list.Items) > 0 {
			s.Update("Deleting Persistent Volume Claim...")

			// Add watcher for waiting for persistent volume clean up
			w, err := pvClient.Watch(ctx, listOptions)
			if err != nil {
				return err
			}

			// Delete the PVCs
			if err = pvClient.DeleteCollection(
				ctx,
				metav1.DeleteOptions{},
				listOptions,
			); err != nil {
				return err
			}

			// Wait until the persistent volumes are cleaned up
			err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
				select {
				case wCh := <-w.ResultChan():
					if wCh.Type == "DELETED" {
						w.Stop()
						return true, nil
					}
					return false, nil
				default:
					return false, nil
				}
			})
			if err != nil {
				return err
			}
		}
		s.Update("Persistent volume claims cleaned up")
		s.Status(terminal.StatusOK)
		s.Done()
	}

	return nil
}

func (i *K8sRunnerInstaller) UninstallFlags(set *flag.Set) {
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
		Name:   "k8s-namespace",
		Target: &i.config.Namespace,
		Usage: "The namespace in the Kubernetes cluster into which the Waypoint " +
			"runner will be installed.",
	})

	set.BoolVar(&flag.BoolVar{
		Name: "k8s-keep-pvc",
		Target: &i.config.KeepPVC,
		Usage: "The Persistent Volume Claim associated with the runner will be " +
			"cleaned up if left unset.",
	})

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
	KeepPVC              bool
}

const (
	defaultHelmChartVersion = "0.1.8"
	defaultRunnerMemory     = "256Mi"
	defaultRunnerCPU        = "250m"
	defaultRunnerImageTag   = "latest"
)
