package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// The port that a service will forward to the pod(s)
const DefaultPort = 80

// Releaser is the ReleaseManager implementation for Google Cloud Run.
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

// Release creates a Kubernetes service configured for the deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *Deployment,
) (*Release, error) {
	var result Release

	st := ui.Status()
	defer st.Close()

	// Get our clientset
	clientset, ns, err := clientset(r.config.KubeconfigPath, r.config.Context)
	if err != nil {
		return nil, err
	}
	serviceclient := clientset.CoreV1().Services(ns)

	// Determine if we have a deployment that we manage already
	create := false
	service, err := serviceclient.Get(ctx, src.App, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		service = result.newService(src.App)
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

	var checkLB bool

	if r.config.LoadBalancer {
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
		checkLB = true
	} else if r.config.NodePort != 0 {
		service.Spec.Type = corev1.ServiceTypeNodePort
		if r.config.NodePort < 0 {
			r.config.NodePort = 0
		}
	} else {
		service.Spec.Type = corev1.ServiceTypeClusterIP
	}

	port := r.config.Port
	if port == 0 {
		port = DefaultPort
	}

	service.Spec.Ports = []corev1.ServicePort{
		{
			Port:       int32(port),
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
			NodePort:   int32(r.config.NodePort),
		},
	}

	// Create/update
	if create {
		st.Update("Creating service...")
		service, err = serviceclient.Create(ctx, service, metav1.CreateOptions{})
	} else {
		st.Update("Updating service...")
		service, err = serviceclient.Update(ctx, service, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, err
	}

	// Wait on the IP
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		service, err = serviceclient.Get(ctx, src.App, metav1.GetOptions{})
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

	st.Step(terminal.StatusOK, "Service succesfully configured!")

	if r.config.LoadBalancer {
		ingress := service.Status.LoadBalancer.Ingress[0]
		result.Url = "http://" + ingress.IP
		if ingress.Hostname != "" {
			result.Url = "http://" + ingress.Hostname
		}

		if port != 80 {
			result.Url = fmt.Sprintf("%s:%d", result.Url, port)
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
	Port int `hcl:"port,optional"`

	// NodePort configures a port to access the service on whichever node
	// is running service.
	NodePort int `hcl:"node_port,optional"`
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
