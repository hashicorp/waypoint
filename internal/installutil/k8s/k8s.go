package k8s

import (
	"fmt"

	"github.com/hashicorp/waypoint/internal/clierrors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sConfig struct {
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	odrImage              string `hcl:"odr_image,optional"`
	odrServiceAccount     string `hcl:"odr_service_account,optional"`
	odrServiceAccountInit bool   `hcl:"odr_service_account_init,optional"`

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
