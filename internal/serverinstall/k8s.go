package serverinstall

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	dockerparser "github.com/novln/docker-parser"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	helminstallutil "github.com/hashicorp/waypoint/internal/installutil/helm"
	k8sinstallutil "github.com/hashicorp/waypoint/internal/installutil/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

//
type K8sInstaller struct {
	Config k8sinstallutil.K8sConfig
	// Config k8sConfig
}

const (
	serviceName                  = "waypoint"
	runnerRoleBindingName        = "waypoint-runner-rolebinding"
	runnerClusterRoleName        = "waypoint-runner"
	runnerClusterRoleBindingName = "waypoint-runner"
)

// newClient creates a new K8S client based on the configured settings.
func (i *K8sInstaller) newClient() (*kubernetes.Clientset, error) {
	// Build our K8S client.
	configOverrides := &clientcmd.ConfigOverrides{}
	if i.Config.K8sContext != "" {
		configOverrides = &clientcmd.ConfigOverrides{
			CurrentContext: i.Config.K8sContext,
		}
	}
	newCmdConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		configOverrides,
	)

	// Discover the current target namespace in the user's config so if they
	// run kubectl commands waypoint will show up. If we use the default namespace
	// they might not see the objects we've created.
	if i.Config.Namespace == "" {
		namespace, _, err := newCmdConfig.Namespace()
		if err != nil {
			return nil, fmt.Errorf(
				"Error getting namespace from client config: %s",
				clierrors.Humanize(err),
			)
		}

		i.Config.Namespace = namespace
	}

	clientconfig, err := newCmdConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf(
			"Error initializing kubernetes client: %s",
			clierrors.Humanize(err),
		)
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		return nil, fmt.Errorf(
			"Error initializing kubernetes client: %s",
			clierrors.Humanize(err),
		)
	}

	return clientset, nil
}

