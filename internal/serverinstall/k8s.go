package serverinstall

import (
	"context"
	"encoding/json"
	"fmt"
	helminstallutil "github.com/hashicorp/waypoint/internal/installutil/helm"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"net"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	k8sinstallutil "github.com/hashicorp/waypoint/internal/installutil/k8s"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	dockerparser "github.com/novln/docker-parser"
)

//
type K8sInstaller struct {
	k8sinstallutil.K8sInstaller
	config k8sConfig
}

type k8sConfig struct {
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	odrImage              string `hcl:"odr_image,optional"`
	odrServiceAccount     string `hcl:"odr_service_account,optional"`
	odrServiceAccountInit bool   `hcl:"odr_service_account_init,optional"`

	advertiseInternal bool   `hcl:"advertise_internal,optional"`
	imagePullPolicy   string `hcl:"image_pull_policy,optional"`
	k8sContext        string `hcl:"k8s_context,optional"`
	cpuRequest        string `hcl:"cpu_request,optional"`
	memRequest        string `hcl:"mem_request,optional"`
	cpuLimit          string `hcl:"cpu_limit,optional"`
	memLimit          string `hcl:"mem_limit,optional"`
	storageClassName  string `hcl:"storageclassname,optional"`
	storageRequest    string `hcl:"storage_request,optional"`
	secretFile        string `hcl:"secret_file,optional"`
	imagePullSecret   string `hcl:"image_pull_secret,optional"`
	kubeConfigPath    string `hcl:"kubeconfig_path,optional"`
	version           string `hcl:"version,optional"`
}

const (
	serviceName                  = "waypoint"
	runnerRoleBindingName        = "waypoint-runner-rolebinding"
	runnerClusterRoleName        = "waypoint-runner"
	runnerClusterRoleBindingName = "waypoint-runner"
)

