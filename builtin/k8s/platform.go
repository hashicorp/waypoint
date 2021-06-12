package k8s

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/internal/clierrors"
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

func (p *Platform) StatusFunc() interface{} {
	return p.Status
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

func (p *Platform) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getClientset),
		resource.WithResource(resource.NewResource(
			resource.WithName("deployment"),
			resource.WithState(&Resource_Deployment{}),
			resource.WithCreate(p.resourceDeploymentCreate),
			resource.WithDestroy(p.resourceDeploymentDestroy),
		)),
	)
	return nil
}

// getClientset is a value provider for our resource manager and provides
// the connection information used by resources to interact with Kubernetes.
func (p *Platform) getClientset() (*clientsetInfo, error) {
	// Get our client
	clientSet, ns, config, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	return &clientsetInfo{
		Clientset: clientSet,
		Namespace: ns,
		Config:    config,
	}, nil
}

// resourceDeploymentCreate creates the Kubernetes deployment.
func (p *Platform) resourceDeploymentCreate(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,

	result *Deployment,
	state *Resource_Deployment,
	csinfo *clientsetInfo,
	sg terminal.StepGroup,
) error {
	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	step := sg.Add("")
	defer func() { step.Abort() }()
	step.Update("Kubernetes client connected to %s with namespace %s", csinfo.Config.Host, ns)
	step.Done()

	step = sg.Add("Preparing deployment...")

	clientSet := csinfo.Clientset
	deployClient := clientSet.AppsV1().Deployments(ns)

	// Determine if we have a deployment that we manage already
	create := false
	deployment, err := deployClient.Get(ctx, result.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		deployment = result.newDeployment(result.Name)
		create = true
		err = nil
	}
	if err != nil {
		return err
	}

	// Setup our port configuration
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
		return fmt.Errorf("Cannot define both 'service_port' and 'ports'. Use" +
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
	var resourceLimits = make(map[corev1.ResourceName]k8sresource.Quantity)
	var resourceRequests = make(map[corev1.ResourceName]k8sresource.Quantity)

	for k, v := range p.config.Resources {
		if strings.HasPrefix(k, "limits_") {
			limitKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(limitKey[1])

			quantity, err := k8sresource.ParseQuantity(v)
			if err != nil {
				return err
			}
			resourceLimits[resourceName] = quantity
		} else if strings.HasPrefix(k, "requests_") {
			reqKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(reqKey[1])

			quantity, err := k8sresource.ParseQuantity(v)
			if err != nil {
				return err
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

	container := corev1.Container{
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
	}

	if p.config.Pod != nil && p.config.Pod.Container != nil {
		containerCfg := p.config.Pod.Container
		if containerCfg.Command != nil {
			container.Command = *containerCfg.Command
		}

		if containerCfg.Args != nil {
			container.Args = *containerCfg.Args
		}
	}

	// Update the deployment with our spec
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{container},
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

	if len(p.config.ScratchSpace) > 0 {
		for idx, scratchSpaceLocation := range p.config.ScratchSpace {
			scratchName := fmt.Sprintf("scratch-%d", idx)
			deployment.Spec.Template.Spec.Volumes = append(
				deployment.Spec.Template.Spec.Volumes,
				corev1.Volume{
					Name: scratchName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			)

			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
				corev1.VolumeMount{
					Name:      scratchName,
					MountPath: scratchSpaceLocation,
				},
			)
		}
	}

	if p.config.Pod != nil {
		// Configure Pod
		podConfig := p.config.Pod
		if podConfig.SecurityContext != nil {
			secCtx := podConfig.SecurityContext
			// Configure Pod Security Context
			deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
				RunAsUser:    secCtx.RunAsUser,
				RunAsNonRoot: secCtx.RunAsNonRoot,
				FSGroup:      secCtx.FsGroup,
			}
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
		saClient := clientSet.CoreV1().ServiceAccounts(ns)
		saCreate := false
		serviceAccount, err := saClient.Get(ctx, p.config.ServiceAccount, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			serviceAccount = newServiceAccount(p.config.ServiceAccount)
			saCreate = true
			err = nil
		}
		if err != nil {
			return err
		}

		if saCreate {
			serviceAccount, err = saClient.Create(ctx, serviceAccount, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}

	dc := clientSet.AppsV1().Deployments(ns)

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
		return err
	}

	// We successfully created or updated, so set the name on our state so
	// that if we error, we'll partially clean up properly. THIS IS IMPORTANT.
	state.Name = result.Name

	step.Done()
	step = sg.Add("Waiting for deployment...")

	ps := clientSet.CoreV1().Pods(ns)
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
		return err
	}

	step.Update("Deployment successfully rolled out!")
	step.Done()

	return nil
}

// Destroy deletes the K8S deployment.
func (p *Platform) resourceDeploymentDestroy(
	ctx context.Context,
	state *Resource_Deployment,
	sg terminal.StepGroup,
	csinfo *clientsetInfo,
) error {
	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	step := sg.Add("")
	defer func() { step.Abort() }()
	step.Update("Kubernetes client connected to %s with namespace %s", csinfo.Config.Host, ns)
	step.Done()

	step = sg.Add("Deleting deployment...")
	deployclient := csinfo.Clientset.AppsV1().Deployments(ns)
	if err := deployclient.Delete(ctx, state.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	step.Done()

	return nil
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
	sg := ui.StepGroup()
	defer sg.Wait()

	// Create our resource manager and create
	rm := p.resourceManager(log)
	if err := rm.CreateAll(
		ctx, log, sg, ui,
		src, img, deployConfig, &result,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	result.ResourceState = rm.State()

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
	defer sg.Wait()

	rm := p.resourceManager(log)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource("deployment").SetState(&Resource_Deployment{
			Name: deployment.Name,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	// Destroy
	return rm.DestroyAll(ctx, log, sg, ui)
}

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Gathering health report for Kubernetes platform...")
	defer func() { step.Abort() }() // Defer in func in case more steps are added to this func in the future

	csInfo, err := p.getClientset()
	if err != nil {
		return nil, err
	}
	clientSet := csInfo.Clientset
	namespace := csInfo.Namespace
	if p.config.Namespace != "" {
		namespace = p.config.Namespace
	}

	podClient := clientSet.CoreV1().Pods(namespace)
	podLabelId := fmt.Sprintf("app=%s", deployment.Name)
	podList, err := podClient.List(ctx, metav1.ListOptions{LabelSelector: podLabelId})
	if err != nil {
		ui.Output(
			"Error listing pods to determine application health: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, err
	}

	// Create our status report
	step.Update("Building status report for running pods...")
	result := buildStatusReport(podList)

	result.GeneratedTime = ptypes.TimestampNow()
	log.Debug("status report complete")

	// update output based on main health state
	step.Update("Finished building report for Kubernetes platform")
	step.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	st.Update("Determining overall container health...")
	if result.Health == sdk.StatusReport_READY {
		st.Step(terminal.StatusOK, fmt.Sprintf("Deployment %q is reporting ready!", deployment.Name))
	} else {
		if result.Health == sdk.StatusReport_PARTIAL {
			st.Step(terminal.StatusWarn, fmt.Sprintf("Deployment %q is reporting partially available!", deployment.Name))
		} else {
			st.Step(terminal.StatusError, fmt.Sprintf("Deployment %q is reporting not ready!", deployment.Name))
		}

		// Extra advisory wording to let user know that the deployment could be still starting up
		// if the report was generated immediately after it was deployed or released.
		st.Step(terminal.StatusWarn, mixedHealthWarn)
	}

	// More UI detail for non-ready resources
	for _, resource := range result.Resources {
		if resource.Health != sdk.StatusReport_READY {
			st.Step(terminal.StatusWarn, fmt.Sprintf("Resource %q is reporting %q", resource.Name, resource.Health.String()))
		}
	}

	return &result, nil
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

	// An array of paths to directories that will be mounted as EmptyDirVolumes in the pod
	// to store temporary data.
	ScratchSpace []string `hcl:"scratch_path,optional"`

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

	// Pod describes the configuration for the pod
	Pod *Pod `hcl:"pod,block"`
}

// Pod describes the configuration for the pod
type Pod struct {
	SecurityContext *PodSecurityContext `hcl:"security_context,block"`
	Container       *Container          `hcl:"container,block"`
}

// Container describes the commands and arguments for a container config
type Container struct {
	Command *[]string `hcl:"command"`
	Args    *[]string `hcl:"args"`
}

// PodSecurityContext describes the security config for the Pod
type PodSecurityContext struct {
	RunAsUser    *int64 `hcl:"run_as_user"`
	RunAsNonRoot *bool  `hcl:"run_as_non_root"`
	FsGroup      *int64 `hcl:"fs_group"`
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
		"pod",
		"the configuration for a pod",
		docs.Summary("Pod describes the configuration for a pod when deploying"),
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"container",
				"container describes the commands and arguments for a container config",
				docs.SubFields(func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"command",
						"An array of strings to run for the container",
					)

					doc.SetField(
						"args",
						"An array of string arguments to pass through to the container",
					)
				}),
			)
			doc.SetField(
				"pod_security_context",
				"holds pod-level security attributes and container settings",
				docs.SubFields(func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"run_as_user",
						"The UID to run the entrypoint of the container process",
					)
					doc.SetField(
						"run_as_non_root",
						"Indicates that the container must run as a non-root user",
					)
					doc.SetField(
						"fs_group",
						"A special supplemental group that applies to all containers in a pod",
					)
				}),
			)
		}),
	)

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
	mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current deployment is not ready, however your application
might be available or still starting up.
`)
)

var (
	_ component.Platform         = (*Platform)(nil)
	_ component.PlatformReleaser = (*Platform)(nil)
	_ component.Configurable     = (*Platform)(nil)
	_ component.Documented       = (*Platform)(nil)
	_ component.Destroyer        = (*Platform)(nil)
	_ component.Status           = (*Platform)(nil)
)