// Install is a method of K8sInstaller and implements the Installer interface to
// register a waypoint-server in a Kubernetes cluster
func (i *K8sInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, string, error) {
	if i.Config.OdrImage == "" {
		var err error
		i.Config.OdrImage, err = installutil.DefaultODRImage(i.Config.ServerImage)
		if err != nil {
			return nil, "", err
		}
	}

	log := opts.Log
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	clientset, err := k8sinstallutil.NewClient(i.Config)
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return nil, "", err
	}

	s := sg.Add("Inspecting Kubernetes cluster")
	// If this is kind, then we want to warn the user that they need
	// to have some loadbalancer system setup or this will not work.
	_, err = clientset.AppsV1().DaemonSets("kube-system").Get(
		ctx, "kindnet", metav1.GetOptions{})
	isKind := err == nil
	if isKind {
		s.Update(warnK8SKind)
		s.Status(terminal.StatusWarn)
		s.Done()
		s = sg.Add("")
	}

	s.Update("Getting Helm configs...")
	defer func() { s.Abort() }()
	settings, err := helminstallutil.SettingsInit()
	if err != nil {
		return nil, "", err
	}

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeConfigPath, i.Config.K8sContext)
	if err != nil {
		return nil, "", err
	}

	if i.Config.Namespace == "" {
		// If all else fails, default the namespace to "default"
		i.Config.Namespace = "default"
	}

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

	var version string
	if i.Config.Version == "" {
		tags, err := helminstallutil.GetLatestHelmChartVersion(ctx)
		if err != nil {
			opts.UI.Output("Error getting latest tag of Waypoint helm chart.", terminal.WithErrorStyle())
			return nil, "", err
		}
		version = *tags[0].Name
	} else {
		version = i.Config.Version
	}
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/"+version+".tar.gz", settings)
	if err != nil {
		return nil, "", err
	}

	c, err := loader.Load(path)
	if err != nil {
		return nil, "", err
	}

	imageRef, err := dockerparser.Parse(i.Config.ServerImage)
	if err != nil {
		ui.Output("Error parsing server image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, "", err
	}

	odrImageRef, err := dockerparser.Parse(i.Config.OdrImage)
	if err != nil {
		ui.Output("Error parsing on-demand runner image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, "", err
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": true,
			"image": map[string]interface{}{
				"repository": imageRef.Repository(),
				"tag":        imageRef.Tag(),
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"memory": i.Config.MemRequest,
					"cpu":    i.Config.CpuRequest,
				},
				"limits": map[string]interface{}{
					"memory": i.Config.MemLimit,
					"cpu":    i.Config.CpuLimit,
				},
			},
		},
		"runner": map[string]interface{}{
			"enabled": false,
			"image": map[string]interface{}{
				"repository": odrImageRef.Repository(),
				"tag":        odrImageRef.Tag(),
			},
			"odr": map[string]interface{}{
				"image": map[string]interface{}{
					"repository": odrImageRef.Repository(),
					"tag":        odrImageRef.Tag(),
				},
			},
		},
	}

	s.Update("Installing Waypoint Helm chart...")
	_, err = client.RunWithContext(ctx, c, values)
	if err != nil {
		return nil, "", err
	}

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	// TODO: Move this to a util function for install and upgrade to use
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		clientset, err := k8sinstallutil.NewClient(i.Config)
		if err != nil {
			return false, err
		}

		s.Update("Getting waypoint-ui service...")
		svc, err := clientset.CoreV1().Services(i.Config.Namespace).Get(
			ctx, "waypoint-ui", metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		ingress := svc.Status.LoadBalancer.Ingress
		if len(ingress) == 0 {
			log.Trace("ingress list is empty, waiting")
			return false, nil
		}

		addr := ingress[0].IP
		if addr == "" {
			addr = ingress[0].Hostname
		}

		// No address, still not ready
		if addr == "" {
			log.Trace("address is empty, waiting")
			return false, nil
		}

		// Get the ports
		var grpcPort int32
		var httpPort int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				grpcPort = spec.Port
			}

			if spec.Name == "https-2" {
				httpPort = spec.Port
			}

			if httpPort != 0 && grpcPort != 0 {
				break
			}
		}
		if grpcPort == 0 || httpPort == 0 {
			// If we didn't find the port, retry...
			log.Trace("no port found on service, retrying")
			return false, nil
		}

		// Set the grpc address
		grpcAddr = fmt.Sprintf("%s:%d", addr, grpcPort)
		log.Info("server service ready", "addr", addr)

		// HTTP address to return
		httpAddr = fmt.Sprintf("%s:%d", addr, httpPort)

		// Ensure the service is ready to use before returning
		s.Update("Checking that the server service is ready...")
		_, err = net.DialTimeout("tcp", httpAddr, 1*time.Second)
		if err != nil {
			// Depending on the platform, this can take a long time. On EKS, it's by far the longest step. Adding an explicit message helps
			s.Update("Service %q exists and is configured, but isn't yet accepting incoming connections. Waiting...", serviceName)
			return false, nil
		}

		s.Update("Service %q is ready", serviceName)
		log.Info("http server ready", "httpAddr", addr)

		// Set our advertise address
		advertiseAddr.Addr = grpcAddr
		advertiseAddr.Tls = true
		advertiseAddr.TlsSkipVerify = true

		// If we want internal or we're a localhost address, we use the internal
		// address. The "localhost" check is specifically for Docker for Desktop
		// since pods can't reach this.
		if i.Config.AdvertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				serverName,
				grpcPort,
			)
		}

		// Set our connection information
		contextConfig = clicontext.Config{
			Server: serverconfig.Client{
				Address:       grpcAddr,
				Tls:           true,
				TlsSkipVerify: true, // always for now
				Platform:      "kubernetes",
			},
		}

		return true, nil
	})
	if err != nil {
		return nil, "", err
	}
	s.Done()

	s = sg.Add("Waiting for the bootstrap process to finish...")
	var bootJobName string
	err = wait.PollImmediate(2*time.Second, 5*time.Minute, func() (bool, error) {
		// the label we use for LabelSelector is set here
		// https://github.com/hashicorp/waypoint-helm/blob/d2f6de6e9010b94da84f37eeaca4a8190a439060/templates/bootstrap-job.yaml#L8
		jobs, err := clientset.BatchV1().Jobs(i.Config.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/instance=waypoint",
		})
		if err != nil {
			return false, nil
		}
		// NOTE(krantzinator): the job we are searching for is prefixed with `waypoint-bootstrap`
		// per our Helm chart; if that naming ever changes, this will also need to be updated
		// https://github.com/hashicorp/waypoint-helm/blob/d2f6de6e9010b94da84f37eeaca4a8190a439060/templates/bootstrap-job.yaml#L5
		jobPrefix := "waypoint-bootstrap-"
		var bootJob *batchv1.Job
		for _, j := range jobs.Items {
			if strings.Contains(j.Name, jobPrefix) {
				bootJob = &j
				bootJobName = j.Name
				break
			}
		}
		if bootJob.Status.Succeeded == 1 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		log.Error("no bootstrap job found, returning", "job_name", bootJobName, "err", err)
		s.Update("No bootstrap job found")
		s.Status(terminal.WarningStyle)
		s.Done()
		return nil, "", err
	}

	secretClient := clientset.CoreV1().Secrets(i.Config.Namespace)

	// Get the secret
	// TODO(briancain): make a flag for this like the cli/login.go command has
	// TODO(briancain): Extract the cli.loginK8S method outside of the CLI package?
	secret, err := secretClient.Get(ctx, "waypoint-server-token", metav1.GetOptions{})
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, "", err
	}

	// Get our token
	tokenB64 := secret.Data["token"]
	if len(tokenB64) == 0 {
		return nil, "", errors.New("failed to read token secret from response")
	}
	bootstrapToken := string(tokenB64)

	s.Update("Server bootstrap complete!")
	s.Done()
	s = sg.Add("")

	s.Update("Waypoint server installed with Helm!")
	s.Status(terminal.StatusOK)
	s.Done()

	return &InstallResults{
			Context:       &contextConfig,
			AdvertiseAddr: &advertiseAddr,
			HTTPAddr:      httpAddr,
		},
		bootstrapToken, nil
}

