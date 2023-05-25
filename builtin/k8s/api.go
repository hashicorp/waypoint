package k8s

import (
	"io/ioutil"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// clientsetInfo is used by some functions to collect the results returned by
// the clientset functions.
type clientsetInfo struct {
	Clientset *kubernetes.Clientset
	Namespace string
	Config    *rest.Config
}

// Clientset returns a K8S clientset and configured namespace. This will
// attempt to use in-cluster auth if available if kubeconfig is not explicitly
// specified. Otherwise, this will fall back to out of cluster auth.
func Clientset(kubeconfig, context string) (*kubernetes.Clientset, string, *rest.Config, error) {
	if kubeconfig == "" {
		cs, ns, c, err := ClientsetInCluster()
		if err == nil {
			return cs, ns, c, nil
		}

		// If we got an error about not being in the cluster, that's okay
		// and fall back to out of cluster auth. If we got any other error
		// though then report an error.
		if err != rest.ErrNotInCluster {
			return nil, "", nil, err
		}
	}

	return ClientsetOutOfCluster(kubeconfig, context)
}

// ClientsetOutOfCluster loads a Kubernetes clientset using only a kubeconfig.
func ClientsetOutOfCluster(kubeconfig, context string) (*kubernetes.Clientset, string, *rest.Config, error) {
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

// ClientsetInCluster returns a K8S clientset and configured namespace for
// in-cluster usage.
func ClientsetInCluster() (*kubernetes.Clientset, string, *rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, "", nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", nil, err
	}

	ns := "default"
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if v := strings.TrimSpace(string(data)); len(v) > 0 {
			ns = v
		}
	}

	return clientset, ns, config, nil
}
