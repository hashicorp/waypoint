package k8s

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// The port that a service will forward to the pod(s)
const DefaultPort = 80

// Releaser is the ReleaseManager implementation for Kubernetes.
type Releaser struct {
	config ReleaserConfig
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// DestroyFunc implements component.Destroyer
func (r *Releaser) DestroyFunc() interface{} {
	return r.Destroy
}

// Release creates a Kubernetes service configured for the deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *Deployment,
) (*Release, error) {
	var result Release
	result.ServiceName = src.App

	sg := ui.StepGroup()
	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Get our clientset
	clientset, ns, config, err := clientset(r.config.KubeconfigPath, r.config.Context)
	if err != nil {
		return nil, err
	}

	// Override namespace if set
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	step.Update("Kubernetes client connected to %s with namespace %s", config.Host, ns)
	step.Done()

	step = sg.Add("Preparing service...")

	serviceclient := clientset.CoreV1().Services(ns)

	// Determine if we have a deployment that we manage already
	create := false
	service, err := serviceclient.Get(ctx, result.ServiceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		service = result.newService(result.ServiceName)
		create = true
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// Update the spec
	service.Spec.Selector = map[string]string{
		"name":  target.Name,
		labelId: target.Id,
	}

	if (r.config.Port != 0 || r.config.NodePort != 0) && r.config.Ports != nil {
		return nil, fmt.Errorf("Cannot define both 'ports' and 'port' or 'node_port'." +
			" Use 'ports' for configuring multiple service ports.")
	} else if r.config.Ports == nil && (r.config.Port != 0 || r.config.NodePort != 0) {
		r.config.Ports = make([]map[string]string, 1)
		r.config.Ports[0] = map[string]string{
			"port":        strconv.Itoa(int(r.config.Port)),
			"target_port": "http",
			"node_port":   strconv.Itoa(int(r.config.NodePort)),
		}
	} else if r.config.Port == 0 && r.config.NodePort == 0 && r.config.Ports == nil {
		// We don't explicitly set nodeport if Port isn't defined, because
		// k8s will automatically assign a nodeport if unspecified
		r.config.Ports = make([]map[string]string, 1)
		r.config.Ports[0] = map[string]string{
			"target_port": "http",
			"port":        strconv.Itoa(int(DefaultPort)),
		}
	}

	var checkLB bool

	if r.config.LoadBalancer {
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
		checkLB = true
	} else if r.config.Ports[0]["node_port"] != "" && r.config.Ports[0]["node_port"] != "0" {
		service.Spec.Type = corev1.ServiceTypeNodePort
	} else {
		service.Spec.Type = corev1.ServiceTypeClusterIP
	}

	servicePorts := make([]corev1.ServicePort, len(r.config.Ports))
	for i, sp := range r.config.Ports {
		nodePort, _ := strconv.ParseInt(sp["node_port"], 10, 32)
		port, _ := strconv.ParseInt(sp["port"], 10, 32)
		if port == 0 {
			// This likely means port was unset and got parsed to 0
			port = DefaultPort
		}

		var target_port int
		if sp["target_port"] == "" {
			sp["target_port"] = "http"
		} else {
			target_port, err = strconv.Atoi(sp["target_port"])
			if err != nil {
				// it's a string label like 'http', not an integer
				target_port = 0
			}
		}

		servicePorts[i] = corev1.ServicePort{
			Name:     sp["name"],
			Port:     int32(port),
			Protocol: corev1.ProtocolTCP,
			NodePort: int32(nodePort),
		}

		// Because of the type TargetPort is expected to be, we can't pass along
		// an int as a string, it expects the int to actually be an int
		if target_port != 0 {
			servicePorts[i].TargetPort = intstr.FromInt(target_port)
		} else {
			servicePorts[i].TargetPort = intstr.FromString(sp["target_port"])
		}
	}

	service.Spec.Ports = servicePorts

	// Create/update
	if create {
		step.Update("Creating service...")
		service, err = serviceclient.Create(ctx, service, metav1.CreateOptions{})
	} else {
		step.Update("Updating service...")
		service, err = serviceclient.Update(ctx, service, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, err
	}

	step.Done()
	step = sg.Add("Waiting for service to become ready...")

	// Wait on the IP
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		service, err = serviceclient.Get(ctx, result.ServiceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if checkLB {
			return len(service.Status.LoadBalancer.Ingress) > 0, nil
		} else {
			return service.Spec.ClusterIP != "", nil
		}
	})
	if err != nil {
		return nil, err
	}

	step.Update("Service is ready!")
	step.Done()

	if r.config.LoadBalancer {
		ingress := service.Status.LoadBalancer.Ingress[0]
		result.Url = "http://" + ingress.IP
		if ingress.Hostname != "" {
			result.Url = "http://" + ingress.Hostname
		}

		if service.Spec.Ports[0].Port != 80 {
			result.Url = fmt.Sprintf("%s:%d", result.Url, service.Spec.Ports[0].Port)
		}
	} else if service.Spec.Ports[0].NodePort > 0 {
		nodeclient := clientset.CoreV1().Nodes()
		nodes, err := nodeclient.List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		nodeIP := nodes.Items[0].Status.Addresses[0].Address
		result.Url = fmt.Sprintf("http://%s:%d", nodeIP, service.Spec.Ports[0].NodePort)
	} else {
		result.Url = fmt.Sprintf("http://%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port)
	}

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (r *Releaser) Destroy(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) error {
	// This is possible if an older version of the Kubernetes plugin was used
	// prior to service name existing. This was only in pre-0.1 releases so
	// we just return nil and pretend the destroy succeeded. We can probably
	// remove this very quickly post-release.
	if release.ServiceName == "" {
		return nil
	}

	sg := ui.StepGroup()
	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Get our client
	clientset, ns, config, err := clientset(r.config.KubeconfigPath, r.config.Context)
	if err != nil {
		return err
	}

	// Override namespace if set
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	step.Update("Kubernetes client connected to %s with namespace %s", config.Host, ns)
	step.Done()
	step = sg.Add("Deleting service...")

	serviceclient := clientset.CoreV1().Services(ns)
	if err := serviceclient.Delete(ctx, release.ServiceName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	step.Done()
	return nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`

	// Load Balancer sets whether or not the service will be a load
	// balancer type service
	LoadBalancer bool `hcl:"load_balancer,optional"`

	// Port configures the port that is used to access the service.
	// The default is 80.
	// Not valid if `Ports` is already defined
	// If defined, will internally be stored into `Ports`
	Port uint `hcl:"port,optional"`

	// A full resource of options to define ports for a service
	Ports []map[string]string `hcl:"ports,optional"`

	// NodePort configures a port to access the service on whichever node
	// is running service.
	// Not valid if `Ports` is already defined
	// If defined, will internally be stored into `Ports`
	NodePort uint `hcl:"node_port,optional"`

	// Namespace is the Kubernetes namespace to target the deployment to.
	Namespace string `hcl:"namespace,optional"`
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Manipulates the Kubernetes Service activate Deployments")

	doc.SetField(
		"kubeconfig",
		"path to the kubeconfig file to use",
		docs.Summary("by default uses from current user's home directory"),
		docs.EnvVar("KUBECONFIG"),
	)

	doc.SetField(
		"context",
		"the kubectl context to use, as defined in the kubeconfig file",
	)

	doc.SetField(
		"load_balancer",
		"indicates if the Kubernetes Service should LoadBalancer type",
		docs.Summary(
			"if the Kubernetes Service is not a LoadBalancer and node_port is not",
			"set, then the Service uses ClusterIP",
		),
	)

	doc.SetField(
		"node_port",
		"the TCP port that the Service should consume as a NodePort",
		docs.Summary(
			"if this is set but load_balancer is not, the service will be NodePort type,",
			"but if load_balancer is also set, it will be LoadBalancer",
		),
	)

	doc.SetField(
		"port",
		"the TCP port that the application is listening on",
		docs.Default(fmt.Sprint(DefaultPort)),
	)

	doc.SetField(
		"ports",
		"a map of ports and options that the application is listening on",
		docs.Summary(
			"used to define and configure multiple ports that the application is",
			"listening on. Available keys are 'port', 'node_port', 'name', and 'target_port'.",
			"If 'node_port' is set but 'load_balancer' is not, the service will be",
			" NodePort type. If 'load_balancer' is also set, it will be LoadBalancer.",
			"Ports defined will be TCP protocol.",
			"Note that 'name' is required if defining more than one port.",
		),
	)

	doc.SetField(
		"namespace",
		"namespace to create Service in",
		docs.Summary(
			"namespace is the name of the Kubernetes namespace to create the deployment in",
			"This is useful to create Services in non-default namespaces without creating kubeconfig contexts for each",
		),
	)

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Destroyer      = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
)
