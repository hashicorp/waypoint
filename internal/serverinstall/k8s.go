package serverinstall

import (
	"context"
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

//
type K8sInstaller struct {
	config k8sConfig
}

type k8sConfig struct {
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	advertiseInternal bool   `hcl:"advertise_internal,optional"`
	imagePullPolicy   string `hcl:"image_pull_policy,optional"`
	k8sContext        string `hcl:"k8s_context,optional"`
	openshift         bool   `hcl:"openshft,optional"`
	cpuRequest        string `hcl:"cpu_request,optional"`
	memRequest        string `hcl:"mem_request,optional"`
	storageClassName  string `hcl:"storageclassname,optional"`
	storageRequest    string `hcl:"storage_request,optional"`
	secretFile        string `hcl:"secret_file,optional"`
	imagePullSecret   string `hcl:"image_pull_secret,optional"`
}

const (
	serviceName = "waypoint"
)

// Install is a method of K8sInstaller and implements the Installer interface to
// register a waypoint-server in a Kubernetes cluster
func (i *K8sInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, error) {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer func() { s.Abort() }()

	clientset, err := i.newClient()
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return nil, err
	}

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

	if i.config.secretFile != "" {
		s.Update("Initializing Kubernetes secret")

		data, err := ioutil.ReadFile(i.config.secretFile)
		if err != nil {
			ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		}

		var secretData apiv1.Secret

		err = yaml.Unmarshal(data, &secretData)
		if err != nil {
			ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		}

		i.config.imagePullSecret = secretData.ObjectMeta.Name

		ui.Output("Installing kubernetes secret...")

		secretsClient := clientset.CoreV1().Secrets(i.config.namespace)
		_, err = secretsClient.Create(context.TODO(), &secretData, metav1.CreateOptions{})
		if err != nil {
			ui.Output(
				"Error creating Kubernetes secret from file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		}

		s.Done()
		s = sg.Add("")
	}

	// Do some probing to see if this is OpenShift. If so, we'll switch the config for the user.
	// Setting the OpenShift flag will short circuit this.
	if !i.config.openshift {
		s.Update("Gathering information about the Kubernetes cluster...")
		namespaceClient := clientset.CoreV1().Namespaces()
		_, err := namespaceClient.Get(context.TODO(), "openshift", metav1.GetOptions{})
		isOpenShift := err == nil

		// Default namespace in OpenShift acts like a regular K8s namespace, so we don't want
		// to remove fsGroup in this case.
		if isOpenShift && i.config.namespace != "default" {
			s.Update("OpenShift detected. Switching configuration...")
			i.config.openshift = true
		}
	}

	// Decode our configuration
	statefulset, err := newStatefulSet(i.config)
	if err != nil {
		ui.Output(
			"Error generating statefulset configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	service, err := newService(i.config)
	if err != nil {
		ui.Output(
			"Error generating service configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	s.Update("Creating Kubernetes resources...")

	serviceClient := clientset.CoreV1().Services(i.config.namespace)
	_, err = serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		ui.Output(
			"Error creating service %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	statefulSetClient := clientset.AppsV1().StatefulSets(i.config.namespace)
	_, err = statefulSetClient.Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err != nil {
		ui.Output(
			"Error creating statefulset %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	s.Done()
	s = sg.Add("Waiting for Kubernetes StatefulSet to be ready...")
	log.Info("waiting for server statefulset to become ready")
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		ss, err := clientset.AppsV1().StatefulSets(i.config.namespace).Get(
			ctx, serverName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if ss.Status.ReadyReplicas != ss.Status.Replicas {
			log.Trace("statefulset not ready, waiting")
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	s.Update("Kubernetes StatefulSet reporting ready")
	s.Done()

	s = sg.Add("Waiting for Kubernetes service to become ready..")

	// Wait for our service to be ready
	log.Info("waiting for server service to become ready")
	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		svc, err := clientset.CoreV1().Services(i.config.namespace).Get(
			ctx, serviceName, metav1.GetOptions{})
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

		endpoints, err := clientset.CoreV1().Endpoints(i.config.namespace).Get(
			ctx, serviceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if len(endpoints.Subsets) == 0 {
			log.Trace("endpoints are empty, waiting")
			return false, nil
		}

		// Get the ports
		var grpcPort int32
		var httpPort int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				grpcPort = spec.Port
			}

			if spec.Name == "http" {
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
		_, err = net.DialTimeout("tcp", httpAddr, 1*time.Second)
		if err != nil {
			return false, nil
		}
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
				serviceName,
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
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer s.Abort()

	clientset, err := i.newClient()
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return nil, err
	}

	// Do some probing to see if this is OpenShift. If so, we'll switch the config for the user.
	// Setting the OpenShift flag will short circuit this.
	if !i.config.openshift {
		s.Update("Gathering information about the Kubernetes cluster...")

		namespaceClient := clientset.CoreV1().Namespaces()
		_, err := namespaceClient.Get(context.TODO(), "openshift", metav1.GetOptions{})
		isOpenShift := err == nil

		// Default namespace in OpenShift acts like a regular K8s namespace, so we don't want
		// to remove fsGroup in this case.
		if isOpenShift && i.config.namespace != "default" {
			s.Update("OpenShift detected. Switching configuration...")
			i.config.openshift = true
		}
	}

	s.Done()

	statefulSetClient := clientset.AppsV1().StatefulSets(i.config.namespace)
	waypointStatefulSet, err := statefulSetClient.Get(ctx, serverName, metav1.GetOptions{})
	if err != nil {
		ui.Output(
			"Error obtaining waypoint statefulset: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	s = sg.Add("Upgrading server to %q", i.config.serverImage)

	// Update pod image to requested serverImage
	podClient := clientset.CoreV1().Pods(i.config.namespace)
	if podList, err := podClient.List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", serverName)}); err != nil {
		ui.Output(
			"Error listing pods: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	} else {
		for _, pod := range podList.Items {
			// patch the pod containers with the new i.config.serverImage
			// Payload should be the updated server config image with the podspec
			for j := range pod.Spec.Containers {
				pod.Spec.Containers[j].Image = i.config.serverImage
			}

			jsonPayload, err := json.Marshal(pod)
			if err != nil {
				return nil, err
			}

			_, err = podClient.Patch(ctx, pod.Name, types.MergePatchType, jsonPayload, metav1.PatchOptions{})
			if err != nil {
				ui.Output(
					"Error submitting patch to update container image: %s", clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return nil, err
			}
		}
	}

	s.Update("Patch update sent to waypoint server pod(s)")

	if waypointStatefulSet.Spec.UpdateStrategy.Type == "OnDelete" {
		s.Update("Deleting pod to refresh image")
		log.Info("Update Strategy is 'OnDelete', deleting pod to refresh image")

		if podList, err := podClient.List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", serverName)}); err != nil {
			ui.Output(
				"Error listing pods: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		} else {
			for _, pod := range podList.Items {
				if err := podClient.Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
					s.Update("Pod deletion failed", terminal.WithErrorStyle)
					s.Done()
					ui.Output(
						"Error deleting pod %q: %s", pod.Name, clierrors.Humanize(err),
						terminal.WithErrorStyle(),
					)
					return nil, err
				}
			}
		}

		log.Info("Pod(s) deleted, k8s will now restart waypoint server ", serverName)
	} else if waypointStatefulSet.Spec.UpdateStrategy.Type == "RollingUpdate" {
		log.Info("Update Strategy is 'RollingUpdate', no further action required")
	} else {
		log.Warn("Update Strategy is not recognized, so no action is taken", "UpdateStrategy",
			waypointStatefulSet.Spec.UpdateStrategy.Type)
	}

	s.Update("Image set to update!")
	s.Done()

	s = sg.Add("Waiting for server to be ready...")
	log.Info("waiting for waypoint server to become ready after image refresh")

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	err = wait.PollImmediate(2*time.Second, 2*time.Minute, func() (bool, error) {
		svc, err := clientset.CoreV1().Services(i.config.namespace).Get(
			ctx, serviceName, metav1.GetOptions{})
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

		endpoints, err := clientset.CoreV1().Endpoints(i.config.namespace).Get(
			ctx, serviceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if len(endpoints.Subsets) == 0 {
			log.Trace("endpoints are empty, waiting")
			return false, nil
		}

		// Get the ports
		var grpcPort int32
		var httpPort int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				grpcPort = spec.Port
			}

			if spec.Name == "http" {
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
		_, err = net.DialTimeout("tcp", httpAddr, 1*time.Second)
		if err != nil {
			return false, nil
		}
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
				serviceName,
				grpcPort,
			)
		}

		// Set our connection information
		contextConfig = clicontext.Config{
			Server: serverconfig.Client{
				Address:       grpcAddr,
				Tls:           true,
				TlsSkipVerify: true, // always for now
			},
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	if waypointStatefulSet.Spec.UpdateStrategy.Type == "RollingUpdate" {
		ui.Output("\nKubernetes is now set to upgrade waypoint server image with its\n" +
			"'RollingUpdate' strategy. This means the pod might not be updated immediately.")
	}
	s.Update("Upgrade complete!")
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
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer func() { s.Abort() }()

	clientset, err := i.newClient()
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

	ssClient := clientset.AppsV1().StatefulSets(i.config.namespace)
	if list, err := ssClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", serverName),
	}); err != nil {
		ui.Output(
			"Error looking up stateful sets: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	} else if len(list.Items) > 0 {
		s.Update("Deleting statefulset and pods...")

		// create our wait channel to later poll for statefulset+pod deletion
		w, err := ssClient.Watch(
			ctx,
			metav1.ListOptions{
				LabelSelector: "app=" + serverName,
			},
		)
		if err != nil {
			ui.Output(
				"Error creating stateful set watcher: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
		// send DELETE to statefulset collection
		if err = ssClient.DeleteCollection(
			ctx,
			metav1.DeleteOptions{},
			metav1.ListOptions{
				LabelSelector: "app=" + serverName,
			},
		); err != nil {
			ui.Output(
				"Error deleting Waypoint statefulset: %s", clierrors.Humanize(err),
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
				log.Trace("statefulset collection not fully removed, waiting")
				return false, nil
			default:
				log.Trace("no message received on watch.ResultChan(), waiting for Event")
				return false, nil
			}
		})
		if err != nil {
			return err
		}
		s.Update("Statefulset and pods deleted")
		s.Done()
		s = sg.Add("")
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(i.config.namespace)
	if list, err := pvcClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", serverName),
	}); err != nil {
		ui.Output(
			"Error looking up persistent volume claims: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	} else if len(list.Items) > 0 {
		s.Update("Deleting Persistent Volume Claim...")

		// create our wait channel to later poll for pvc deletion
		w, err := pvcClient.Watch(
			ctx,
			metav1.ListOptions{
				LabelSelector: "app=" + serverName,
			},
		)
		if err != nil {
			ui.Output(
				"Error creating persistent volume claims watcher: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
		// delete persistent volume claims
		if err = pvcClient.DeleteCollection(
			ctx,
			metav1.DeleteOptions{},
			metav1.ListOptions{
				LabelSelector: "app=" + serverName,
			},
		); err != nil {
			ui.Output(
				"Error deleting Waypoint pvc: %s", clierrors.Humanize(err),
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
				log.Trace("persistent volume claims collection not fully removed, waiting")
				return false, nil
			default:
				log.Trace("no message received on watch.ResultChan(), waiting for Event")
				return false, nil
			}
		})
		if err != nil {
			return err
		}

		s.Update("Persistent Volume Claim deleted")
		s.Done()
		s = sg.Add("")
	}

	svcClient := clientset.CoreV1().Services(i.config.namespace)
	if list, err := svcClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", serverName),
	}); err != nil {
		ui.Output(
			"Error looking up services: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	} else if len(list.Items) > 0 {
		s.Update("Deleting service...")

		// create our wait channel to later poll for service deletion
		w, err := svcClient.Watch(
			ctx,
			metav1.ListOptions{
				LabelSelector: "app=" + serverName,
			},
		)
		if err != nil {
			ui.Output(
				"Error creating service client watcher: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
		// delete waypoint service
		if err = svcClient.Delete(
			ctx,
			serviceName,
			metav1.DeleteOptions{},
		); err != nil {
			ui.Output(
				"Error deleting Waypoint service: %s", clierrors.Humanize(err),
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
				log.Trace("no message received on watch.ResultChan(), waiting for Event")
				return false, nil
			default:
				log.Trace("persistent volume claims not fully removed, waiting")
				return false, nil
			}
		})
		if err != nil {
			return err
		}

		s.Update("Service deleted")
		s.Done()
	}

	s.Done()

	return nil
}

// InstallRunner implements Installer.
func (i *K8sInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting Kubernetes cluster...")
	defer func() { s.Abort() }()

	clientset, err := i.newClient()
	if err != nil {
		ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

	// Decode our configuration
	deployment, err := newDeployment(i.config, opts)
	if err != nil {
		ui.Output(
			"Error generating deployment configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}

	s.Update("Creating Deployment for Runner")

	deploymentClient := clientset.AppsV1().Deployments(i.config.namespace)
	_, err = deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		ui.Output(
			"Error creating deployment %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}

	s.Done()
	s = sg.Add("Waiting for Kubernetes Deployment to be ready...")
	log.Info("waiting for server deployment to become ready")
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		ss, err := clientset.AppsV1().Deployments(i.config.namespace).Get(
			ctx, runnerName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if ss.Status.ReadyReplicas > 0 {
			return true, nil
		}

		log.Trace("deployment not ready, waiting")
		return false, nil
	})
	if err != nil {
		return err
	}

	s.Update("Kubernetes Deployment for Waypoint runner reporting ready")
	s.Done()

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

	clientset, err := i.newClient()
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
	clientset, err := i.newClient()
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

// newDeployment takes in a k8sConfig and creates a new Waypoint Deployment for
// deploying Waypoint runners.
func newDeployment(c k8sConfig, opts *InstallRunnerOpts) (*appsv1.Deployment, error) {
	// This is the port we'll use for the liveness check with the
	// runner. This isn't exposed outside the pod so it doesn't really
	// matter what it is.
	const livenessPort = 1234

	cpuRequest, err := resource.ParseQuantity(c.cpuRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse cpu request resource %s: %s", c.cpuRequest, err)
	}

	memRequest, err := resource.ParseQuantity(c.memRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse memory request resource %s: %s", c.memRequest, err)
	}

	securityContext := &apiv1.PodSecurityContext{}
	if !c.openshift {
		securityContext.FSGroup = int64Ptr(1000)
	}

	// Build our env vars so we can connect back to the Waypoint server.
	var envs []apiv1.EnvVar
	for _, line := range opts.AdvertiseClient.Env() {
		idx := strings.Index(line, "=")
		if idx == -1 {
			// Should never happen but let's not crash.
			continue
		}

		key := line[:idx]
		value := line[idx+1:]
		envs = append(envs, apiv1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runnerName,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app": runnerName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": runnerName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": runnerName,
					},

					Annotations: map[string]string{
						// These annotations are required for `img` to work
						// properly within Kubernetes.
						"container.apparmor.security.beta.kubernetes.io/runner": "unconfined",
						"container.seccomp.security.alpha.kubernetes.io/runner": "unconfined",
					},
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: c.imagePullSecret,
						},
					},
					SecurityContext: securityContext,
					Containers: []apiv1.Container{
						{
							Name:            "runner",
							Image:           c.serverImage,
							ImagePullPolicy: apiv1.PullPolicy(c.imagePullPolicy),
							Env:             envs,
							Command:         []string{serviceName},
							Args: []string{
								"runner",
								"agent",
								"-vvv",
								"-liveness-tcp-addr=:" + strconv.Itoa(livenessPort),
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									TCPSocket: &apiv1.TCPSocketAction{
										Port: intstr.FromInt(livenessPort),
									},
								},
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceMemory: memRequest,
									apiv1.ResourceCPU:    cpuRequest,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// newStatefulSet takes in a k8sConfig and creates a new Waypoint Statefulset
// for deployment in Kubernetes.
func newStatefulSet(c k8sConfig) (*appsv1.StatefulSet, error) {
	cpuRequest, err := resource.ParseQuantity(c.cpuRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse cpu request resource %s: %s", c.cpuRequest, err)
	}

	memRequest, err := resource.ParseQuantity(c.memRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse memory request resource %s: %s", c.memRequest, err)
	}

	storageRequest, err := resource.ParseQuantity(c.storageRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse storage request resource %s: %s", c.storageRequest, err)
	}

	securityContext := &apiv1.PodSecurityContext{}
	if !c.openshift {
		securityContext.FSGroup = int64Ptr(1000)
	}

	volumeClaimTemplates := []apiv1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "data",
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceStorage: storageRequest,
					},
				},
			},
		},
	}

	if c.storageClassName != "" {
		volumeClaimTemplates[0].Spec.StorageClassName = &c.storageClassName
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serverName,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app": serverName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": serverName,
				},
			},
			ServiceName: serviceName,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": serverName,
					},
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: c.imagePullSecret,
						},
					},
					SecurityContext: securityContext,
					Containers: []apiv1.Container{
						{
							Name:            "server",
							Image:           c.serverImage,
							ImagePullPolicy: apiv1.PullPolicy(c.imagePullPolicy),
							Env: []apiv1.EnvVar{
								{
									Name:  "HOME",
									Value: "/data",
								},
							},
							Command: []string{serviceName},
							Args: []string{
								"server",
								"run",
								"-accept-tos",
								"-vvv",
								"-db=/data/data.db",
								"-listen-grpc=0.0.0.0:9701",
								"-listen-http=0.0.0.0:9702",
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "grpc",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 9701,
								},
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 9702,
								},
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path:   "/",
										Port:   intstr.FromString("http"),
										Scheme: "HTTPS",
									},
								},
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceMemory: memRequest,
									apiv1.ResourceCPU:    cpuRequest,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/data",
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: volumeClaimTemplates,
		},
	}, nil
}

// newService takes in a k8sConfig and creates a new Waypoint LoadBalancer
// for deployment in Kubernetes.
func newService(c k8sConfig) (*apiv1.Service, error) {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app": serverName,
			},
			Annotations: c.serviceAnnotations,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port: 9701,
					Name: "grpc",
				},
				{
					Port: 9702,
					Name: "http",
				},
			},
			Selector: map[string]string{
				"app": serverName,
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}, nil
}

func (i *K8sInstaller) InstallFlags(set *flag.Set) {
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
		Name:    "k8s-namespace",
		Target:  &i.config.namespace,
		Usage:   "Namespace to install the Waypoint server into for Kubernetes.",
		Default: "",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-openshift",
		Target:  &i.config.openshift,
		Default: false,
		Usage:   "Enables installing the Waypoint server on Kubernetes on Red Hat OpenShift. If set, auto-configures the installation.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-pull-policy",
		Target:  &i.config.imagePullPolicy,
		Usage:   "Set the pull policy for the Waypoint server image.",
		Default: "",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-pull-secret",
		Target:  &i.config.imagePullSecret,
		Usage:   "Secret to use to access the Waypoint server image on Kubernetes.",
		Default: "github",
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

	set.BoolVar(&flag.BoolVar{
		Name:    "k8s-openshift",
		Target:  &i.config.openshift,
		Default: false,
		Usage:   "Enables installing the Waypoint server on Kubernetes on Red Hat OpenShift. If set, auto-configures the installation.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
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
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

// newClient creates a new K8S client based on the configured settings.
func (i *K8sInstaller) newClient() (*kubernetes.Clientset, error) {
	// Build our K8S client.
	configOverrides := &clientcmd.ConfigOverrides{}
	if i.config.k8sContext != "" {
		configOverrides = &clientcmd.ConfigOverrides{
			CurrentContext: i.config.k8sContext,
		}
	}
	newCmdConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		configOverrides,
	)

	// Discover the current target namespace in the user's config so if they
	// run kubectl commands waypoint will show up. If we use the default namespace
	// they might not see the objects we've created.
	if i.config.namespace == "" {
		namespace, _, err := newCmdConfig.Namespace()
		if err != nil {
			return nil, fmt.Errorf(
				"Error getting namespace from client config: %s",
				clierrors.Humanize(err),
			)
		}

		i.config.namespace = namespace
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

var (
	warnK8SKind = strings.TrimSpace(`
Kind cluster detected!

Installing Waypoint to a Kind cluster requires that the cluster has
LoadBalancer capabilities (such as metallb). If Kind isn't configured
in this way, then the install may hang. If this happens, please delete
all the Waypoint resources and try again.
`)
)
