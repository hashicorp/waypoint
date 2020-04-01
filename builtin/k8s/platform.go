package k8s

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/mitchellh/devflow/builtin/docker"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/terminal"
)

const (
	labelId    = "devflow.hashicorp.com/id"
	labelNonce = "devflow.hashicorp.com/nonce"
)

// Platform is the Platform implementation for Google Cloud Run.
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// Deploy deploys an image to GCR.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	dir *datadir.Component,
	ui terminal.UI,
) (*Deployment, error) {
	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	result.Name = src.App

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get our client
	clientset, ns, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	deployclient := clientset.AppsV1().Deployments(ns)

	// Determine if we have a deployment that we manage already
	create := false
	deployment, err := deployclient.Get(ctx, result.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		deployment = result.newDeployment(result.Name)
		create = true
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// Set our ID on the label. We use this ID so that we can have a key
	// to route to multiple versions during release management.
	deployment.Spec.Template.Labels[labelId] = result.Id

	// Update the deployment with our spec
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			corev1.Container{
				Name:            result.Name,
				Image:           img.Name(),
				ImagePullPolicy: corev1.PullAlways,
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Name:          "http",
						ContainerPort: 3000,
					},
				},
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "PORT",
						Value: "3000",
					},
				},
			},
		},
	}
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations[labelNonce] =
		time.Now().UTC().Format(time.RFC3339Nano)

	// Create/update
	if create {
		st.Update("Creating deployment...")
		deployment, err = clientset.AppsV1().Deployments(ns).Create(
			ctx, deployment, metav1.CreateOptions{})
	} else {
		st.Update("Updating deployment...")
		deployment, err = clientset.AppsV1().Deployments(ns).Update(
			ctx, deployment, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
