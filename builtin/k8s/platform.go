package k8s

import (
	"context"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/go-homedir"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/mitchellh/devflow/builtin/docker"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/terminal"
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
	var result Deployment

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get our client
	clientset, ns, err := p.clientset()
	if err != nil {
		return nil, err
	}

	deployclient := clientset.AppsV1().Deployments(ns)

	// Determine if we have a deployment that we manage already
	create := false
	deployment, err := deployclient.Get(ctx, src.App, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		deployment = result.newDeployment(src.App)
		create = true
	}
	if err != nil {
		return nil, err
	}

	// Update the deployment with our spec
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			corev1.Container{
				Name:            src.App,
				Image:           img.Name(),
				ImagePullPolicy: corev1.PullAlways,
			},
		},
	}
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["devflow.hashicorp.com/nonce"] =
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

	return &Deployment{
		Name: deployment.Name,
	}, nil
}

// clientset returns a K8S clientset and configured namespace.
func (p *Platform) clientset() (*kubernetes.Clientset, string, error) {
	// Path to the kube config file
	kubeconfig := p.config.KubeconfigPath
	if kubeconfig == "" {
		dir, err := homedir.Dir()
		if err != nil {
			return nil, "", status.Errorf(codes.Aborted,
				"failed to load home directory: %s",
				err)
		}

		kubeconfig = filepath.Join(dir, ".kube", "config")
	}

	// Build our config and client
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: p.config.Context,
		},
	)

	// Get our configured namespace
	ns, _, err := config.Namespace()
	if err != nil {
		return nil, "", status.Errorf(codes.Aborted,
			"failed to initialize K8S client configuration: %s", err)
	}

	clientconfig, err := config.ClientConfig()
	if err != nil {
		return nil, "", status.Errorf(codes.Aborted,
			"failed to initialize K8S client configuration: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		return nil, "", status.Errorf(codes.Aborted,
			"failed to initialize K8S client: %s", err)
	}

	return clientset, ns, nil
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
