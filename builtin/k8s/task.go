// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package k8s

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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

// WatchTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) WatchTaskFunc() interface{} {
	return p.WatchTask
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

	// Optionally define various cpu resource limits and requests for kubernetes pod containers
	CPU *ResourceConfig `hcl:"cpu,block"`

	// Optionally define various memory resource limits and requests for kubernetes pod containers
	Memory *ResourceConfig `hcl:"memory,block"`

	// Optionally define various ephemeral storage resource limits and requests for kubernetes pod containers
	EphemeralStorage *ResourceConfig `hcl:"ephemeral_storage,block"`

	// How long WatchTask should wait for a pod to startup. This option is specifically
	// wordy because it's only for the WatchTask timing out waiting for the pod
	// its watching to start up before it attempts to stream its logs.
	WatchTaskStartupTimeoutSeconds int `hcl:"watchtask_startup_timeout_seconds,optional"`

	// The PodSecurityContext to apply to the pod
	SecurityContext *PodSecurityContext `hcl:"security_context,block"`
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
		"memory",
		"memory resource request to be added to the task container",
	)

	doc.SetField(
		"cpu",
		"cpu resource request to be added to the task container",
	)

	doc.SetField(
		"ephemeral_storage",
		"ephemeral_storage resource request to be added to the task container",
	)

	doc.SetField(
		"image_pull_policy",
		"pull policy to use for the task container image",
	)

	doc.SetField(
		"namespace",
		"namespace in which to launch task",
	)

	doc.SetField(
		"watchtask_startup_timeout_seconds",
		"This option configures how long the WatchTask should wait for a task pod to start-up "+
			"before attempting to stream its logs. If the pod does not start up within "+
			"the given timeout, WatchTask will exit.",
		docs.Default("30"),
	)

	return doc, nil
}

