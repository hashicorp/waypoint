package serverinstall

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	serverName        string `hcl:"server_name,optional"`
	serviceName       string `hcl:"service_name,optional"`
	openshift         bool   `hcl:"openshft,optional"`
	cpuRequest        string `hcl:"cpu_request,optional"`
	memRequest        string `hcl:"mem_request,optional"`
	storageRequest    string `hcl:"storage_request,optional"`
	secretFile        string `hcl:"secret_file,optional"`
	imagePullSecret   string `hcl:"image_pull_secret,optional"`
}

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
	defer s.Abort()

	// Build our K8S client.
	newCmdConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Discover the current target namespace in the user's config so if they
	// run kubectl commands waypoint will show up. If we use the default namespace
	// they might not see the objects we've created.
	if i.config.namespace == "" {
		namespace, _, err := newCmdConfig.Namespace()
		if err != nil {
			ui.Output(
				"Error getting namespace from client config: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		}
		i.config.namespace = namespace
	}

	clientconfig, err := newCmdConfig.ClientConfig()
	if err != nil {
		ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
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

		var secretData struct {
			Metadata struct {
				Name string `yaml:"name"`
			} `yaml:"metadata"`
		}

		err = yaml.Unmarshal(data, &secretData)
		if err != nil {
			ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		}

		if secretData.Metadata.Name == "" {
			ui.Output(
				"Invalid secret, no metadata.name",
				terminal.WithErrorStyle(),
			)
			return nil, err
		}

		i.config.imagePullSecret = secretData.Metadata.Name

		ui.Output("Installing kubernetes secret...")

		cmd := exec.Command("kubectl", "create", "-f", "-")
		cmd.Stdin = bytes.NewReader(data)
		cmd.Stdout = s.TermOutput()
		cmd.Stderr = cmd.Stdout

		if err = cmd.Run(); err != nil {
			ui.Output(
				"Error executing kubectl to install secret: %s", clierrors.Humanize(err),
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
	}

	statefulSetClient := clientset.AppsV1().StatefulSets(i.config.namespace)
	_, err = statefulSetClient.Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err != nil {
		ui.Output(
			"Error creating statefulset %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
	}

	s.Done()
	s = sg.Add("Waiting for Kubernetes StatefulSet to be ready...")
	log.Info("waiting for server statefulset to become ready")
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		ss, err := clientset.AppsV1().StatefulSets(i.config.namespace).Get(
			ctx, "waypoint-server", metav1.GetOptions{})
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
			ctx, i.config.serviceName, metav1.GetOptions{})
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
			ctx, i.config.serviceName, metav1.GetOptions{})
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
				i.config.serviceName,
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

	s.Done()

	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
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

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.serverName,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app": c.serverName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": c.serverName,
				},
			},
			ServiceName: c.serviceName,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": c.serverName,
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
							Command: []string{"waypoint"},
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
			VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{
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
			},
		},
	}, nil
}

// newService takes in a k8sConfig and creates a new Waypoint LoadBalancer
// for deployment in Kubernetes.
func newService(c k8sConfig) (*apiv1.Service, error) {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.serviceName,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app": c.serverName,
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
				"app": c.serverName,
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}, nil
}

// InstallRunner implements Installer.
func (i *K8sInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	// TODO
	return nil
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
		Default: "hashicorp/waypoint:latest",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-server-name",
		Target:  &i.config.serverName,
		Usage:   "Name of the Waypoint server for Kubernetes.",
		Default: "waypoint-server",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-service-name",
		Target:  &i.config.serviceName,
		Usage:   "Name of the Kubernetes service for the Waypoint server.",
		Default: "waypoint",
	})

	set.StringVar(&flag.StringVar{
		Name:    "k8s-storage-request",
		Target:  &i.config.storageRequest,
		Usage:   "Configures the requested persistent volume size for the Waypoint server in Kubernetes.",
		Default: "1Gi",
	})
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
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
