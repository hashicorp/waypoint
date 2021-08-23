package k8s

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
)

// TaskLauncher implements the TaskLauncher plugin interface to support
// launching on-demand tasks for the Waypoint server.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// StartTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StartTaskFunc() interface{} {
	return p.StartTask
}

// StopTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StopTaskFunc() interface{} {
	return p.StopTask
}

// TaskLauncherConfig is the configuration structure for the task plugin.
type TaskLauncherConfig struct {
	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`

	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory. If we are running within
	// a pod, we will use the service account authentication if available if
	// this isn't set.
	KubeconfigPath string `hcl:"kubeconfig,optional"`
}

func (p *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(p.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`
Launch a Kubernetes pod for on-demand tasks from the Waypoint server.

This will use the standard Kubernetes environment variables to source
authentication information for Kubernetes. If this is running within Kubernetes
itself (typical for a Kubernetes-based installation), it will use the pod's
service account unless other auth is explicitly given. This allows the task
launcher to work by default.
`)

	doc.Example(`
task {
	use "kubernetes" {}
}
`)

	doc.SetField(
		"kubeconfig",
		"path to the kubeconfig file to use",
		docs.Summary("by default uses from current user's home directory"),
		docs.EnvVar("KUBECONFIG"),
	)

	doc.SetField(
		"context",
		"the kubectl context to use, as defined in the kubeconfig file",
	)

	return doc, nil
}

// TaskLauncher implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to docker to stop the container created previously
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *Task,
) error {
	// Get our client
	clientSet, ns, _, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return err
	}

	// Get our pods client
	podsClient := clientSet.CoreV1().Pods(ns)
	err = podsClient.Delete(ctx, ti.Id, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		err = nil
	}

	return err
}

// StartTask creates a docker container for the task.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*Task, error) {
	// Get our client
	clientSet, ns, _, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	// Generate an ID for our pod name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("waypoint-task-%s", id.String())

	// Parse our image to determine some details later.
	named, err := reference.ParseNormalizedNamed(tli.OciUrl)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unable to parse image name: %s", tli.OciUrl)
	}
	// This ensures that the image has a tag associated with it.
	named = reference.TagNameOnly(named)

	// Build our env vars
	env := []corev1.EnvVar{}
	for k, v := range tli.EnvironmentVariables {
		env = append(env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// If the user is using the latest tag, then don't specify an overriding pull policy.
	// This by default means kubernetes will always pull so that latest is useful.
	pullPolicy := corev1.PullIfNotPresent
	if t, ok := named.(reference.Tagged); ok && t.Tag() == "latest" {
		pullPolicy = ""
	}

	// Get container resource limits and requests
	var resourceLimits = make(map[corev1.ResourceName]k8sresource.Quantity)
	var resourceRequests = make(map[corev1.ResourceName]k8sresource.Quantity)
	resourceRequirements := corev1.ResourceRequirements{
		Limits:   resourceLimits,
		Requests: resourceRequests,
	}

	// Build our container
	container := corev1.Container{
		Name:            name,
		Image:           tli.OciUrl,
		ImagePullPolicy: pullPolicy,
		Command:         []string{tli.Arguments[0]},
		Args:            tli.Arguments,
		Env:             env,
		Resources:       resourceRequirements,
	}

	// Get our pods client and create our pod.
	podsClient := clientSet.CoreV1().Pods(ns)
	_, err = podsClient.Create(ctx, &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{container},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &Task{
		Id: name,
	}, nil
}