// Upgrade is a method of K8sInstaller and implements the Installer interface to
// upgrade a waypoint-server in a Kubernetes cluster
func (i *K8sInstaller) Upgrade(
	ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (
	*InstallResults, error,
) {
	if i.Config.OdrImage == "" {
		var err error
		i.Config.OdrImage, err = installutil.DefaultODRImage(i.Config.ServerImage)
		if err != nil {
			return nil, err
		}
	}

	log := opts.Log
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Getting Helm configs...")
	defer func() { s.Abort() }()
	settings, err := helminstallutil.SettingsInit()
	if err != nil {
		return nil, err
	}

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeConfigPath, i.Config.K8sContext)
	if err != nil {
		return nil, err
	}

	if i.Config.Namespace == "" {
		// If all else fails, default the namespace to "default"
		i.Config.Namespace = "default"
	}

	client := action.NewUpgrade(actionConfig)
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.WaitForJobs = false
	client.Devel = true
	client.DependencyUpdate = false
	client.Timeout = 300 * time.Second
	client.Namespace = i.Config.Namespace
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = true
	client.DisableOpenAPIValidation = false
	client.Description = ""

	var version string
	if i.Config.Version == "" {
		tags, err := helminstallutil.GetLatestHelmChartVersion(ctx)
		if err != nil {
			opts.UI.Output("Error getting latest tag of Waypoint helm chart.", terminal.WithErrorStyle())
			return nil, err
		}
		version = *tags[0].Name
	} else {
		version = i.Config.Version
	}

	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/"+version+".tar.gz", settings)
	if err != nil {
		return nil, err
	}

	c, err := loader.Load(path)
	if err != nil {
		return nil, err
	}

	imageRef, err := dockerparser.Parse(i.Config.ServerImage)
	if err != nil {
		ui.Output("Error parsing server image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, err
	}

	odrImageRef, err := dockerparser.Parse(i.Config.OdrImage)
	if err != nil {
		ui.Output("Error parsing on-demand runner image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, err
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": true,
			"image": map[string]interface{}{
				"repository": imageRef.Repository(),
				"tag":        imageRef.Tag(),
			},
		},
		"runner": map[string]interface{}{
			"enabled": false,
			"image": map[string]interface{}{
				"repository": odrImageRef.Repository(),
				"tag":        odrImageRef.Tag(),
			},
			"odr": map[string]interface{}{
				"image": map[string]interface{}{
					"repository": odrImageRef.Repository(),
					"tag":        odrImageRef.Tag(),
				},
			},
		},
	}

	s.Update("Installing Waypoint Helm chart...")
	_, err = client.RunWithContext(ctx, "waypoint", c, values)
	if err != nil {
		return nil, err
	}

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		clientset, err := k8sinstallutil.NewClient(i.Config)
		if err != nil {
			return false, err
		}

		s.Update("Getting waypoint-ui service...")
		svc, err := clientset.CoreV1().Services(i.Config.Namespace).Get(
			ctx, "waypoint-ui", metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		ingress := svc.Status.LoadBalancer.Ingress
		if len(ingress) == 0 {
			log.Trace("ingress list is empty, waiting")
			return false, nil
		}

		addr := ingress[0].IP
		if addr == "" {
			addr = ingress[0].Hostname
		}

		// No address, still not ready
		if addr == "" {
			log.Trace("address is empty, waiting")
			return false, nil
		}

		// Get the ports
		var grpcPort int32
		var httpPort int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				grpcPort = spec.Port
			}

			if spec.Name == "https-2" {
				httpPort = spec.Port
			}

			if httpPort != 0 && grpcPort != 0 {
				break
			}
		}
		if grpcPort == 0 || httpPort == 0 {
			// If we didn't find the port, retry...
			log.Trace("no port found on service, retrying")
			return false, nil
		}

		// Set the grpc address
		grpcAddr = fmt.Sprintf("%s:%d", addr, grpcPort)
		log.Info("server service ready", "addr", addr)

		// HTTP address to return
		httpAddr = fmt.Sprintf("%s:%d", addr, httpPort)

		// Ensure the service is ready to use before returning
		s.Update("Checking that the server service is ready...")
		_, err = net.DialTimeout("tcp", httpAddr, 1*time.Second)
		if err != nil {
			// Depending on the platform, this can take a long time. On EKS, it's by far the longest step. Adding an explicit message helps
			s.Update("Service %q exists and is configured, but isn't yet accepting incoming connections. Waiting...", serviceName)
			return false, nil
		}

		s.Update("Service %q is ready", serviceName)
		s.Status(terminal.StatusOK)
		s.Done()
		log.Info("http server ready", "httpAddr", addr)

		// Set our advertise address
		advertiseAddr.Addr = grpcAddr
		advertiseAddr.Tls = true
		advertiseAddr.TlsSkipVerify = true

		// If we want internal or we're a localhost address, we use the internal
		// address. The "localhost" check is specifically for Docker for Desktop
		// since pods can't reach this.
		if i.Config.AdvertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				serverName,
				grpcPort,
			)
		}

		// Set our connection information
		contextConfig = clicontext.Config{
			Server: serverconfig.Client{
				Address:       grpcAddr,
				Tls:           true,
				TlsSkipVerify: true, // always for now
				Platform:      "kubernetes",
			},
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	s = sg.Add("Upgrade complete!")
	s.Done()

	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Uninstall is a method of K8sInstaller and implements the Installer interface to
// remove a waypoint-server statefulset and the associated PVC and service from
// a Kubernetes cluster
func (i *K8sInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Getting Helm action configuration...")
	defer func() { s.Abort() }()

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.Config.KubeConfigPath, i.Config.K8sContext)
	if err != nil {
		return err
	}
	s.Update("Helm action initialized, creating new Helm uninstall object...")

	chartNS := ""
	if v := i.Config.Namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

	client := action.NewUninstall(actionConfig)
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.Timeout = 300 * time.Second
	client.Description = ""
	s.Update("Helm uninstall created; uninstalling Helm chart...")

	_, err = client.Run("waypoint")
	if err != nil {
		return err
	}
	s.Update("Waypoint uninstalled with Helm!")
	s.Status(terminal.StatusOK)
	s.Done()

	// TODO: Clean-up waypoint server PVC
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s,component=server", "waypoint"),
	}
	err = k8sinstallutil.CleanPVC(ctx, ui, opts.Log, listOptions, i.Config)
	if err != nil {
		return err
	}
	// TODO: Delete runner (or all runners?)

	return nil
}