// Config implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to Kubernetes to stop the container created previously
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	// If a job completes and the corresponding pod exits with a "completed"
	// status, we purposely do nothing here. We leverage the job TTL feature in
	// Kube 1.19+ so that Kubernetes automatically deletes old jobs and pods
	// after they complete running.
	//
	// If a Waypoint job is cancelled or otherwise times out, we check for
	// existing Kubernetes jobs and delete them, and clean up any Pending
	// containers.
	clientSet, ns, _, err := Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return err
	}
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	// List pods with this job label
	podsClient := clientSet.CoreV1().Pods(ns)
	pods, err := podsClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", ti.Id),
	})
	// It's not clear from the documentation if an error is returned from the
	// List API call if no jobs are found, so we guard here just in case
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if pods == nil {
		log.Info("no pods found for job, returning", "job_id", ti.Id)
		return nil
	}

	// Find any pods stuck in pending
	var pendingPods []string
	for _, p := range pods.Items {
		if p.Status.Phase == corev1.PodPending {
			pendingPods = append(pendingPods, p.Name)
		}
	}

	// If we've found pending/stuck pods, attempt to clean up
	if len(pendingPods) > 0 {
		// Delete the job. This does *not* delete any running pods that the job
		// created.
		jobsClient := clientSet.BatchV1().Jobs(ns)
		if err := jobsClient.Delete(ctx, ti.Id, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
		for _, name := range pendingPods {
			log.Warn("job pod is in pending phase in StopTask operation, cancelling", "job_id", ti.Id)
			if err := podsClient.Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
				if !errors.IsNotFound(err) {
					return err
				}
			}
		}
	}

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

	// NOTE(briancain): This is here to help kaniko detect that this is a docker container.
	// See https://github.com/GoogleContainerTools/kaniko/blob/7e3954ac734534ce5ce68ad6300a2d3143d82f40/vendor/github.com/genuinetools/bpfd/proc/proc.go#L138
	// for more info.
	log.Warn("temporarily setting 'container=docker' environment variable to patch Kaniko working on Kubernetes 1.23")
	env = append(env, corev1.EnvVar{
		Name:  "container",
		Value: "docker",
	})

	// If the user is using the latest tag, then don't specify an overriding pull policy.
	// This by default means kubernetes will always pull so that latest is used.
	pullPolicy := corev1.PullIfNotPresent
	if v := p.config.PullPolicy; v != "" {
		pullPolicy = corev1.PullPolicy(v)
	} else if t, ok := named.(reference.Tagged); ok && t.Tag() == "latest" {
		pullPolicy = ""
	}

	// Get container resource limits and requests
	resourceLimits := make(map[corev1.ResourceName]k8sresource.Quantity)
	resourceRequests := make(map[corev1.ResourceName]k8sresource.Quantity)

	if p.config.CPU != nil {
		if p.config.CPU.Request != "" {
			q, err := k8sresource.ParseQuantity(p.config.CPU.Request)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse cpu request %q to k8s quantity: %s", p.config.CPU.Request, err)
			}
			resourceRequests[corev1.ResourceCPU] = q
		}

		if p.config.CPU.Limit != "" {
			q, err := k8sresource.ParseQuantity(p.config.CPU.Limit)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse cpu limit %q to k8s quantity: %s", p.config.CPU.Limit, err)
			}
			resourceLimits[corev1.ResourceCPU] = q
		}
	}

	if p.config.Memory != nil {
		if p.config.Memory.Request != "" {
			q, err := k8sresource.ParseQuantity(p.config.Memory.Request)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse memory requested %q to k8s quantity: %s", p.config.Memory.Request, err)
			}
			resourceRequests[corev1.ResourceMemory] = q
		}

		if p.config.Memory.Limit != "" {
			q, err := k8sresource.ParseQuantity(p.config.Memory.Limit)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse memory limit %q to k8s quantity: %s", p.config.Memory.Limit, err)
			}
			resourceLimits[corev1.ResourceMemory] = q
		}
	}

	if p.config.EphemeralStorage != nil {
		if p.config.EphemeralStorage.Request != "" {
			q, err := k8sresource.ParseQuantity(p.config.EphemeralStorage.Request)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse ephemeral-storage requested %q to k8s quantity: %s", p.config.EphemeralStorage.Request, err)
			}
			resourceRequests[corev1.ResourceEphemeralStorage] = q
		}

		if p.config.EphemeralStorage.Limit != "" {
			q, err := k8sresource.ParseQuantity(p.config.EphemeralStorage.Limit)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse ephemeral-storage limit %q to k8s quantity: %s", p.config.EphemeralStorage.Limit, err)
			}
			resourceLimits[corev1.ResourceEphemeralStorage] = q
		}
	}

	var securityContext *corev1.PodSecurityContext = nil
	podSc := p.config.SecurityContext
	if podSc != nil {
		securityContext = &corev1.PodSecurityContext{
			RunAsUser:    podSc.RunAsUser,
			RunAsGroup:   podSc.RunAsGroup,
			RunAsNonRoot: podSc.RunAsNonRoot,
			FSGroup:      podSc.FsGroup,
		}
	}

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
					SecurityContext:    securityContext,
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