// Install is a method of K8sInstaller and implements the Installer interface to
// register a waypoint-server in a Kubernetes cluster
func (i *K8sInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, error) {
	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = defaultODRImage(i.config.serverImage)
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
	s.Update("Helm settings retrieved")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Getting Helm action configuration...")
	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.config.kubeConfigPath, i.config.k8sContext)
	if err != nil {
		return nil, err
	}
	s.Update("Helm action initialized")
	s.Status(terminal.StatusOK)
	s.Done()

	chartNS := ""
	if v := i.config.namespace; v != "" {
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
	s.Update("Helm install created")
	s.Status(terminal.StatusOK)
	s.Done()

	var version string
	if i.config.version == "" {
		version = helminstallutil.DefaultHelmChartVersion
	} else {
		version = i.config.version
	}

	s = sg.Add("Locating chart...")
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/v"+version+".tar.gz", settings)
	if err != nil {
		return nil, err
	}
	s.Update("Helm chart located")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Loading Helm chart...")
	c, err := loader.Load(path)
	if err != nil {
		return nil, err
	}
	s.Update("Helm chart loaded")
	s.Status(terminal.StatusOK)
	s.Done()

	imageRef, err := dockerparser.Parse(i.config.serverImage)
	if err != nil {
		ui.Output("Error parsing image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, err
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": true,
			"image": map[string]interface{}{
				"repository": imageRef.ShortName(),
				"tag":        imageRef.Tag(),
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"memory": i.config.memRequest,
					"cpu":    i.config.cpuRequest,
				},
				"limits": map[string]interface{}{
					"memory": i.config.memLimit,
					"cpu":    i.config.cpuLimit,
				},
			},
		},
		"runner": map[string]interface{}{
			"enabled": false,
		},
	}
	s = sg.Add("Installing Waypoint Helm chart...")
	_, err = client.RunWithContext(ctx, c, values)
	if err != nil {
		return nil, err
	}

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	// TODO: Move this to a util function for install and upgrade to use
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		clientset, err := i.NewClient()
		if err != nil {
			return false, err
		}

		s.Update("Getting waypoint-ui service...")
		svc, err := clientset.CoreV1().Services(i.config.namespace).Get(
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
		log.Info("server service ready: %s", addr)

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
		if i.config.advertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				"waypoint-server",
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
	s.Done()

	s.Update("Waypoint server installed with Helm!")
	s.Status(terminal.StatusOK)
	s.Done()

	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Upgrade is a method of K8sInstaller and implements the Installer interface to
// upgrade a waypoint-server in a Kubernetes cluster
func (i *K8sInstaller) Upgrade(
	ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (
	*InstallResults, error,
) {
	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = defaultODRImage(i.config.serverImage)
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
	s.Update("Helm settings retrieved")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Getting Helm action configuration...")
	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.config.kubeConfigPath, i.config.k8sContext)
	if err != nil {
		return nil, err
	}
	s.Update("Helm action initialized")
	s.Status(terminal.StatusOK)
	s.Done()

	chartNS := ""
	if v := i.config.namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

	s = sg.Add("Creating new Helm upgrade object...")
	client := action.NewUpgrade(actionConfig)
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.WaitForJobs = false
	client.Devel = true
	client.DependencyUpdate = false
	client.Timeout = 300 * time.Second
	client.Namespace = chartNS
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = true
	client.DisableOpenAPIValidation = false
	client.Description = ""
	s.Update("Helm upgrade created")
	s.Status(terminal.StatusOK)
	s.Done()

	var version string
	if i.config.version == "" {
		version = helminstallutil.DefaultHelmChartVersion
	} else {
		version = i.config.version
	}

	s = sg.Add("Locating chart...")
	path, err := client.LocateChart("https://github.com/hashicorp/waypoint-helm/archive/refs/tags/v"+version+".tar.gz", settings)
	if err != nil {
		return nil, err
	}
	s.Update("Helm chart located")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Loading Helm chart...")
	c, err := loader.Load(path)
	if err != nil {
		return nil, err
	}
	s.Update("Helm chart loaded")
	s.Status(terminal.StatusOK)
	s.Done()

	imageRef, err := dockerparser.Parse(i.config.serverImage)
	if err != nil {
		ui.Output("Error parsing image ref: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return nil, err
	}

	values := map[string]interface{}{
		"server": map[string]interface{}{
			"enabled": true,
			"image": map[string]interface{}{
				"repository": imageRef.ShortName(),
				"tag":        imageRef.Tag(),
			},
		},
		"runner": map[string]interface{}{
			"enabled": false,
		},
	}
	s = sg.Add("Installing Waypoint Helm chart...")
	_, err = client.RunWithContext(ctx, "waypoint", c, values)
	if err != nil {
		return nil, err
	}

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		clientset, err := i.NewClient()
		if err != nil {
			return false, err
		}

		s.Update("Getting waypoint-ui service...")
		svc, err := clientset.CoreV1().Services(i.config.namespace).Get(
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
		log.Info("server service ready: %s", addr)

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
		if i.config.advertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				"waypoint-server",
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

	actionConfig, err := helminstallutil.ActionInit(opts.Log, i.config.kubeConfigPath, i.config.k8sContext)
	if err != nil {
		return err
	}
	s.Update("Helm action initialized")
	s.Status(terminal.StatusOK)
	s.Done()

	chartNS := ""
	if v := i.config.namespace; v != "" {
		chartNS = v
	}
	if chartNS == "" {
		// If all else fails, default the namespace to "default"
		chartNS = "default"
	}

	s = sg.Add("Creating new Helm uninstall object...")
	client := action.NewUninstall(actionConfig)
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.Timeout = 300 * time.Second
	client.Description = ""
	s.Update("Helm uninstall created")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Uninstalling Helm chart...")
	_, err = client.Run("waypoint")
	if err != nil {
		return err
	}
	s.Update("Waypoint uninstalled with Helm!")
	s.Status(terminal.StatusOK)
	s.Done()

	// TODO: Delete runner (or all runners?)

	return nil
}

// InstallRunner implements Installer.
func (i *K8sInstaller) InstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	ref, err := dockerparser.Parse(i.config.serverImage)
	if err != nil {
		opts.UI.Output("Error parsing image name: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}
	runnerInstaller := runnerinstall.K8sRunnerInstaller{
		Config: runnerinstall.K8sConfig{
			KubeconfigPath:       "",
			K8sContext:           i.config.k8sContext,
			Version:              helminstallutil.DefaultHelmChartVersion,
			Namespace:            i.config.namespace,
			RunnerImage:          ref.ShortName(),
			RunnerImageTag:       ref.Tag(),
			CpuRequest:           i.config.cpuRequest,
			MemRequest:           i.config.memRequest,
			CreateServiceAccount: true,
		},
	}

	err = runnerInstaller.Install(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// UninstallRunner implements Installer.
func (i *K8sInstaller) UninstallRunner(
	ctx context.Context,
	opts *InstallOpts,
) error {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer func() { s.Abort() }()

	clientset, err := i.NewClient()
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

	deploymentClient := clientset.AppsV1().Deployments(i.config.namespace)
	if list, err := deploymentClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", runnerName),
	}); err != nil {
		ui.Output(
			"Error looking up deployments: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	} else if len(list.Items) > 0 {
		s.Update("Deleting any automatically installed runners...")

		// Record various settings we can reuse for runner reinstallation
		// if we're doing an upgrade. We need to do this because the upgrade
		// flags don't contain the installation settings, and we prefer them
		// not to; instead we just retain the old settings.
		//
		// Note we have lots of conditionals here to try to avoid weird
		// panic situations if the remote side doesn't have the fields we
		// expect.
		podSpec := list.Items[0].Spec.Template.Spec
		if secrets := podSpec.ImagePullSecrets; len(secrets) > 0 {
			i.config.imagePullSecret = secrets[0].Name
		}
		if v := podSpec.Containers; len(v) > 0 {
			c := v[0]

			i.config.imagePullPolicy = string(c.ImagePullPolicy)
			if m := c.Resources.Requests; len(m) > 0 {
				if v, ok := m[apiv1.ResourceMemory]; ok {
					i.config.memRequest = v.String()
				}
				if v, ok := m[apiv1.ResourceCPU]; ok {
					i.config.cpuRequest = v.String()
				}
			}
			if m := c.Resources.Limits; len(m) > 0 {
				if v, ok := m[apiv1.ResourceMemory]; ok {
					i.config.memLimit = v.String()
				}
				if v, ok := m[apiv1.ResourceCPU]; ok {
					i.config.cpuLimit = v.String()
				}
			}
		}

		// create our wait channel to later poll for statefulset+pod deletion
		w, err := deploymentClient.Watch(
			ctx,
			metav1.ListOptions{
				LabelSelector: "app=" + runnerName,
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
				LabelSelector: "app=" + runnerName,
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
	} else {
		s.Update("No runners installed.")
		s.Done()
	}

	return nil
}

// HasRunner implements Installer.
func (i *K8sInstaller) HasRunner(
	ctx context.Context,
	opts *InstallOpts,
) (bool, error) {
	clientset, err := i.NewClient()
	if err != nil {
		return false, err
	}

	deploymentClient := clientset.AppsV1().Deployments(i.config.namespace)
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
	if v := i.config.imagePullSecret; v != "" {
		cfgMap["image_secret"] = v
	}
	if v := i.config.odrServiceAccount; v != "" {
		cfgMap["service_account"] = v
	}
	if v := i.config.imagePullPolicy; v != "" {
		cfgMap["image_pull_policy"] = v
	}

	var cpuConfig k8s.ResourceConfig
	var memConfig k8s.ResourceConfig
	if v := i.config.cpuRequest; v != "" {
		cpuConfig.Requested = v
	}
	if v := i.config.memRequest; v != "" {
		memConfig.Requested = v
	}
	if v := i.config.cpuLimit; v != "" {
		cpuConfig.Limit = v
	}
	if v := i.config.memLimit; v != "" {
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
		OciUrl:       i.config.odrImage,
		PluginType:   "kubernetes",
		Default:      true,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
	}
}

// newServiceAccount takes in a k8sConfig and creates the ServiceAccount
// definition for the ODR.
func newServiceAccount(c k8sConfig) (*apiv1.ServiceAccount, error) {
	return &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.odrServiceAccount,
			Namespace: c.namespace,
		},
	}, nil
}

// newServiceAccountClusterRoleWithBinding creates the cluster role and binding necessary to create and verify
// a nodeport type services.
func newServiceAccountClusterRoleWithBinding(c k8sConfig) (*rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding, error) {
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
					Name:      c.odrServiceAccount,
					Namespace: c.namespace,
				},
			},
		}, nil
}

// newServiceAccountRoleBinding creates the role binding necessary to
// map the ODR role to the service account.
func newServiceAccountRoleBinding(c k8sConfig) (*rbacv1.RoleBinding, error) {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runnerRoleBindingName,
			Namespace: c.namespace,
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
				Name:      c.odrServiceAccount,
				Namespace: c.namespace,
			},
		},
	}, nil
}

