// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runnerinstall

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	dockerparser "github.com/novln/docker-parser"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	helminstallutil "github.com/hashicorp/waypoint/internal/installutil/helm"
	k8sinstallutil "github.com/hashicorp/waypoint/internal/installutil/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type K8sRunnerInstaller struct {
	Config k8sinstallutil.K8sConfig
}

const (
	defaultRunnerMemory          = "256Mi"
	defaultRunnerCPU             = "250m"
	defaultOdrServiceAccountName = "waypoint-runner-odr"
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
	settings, err := helminstallutil.SettingsInit(i.Config.Namespace)
	if err != nil {
		opts.UI.Output("Unable to retrieve Helm configuration.", terminal.WithErrorStyle())
		return err
	}

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeconfigPath, i.Config.K8sContext, i.Config.Namespace)
	if err != nil {
		opts.UI.Output("Unable to initialize Helm.", terminal.WithErrorStyle())
		return err
	}

	// If all else fails, default the namespace to "default"
	if i.Config.Namespace == "" {
		i.Config.Namespace = "default"
	}

	// This setup for Helm install matches the setup for the Helm platform plugin
	client := action.NewInstall(actionConfig)
	client.ClientOnly = false
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.WaitForJobs = false
	client.Devel = true
	client.DependencyUpdate = false
	client.Timeout = 300 * time.Second
	client.Namespace = i.Config.Namespace
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

	version := i.Config.Version
	if version == "" {
		tags, err := helminstallutil.GetLatestHelmChartVersion(ctx)
		if err != nil {
			opts.UI.Output("Error getting latest tag of Waypoint helm chart.", terminal.WithErrorStyle())
			return err
		}
		version = *tags[0].Name
	}

	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/"+version+".tar.gz", settings)
	if err != nil {
		opts.UI.Output("Unable to locate Waypoint helm chart.", terminal.WithErrorStyle())
		return err
	}

	c, err := loader.Load(path)
	if err != nil {
		opts.UI.Output("Unable to load Waypoint helm chart.", terminal.WithErrorStyle())
		return err
	}

	runnerImageRef, err := dockerparser.Parse(i.Config.RunnerImage)
	if err != nil {
		opts.UI.Output("Error parsing runner image name: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	clientSet, err := k8sinstallutil.NewClient(i.Config)
	if err != nil {
		opts.UI.Output("Error creating k8s clientset: %s", clierrors.Humanize(err), terminal.StatusError)
		return err
	}
	// Determine if we need to make a service account
	if i.Config.CreateServiceAccount {
		saClient := clientSet.CoreV1().ServiceAccounts(i.Config.Namespace)
		_, err = saClient.Get(ctx, defaultRunnerTagName, metav1.GetOptions{})
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				err = nil
			} else {
				opts.UI.Output(fmt.Sprintf("Error getting service account: %s", clierrors.Humanize(err)), terminal.StatusError)
				return err
			}
		} else {
			opts.UI.Output("Waypoint runner service account already exists - a new service account will not be created",
				terminal.WithInfoStyle())
			i.Config.CreateServiceAccount = false
		}
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": false,
		},
		"runner": map[string]interface{}{
			"agentArgs": opts.RunnerAgentFlags,
			"id":        opts.Id,
			"image": map[string]interface{}{
				"repository": runnerImageRef.Repository(),
				"tag":        runnerImageRef.Tag(),
				"pullPolicy": i.Config.ImagePullPolicy,
			},
			"odr": map[string]interface{}{
				// odr image stanza not specified - this is used by the helm chart to
				// give to the bootstrap job to populate the ODR profile, but only
				// during a server install. For runner installs, we'll create the
				// ODR profile ourselves later.
				"serviceAccount": map[string]interface{}{
					"create": i.Config.CreateServiceAccount,
					"name":   defaultOdrServiceAccountName,
				},
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"memory": i.Config.MemRequest,
					"cpu":    i.Config.CpuRequest,
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
				"name":   defaultRunnerTagName,
			},

			"pullPolicy": "always",
		},
	}

	s.Update("Installing Waypoint Helm chart with runner options: " + c.Name())
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
		Usage:  "Path to the kubeconfig file to use.",
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
		Usage: "The version of the Helm chart to use for the Waypoint runner install. " +
			"The required version number format is: 'vX.Y.Z'.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Default: "default",
		Usage: "The namespace in the Kubernetes cluster into which the Waypoint " +
			"runner will be installed.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-runner-image",
		Target: &i.Config.RunnerImage,
		// This is the static (non-odr) runner, and therefore needs to use the non-ODR
		// image. The server and the static runner use the same image.
		Default: installutil.DefaultServerImage,
		Usage:   "Docker image for the Waypoint runner.",
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
		Usage:   "Create the service account if it does not exist.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-image-pull-policy",
		Target:  &i.Config.ImagePullPolicy,
		Default: "",
		Usage:   "Set the pull policy for the Waypoint runner image",
	})
}

