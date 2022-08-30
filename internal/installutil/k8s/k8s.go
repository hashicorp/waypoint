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
				"error getting namespace from client config: %s",
				clierrors.Humanize(err),
			)
		}

		i.config.namespace = namespace
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
func (i *K8sInstaller) CleanPVC(ctx context.Context, ui terminal.UI, log hclog.Logger, listOptions metav1.ListOptions) error {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Deleting PVCs...")
	defer func() { s.Abort() }()

	clientset, err := i.NewClient()
	if err != nil {
		return err
	}
	pvcClient := clientset.CoreV1().PersistentVolumeClaims(i.config.namespace)
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
			s.Update("Unable to delete PVCs")
			s.Abort()
			return err
		}
	}
	s.Update("Persistent volume claims cleaned up")
	s.Status(terminal.StatusOK)
	s.Done()
	return nil
}
