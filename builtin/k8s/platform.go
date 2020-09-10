package k8s

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"
)

// Platform is the Platform implementation for Kubernetes.
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return p.Destroy
}

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return p.ValidateAuth
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return p.Auth
}

func (p *Platform) Auth() error {
	return nil
}

func (p *Platform) ValidateAuth() error {
	return nil
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	c, ok := config.(*Config)
	if !ok {
		// this should never happen
		return fmt.Errorf("Invalid configuration, expected *cloudrun.Config, got %s", reflect.TypeOf(config))
	}

	// set defaults
	if c.ContainerPort < 0 && c.ContainerPort < 65535 {
		c.ContainerPort = 3000
	}

	return nil
}

// DefaultReleaserFunc implements component.PlatformReleaser
func (p *Platform) DefaultReleaserFunc() interface{} {
	return func() *Releaser { return &Releaser{} }
}

// Deploy deploys an image to Kubernetes.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	result.Name = strings.ToLower(fmt.Sprintf("%s-%s", src.App, id))

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get our client
	clientset, ns, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	deployclient := clientset.AppsV1().Deployments(ns)

	// Determine if we have a deployment that we manage already
	create := false
	deployment, err := deployclient.Get(ctx, result.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		deployment = result.newDeployment(result.Name)
		create = true
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// Build our env vars
	env := []corev1.EnvVar{}

	for k, v := range p.config.StaticEnvVars {
		env = append(env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	for k, v := range deployConfig.Env() {
		env = append(env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// If no count is specified, presume that the user is managing the replica
	// count some other way (perhaps manual scaling, perhaps a pod autoscaler).
	// Either way if they don't specify a count, we should be sure we don't send one.
	if p.config.Count > 0 {
		deployment.Spec.Replicas = &p.config.Count
	}

	// Set our ID on the label. We use this ID so that we can have a key
	// to route to multiple versions during release management.
	deployment.Spec.Template.Labels[labelId] = result.Id

	// If the user is using the latest tag, then don't specify an overriding pull policy.
	// This by default means kubernetes will always pull so that latest is useful.
	pullPolicy := corev1.PullIfNotPresent
	if img.Tag == "latest" {
		pullPolicy = ""
	}

	// Update the deployment with our spec
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            result.Name,
				Image:           img.Name(),
				ImagePullPolicy: pullPolicy,
				Ports: []corev1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: int32(p.config.ContainerPort),
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(p.config.ContainerPort),
						},
					},
					InitialDelaySeconds: 5,
					TimeoutSeconds:      5,
					FailureThreshold:    5,
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(p.config.ContainerPort),
						},
					},
					InitialDelaySeconds: 5,
					TimeoutSeconds:      5,
				},
				Env: env,
			},
		},
	}

	// Override the default TCP socket checks if we have a probe path
	if p.config.ProbePath != "" {
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: p.config.ProbePath,
					Port: intstr.FromInt(p.config.ContainerPort),
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			FailureThreshold:    5,
		}

		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: p.config.ProbePath,
					Port: intstr.FromInt(p.config.ContainerPort),
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
		}
	}

	if p.config.ScratchSpace != "" {
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "scratch",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}

		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "scratch",
				MountPath: p.config.ScratchSpace,
			},
		}
	}

	if p.config.ImageSecret != "" {
		deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: p.config.ImageSecret,
			},
		}
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.Annotations[labelNonce] =
		time.Now().UTC().Format(time.RFC3339Nano)

	dc := clientset.AppsV1().Deployments(ns)

	// Create/update
	if create {
		st.Update("Creating deployment...")
		deployment, err = dc.Create(ctx, deployment, metav1.CreateOptions{})
	} else {
		st.Update("Updating deployment...")
		deployment, err = dc.Update(ctx, deployment, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, err
	}

	var lastStatus time.Time

	// Wait on the Pod to start
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		dep, err := dc.Get(ctx, result.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if time.Since(lastStatus) > 10*time.Second {
			st.Update(fmt.Sprintf(
				"Waiting on deployment to become available: %d/%d/%d",
				*dep.Spec.Replicas,
				dep.Status.UnavailableReplicas,
				dep.Status.AvailableReplicas,
			))
			lastStatus = time.Now()
		}

		if dep.Status.AvailableReplicas > 0 {
			return true, nil
		}

		// TODO: Report the statuses and events of the pods that are starting
		// here so that users know why stuff isn't starting. Most commonly here
		// it's going to be an error pulling the image.

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	st.Step(terminal.StatusOK, "Deployment succesfully rolled out!")

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	clientset, ns, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return err
	}

	st.Update("Deleting deployment...")
	deployclient := clientset.AppsV1().Deployments(ns)
	return deployclient.Delete(ctx, deployment.Name, metav1.DeleteOptions{})
}

// Config is the configuration structure for the Platform.
type Config struct {
	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`

	// The number of replicas of the service to maintain. If this number is maintained
	// outside waypoint, for instance by a pod autoscaler, do not set this variable.
	Count int32 `hcl:"replicas,optional"`

	// If set, this is the HTTP path to request to test that the application
	// is up and running. Without this, we only test that a connection can be
	// made to the port.
	ProbePath string `hcl:"probe_path,optional"`

	// A path to a directory that will be created for the service to store
	// temporary data.
	ScratchSpace string `hcl:"scratch_path,optional"`

	// The name of the Kubernetes secret to use to pull the image stored
	// in the registry.
	// TODO This maybe should be required because the vast majority of deployments
	// will be against private images.
	ImageSecret string `hcl:"image_secret,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to port 3000. 
	// TODO Evaluate if this should remain as a default 3000, should be a required field,
	// or default to another port. 
	ContainerPort int `hcl:"container_port,optional"`
}

var (
	_ component.Platform         = (*Platform)(nil)
	_ component.PlatformReleaser = (*Platform)(nil)
	_ component.Configurable     = (*Platform)(nil)
	_ component.Destroyer        = (*Platform)(nil)
)