// WatchTask implements TaskLauncher
func (p *TaskLauncher) WatchTask(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	ti *TaskInfo,
) (*component.TaskResult, error) {
	// Get our client
	clientSet, ns, _, err := Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	// List pods with this job label
	podsClient := clientSet.CoreV1().Pods(ns)
	pods, err := podsClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", ti.Id),
	})
	// It's not clear from the documentation if an error is returned from the
	// List API call if no jobs are found, so we guard here just in case
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if pods == nil {
		log.Info("no pods found for job label, returning", "job_id", ti.Id)
		ui.Output("no pods found for job-name %q, cannot watch task job...", ti.Id, terminal.WithErrorStyle())

		return nil, nil
	}

	if len(pods.Items) < 1 {
		// This is an error. We found the pods by job name but for some reason there
		// are no actual pods inside the job. This might happen if Waypoint server
		// encounters an internal error or panic mid-task launch.
		return nil, status.Errorf(codes.NotFound, "No pods found in job %q for WatchTask", ti.Id)
	}

	// Assume first one exists for now? Our task launcher for k8s only launches
	// one pod right now.
	pod := pods.Items[0]

	// How long to wait for the pod in question to start before attempting to stream
	// its logs.
	podStartUpTimeout := time.Duration(30 * time.Second)
	if p.config.WatchTaskStartupTimeoutSeconds != 0 {
		podStartUpTimeout = time.Duration(p.config.WatchTaskStartupTimeoutSeconds) * time.Second
	}

	log.Info("waiting for pod to start", "name", pod.Name)

	// Ensure the pod exists before attempting to stream its logs
	err = wait.PollImmediate(time.Second, podStartUpTimeout, func() (bool, error) {
		p, err := clientSet.CoreV1().Pods(ns).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		switch p.Status.Phase {
		case v1.PodRunning, v1.PodFailed, v1.PodSucceeded:
			return true, nil
		case v1.PodPending, v1.PodUnknown:
			return false, nil
		}
		return false, nil
	})
	if err != nil {
		log.Error("pod failed to start before timeout", "pod", pod.Name, "timeout", podStartUpTimeout, "err", err)
		ui.Output("pod %q failed to start before WatchTask timeout %q. err: %s", pod.Name, podStartUpTimeout, err, terminal.WithErrorStyle())

		return nil, err
	}

	log.Info("attempting to stream pod logs")

	// Accumulate our result on this
	var result component.TaskResult

	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		log.Error("failed to read pod log stream", "err", err)
		ui.Output("WatchTask failed to read pod %q log stream: err: %s", pod.Name, err, terminal.WithErrorStyle())

		return nil, err
	}
	defer podLogs.Close()

	log.Info("reading pod logs", "name", pod.Name)

	// Start a func to watch the pod phase. If complete, we no longer need to
	// attempt to stream the pod logs.
	logsDoneCh := make(chan bool)
	go func() {
		defer close(logsDoneCh)
		for {
			p, err := clientSet.CoreV1().Pods(ns).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				log.Warn("error getting pod status", "pod", pod.Name, "err", err)
				ui.Output("Error getting pod %q status: %s", pod.Name, err, terminal.WithErrorStyle())

				logsDoneCh <- true
				return
			}

			switch p.Status.Phase {
			case v1.PodRunning:
				// Pod is still running, so wait
				logsDoneCh <- false
			case v1.PodFailed:
				// Pod has finished
				result.ExitCode = 1
				logsDoneCh <- true
			case v1.PodSucceeded:
				// Pod has finished
				result.ExitCode = 0
				logsDoneCh <- true
				return
			case v1.PodPending, v1.PodUnknown:
				// Unknown state, still wait
				logsDoneCh <- false
			}

			// Sleep a bit so we don't hammer the k8s cluster
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Read the log stream and send to the UI
	for {
		buf := make([]byte, 2000)
		numBytes, err := podLogs.Read(buf)

		if numBytes == 0 {
			// This is here because it doesn't look like the pod log reader ever
			// sends an io.EOF when the logstream is finished. We instead can look
			// at the pod phase for when it's no longer running and the streamer
			// hasn't sent any log bytes.
			select {
			case <-logsDoneCh:
				log.Trace("pod is finished", "pod", pod.Name)

				return &result, nil
			default:
				// Pod is still running, but no log output
				continue
			}

			continue
		}
		// NOTE(briancain): it doesn't seem like the k8s API is sending an EOF
		// Maybe some day this will work.
		if err == io.EOF {
			log.Info("end of stream")
			result.ExitCode = 0
			break
		}
		if err != nil {
			log.Error("got an error streaming pod logs", "pod", pod.Name, "err", err)
			return nil, err
		}

		// stream message to the ui
		message := string(buf[:numBytes])
		log.Info(message)
		ui.Output(message)
	}

	return &result, nil
}

var _ component.TaskLauncher = (*TaskLauncher)(nil)