func (i *K8sRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI
	// Our checks here follow the logic of:
	// Up until v0.8.2, we installed runners with the k8s client,
	// and the Label was "app=waypoint-runner" and the Name "waypoint-runner-random-id"
	// As of 0.9.0, we install runners with helm, with a Label following the
	// pattern ("app.kubernetes.io/instance=waypoint-%s", runnerId)
	// and the Name ("waypoint-"+strings.ToLower(runnerId))
	//
	// Therefore we need to ascertain A) if the runner exists on the cluster at
	// all (it might not be if the user is auth'd to the wrong cluster), and then B)
	// if the name/label matches the Helm pattern or the k8s client pattern so
	// we know how to uninstall it

	// A) Is runner on the cluster at all?
	clientset, err := k8sinstallutil.NewClient(i.Config)
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

	// Search for a runner with 0.8.x tag format, installed with k8s client
	deploymentClient := clientset.AppsV1().Deployments(i.Config.Namespace)
	k8sClientList, err := deploymentClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", defaultRunnerTagName),
	})
	if err != nil {
		return fmt.Errorf("could not list deployments in namespace %q with current context: %s", i.Config.Namespace, err)
	}

	// Search for runner with 0.9+ tag format, installed with helm
	podClient := clientset.CoreV1().Pods(i.Config.Namespace)
	helmClientList, err := podClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=waypoint-%s", strings.ToLower(opts.Id)),
	})
	if err != nil {
		return fmt.Errorf("could not list pods in namespace %q with current context: %s", i.Config.Namespace, err)
	}

	// If both lists are empty, the runner is not here at all
	// Move to: B) Decide which uninstall path we use based on if there is a runner
	// with the naming patterns we get with our 0.9.0+ helm installer
	if len(k8sClientList.Items) == 0 && len(helmClientList.Items) == 0 {
		return fmt.Errorf("runner with ID %q not found in namespace %q with current context", opts.Id, i.Config.Namespace)
	} else if len(helmClientList.Items) > 0 {
		err = i.uninstallWithHelm(ctx, opts)
	} else {
		// Once we're here, we know that there is >0 K8S runners, and 0 Helm runners,
		// so we can proceed with the k8s uninstall
		// This should only include default runners installed on 0.8.x
		err = i.uninstallWithK8s(ctx, opts, k8sClientList)
	}
	return err
}

