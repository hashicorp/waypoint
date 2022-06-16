package k8s

import (
	"fmt"

	"github.com/hashicorp/waypoint/internal/clierrors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sConfig struct {
	namespace string `hcl:"namespace,optional"`

	k8sContext string `hcl:"k8s_context,optional"`
}

type K8sInstaller struct {
	config K8sConfig
}

// newClient creates a new K8S client based on the configured settings.
func (i *K8sInstaller) NewClient() (*kubernetes.Clientset, error) {
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