func (i *K8sInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-config-path",
		Usage:  "Path to the kubeconfig file to use,",
		Target: &i.config.kubeConfigPath,
	})

	set.BoolVar(&flag.BoolVar{
		Name:   "k8s-advertise-internal",
		Target: &i.config.advertiseInternal,
		Usage: "Advertise the internal service address rather than the external. " +
			"This is useful if all your deployments will be able to access the private " +
			"service address. This will default to false but will be automatically set to " +
			"true if the external host is detected to be localhost.",
	})

	set.StringMapVar(&flag.StringMapVar{
		Name:   "k8s-annotate-service",
		Target: &i.config.serviceAnnotations,
		Usage:  "Annotations for the Service generated.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.config.k8sContext,
		Usage: "The Kubernetes context to install the Waypoint server to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-helm-version",
		Target: &i.config.version,
		Usage:  "The version of the Helm chart to use for the Waypoint runner install.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-request",
		Target:  &i.config.cpuRequest,
		Usage:   "Configures the requested CPU amount for the Waypoint server in Kubernetes.",
		Default: "100m",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-request",
		Target:  &i.config.memRequest,
		Usage:   "Configures the requested memory amount for the Waypoint server in Kubernetes.",
		Default: "256Mi",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-cpu-limit",
		Target:  &i.config.cpuLimit,
		Usage:   "Configures the CPU limit for the Waypoint server in Kubernetes.",
		Default: "100m",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-mem-limit",
		Target:  &i.config.memLimit,
		Usage:   "Configures the memory limit for the Waypoint server in Kubernetes.",
		Default: "256Mi",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.config.namespace,
		Usage:   "Namespace to install the Waypoint server into for Kubernetes.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-pull-policy",
		Target:  &i.config.imagePullPolicy,
		Usage:   "Set the pull policy for the Waypoint server image.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-pull-secret",
		Target: &i.config.imagePullSecret,
		Usage:  "Secret to use to access the Waypoint server image on Kubernetes.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-secret-file",
		Target: &i.config.secretFile,
		Usage:  "Use the Kubernetes Secret in the given path to access the Waypoint server image.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-odr-image",
		Target: &i.config.odrImage,
		Usage:  "Docker image for the Waypoint On-Demand Runners",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-runner-service-account",
		Target: &i.config.odrServiceAccount,
		Usage: "Service account to assign to the on-demand runner. If this is blank, " +
			"a service account will be created automatically with the correct permissions.",
		Default: "waypoint-runner",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.config.odrServiceAccountInit,
		Usage:   "Create the service account if it does not exist.",
		Default: true,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-storageclassname",
		Target: &i.config.storageClassName,
		Usage:  "Name of the StorageClass required by the volume claim to install the Waypoint server image to.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-storage-request",
		Target:  &i.config.storageRequest,
		Usage:   "Configures the requested persistent volume size for the Waypoint server in Kubernetes.",
		Default: "1Gi",
	})
}