// InstallRunner implements Installer.
func (i *K8sInstaller) InstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.K8sRunnerInstaller{
		Config: k8sinstallutil.K8sConfig{
			K8sContext:           i.Config.K8sContext,
			Namespace:            i.Config.Namespace,
			RunnerImage:          i.Config.ServerImage,
			CpuRequest:           i.Config.CpuRequest,
			MemRequest:           i.Config.MemRequest,
			CreateServiceAccount: true,
			OdrImage:             i.Config.OdrImage,
		},
	}
	// parachute in case we remove the flag defaults one day
	if runnerInstaller.Config.Namespace == "" {
		runnerInstaller.Config.Namespace = "default"
	}
	err := runnerInstaller.Install(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// UninstallRunner implements Installer.
func (i *K8sInstaller) UninstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerUninstaller := runnerinstall.K8sRunnerInstaller{
		Config: k8sinstallutil.K8sConfig{
			KubeconfigPath: "",
			K8sContext:     i.Config.K8sContext,
			Namespace:      i.Config.Namespace,
		},
	}

	// parachute in case we remove the flag defaults one day
	if runnerUninstaller.Config.Namespace == "" {
		runnerUninstaller.Config.Namespace = "default"
	}

	err := runnerUninstaller.Uninstall(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// HasRunner implements Installer.
func (i *K8sInstaller) HasRunner(
	ctx context.Context,
	opts *InstallOpts,
) (bool, error) {
	clientset, err := k8sinstallutil.NewClient(i.Config)
	if err != nil {
		return false, err
	}

	deploymentClient := clientset.AppsV1().Deployments(i.Config.Namespace)
	list, err := deploymentClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", runnerName),
	})
	if err != nil {
		return false, err
	}

	return len(list.Items) > 0, nil
}

