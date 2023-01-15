package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sConfig struct {
	KubeconfigPath       string `hcl:"kubeconfig,optional"`
	K8sContext           string `hcl:"context,optional"`
	Version              string `hcl:"version,optional"`
	Timeout              int    `hcl:"timeout,optional"`
	Namespace            string `hcl:"namespace,optional"`
	RunnerImage          string `hcl:"runner_image,optional"`
	CpuRequest           string `hcl:"runner_cpu_request,optional"`
	MemRequest           string `hcl:"runner_mem_request,optional"`
	CreateServiceAccount bool   `hcl:"odr_service_account_init,optional"`
	OdrImage             string `hcl:"odr_image"`

	// Required for backwards compatibility
	ImagePullPolicy string `hcl:"image_pull_policy,optional"`
	CpuLimit        string `hcl:"cpu_limit,optional"`
	MemLimit        string `hcl:"mem_limit,optional"`
	ImagePullSecret string `hcl:"image_pull_secret,optional"`

	// Used for serverinstall
	ServerImage        string            `hcl:"server_image,optional"`
	ServiceAnnotations map[string]string `hcl:"service_annotations,optional"`

	OdrServiceAccount     string `hcl:"odr_service_account,optional"`
	OdrServiceAccountInit bool   `hcl:"odr_service_account_init,optional"`

	AdvertiseInternal bool   `hcl:"advertise_internal,optional"`
	StorageClassName  string `hcl:"storageclassname,optional"`
	StorageRequest    string `hcl:"storage_request,optional"`
	SecretFile        string `hcl:"secret_file,optional"`
	KubeConfigPath    string `hcl:"kubeconfig_path,optional"`
}

// newClient creates a new K8S client based on the configured settings.
func NewClient(config K8sConfig) (*kubernetes.Clientset, error) {
	// Build our K8S client.
	configOverrides := &clientcmd.ConfigOverrides{}
	if config.K8sContext != "" {
		configOverrides = &clientcmd.ConfigOverrides{
			CurrentContext: config.K8sContext,
		}
	}
	newCmdConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		configOverrides,
	)

	// Discover the current target namespace in the user's config so if they
	// run kubectl commands waypoint will show up. If we use the default namespace
	// they might not see the objects we've created.
	if config.Namespace == "" {
		namespace, _, err := newCmdConfig.Namespace()
		if err != nil {
			return nil, fmt.Errorf(
				"error getting namespace from client config: %s",
				clierrors.Humanize(err),
			)
		}

		config.Namespace = namespace
	}

	clientconfig, err := newCmdConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf(
			"error initializing kubernetes client: %s",
			clierrors.Humanize(err),
		)
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		return nil, fmt.Errorf(
			"error initializing kubernetes client: %s",
			clierrors.Humanize(err),
		)
	}

	return clientset, nil
}

// Takes list options and cleans up any PVCs found in the query of resources from the kubernetes api.
// Useful for cleaning up PVCs left behind by statefulsets deployed via helm.
func CleanPVC(ctx context.Context, ui terminal.UI, log hclog.Logger, listOptions metav1.ListOptions, config K8sConfig) error {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Deleting PVCs...")
	defer func() { s.Abort() }()

	clientset, err := NewClient(config)
	if err != nil {
		return err
	}
	pvcClient := clientset.CoreV1().PersistentVolumeClaims(config.Namespace)
	if list, err := pvcClient.List(ctx, listOptions); err != nil {
		return err
	} else if len(list.Items) > 0 {
		// Add watcher for waiting for persistent volume clean up
		w, err := pvcClient.Watch(ctx, listOptions)
		if err != nil {
			return err
		}

		// Delete the PVCs
		if err = pvcClient.DeleteCollection(
			ctx,
			metav1.DeleteOptions{},
			listOptions,
		); err != nil {
			s.Update("Unable to delete PVCs")
			s.Abort()
			return err
		}

		// Wait until the persistent volumes are cleaned up
		err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
			select {
			case wCh := <-w.ResultChan():
				if wCh.Type == watch.Deleted {
					w.Stop()
					return true, nil
				}
				return false, nil
			default:
				return false, nil
			}
		})
		if err != nil {
			s.Update("Deleted PVCs not cleaned up after 10 minutes")
			s.Abort()
			return err
		}
	}
	s.Update("Persistent volume claims cleaned up")
	s.Status(terminal.StatusOK)
	s.Done()
	return nil
}
