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

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/terminal"
)

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

// Release deploys an image to GCR.
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	targets []component.ReleaseTarget,
) (*Release, error) {
	if len(targets) > 1 {
		return nil, fmt.Errorf(
			"The 'kubernetes' release manager does not support traffic splitting.")
	}

	var result Release

	// Get the deployment
	var deploy Deployment
	if err := component.ProtoAnyUnmarshal(targets[0].Deployment, &deploy); err != nil {
		return nil, err
	}

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
		"name":  deploy.Name,
		labelId: deploy.Id,
	}
	service.Spec.Type = corev1.ServiceTypeLoadBalancer
	service.Spec.Ports = []corev1.ServicePort{
		corev1.ServicePort{
			Port:       80,
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
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

		return len(service.Status.LoadBalancer.Ingress) > 0, nil
	})
	if err != nil {
		return nil, err
	}

	ingress := service.Status.LoadBalancer.Ingress[0]
	result.Url = "http://" + ingress.IP
	if ingress.Hostname != "" {
		result.Url = "http://" + ingress.Hostname
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
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