func (i *K8sInstaller) UpgradeFlags(set *flag.Set) {
	set.BoolVar(&flag.BoolVar{
		Name:   "k8s-advertise-internal",
		Target: &i.config.advertiseInternal,
		Usage: "Advertise the internal service address rather than the external. " +
			"This is useful if all your deployments will be able to access the private " +
			"service address. This will default to false but will be automatically set to " +
			"true if the external host is detected to be localhost.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.config.k8sContext,
		Usage: "The Kubernetes context to upgrade the Waypoint server to. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.config.namespace,
		Usage:   "Namespace to install the Waypoint server into for Kubernetes.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-odr-image",
		Target: &i.config.odrImage,
		Usage:  "Docker image for the Waypoint On-Demand Runners",
	})

	set.StringVar(&flag.StringVar{
		Name:   "k8s-runner-service-account",
		Target: &i.config.odrServiceAccount,
		Usage: "Service account to assign to the on-demand runner. If this is blank, " +
			"a service account will be created automatically with the correct permissions.",
		Default: "waypoint-runner",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-runner-service-account-init",
		Target:  &i.config.odrServiceAccountInit,
		Usage:   "Create the service account if it does not exist.",
		Default: true,
	})
}

func (i *K8sInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "k8s-context",
		Target: &i.config.k8sContext,
		Usage: "The Kubernetes context to unisntall the Waypoint server from. If left" +
			" unset, Waypoint will use the current Kubernetes context.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-namespace",
		Target:  &i.config.namespace,
		Usage:   "Namespace in Kubernetes to uninstall the Waypoint server from.",
		Default: "",
	})
}

var warnK8SKind = strings.TrimSpace(`
Kind cluster detected!

Installing Waypoint to a Kind cluster requires that the cluster has
LoadBalancer capabilities (such as metallb). If Kind isn't configured
in this way, then the install may hang. If this happens, please delete
all the Waypoint resources and try again.
`)
