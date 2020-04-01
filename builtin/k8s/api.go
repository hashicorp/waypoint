package k8s

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// clientset returns a K8S clientset and configured namespace.
func clientset(kubeconfig, context string) (*kubernetes.Clientset, string, error) {
	// Path to the kube config file
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
			CurrentContext: context,
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