// OnDemandRunnerConfig implements OnDemandRunnerConfigProvider
func (i *K8sInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	// Generate some configuration
	cfgMap := map[string]interface{}{}
	if v := i.Config.ImagePullSecret; v != "" {
		cfgMap["image_secret"] = v
	}
	if v := i.Config.OdrServiceAccount; v != "" {
		cfgMap["service_account"] = v
	}
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
		OciUrl:       i.Config.OdrImage,
		PluginType:   "kubernetes",
		Default:      true,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
	}
}

// newServiceAccount takes in a k8sConfig and creates the ServiceAccount
// definition for the ODR.
func newServiceAccount(c k8sinstallutil.K8sConfig) (*apiv1.ServiceAccount, error) {
	return &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.OdrServiceAccount,
			Namespace: c.Namespace,
		},
	}, nil
}

// newServiceAccountClusterRoleWithBinding creates the cluster role and binding necessary to create and verify
// a nodeport type services.
func newServiceAccountClusterRoleWithBinding(c k8sinstallutil.K8sConfig) (*rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding, error) {
	return &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: runnerClusterRoleName,
			},
			Rules: []rbacv1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"get", "list"},
			}},
		}, &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: runnerClusterRoleBindingName,
			},

			// Our default runner role is just the default "edit" role. This
			// gives access to read/write most things in this namespace but
			// disallows modifying roles and rolebindings.
			RoleRef: rbacv1.RoleRef{
				APIGroup: "",
				Kind:     "ClusterRole",
				Name:     runnerClusterRoleName,
			},

			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      c.OdrServiceAccount,
					Namespace: c.Namespace,
				},
			},
		}, nil
}

// newServiceAccountRoleBinding creates the role binding necessary to
// map the ODR role to the service account.
func newServiceAccountRoleBinding(c k8sinstallutil.K8sConfig) (*rbacv1.RoleBinding, error) {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runnerRoleBindingName,
			Namespace: c.Namespace,
		},

		// Our default runner role is just the default "edit" role. This
		// gives access to read/write most things in this namespace but
		// disallows modifying roles and rolebindings.
		RoleRef: rbacv1.RoleRef{
			APIGroup: "",
			Kind:     "ClusterRole",
			Name:     "edit",
		},

		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.OdrServiceAccount,
				Namespace: c.Namespace,
			},
		},
	}, nil
}

