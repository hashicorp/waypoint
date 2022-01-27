package k8s

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

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

	// The name of the Kubernetes secret to use to pull images started by
	// this task.
	ImageSecret string `hcl:"image_secret,optional"`

	// ServiceAccount is the name of the Kubernetes service account to apply to the
	// application deployment. This is useful to apply Kubernetes RBAC to the pod.
	ServiceAccount string `hcl:"service_account,optional"`

	// Set an explicit pull policy for this task launching. By default
	// we use "PullIfNotPresent" unless the image tag is "latest" when we
	// use "Always".
	PullPolicy string `hcl:"image_pull_policy,optional"`

	// The namespace to use for launching this task in Kubernetes
	Namespace string `hcl:"namespace,optional"`
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

	doc.SetField(
		"image_secret",
		"name of the Kubernetes secret to use for the image",
		docs.Summary(
			"this references an existing secret; Waypoint does not create this secret",
		),
	)

	doc.SetField(
		"service_account",
		"service account name to be added to the application pod",
		docs.Summary(
			"service account is the name of the Kubernetes service account to add to the pod.",
			"This is useful to apply Kubernetes RBAC to the application.",
		),
	)

	doc.SetField(
		"image_pull_policy",
		"pull policy to use for the task container image",
	)

	doc.SetField(
		"namespace",
		"namespace in which to launch task",
	)

	return doc, nil
}

// Config implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to docker to stop the container created previously
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	// Purposely do nothing. We leverage the job TTL feature in Kube 1.19+
	// so that Kubernetes automatically deletes old jobs after they complete
	// running.
	//
	// In the future, we may want to get more clever about this and explicitly
	// delete jobs under certain conditions, but for now we leave them around
	// and let K8S clean it up
	return nil
}

// StartTask creates a docker container for the task.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	// Get our client
	clientSet, ns, _, err := Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	// Generate an ID for our pod name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// This must be lowercase because K8S enforces that resource names
	// are lowercase.
	name := strings.ToLower(fmt.Sprintf("waypoint-task-%s", id.String()))

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
	// This by default means kubernetes will always pull so that latest is used.
	pullPolicy := corev1.PullIfNotPresent
	if v := p.config.PullPolicy; v != "" {
		pullPolicy = corev1.PullPolicy(v)
	} else if t, ok := named.(reference.Tagged); ok && t.Tag() == "latest" {
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
		Command:         tli.Entrypoint,
		Args:            tli.Arguments,
		Env:             env,
		Resources:       resourceRequirements,
	}

	// Determine our image pull secret
	var pullSecrets []corev1.LocalObjectReference
	if p.config.ImageSecret != "" {
		pullSecrets = []corev1.LocalObjectReference{
			{
				Name: p.config.ImageSecret,
			},
		}
	}

	// Get our jobs client and create our job
	jobsClient := clientSet.BatchV1().Jobs(ns)
	_, err = jobsClient.Create(ctx, &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		Spec: batchv1.JobSpec{
			Parallelism:             pointer.Int32(1),
			Completions:             pointer.Int32(1),
			BackoffLimit:            pointer.Int32(3),
			TTLSecondsAfterFinished: pointer.Int32(600),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: p.config.ServiceAccount,
					Containers:         []corev1.Container{container},
					ImagePullSecrets:   pullSecrets,
					RestartPolicy:      corev1.RestartPolicyOnFailure,
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// NOTE(mitchellh): In the future, we can probably do some waiting
	// here to check that the pods are successfully starting. This will
	// result in a more immediate error message.

	return &TaskInfo{
		Id: name,
	}, nil
}
