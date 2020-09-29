package k8s

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// clientset returns a K8S clientset and configured namespace.
func clientset(kubeconfig, context string) (*kubernetes.Clientset, string, *rest.Config, error) {
	loader := clientcmd.NewDefaultClientConfigLoadingRules()

	// Path to the kube config file
	if kubeconfig != "" {
		loader.ExplicitPath = kubeconfig
	}

	// Build our config and client
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loader,
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)

	// Get our configured namespace
	ns, _, err := config.Namespace()
	if err != nil {
		return nil, "", nil, status.Errorf(codes.Aborted,
			"failed to initialize K8S client configuration: %s", err)
	}

	clientconfig, err := config.ClientConfig()
	if err != nil {
		return nil, "", nil, status.Errorf(codes.Aborted,
			"failed to initialize K8S client configuration: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		return nil, "", nil, status.Errorf(codes.Aborted,
			"failed to initialize K8S client: %s", err)
	}

	return clientset, ns, clientconfig, nil
}