func (i *K8sInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-config-path",
		Usage:  "Path to the kubeconfig file to use.",
		Target: &i.Config.KubeConfigPath,
	})

	set.BoolVar(&flag.BoolVar{
		Name:   "k8s-advertise-internal",
		Target: &i.Config.AdvertiseInternal,
		Usage: "Advertise the internal service address rather than the external. " +
			"This is useful if all your deployments will be able to access the private " +
			"service address. This will default to false but will be automatically set to " +
			"true if the external host is detected to be localhost.",
	})

	set.StringMapVar(&flag.StringMapVar{
		Name:   "k8s-annotate-service",
		Target: &i.Config.ServiceAnnotations,
		Usage:  "Annotations for the Service generated.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.Config.K8sContext,
		Usage: "The Kubernetes context to install the Waypoint server to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	// NOTE(briancain): We set the default for these values to 0. In the k8s API,
	// setting the limits and requests values to 0 is the same as not setting it all.
	// This is the expected behavior we'll want. If someone _does_ set a value using
	// these flags, they will be parsed and used instead.

	set.StringVar(&flag.StringVar{
		Name:   "k8s-helm-version",
		Target: &i.Config.Version,
		Usage: "The version of the Helm chart to use for the Waypoint runner install. " +
			"The required version number format is: 'vX.Y.Z'.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-request",
		Target:  &i.Config.CpuRequest,
		Usage:   "Configures the requested CPU amount for the Waypoint server in Kubernetes.",
		Default: "0",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-request",
		Target:  &i.Config.MemRequest,
		Usage:   "Configures the requested memory amount for the Waypoint server in Kubernetes.",
		Default: "0",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-limit",
		Target:  &i.Config.CpuLimit,
		Usage:   "Configures the CPU limit for the Waypoint server in Kubernetes.",
		Default: "0",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-limit",
		Target:  &i.Config.MemLimit,
		Usage:   "Configures the memory limit for the Waypoint server in Kubernetes.",
		Default: "0",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Usage:   "Namespace to install the Waypoint server into for Kubernetes.",
		Default: "default",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-pull-policy",
		Target:  &i.Config.ImagePullPolicy,
		Usage:   "Set the pull policy for the Waypoint server image.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-pull-secret",
		Target: &i.Config.ImagePullSecret,
		Usage:  "Secret to use to access the Waypoint server image on Kubernetes.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-secret-file",
		Target: &i.Config.SecretFile,
		Usage:  "Use the Kubernetes Secret in the given path to access the Waypoint server image.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-image",
		Target:  &i.Config.ServerImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-odr-image",
		Target: &i.Config.OdrImage,
		Usage:  "Docker image for the Waypoint On-Demand Runners",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-runner-service-account",
		Target: &i.Config.OdrServiceAccount,
		Usage: "Service account to assign to the on-demand runner. If this is blank, " +
			"a service account will be created automatically with the correct permissions.",
		Default: "waypoint-runner",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.Config.OdrServiceAccountInit,
		Usage:   "Create the service account if it does not exist.",
		Default: true,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-storageclassname",
		Target: &i.Config.StorageClassName,
		Usage:  "Name of the StorageClass required by the volume claim to install the Waypoint server image to.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-storage-request",
		Target:  &i.Config.StorageRequest,
		Usage:   "Configures the requested persistent volume size for the Waypoint server in Kubernetes.",
		Default: "1Gi",
	})
}

func (i *K8sInstaller) UpgradeFlags(set *flag.Set) {
	set.BoolVar(&flag.BoolVar{
		Name:   "k8s-advertise-internal",
		Target: &i.Config.AdvertiseInternal,
		Usage: "Advertise the internal service address rather than the external. " +
			"This is useful if all your deployments will be able to access the private " +
			"service address. This will default to false but will be automatically set to " +
			"true if the external host is detected to be localhost.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.Config.K8sContext,
		Usage: "The Kubernetes context to upgrade the Waypoint server to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Usage:   "Namespace to install the Waypoint server into for Kubernetes.",
		Default: "default",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-image",
		Target:  &i.Config.ServerImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-odr-image",
		Target: &i.Config.OdrImage,
		Usage:  "Docker image for the Waypoint On-Demand Runners",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-runner-service-account",
		Target: &i.Config.OdrServiceAccount,
		Usage: "Service account to assign to the on-demand runner. If this is blank, " +
			"a service account will be created automatically with the correct permissions.",
		Default: "waypoint-runner",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.Config.OdrServiceAccountInit,
		Usage:   "Create the service account if it does not exist.",
		Default: true,
	})
}

func (i *K8sInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.Config.K8sContext,
		Usage: "The Kubernetes context to unisntall the Waypoint server from. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.Config.Namespace,
		Usage:   "Namespace in Kubernetes to uninstall the Waypoint server from.",
		Default: "default",
	})
}

var warnK8SKind = strings.TrimSpace(`
Kind cluster detected!

Installing Waypoint to a Kind cluster requires that the cluster has
LoadBalancer capabilities (such as metallb). If Kind isn't configured
in this way, then the install may hang. If this happens, please delete
all the Waypoint resources and try again.
`)
