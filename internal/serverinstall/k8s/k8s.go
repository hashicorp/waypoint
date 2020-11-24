package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/hashicorp/waypoint/internal/serverinstall/config"
)

type Platform struct {
	Config *config.BaseConfig
}

func (p *Platform) Install(
	ctx context.Context, ui terminal.UI, log hclog.Logger,
) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, error) {
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
	if p.Config.Namespace == "" {
		namespace, _, err := newCmdConfig.Namespace()
		if err != nil {
			ui.Output(
				"Error getting namespace from client config: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", err
		}
		p.Config.Namespace = namespace
	}

	clientconfig, err := newCmdConfig.ClientConfig()
	if err != nil {
		ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", err
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", err
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

	if p.Config.SecretFile != "" {
		s.Update("Initializing Kubernetes secret")

		data, err := ioutil.ReadFile(p.Config.SecretFile)
		if err != nil {
			ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", err
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
			return nil, nil, "", err
		}

		if secretData.Metadata.Name == "" {
			ui.Output(
				"Invalid secret, no metadata.name",
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", err
		}

		p.Config.ImagePullSecret = secretData.Metadata.Name

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

			return nil, nil, "", err
		}

		s.Done()
		s = sg.Add("")
	}

	// Do some probing to see if this is OpenShift. If so, we'll switch the config for the user.
	// Setting the OpenShift flag will short circuit this.
	if !p.Config.OpenShift {
		s.Update("Gathering information about the Kubernetes cluster...")
		namespaceClient := clientset.CoreV1().Namespaces()
		_, err := namespaceClient.Get(context.TODO(), "openshift", metav1.GetOptions{})
		isOpenShift := err == nil

		// Default namespace in OpenShift acts like a regular K8s namespace, so we don't want
		// to remove fsGroup in this case.
		if isOpenShift && p.Config.Namespace != "default" {
			s.Update("OpenShift detected. Switching configuration...")
			p.Config.OpenShift = true
		}
	}

	// Decode our configuration
	statefulset, err := newStatefulSet(p.Config)
	if err != nil {
		ui.Output(
			"Error generating statefulset configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", err
	}

	service, err := newService(p.Config)
	if err != nil {
		ui.Output(
			"Error generating service configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", err
	}

	s.Update("Creating Kubernetes resources...")

	serviceClient := clientset.CoreV1().Services(p.Config.Namespace)
	_, err = serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		ui.Output(
			"Error creating service %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
	}

	statefulSetClient := clientset.AppsV1().StatefulSets(p.Config.Namespace)
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
		ss, err := clientset.AppsV1().StatefulSets(p.Config.Namespace).Get(
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
		return nil, nil, "", fmt.Errorf(
			"error waiting for statefulset ready: %s",
			err,
		)
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
		svc, err := clientset.CoreV1().Services(p.Config.Namespace).Get(
			ctx, p.Config.ServiceName, metav1.GetOptions{})
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

		endpoints, err := clientset.CoreV1().Endpoints(p.Config.Namespace).Get(
			ctx, p.Config.ServiceName, metav1.GetOptions{})
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
		if p.Config.AdvertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				p.Config.ServiceName,
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
		return nil, nil, "", fmt.Errorf(
			"error waiting for service ready: %s",
			err,
		)
	}

	s.Done()

	return &contextConfig, &advertiseAddr, httpAddr, err
}

// NewStatefulSet creates a new Waypoint Statefulset for deployment in Kubernetes.
func newStatefulSet(c *config.BaseConfig) (*appsv1.StatefulSet, error) {
	cpuRequest, err := resource.ParseQuantity(c.CPURequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse cpu request resource %s: %s", c.CPURequest, err)
	}

	memRequest, err := resource.ParseQuantity(c.MemRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse memory request resource %s: %s", c.MemRequest, err)
	}

	storageRequest, err := resource.ParseQuantity(c.StorageRequest)
	if err != nil {
		return nil, fmt.Errorf("could not parse storage request resource %s: %s", c.StorageRequest, err)
	}

	securityContext := &apiv1.PodSecurityContext{}
	if !c.OpenShift {
		securityContext.FSGroup = int64Ptr(1000)
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.ServerName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"app": c.ServerName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": c.ServerName,
				},
			},
			ServiceName: c.ServiceName,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": c.ServerName,
					},
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: c.ImagePullSecret,
						},
					},
					SecurityContext: securityContext,
					Containers: []apiv1.Container{
						{
							Name:            "server",
							Image:           c.ServerImage,
							ImagePullPolicy: apiv1.PullAlways,
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

// NewService creates a new Waypoint LoadBalancer for deployment in Kubernetes.
func newService(c *config.BaseConfig) (*apiv1.Service, error) {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.ServiceName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"app": c.ServerName,
			},
			Annotations: c.ServiceAnnotations,
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
				"app": c.ServerName,
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}, nil
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
