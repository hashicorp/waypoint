package k8s

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"

	// TODO Evaluate if this should remain as a default 3000 to another port.
	DefaultServicePort = 3000
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

// DefaultReleaserFunc implements component.PlatformReleaser
func (p *Platform) DefaultReleaserFunc() interface{} {
	var rc ReleaserConfig
	if err := mapstructure.Decode(p.config, &rc); err != nil {
		// shouldn't happen
		panic("error decoding config: " + err.Error())
	}

	return func() *Releaser {
		return &Releaser{
			config: rc,
		}
	}
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

	sg := ui.StepGroup()
	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Get our client
	clientset, ns, config, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	// Override namespace if set
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	step.Update("Kubernetes client connected to %s with namespace %s", config.Host, ns)
	step.Done()

	step = sg.Add("Preparing deployment...")

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

	if p.config.ServicePort == 0 && p.config.Ports == nil {
		// nothing defined, set up the defaults
		p.config.Ports = make([]map[string]string, 1)
		p.config.Ports[0] = map[string]string{"port": strconv.Itoa(DefaultServicePort), "name": "http"}
	} else if p.config.ServicePort > 0 && p.config.Ports == nil {
		// old ServicePort var is used, so set it up in our Ports map to be used
		p.config.Ports = make([]map[string]string, 1)
		p.config.Ports[0] = map[string]string{"port": strconv.Itoa(int(p.config.ServicePort)), "name": "http"}
	} else if p.config.ServicePort > 0 && len(p.config.Ports) > 0 {
		// both defined, this is an error
		return nil, fmt.Errorf("Cannot define both 'service_port' and 'ports'. Use" +
			" 'ports' for configuring multiple container ports.")
	}

	// Build our env vars
	env := []corev1.EnvVar{
		{
			Name:  "PORT",
			Value: fmt.Sprint(p.config.Ports[0]["port"]),
		},
	}

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
	// Version label duplicates "labelId" to support services like Istio that
	// expect pods to be labled with 'version'
	deployment.Spec.Template.Labels["version"] = result.Id

	// Apply user defined labels
	for k, v := range p.config.Labels {
		deployment.Spec.Template.Labels[k] = v
	}

	// If the user is using the latest tag, then don't specify an overriding pull policy.
	// This by default means kubernetes will always pull so that latest is useful.
	pullPolicy := corev1.PullIfNotPresent
	if img.Tag == "latest" {
		pullPolicy = ""
	}

	// Get container resource limits and requests
	var resourceLimits = make(map[corev1.ResourceName]resource.Quantity)
	var resourceRequests = make(map[corev1.ResourceName]resource.Quantity)

	for k, v := range p.config.Resources {
		if strings.HasPrefix(k, "limits_") {
			limitKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(limitKey[1])

			quantity, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			resourceLimits[resourceName] = quantity
		} else if strings.HasPrefix(k, "requests_") {
			reqKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(reqKey[1])

			quantity, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			resourceRequests[resourceName] = quantity
		} else {
			log.Warn("ignoring unrecognized k8s resources key: %q", k)
		}
	}

	resourceRequirements := corev1.ResourceRequirements{
		Limits:   resourceLimits,
		Requests: resourceRequests,
	}

	containerPorts := make([]corev1.ContainerPort, len(p.config.Ports))
	for i, cp := range p.config.Ports {
		hostPort, _ := strconv.ParseInt(cp["host_port"], 10, 32)
		port, _ := strconv.ParseInt(cp["port"], 10, 32)

		containerPorts[i] = corev1.ContainerPort{
			Name:          cp["name"],
			ContainerPort: int32(port),
			HostPort:      int32(hostPort),
			HostIP:        cp["host_ip"],
			Protocol:      corev1.ProtocolTCP,
		}
	}

	// assume the first port defined is the 'main' port to use
	defaultPort := int(containerPorts[0].ContainerPort)

	initialDelaySeconds := int32(5)
	timeoutSeconds := int32(5)
	failureThreshold := int32(5)
	if p.config.Probe != nil {
		if p.config.Probe.InitialDelaySeconds != 0 {
			initialDelaySeconds = int32(p.config.Probe.InitialDelaySeconds)
		}
		if p.config.Probe.TimeoutSeconds != 0 {
			timeoutSeconds = int32(p.config.Probe.TimeoutSeconds)
		}
		if p.config.Probe.FailureThreshold != 0 {
			failureThreshold = int32(p.config.Probe.FailureThreshold)
		}
	}

	// Update the deployment with our spec
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            result.Name,
				Image:           img.Name(),
				ImagePullPolicy: pullPolicy,
				Ports:           containerPorts,
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(defaultPort),
						},
					},
					InitialDelaySeconds: initialDelaySeconds,
					TimeoutSeconds:      timeoutSeconds,
					FailureThreshold:    failureThreshold,
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(defaultPort),
						},
					},
					InitialDelaySeconds: initialDelaySeconds,
					TimeoutSeconds:      timeoutSeconds,
				},
				Env:       env,
				Resources: resourceRequirements,
			},
		},
	}

	// Override the default TCP socket checks if we have a probe path
	if p.config.ProbePath != "" {
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: p.config.ProbePath,
					Port: intstr.FromInt(defaultPort),
				},
			},
			InitialDelaySeconds: initialDelaySeconds,
			TimeoutSeconds:      timeoutSeconds,
			FailureThreshold:    failureThreshold,
		}

		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: p.config.ProbePath,
					Port: intstr.FromInt(defaultPort),
				},
			},
			InitialDelaySeconds: initialDelaySeconds,
			TimeoutSeconds:      timeoutSeconds,
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

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.Annotations = p.config.Annotations

	if p.config.ServiceAccount != "" {
		deployment.Spec.Template.Spec.ServiceAccountName = p.config.ServiceAccount

		// Determine if we need to make a service account
		saClient := clientset.CoreV1().ServiceAccounts(ns)
		saCreate := false
		serviceAccount, err := saClient.Get(ctx, p.config.ServiceAccount, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			serviceAccount = newServiceAccount(p.config.ServiceAccount)
			saCreate = true
			err = nil
		}
		if err != nil {
			return nil, err
		}

		if saCreate {
			serviceAccount, err = saClient.Create(ctx, serviceAccount, metav1.CreateOptions{})
			if err != nil {
				return nil, err
			}
		}
	}

	dc := clientset.AppsV1().Deployments(ns)

	// Create/update
	if create {
		log.Debug("no existing deployment, creating a new one")
		step.Update("Creating deployment...")
		deployment, err = dc.Create(ctx, deployment, metav1.CreateOptions{})
	} else {
		log.Debug("updating deployment")
		step.Update("Updating deployment...")
		deployment, err = dc.Update(ctx, deployment, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, err
	}

	step.Done()
	step = sg.Add("Waiting for deployment...")

	ps := clientset.CoreV1().Pods(ns)
	podLabelId := fmt.Sprintf("%s=%s", labelId, result.Id)

	var (
		lastStatus    time.Time
		detectedError string
		k8error       string
		reportedError bool
	)

	timeout := 10 * time.Minute

	// Wait on the Pod to start
	err = wait.PollImmediate(2*time.Second, timeout, func() (bool, error) {
		dep, err := dc.Get(ctx, result.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if time.Since(lastStatus) > 10*time.Second {
			step.Update(fmt.Sprintf(
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

		pods, err := ps.List(ctx, metav1.ListOptions{
			LabelSelector: podLabelId,
		})

		if err != nil {
			return false, nil
		}

		for _, p := range pods.Items {
			for _, cs := range p.Status.ContainerStatuses {
				if cs.Ready {
					continue
				}

				if cs.State.Waiting != nil {
					// TODO: handle other pod failures here
					if cs.State.Waiting.Reason == "ImagePullBackOff" ||
						cs.State.Waiting.Reason == "ErrImagePull" {
						detectedError = "Pod unable to access Docker image"
						k8error = cs.State.Waiting.Message
					}
				}
			}
		}

		if detectedError != "" && !reportedError {
			// we use ui output here instead of a step group, otherwise the warning
			// gets swallowed up on the next poll iteration
			ui.Output("Detected pods having an issue starting - %s: %s",
				detectedError, k8error, terminal.WithWarningStyle())
			reportedError = true

			// force a faster rerender
			lastStatus = time.Time{}
		}

		return false, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			err = fmt.Errorf("Deployment was not able to start pods after %s", timeout)
		}
		return nil, err
	}

	step.Update("Deployment successfully rolled out!")
	step.Done()

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Get our client
	clientset, ns, config, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return err
	}

	// Override namespace if set
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	step.Update("Kubernetes client connected to %s with namespace %s", config.Host, ns)
	step.Done()
	step = sg.Add("Deleting deployment...")

	deployclient := clientset.AppsV1().Deployments(ns)
	if err := deployclient.Delete(ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	step.Done()
	return nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// Annotations are added to the pod spec of the deployed application.  This is
	// useful when using mutating webhook admission controllers to further process
	// pod events.
	Annotations map[string]string `hcl:"annotations,optional"`

	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`

	// The number of replicas of the service to maintain. If this number is maintained
	// outside waypoint, for instance by a pod autoscaler, do not set this variable.
	Count int32 `hcl:"replicas,optional"`

	// The name of the Kubernetes secret to use to pull the image stored
	// in the registry.
	// TODO This maybe should be required because the vast majority of deployments
	// will be against private images.
	ImageSecret string `hcl:"image_secret,optional"`

	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// A map of key vals to label the deployed Pod and Deployment with.
	Labels map[string]string `hcl:"labels,optional"`

	// Namespace is the Kubernetes namespace to target the deployment to.
	Namespace string `hcl:"namespace,optional"`

	// A full resource of options to define ports for your service running on the container
	// Defaults to port 3000.
	Ports []map[string]string `hcl:"ports,optional"`

	// If set, this is the HTTP path to request to test that the application
	// is up and running. Without this, we only test that a connection can be
	// made to the port.
	ProbePath string `hcl:"probe_path,optional"`

	// Probe details for describing a health check to be performed against a container.
	Probe *Probe `hcl:"probe,block"`

	// Optionally define various resources limits for kubernetes pod containers
	// such as memory and cpu.
	Resources map[string]string `hcl:"resources,optional"`

	// A path to a directory that will be created for the service to store
	// temporary data.
	ScratchSpace string `hcl:"scratch_path,optional"`

	// ServiceAccount is the name of the Kubernetes service account to apply to the
	// application deployment. This is useful to apply Kubernetes RBAC to the pod.
	ServiceAccount string `hcl:"service_account,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to DefaultServicePort const.
	// NOTE: Ports and ServicePort cannot both be defined
	ServicePort uint `hcl:"service_port,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// Time in seconds to wait before performing the initial liveness and readiness probes.
	// Defaults to 5 seconds.
	InitialDelaySeconds uint `hcl:"initial_delay,optional"`

	// Time in seconds before the probe fails.
	// Defaults to 5 seconds.
	TimeoutSeconds uint `hcl:"timeout,optional"`

	// Number of times a liveness probe can fail before the container is killed.
	// FailureThreshold * TimeoutSeconds should be long enough to cover your worst
	// case startup times. Defaults to 5 failures.
	FailureThreshold uint `hcl:"failure_threshold,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy the application into a Kubernetes cluster using Deployment objects")

	doc.Example(`
deploy "kubernetes" {
	image_secret = "registry_secret"
	replicas = 3
	probe_path = "/_healthz"
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
		"replicas",
		"the number of replicas to maintain",
		docs.Summary(
			"if the replica count is maintained outside waypoint,",
			"for instance by a pod autoscaler, do not set this variable",
		),
	)

	doc.SetField(
		"resources",
		"a map of resource limits and requests to apply to a pod on deploy",
		docs.Summary(
			"resource limits and requests for a pod. limits and requests options "+
				"must start with either 'limits\\_' or 'requests\\_'. Any other options "+
				"will be ignored.",
		),
	)

	doc.SetField(
		"ports",
		"a map of ports and options that the application is listening on",
		docs.Summary(
			"used to define and expose multiple ports that the application is",
			"listening on for the container in use. Available keys are 'port', 'name'",
			", 'host_port', and 'host_ip'. Ports defined will be TCP protocol.",
		),
	)

	doc.SetField(
		"probe_path",
		"the HTTP path to request to test that the application is running",
		docs.Summary(
			"without this, the test will simply be that the application has bound to the port",
		),
	)

	doc.SetField(
		"probe",
		"configuration to control liveness and readiness probes",
		docs.Summary("Probe describes a health check to be performed against a ",
			"container to determine whether it is alive or ready to receive traffic."),
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"initial_delay",
				"time in seconds to wait before performing the initial liveness and readiness probes",
				docs.Default("5"),
			)

			doc.SetField(
				"timeout",
				"time in seconds before the probe fails",
				docs.Default("5"),
			)

			doc.SetField(
				"failure_threshold",
				"number of times a liveness probe can fail before the container is killed",
				docs.Summary(
					"failureThreshold * TimeoutSeconds should be long enough to cover your worst case startup times",
				),
				docs.Default("5"),
			)

		}),
	)

	doc.SetField(
		"scratch_path",
		"a path for the service to store temporary data",
		docs.Summary(
			"a path to a directory that will be created for the service to store temporary data using tmpfs",
		),
	)

	doc.SetField(
		"image_secret",
		"name of the Kubernetes secrete to use for the image",
		docs.Summary(
			"this references an existing secret, waypoint does not create this secret",
		),
	)

	doc.SetField(
		"static_environment",
		"environment variables to control broad modes of the application",
		docs.Summary(
			"environment variables that are meant to configure the application in a static",
			"way. This might be control an image that has multiple modes of operation,",
			"selected via environment variable. Most configuration should use the waypoint",
			"config commands",
		),
	)

	doc.SetField(
		"service_port",
		"the TCP port that the application is listening on",
		docs.Default(fmt.Sprint(DefaultServicePort)),
		docs.Summary(
			"by default, this config variable is used for exposing a single port for",
			"the container in use. For multi-port configuration, use 'ports' instead.",
		),
	)

	doc.SetField(
		"annotations",
		"annotations to be added to the application pod",
		docs.Summary(
			"annotations are added to the pod spec of the deployed application. This is",
			"useful when using mutating webhook admission controllers to further process",
			"pod events.",
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
		"labels",
		"a map of key value labels to apply to the deployment pod",
	)

	doc.SetField(
		"namespace",
		"namespace to target deployment into",
		docs.Summary(
			"namespace is the name of the Kubernetes namespace to apply the deployment in.",
			"This is useful to create deployments in non-default namespaces without creating kubeconfig contexts for each",
		),
	)

	return doc, nil
}

var (
	_ component.Platform         = (*Platform)(nil)
	_ component.PlatformReleaser = (*Platform)(nil)
	_ component.Configurable     = (*Platform)(nil)
	_ component.Documented       = (*Platform)(nil)
	_ component.Destroyer        = (*Platform)(nil)
)