// Uninstall is a method of K8sInstaller and implements the Installer interface to
// remove a waypoint-server statefulset and the associated PVC and service from
// a Kubernetes cluster
func (i *K8sRunnerInstaller) uninstallWithK8s(ctx context.Context, opts *InstallOpts, listK8sClient *v1.DeploymentList) error {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer func() { s.Abort() }()

	clientset, err := k8sinstallutil.NewClient(i.Config)
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

	deploymentClient := clientset.AppsV1().Deployments(i.Config.Namespace)
	s.Update("Deleting any automatically installed runners...")

	// Record various settings we can reuse for runner reinstallation
	// if we're doing an upgrade. We need to do this because the upgrade
	// flags don't contain the installation settings, and we prefer them
	// not to; instead we just retain the old settings.
	//
	// Note we have lots of conditionals here to try to avoid weird
	// panic situations if the remote side doesn't have the fields we
	// expect.
	podSpec := listK8sClient.Items[0].Spec.Template.Spec
	if secrets := podSpec.ImagePullSecrets; len(secrets) > 0 {
		i.Config.ImagePullSecret = secrets[0].Name
	}
	if v := podSpec.Containers; len(v) > 0 {
		c := v[0]

		i.Config.ImagePullPolicy = string(c.ImagePullPolicy)
		if m := c.Resources.Requests; len(m) > 0 {
			if v, ok := m[apiv1.ResourceMemory]; ok {
				i.Config.MemRequest = v.String()
			}
			if v, ok := m[apiv1.ResourceCPU]; ok {
				i.Config.CpuRequest = v.String()
			}
		}
		if m := c.Resources.Limits; len(m) > 0 {
			if v, ok := m[apiv1.ResourceMemory]; ok {
				i.Config.MemLimit = v.String()
			}
			if v, ok := m[apiv1.ResourceCPU]; ok {
				i.Config.CpuLimit = v.String()
			}
		}
	}

	// create our wait channel to later poll for statefulset+pod deletion
	w, err := deploymentClient.Watch(
		ctx,
		metav1.ListOptions{
			LabelSelector: "app=" + defaultRunnerTagName,
		},
	)
	if err != nil {
		ui.Output(
			"Error creating deployments watcher %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err

	}
	// send DELETE to statefulset collection
	if err = deploymentClient.DeleteCollection(
		ctx,
		metav1.DeleteOptions{},
		metav1.ListOptions{
			LabelSelector: "app=" + defaultRunnerTagName,
		},
	); err != nil {
		ui.Output(
			"Error deleting Waypoint deployment: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}

	// wait for deletion to complete
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		select {
		case wCh := <-w.ResultChan():
			if wCh.Type == "DELETED" {
				w.Stop()
				return true, nil
			}
			log.Trace("deployment collection not fully removed, waiting")
			return false, nil
		default:
			log.Trace("no message received on watch.ResultChan(), waiting for Event")
			return false, nil
		}
	})
	if err != nil {
		return err
	}
	s.Update("Runner deployment deleted")
	s.Done()

	return nil
}

func (i *K8sRunnerInstaller) uninstallWithHelm(ctx context.Context, opts *InstallOpts) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Preparing Helm...")
	defer func() { s.Abort() }()

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeconfigPath, i.Config.K8sContext, i.Config.Namespace)
	if err != nil {
		s.Update("Unable to setup Helm.")
		s.Status(terminal.StatusError)
		return err
	}

	s.Update("Uninstallation Pre-check...")
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
		return fmt.Errorf("found runner with id %q does not match given id %q", runnerCfg.Id, opts.Id)
	}
	s.Update("Runner %q found; uninstalling runner...", opts.Id)
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
	s.Update("Runner %q uninstalled", opts.Id)
	s.Status(terminal.StatusOK)
	s.Done()

	// Delete left over runner persistent volume claim
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", helmRunnerId),
	}
	err = k8sinstallutil.CleanPVC(ctx, opts.UI, opts.Log, listOptions, i.Config)

	return err
}

func (i *K8sRunnerInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-config-path",
		Usage:  "Path to the kubeconfig file to use",
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

// OnDemandRunnerConfig implements OnDemandRunnerConfigProvider
func (i *K8sRunnerInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	// Generate some configuration
	cfgMap := map[string]interface{}{}
	if v := i.Config.ImagePullSecret; v != "" {
		cfgMap["image_secret"] = v
	}
	// TODO: Enable specification of service account name
	if v := i.Config.ImagePullPolicy; v != "" {
		cfgMap["image_pull_policy"] = v
	}

	var cpuConfig k8s.ResourceConfig
	var memConfig k8s.ResourceConfig
	if v := i.Config.CpuRequest; v != "" {
		cpuConfig.Request = v
	}
	if v := i.Config.MemRequest; v != "" {
		memConfig.Request = v
	}
	if v := i.Config.CpuLimit; v != "" {
		cpuConfig.Limit = v
	}
	if v := i.Config.MemLimit; v != "" {
		memConfig.Limit = v
	}
	cfgMap["cpu"] = cpuConfig
	cfgMap["memory"] = memConfig

	cfgMap["service_account"] = defaultOdrServiceAccountName

	// Marshal our config
	cfgJson, err := json.MarshalIndent(cfgMap, "", "\t")
	if err != nil {
		// This shouldn't happen cause we control our input. If it does,
		// just panic cause this will be in a `server install` CLI and
		// we want the user to report a bug.
		panic(err)
	}

	return &pb.OnDemandRunnerConfig{
		Name:         "kubernetes",
		PluginType:   "kubernetes",
		Default:      false,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
		// Can't use i.Config.OdrImage here, because it hasn't been initalized.
		// This doesn't matter in practice - this is used in the `runner install` command,
		// which has its own -odr-image flag that it's going to use to overwrite
		// this value.
		OciUrl: installutil.DefaultODRImage,
	}
}
