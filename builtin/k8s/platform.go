package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"

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

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	c, ok := config.(*Config)
	if !ok {
		// this should never happen
		return status.Errorf(codes.FailedPrecondition, "invalid configuration, expected *k8s.Config, got %T", config)
	}

	if len(c.DeprecatedPorts) > 0 {
		return status.Errorf(codes.InvalidArgument, "invalid kubernetes platform config - the 'ports' field has been deprecated and removed "+
			"in favor of Pod.Container.Port. Refer to Port documentation here: https://www.waypointproject.io/plugins/kubernetes#port")
	}

	// Some fields can be specified on pod.Container and at the top level, for convenience and for
	// historical reasons. Validate that both are not set at once.
	if c.Pod != nil && c.Pod.Container != nil {
		containerOverlayErrStr := "%s defined multiple times - in top level config and in Pod.Container"
		container := c.Pod.Container
		err := utils.Error(validation.ValidateStruct(c,
			validation.Field(&c.Probe,
				validation.Empty.When(container.Probe != nil).Error(fmt.Sprintf(containerOverlayErrStr, "Probe")),
			),
			validation.Field(&c.ProbePath,
				validation.Empty.When(container.ProbePath != "").Error(fmt.Sprintf(containerOverlayErrStr, "ProbePath")),
			),
			validation.Field(&c.Resources,
				validation.Empty.When(container.Resources != nil).Error(fmt.Sprintf(containerOverlayErrStr, "Resources")),
			),
			validation.Field(&c.CPU,
				validation.Empty.When(container.CPU != nil).Error(fmt.Sprintf(containerOverlayErrStr, "CPU")),
			),
			validation.Field(&c.Memory,
				validation.Empty.When(container.Memory != nil).Error(fmt.Sprintf(containerOverlayErrStr, "Memory")),
			),
			validation.Field(&c.StaticEnvVars,
				validation.Empty.When(container.StaticEnvVars != nil).Error(fmt.Sprintf(containerOverlayErrStr, "StaticEnvVars")),
			),
			validation.Field(&c.ServicePort,
				validation.Empty.When(len(container.Ports) > 0).Error("Cannot define both 'service_port' and container 'port'. Use"+
					" container 'port' multiple times for configuring multiple container ports"),
			),
		))
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid kubernetes platform plugin config: %s", err.Error())
		}
	}

	return nil
}

func (p *Platform) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getClientset),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithResource(resource.NewResource(
			resource.WithName(platformName),
			resource.WithState(&Resource_Deployment{}),
			resource.WithCreate(p.resourceDeploymentCreate),
			resource.WithDestroy(p.resourceDeploymentDestroy),
			resource.WithStatus(p.resourceDeploymentStatus),
			resource.WithPlatform(platformName),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("autoscaler"),
			resource.WithState(&Resource_Autoscale{}),
			resource.WithCreate(p.resourceAutoscalerCreate),
			resource.WithDestroy(p.resourceAutoscalerDestroy),
			resource.WithStatus(p.resourceAutoscalerStatus),
			resource.WithPlatform(platformName),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
	)
}

// getClientset is a value provider for our resource manager and provides
// the connection information used by resources to interact with Kubernetes.
func (p *Platform) getClientset() (*clientsetInfo, error) {
	// Get our client
	clientSet, ns, config, err := Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	return &clientsetInfo{
		Clientset: clientSet,
		Namespace: ns,
		Config:    config,
	}, nil
}

func (p *Platform) resourceDeploymentStatus(
	ctx context.Context,
	sg terminal.StepGroup,
	deploymentState *Resource_Deployment,
	clientset *clientsetInfo,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Checking status of the Kubernetes deployment resource...")
	defer s.Abort()

	namespace := p.config.Namespace
	if namespace == "" {
		namespace = clientset.Namespace
	}

	// Get deployment status

	deploymentResource := sdk.StatusReport_Resource{
		Type:                "deployment",
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
	}
	sr.Resources = append(sr.Resources, &deploymentResource)

	deployResp, err := clientset.Clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentState.Name, metav1.GetOptions{})
	if deployResp == nil {
		return status.Errorf(codes.FailedPrecondition, "kubernetes deployment response cannot be empty")
	} else if err != nil {
		if !errors.IsNotFound(err) {
			return status.Errorf(codes.FailedPrecondition, "error getting kubernetes deployment %s: %s", deploymentState.Name, err)
		} else {
			deploymentResource.Name = deploymentState.Name
			deploymentResource.Health = sdk.StatusReport_MISSING
			deploymentResource.HealthMessage = sdk.StatusReport_MISSING.String()

			// Continue on with getting statuses for the rest of our resources
		}
	} else {
		// Found the deployment, and can use it to populate our resource

		var mostRecentCondition v1.DeploymentCondition
		for _, condition := range deployResp.Status.Conditions {
			if condition.LastUpdateTime.Time.After(mostRecentCondition.LastUpdateTime.Time) {
				mostRecentCondition = condition
			}
		}

		// The most recently updated condition isn't always the most pertinent - a healthy deployment
		// can have a "Progressing" most recently updated condition at steady-state.
		// If the deployment exists, we'll mark it as "Ready", and rely on our pod status checks
		// to give more detailed status.
		deployHealth := sdk.StatusReport_READY

		// Redact env vars from containers - they can contain secrets
		for i := 0; i < len(deployResp.Spec.Template.Spec.Containers); i++ {
			deployResp.Spec.Template.Spec.Containers[i].Env = []corev1.EnvVar{}
		}

		deployStateJson, err := json.Marshal(map[string]interface{}{
			"deployment": deployResp,
		})
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "failed to marshal deployment to json: %s", err)
		}

		deploymentResource.Name = deployResp.Name
		deploymentResource.Id = fmt.Sprintf("%s", deployResp.UID)
		deploymentResource.CreatedTime = timestamppb.New(deployResp.CreationTimestamp.Time)
		deploymentResource.Health = deployHealth
		deploymentResource.HealthMessage = fmt.Sprintf("%s: %s", mostRecentCondition.Type, mostRecentCondition.Message)
		deploymentResource.StateJson = string(deployStateJson)
	}

	// Get pod status

	podClient := clientset.Clientset.CoreV1().Pods(namespace)
	podLabelId := fmt.Sprintf("app=%s", deploymentState.Name)
	podList, err := podClient.List(ctx, metav1.ListOptions{LabelSelector: podLabelId})
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "error listing pods to determine application health: %s", err)
	}
	if podList == nil {
		return status.Errorf(codes.FailedPrecondition, "kubernetes pod list cannot be nil")
	}

	for _, pod := range podList.Items {
		// Redact env vars because they can contain secrets
		for i := 0; i < len(pod.Spec.Containers); i++ {
			pod.Spec.Containers[i].Env = []corev1.EnvVar{}
		}

		podJson, err := json.Marshal(map[string]interface{}{
			"pod":       pod,
			"hostIP":    pod.Status.HostIP,
			"ipAddress": pod.Status.PodIP,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal k8s pod definition to json: %s", podJson)
		}

		var health sdk.StatusReport_Health
		var healthMessage string

		switch pod.Status.Phase {
		case corev1.PodPending:
			health = sdk.StatusReport_ALIVE
		case corev1.PodRunning:
			health = sdk.StatusReport_ALIVE
			// Extra checks on the latest condition to ensure pod is reporting ready and running
			for _, c := range pod.Status.Conditions {
				if c.Status == corev1.ConditionTrue && c.Type == corev1.PodReady {
					health = sdk.StatusReport_READY
					healthMessage = fmt.Sprintf("ready")
					break
				}
			}
		case corev1.PodSucceeded:
			// kind of a weird one - in our current model pods are always assumed to be long-lived. If a pod exits at all, it's Down.
			health = sdk.StatusReport_DOWN
		case corev1.PodFailed:
			health = sdk.StatusReport_DOWN
		case corev1.PodUnknown:
			health = sdk.StatusReport_UNKNOWN
		default:
			health = sdk.StatusReport_UNKNOWN
		}

		// If we don't have anything better, the pod status phase is an OK health message
		// NOTE(izaak): An alternative health message could be the "type" of all conditions tied for the most recent
		// latestTransitionTime concatenated together.
		if healthMessage == "" {
			healthMessage = fmt.Sprintf("%s", pod.Status.Phase)
		}

		sr.Resources = append(sr.Resources, &sdk.StatusReport_Resource{
			Name:                pod.ObjectMeta.Name,
			Id:                  fmt.Sprintf("%s", pod.UID),
			Type:                "pod",
			ParentResourceId:    deploymentResource.Id,
			Health:              health,
			HealthMessage:       healthMessage,
			CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE,
			StateJson:           string(podJson),
			CreatedTime:         timestamppb.New(pod.CreationTimestamp.Time),
		})
	}
	s.Update("Finished building report for Kubernetes deployment resource")
	s.Done()
	return nil
}

func configureContainer(
	c *Container,
	image string,
	envVars map[string]string,
	scratchSpace []string,
	volumes []corev1.Volume,
	autoscaleConfig *AutoscaleConfig,
	log hclog.Logger,
	ui terminal.UI,
) (*corev1.Container, error) {
	// If the user is using the latest tag, then don't specify an overriding pull policy.
	// This by default means kubernetes will always pull so that latest is useful.

	var pullPolicy corev1.PullPolicy
	imageReference, err := reference.Parse(image)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "image %q is not a valid OCI reference: %q", image, err)
	}
	taggedImageReference, ok := imageReference.(reference.Tagged)
	if !ok || taggedImageReference.Tag() == "latest" {
		// If no tag is set, docker will default to "latest", and we always want to pull.
		pullPolicy = corev1.PullAlways
	} else {
		// A tag is present, we can use k8s worker caching
		pullPolicy = corev1.PullIfNotPresent
	}

	var k8sPorts []corev1.ContainerPort
	for _, port := range c.Ports {
		// Default the port protocol to TCP
		if port.Protocol == "" {
			port.Protocol = "TCP"
		}

		k8sPorts = append(k8sPorts, corev1.ContainerPort{
			Name:          port.Name,
			ContainerPort: int32(port.Port),
			HostPort:      int32(port.HostPort),
			HostIP:        port.HostIP,
			Protocol:      corev1.Protocol(strings.TrimSpace(strings.ToUpper(port.Protocol))),
		})
	}

	if envVars == nil {
		envVars = make(map[string]string)
	}

	// assume the first port defined is the 'main' port to use
	var defaultPort int
	if len(k8sPorts) != 0 {
		defaultPort = int(k8sPorts[0].ContainerPort)
		envVars["PORT"] = fmt.Sprintf("%d", defaultPort)
	}
	var k8sEnvVars []corev1.EnvVar
	for k, v := range envVars {
		k8sEnvVars = append(k8sEnvVars, corev1.EnvVar{Name: k, Value: v})
	}

	initialDelaySeconds := int32(5)
	timeoutSeconds := int32(5)
	failureThreshold := int32(30)

	if c.Probe != nil {
		if c.Probe.InitialDelaySeconds != 0 {
			initialDelaySeconds = int32(c.Probe.InitialDelaySeconds)
		}
		if c.Probe.TimeoutSeconds != 0 {
			timeoutSeconds = int32(c.Probe.TimeoutSeconds)
		}
		if c.Probe.FailureThreshold != 0 {
			failureThreshold = int32(c.Probe.FailureThreshold)
		}
	}

	// Get container resource limits and requests
	resourceLimits := make(map[corev1.ResourceName]k8sresource.Quantity)
	resourceRequests := make(map[corev1.ResourceName]k8sresource.Quantity)

	if c.CPU != nil {
		if c.CPU.Request != "" {
			q, err := k8sresource.ParseQuantity(c.CPU.Request)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse cpu request %s to k8s quantity: %s", c.CPU.Request, err)
			}
			resourceRequests[corev1.ResourceCPU] = q
		}

		if c.CPU.Limit != "" {
			q, err := k8sresource.ParseQuantity(c.CPU.Limit)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse cpu limit %s to k8s quantity: %s", c.CPU.Limit, err)
			}
			resourceLimits[corev1.ResourceCPU] = q
		}
	}

	if c.Memory != nil {
		if c.Memory.Request != "" {
			q, err := k8sresource.ParseQuantity(c.Memory.Request)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse memory requested %s to k8s quantity: %s", c.Memory.Request, err)
			}
			resourceRequests[corev1.ResourceMemory] = q
		}

		if c.Memory.Limit != "" {
			q, err := k8sresource.ParseQuantity(c.Memory.Limit)
			if err != nil {
				return nil,
					status.Errorf(codes.InvalidArgument, "failed to parse memory limit %s to k8s quantity: %s", c.Memory.Limit, err)
			}
			resourceLimits[corev1.ResourceMemory] = q
		}
	}

	for k, v := range c.Resources {
		if strings.HasPrefix(k, "limits_") {
			limitKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(limitKey[1])

			q, err := k8sresource.ParseQuantity(v)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "failed to parse resource %s to k8s quantity: %s", v, err)
			}
			resourceLimits[resourceName] = q
		} else if strings.HasPrefix(k, "requests_") {
			reqKey := strings.Split(k, "_")
			resourceName := corev1.ResourceName(reqKey[1])

			q, err := k8sresource.ParseQuantity(v)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "failed to parse resource %s to k8s quantity: %s", v, err)
			}
			resourceRequests[resourceName] = q
		} else {
			log.Warn("ignoring unrecognized k8s resources key", "key", k)
		}
	}

	_, cpuLimit := resourceLimits[corev1.ResourceCPU]
	_, cpuRequest := resourceRequests[corev1.ResourceCPU]

	// Check autoscaling
	if autoscaleConfig != nil && !(cpuLimit || cpuRequest) {
		ui.Output("For autoscaling in Kubernetes to work, a deployment must specify "+
			"cpu resource limits and requests. Otherwise the metrics-server will not properly be able "+
			"to scale your deployment.", terminal.WithWarningStyle())
	}

	resourceRequirements := corev1.ResourceRequirements{
		Limits:   resourceLimits,
		Requests: resourceRequests,
	}

	var volumeMounts []corev1.VolumeMount
	for idx, scratchSpaceLocation := range scratchSpace {
		volumeMounts = append(
			volumeMounts,
			corev1.VolumeMount{
				// We know all the volumes are identical
				Name:      volumes[idx].Name,
				MountPath: scratchSpaceLocation,
			},
		)
	}

	container := corev1.Container{
		Name:            c.Name,
		Image:           image,
		ImagePullPolicy: pullPolicy,
		Env:             k8sEnvVars,
		Resources:       resourceRequirements,
		VolumeMounts:    volumeMounts,
	}

	if len(k8sPorts) > 0 {
		container.Ports = k8sPorts
	}
	if c.Command != nil {
		container.Command = *c.Command
	}
	if c.Args != nil {
		container.Args = *c.Args
	}

	// Only define liveliness & readiness checks if container binds to a port
	if defaultPort > 0 {
		var handler corev1.ProbeHandler
		if c.ProbePath != "" {
			handler = corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: c.ProbePath,
					Port: intstr.FromInt(defaultPort),
				},
			}
		} else {
			// If no probe path is defined, assume app will bind to default TCP port
			// TODO: handle apps that aren't socket listeners
			handler = corev1.ProbeHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt(defaultPort),
				},
			}
		}

		container.LivenessProbe = &corev1.Probe{
			ProbeHandler:             handler,
			InitialDelaySeconds: initialDelaySeconds,
			TimeoutSeconds:      timeoutSeconds,
			FailureThreshold:    failureThreshold,
		}
		container.ReadinessProbe = &corev1.Probe{
			ProbeHandler:             handler,
			InitialDelaySeconds: initialDelaySeconds,
			TimeoutSeconds:      timeoutSeconds,
			FailureThreshold:    failureThreshold,
		}
	}

	return &container, nil
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

	var overlayTarget *Container
	if p.config.Pod != nil && p.config.Pod.Container != nil {
		overlayTarget = p.config.Pod.Container
	} else {
		overlayTarget = &Container{}
	}

	appContainerSpec, err := overlayTopLevelProperties(p.config, overlayTarget)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Failed to parse container config: %s", err)
	}

	// App container must have some kind of port
	if len(appContainerSpec.Ports) == 0 {
		log.Warn("No ports defined in waypoint.hcl - defaulting to http on port", "port", DefaultServicePort)
		appContainerSpec.Ports = append(appContainerSpec.Ports, &Port{Port: DefaultServicePort, Name: "http"})
	}

	portStep := sg.Add("")
	defer func() { portStep.Abort() }()
	// we dont use %q to save us convering a uint Port to a string and handling the error
	portStep.Update("Expected %q port \"%d\" for app %q",
		appContainerSpec.Ports[0].Name,
		appContainerSpec.Ports[0].Port,
		result.Name)
	portStep.Done()

	envVars := make(map[string]string)
	// Add deploy config environment to container env vars
	for k, v := range p.config.StaticEnvVars {
		envVars[k] = v
	}
	if p.config.Pod != nil && p.config.Pod.Container != nil {
		for k, v := range p.config.Pod.Container.StaticEnvVars {
			envVars[k] = v
		}
	}
	for k, v := range deployConfig.Env() {
		envVars[k] = v
	}

	// Create scratch space volumes
	var volumes []corev1.Volume
	for idx := range p.config.ScratchSpace {
		scratchName := fmt.Sprintf("scratch-%d", idx)
		volumes = append(volumes,
			corev1.Volume{
				Name: scratchName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		)
	}

	appImage := fmt.Sprintf("%s:%s", img.Image, img.Tag)
	appContainerSpec.Name = src.App

	appContainer, err := configureContainer(
		appContainerSpec,
		appImage,
		envVars,
		p.config.ScratchSpace,
		volumes,
		p.config.AutoscaleConfig,
		log,
		ui,
	)
	if err != nil {
		return status.Errorf(status.Code(err),
			"Failed to define app container: %s", err)
	}

	var sidecarContainers []corev1.Container
	if p.config.Pod != nil {
		for _, sidecarConfig := range p.config.Pod.Sidecars {
			envVars := make(map[string]string)
			// Add deploy config environment to container env vars
			for k, v := range sidecarConfig.Container.StaticEnvVars {
				envVars[k] = v
			}
			for k, v := range deployConfig.Env() {
				envVars[k] = v
			}

			sidecarContainer, err := configureContainer(
				sidecarConfig.Container,
				sidecarConfig.Image,
				envVars,
				p.config.ScratchSpace,
				volumes,
				p.config.AutoscaleConfig,
				log,
				ui,
			)
			if err != nil {
				return status.Errorf(status.Code(err),
					"Failed to define sidecar container %s: %s", sidecarConfig.Container.Name, err)
			}
			sidecarContainers = append(sidecarContainers, *sidecarContainer)
		}
	}

	// Update the deployment with our spec
	containers := []corev1.Container{*appContainer}
	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: append(containers, sidecarContainers...),
		Volumes:    volumes,
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
	// expect pods to be labeled with 'version'
	deployment.Spec.Template.Labels["version"] = result.Id

	// Apply user defined labels
	for k, v := range p.config.Labels {
		deployment.Spec.Template.Labels[k] = v
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
		deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{
			Name: p.config.ImageSecret,
		}}
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
		return status.Errorf(codes.Internal, "failed to create or update deployment: %s", err)
	}

	ev := clientSet.CoreV1().Events(ns)

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

	var timeoutSeconds int
	var failureThreshold int
	var initialDelaySeconds int

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if int(container.ReadinessProbe.TimeoutSeconds) > timeoutSeconds {
			timeoutSeconds = int(container.ReadinessProbe.TimeoutSeconds)
		}

		if int(container.ReadinessProbe.FailureThreshold) > failureThreshold {
			failureThreshold = int(container.ReadinessProbe.FailureThreshold)
		}

		if int(container.ReadinessProbe.TimeoutSeconds) > initialDelaySeconds {
			initialDelaySeconds = int(container.ReadinessProbe.InitialDelaySeconds)
		}
	}

	// We wait the maximum amount of time that the deployment controller would wait for a pod
	// to start before exiting. We double the time to allow for various Kubernetes based
	// delays in startup, detection, and reporting.
	timeout := time.Duration((timeoutSeconds*failureThreshold)+initialDelaySeconds) * 2 * time.Second

	podsSeen := make(map[types.UID]string)

	// Wait on the Pod to start
	err = wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		dep, err := dc.Get(ctx, result.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if time.Since(lastStatus) > 10*time.Second {
			step.Update(fmt.Sprintf(
				"Waiting on deployment to become available: requested=%d running=%d ready=%d",
				*dep.Spec.Replicas,
				dep.Status.UnavailableReplicas+dep.Status.AvailableReplicas,
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
			podsSeen[p.UID] = p.Name

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
		step.Update("Error detected waiting for Deployment to start.")
		step.Status(terminal.StatusError)
		step.Abort()

		ui.Output("The following is events for pods observed while attempting to start the Deployment", terminal.WithWarningStyle())

		for uid, name := range podsSeen {
			sel := ev.GetFieldSelector(nil, nil, nil, (*string)(&uid))

			events, err := ev.List(ctx, metav1.ListOptions{
				FieldSelector: sel.String(),
			})
			if err == nil {
				ui.Output("Events for %s", name, terminal.WithHeaderStyle())
				for _, ev := range events.Items {
					if ev.Type == "Normal" {
						continue
					}
					ui.Output("  %s: %s (%s)", ev.Type, ev.Message, ev.Reason)
				}
			}
		}

		if err == wait.ErrWaitTimeout {
			err = fmt.Errorf("Deployment was not able to start pods after %s", timeout)
		}
		return err
	}

	step.Update("Deployment successfully rolled out!")
	step.Done()

	return nil
}

var deleteGrace = int64(120)

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

	del := metav1.DeletePropagationBackground

	msg := "Deployment deleted"
	deployclient := csinfo.Clientset.AppsV1().Deployments(ns)
	if err := deployclient.Delete(ctx, state.Name, metav1.DeleteOptions{
		GracePeriodSeconds: &deleteGrace,
		PropagationPolicy:  &del,
	}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		msg = fmt.Sprintf("Deployment (%s) not found, continuing..", state.Name)
	}
	step.Update(msg)
	step.Done()

	return nil
}

// resourceAutoscalerCreate creates the Kubernetes deployment.
func (p *Platform) resourceAutoscalerCreate(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,

	result *Deployment,
	deployState *Resource_Deployment, // a deployment must exist before creating an autoscaler
	state *Resource_Autoscale,
	csinfo *clientsetInfo,
	sg terminal.StepGroup,
) error {
	if p.config.AutoscaleConfig == nil {
		log.Trace("no autoscale config detected, will not create one")
		return nil
	}

	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	s := sg.Add("Preparing horizontal pod autoscaler...")
	defer func() { s.Abort() }() // Defer in func in case more steps are added to this func in the future

	clientSet := csinfo.Clientset
	autoscaleClient := clientSet.AutoscalingV1().HorizontalPodAutoscalers(ns)

	// Attempt to detect an existing metrics-server and display warning if none found
	metricsServerLabel := "k8s-app=metrics-server"
	metricsPods, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: metricsServerLabel,
	})
	if err != nil {
		// we don't return the error, this was mostly to provide a helpful warning
		log.Info("receieved error while listing pods in attempt to detect existing metrics-server: %s", err)
		err = nil

		// The apis return an non-nil but empty value when observing an error, so we need to be sure to
		// skip the code below.
		metricsPods = nil
	}

	if metricsPods != nil && len(metricsPods.Items) == 0 {
		ui.Output("There were no pods recognized in the cluster as a metrics-server "+
			"with the label '%s'. This means your autoscaler might "+
			"not be functional until a metrics-server exists and can provide metrics to "+
			"the horizontal pod autoscaler.",
			metricsServerLabel,
			terminal.WithWarningStyle())
		ui.Output("If you have not yet setup a metrics-server inside your Kubernetes cluster, "+
			"please refer to the metrics-server project documentation for properly "+
			"installing one: https://github.com/kubernetes-sigs/metrics-server", terminal.WithWarningStyle())
		ui.Output("Waypoint will continue to configure horizontal pod autoscaler ...",
			terminal.WithWarningStyle())
	}

	state.Name = deployState.Name

	var autoscaler *autoscalingv1.HorizontalPodAutoscaler
	create := false
	autoscaler, err = autoscaleClient.Get(ctx, state.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Debug("no horizontal pod autoscaler found, will create one")

		err = nil
		create = true

		autoscaler = &autoscalingv1.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name: state.Name,
			},
			Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       state.Name,
				},
			},
		}
	}
	if err != nil {
		return err
	}

	if p.config.AutoscaleConfig.MaxReplicas == 0 {
		return status.Error(codes.FailedPrecondition,
			"No max replica config value set for autoscale")
	}

	// Cannot be smaller than min replicas, but k8s will return a good error
	// message if it is when we go to create an HPA
	autoscaler.Spec.MaxReplicas = p.config.AutoscaleConfig.MaxReplicas

	if p.config.AutoscaleConfig.MinReplicas > 0 {
		autoscaler.Spec.MinReplicas = &p.config.AutoscaleConfig.MinReplicas
	}

	if p.config.AutoscaleConfig.TargetCPU >= 0 {
		autoscaler.Spec.TargetCPUUtilizationPercentage = &p.config.AutoscaleConfig.TargetCPU
	}

	// create it
	action := "created"
	if create {
		_, err = autoscaleClient.Create(ctx, autoscaler, metav1.CreateOptions{})
	} else {
		_, err = autoscaleClient.Update(ctx, autoscaler, metav1.UpdateOptions{})
		action = "updated"
	}
	if err != nil {
		return err
	}

	s.Update("Horizontal Pod Autoscaler has been %s", action)
	s.Done()

	return nil
}

// resourceAutoscalerDestroy deletes the K8S horizontal pod autoscaler
func (p *Platform) resourceAutoscalerDestroy(
	ctx context.Context,
	state *Resource_Autoscale,
	sg terminal.StepGroup,
	csinfo *clientsetInfo,
) error {
	if state.Name == "" {
		// No autoscale config, so don't destroy one
		return nil
	}

	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	s := sg.Add("Deleting horizontal pod autoscaler...")
	defer func() { s.Abort() }()

	clientSet := csinfo.Clientset
	autoscaleClient := clientSet.AutoscalingV1().HorizontalPodAutoscalers(ns)

	if err := autoscaleClient.Delete(ctx, state.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	s.Update("Horizontal Pod Autoscaler deleted")
	s.Done()
	return nil
}

// resourceAutoscalerStatus gets information about the autoscaler for a status report
func (p *Platform) resourceAutoscalerStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Autoscale,
	clientsetInfo *clientsetInfo,
	sr *resource.StatusResponse,
) error {
	if p.config.AutoscaleConfig == nil && state.Name == "" {
		// No autoscale config, so don't take the status of one
		return nil
	}

	// Prepare our namespace and override if set.
	ns := clientsetInfo.Namespace
	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	s := sg.Add("Checking status of Kubernetes horizontal pod autoscaler %q...", state.Name)
	defer s.Abort()

	clientSet := clientsetInfo.Clientset
	autoscaleClient := clientSet.AutoscalingV1().HorizontalPodAutoscalers(ns)

	hpaResource := sdk.StatusReport_Resource{
		Type:                "horizontal pod autoscaler",
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
	}
	sr.Resources = append(sr.Resources, &hpaResource)

	autoscalerResp, err := autoscaleClient.Get(ctx, state.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return status.Errorf(codes.FailedPrecondition,
				"error getting kubernetes horizontal pod autoscaler %s: %s", state.Name, err)
		} else {
			hpaResource.Name = state.Name
			hpaResource.Health = sdk.StatusReport_MISSING
			hpaResource.HealthMessage = sdk.StatusReport_MISSING.String()
		}
	} else if autoscalerResp == nil {
		return status.Errorf(codes.FailedPrecondition,
			"kubernetes horizontal pod autoscaler response cannot be empty")
	} else {

		hpaResource.Name = state.Name
		hpaResource.Id = fmt.Sprintf("%s", autoscalerResp.ObjectMeta.UID)
		hpaResource.CreatedTime = timestamppb.New(autoscalerResp.ObjectMeta.CreationTimestamp.Time)
		// the existence of the resource means it's ready. It has no other status
		hpaResource.Health = sdk.StatusReport_READY
		hpaResource.HealthMessage = "The HPA resource is ready"

		hpaStateJson, err := json.Marshal(map[string]interface{}{
			"horizontalPodAutoscaler": &hpaResource,
		})
		if err != nil {
			return status.Errorf(codes.FailedPrecondition,
				"failed to marshal horizontal pod autoscaler to json: %s", err)
		}

		hpaResource.StateJson = string(hpaStateJson)
	}

	s.Update("Finished building report for Kubernetes horizontal pod autoscaler")
	s.Done()
	return nil
}

// Deploy deploys an image to Kubernetes.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	dcr *component.DeclaredResourcesResp,
	ui terminal.UI,
) (*Deployment, error) {
	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	seq := deployConfig.Sequence
	result.Id = id
	result.Name = strings.ToLower(fmt.Sprintf("%s-v%d", src.App, seq))

	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()

	// Create our resource manager and create
	rm := p.resourceManager(log, dcr)
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

	rm := p.resourceManager(log, nil)

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

	step := sg.Add("Gathering health report for Kubernetes deployment...")
	defer func() { step.Abort() }() // Defer in func in case more steps are added to this func in the future

	rm := p.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource("deployment").SetState(&Resource_Deployment{
			Name: deployment.Name,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return nil, err
		}
	}

	result, err := rm.StatusReport(ctx, log, sg, ui)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resource manager failed to generate resource statuses: %s", err)
	}

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()
	st.Update("Determining overall container health...")

	log.Debug("status report complete")

	// update output based on main health state
	step.Update("Finished building report for Kubernetes platform")
	step.Done()

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
			st.Step(terminal.StatusWarn, fmt.Sprintf("%s %q is reporting %q", resource.Type, resource.Name, resource.Health.String()))
		}
	}

	return result, nil
}

// overlayDefaultProperties overlays the top level container properties from config onto the
// more detailed container properties in container.
// ConfigSet has already validated that both are not set, so we don't have to check here.
func overlayTopLevelProperties(config Config, container *Container) (*Container, error) {
	var overlaidContainer *Container
	i, err := copystructure.Copy(container)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to initialize container spec: %s", err)
	}
	overlaidContainer, _ = i.(*Container)

	if config.ProbePath != "" {
		overlaidContainer.ProbePath = config.ProbePath
	}
	if config.Probe != nil {
		overlaidContainer.ProbePath = config.ProbePath
	}
	if config.Resources != nil {
		overlaidContainer.Resources = config.Resources
	}
	if config.CPU != nil {
		overlaidContainer.CPU = config.CPU
	}
	if config.Memory != nil {
		overlaidContainer.Memory = config.Memory
	}
	if config.StaticEnvVars != nil {
		overlaidContainer.StaticEnvVars = config.StaticEnvVars
	}
	if config.ServicePort != nil {
		// We've already validated that ports is nil in ConfigSet - they cannot both be set at once.
		overlaidContainer.Ports = []*Port{{Port: *config.ServicePort, Name: "http"}}
	}

	return overlaidContainer, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// Annotations are added to the pod spec of the deployed application.  This is
	// useful when using mutating webhook admission controllers to further process
	// pod events.
	Annotations map[string]string `hcl:"annotations,optional"`

	// AutoscaleConfig will create a horizontal pod autoscaler for a given
	// deployment and scale the replica pods up or down based on a given
	// load metric, such as CPU utilization
	AutoscaleConfig *AutoscaleConfig `hcl:"autoscale,block"`

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

	// If set, this is the HTTP path to request to test that the application
	// is up and running. Without this, we only test that a connection can be
	// made to the port.
	ProbePath string `hcl:"probe_path,optional"`

	// Probe details for describing a health check to be performed against a container.
	Probe *Probe `hcl:"probe,block"`

	// Optionally define various resources limits for kubernetes pod containers
	// such as memory and cpu.
	Resources map[string]string `hcl:"resources,optional"`

	// Optionally define various cpu resource limits and requests for kubernetes pod containers
	CPU *ResourceConfig `hcl:"cpu,block"`

	// Optionally define various memory resource limits and requests for kubernetes pod containers
	Memory *ResourceConfig `hcl:"memory,block"`

	// An array of paths to directories that will be mounted as EmptyDirVolumes in the pod
	// to store temporary data.
	ScratchSpace []string `hcl:"scratch_path,optional"`

	// ServiceAccount is the name of the Kubernetes service account to apply to the
	// application deployment. This is useful to apply Kubernetes RBAC to the pod.
	ServiceAccount string `hcl:"service_account,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to DefaultServicePort const.
	// NOTE: Ports and ServicePort cannot both be defined
	ServicePort *uint `hcl:"service_port,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Pod describes the configuration for the pod
	Pod *Pod `hcl:"pod,block"`

	// Deprecated field, previous definition of ports
	DeprecatedPorts []map[string]string `hcl:"ports,optional" docs:"hidden"`
}

// ResourceConfig describes the request and limit of a resource. Used for
// cpu and memory resource configuration.
type ResourceConfig struct {
	Request string `hcl:"request,optional" json:"request"`
	Limit   string `hcl:"limit,optional" json:"limit"`
}

// AutoscaleConfig describes the possible configuration for creating a
// horizontal pod autoscaler
type AutoscaleConfig struct {
	MinReplicas int32 `hcl:"min_replicas,optional"`
	MaxReplicas int32 `hcl:"max_replicas,optional"`
	// TargetCPU will determine the max load before the autoscaler will increase
	// a replica
	TargetCPU int32 `hcl:"cpu_percent,optional"`
}

// Pod describes the configuration for the pod
type Pod struct {
	SecurityContext *PodSecurityContext `hcl:"security_context,block"`
	Container       *Container          `hcl:"container,block"`
	Sidecars        []*Sidecar          `hcl:"sidecar,block"`
}

type Sidecar struct {
	// Specifying Image in Container would make it visible on the main Pod config,
	// which isn't the right way to specify the app image.
	Image string `hcl:"image"`

	Container *Container `hcl:"container,block"`
}

type Port struct {
	Name     string `hcl:"name"`
	Port     uint   `hcl:"port"`
	HostPort uint   `hcl:"host_port,optional"`
	HostIP   string `hcl:"host_ip,optional"`
	Protocol string `hcl:"protocol,optional"`
}

// Container describes the detailed parameters to declare a kubernetes container
type Container struct {
	Name          string            `hcl:"name,optional"`
	Ports         []*Port           `hcl:"port,block"`
	ProbePath     string            `hcl:"probe_path,optional"`
	Probe         *Probe            `hcl:"probe,block"`
	CPU           *ResourceConfig   `hcl:"cpu,block"`
	Memory        *ResourceConfig   `hcl:"memory,block"`
	Resources     map[string]string `hcl:"resources,optional"`
	Command       *[]string         `hcl:"command,optional"`
	Args          *[]string         `hcl:"args,optional"`
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

// PodSecurityContext describes the security config for the Pod
type PodSecurityContext struct {
	RunAsUser    *int64 `hcl:"run_as_user"`
	RunAsGroup   *int64 `hcl:"run_as_group"`
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
	// case startup times. Defaults to 30 failures.
	FailureThreshold uint `hcl:"failure_threshold,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy the application into a Kubernetes cluster using Deployment objects")
	doc.Input("docker.Image")
	doc.Output("k8s.Deployment")

	doc.Example(`
use "kubernetes" {
	image_secret = "registry_secret"
	replicas = 3
	probe_path = "/_healthz"
}
`)

	commonSubFields := map[string]func(doc *docs.SubFieldDoc){
		"port": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"port",
				"a port and options that the application is listening on",
				docs.Summary(
					"used to define and expose multiple ports that the application or process is",
					"listening on for the container in use. Can be specified multiple times for many ports.",
				),
				doc.SubFields("port", func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"name",
						"name of the port",
						docs.Summary("If specified, this must be an IANA_SVC_NAME and unique within the pod. Each",
							"named port in a pod must have a unique name. Name for the port that can be",
							"referred to by services.",
						),
					)
					doc.SetField(
						"port",
						"the port number",
						docs.Summary("Number of port to expose on the pod's IP address.",
							"This must be a valid port number, 0 < x < 65536.",
						),
					)
					doc.SetField(
						"host_port",
						"the corresponding worker node port",
						docs.Summary("Number of port to expose on the host.",
							"If specified, this must be a valid port number, 0 < x < 65536.",
							"If HostNetwork is specified, this must match ContainerPort.",
							"Most containers do not need this.",
						),
					)
					doc.SetField(
						"host_ip",
						"what host IP to bind the external port to",
					)
					doc.SetField(
						"protocol",
						"protocol for port. Must be UDP, TCP, or SCTP",
						docs.Default("TCP"),
					)
				}),
			)
		},
		"static_environment": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"static_environment",
				"environment variables to control broad modes of the application",
				docs.Summary(
					"environment variables that are meant to configure the container in a static",
					"way. This might be control an image that has multiple modes of operation,",
					"selected via environment variable. Most configuration should use the waypoint",
					"config commands",
				),
			)
		},
		"cpu": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"cpu",
				"cpu resource configuration",
				docs.Summary("CPU lets you define resource limits and requests for a container in "+
					"a deployment."),
				doc.SubFields("cpu", func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"request",
						"how much cpu to give the container in cpu cores. Supports m to indicate milli-cores",
					)
					doc.SetField(
						"limit",
						"maximum amount of cpu to give the container. Supports m to indicate milli-cores",
					)
				}),
			)
		},
		"memory": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"memory",
				"memory resource configuration",
				docs.Summary("Memory lets you define resource limits and requests for a container in "+
					"a deployment."),
				doc.SubFields("memory", func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"request",
						"how much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes",
					)

					doc.SetField(
						"limit",
						"maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes",
					)
				}),
			)
		},
		"resources": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"resources",
				"a map of resource limits and requests to apply to a container on deploy",
				docs.Summary(
					"resource limits and requests for a container. This exists to allow any possible "+
						"resources. For cpu and memory, use those relevant settings instead. "+
						"Keys must start with either `limits_` or `requests_`. Any other options "+
						"will be ignored.",
				),
			)
		},
		"probe_path": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"probe_path",
				"the HTTP path to request to test that the application is running",
				docs.Summary(
					"without this, the test will simply be that the application has bound to the port",
				),
			)
		},
		"probe": func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"probe",
				"configuration to control liveness and readiness probes",
				docs.Summary("Probe describes a health check to be performed against a ",
					"container to determine whether it is alive or ready to receive traffic."),
				doc.SubFields("probe", func(doc *docs.SubFieldDoc) {
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
		},
	}
	commonSubFields["container"] = func(doc *docs.SubFieldDoc) {
		doc.SetField(
			"container",
			"container describes the commands and arguments for a container config",
			doc.SubFields("container", func(doc *docs.SubFieldDoc) {
				doc.SetField(
					"name",
					"name of the container",
				)

				commonSubFields["cpu"](doc)
				commonSubFields["memory"](doc)
				commonSubFields["resources"](doc)
				commonSubFields["probe_path"](doc)
				commonSubFields["probe"](doc)
				commonSubFields["port"](doc)
				commonSubFields["static_environment"](doc)

				doc.SetField(
					"command",
					"an array of strings to run for the container",
				)

				doc.SetField(
					"args",
					"an array of string arguments to pass through to the container",
				)
			}),
		)
	}

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
		"autoscale",
		"sets up a horizontal pod autoscaler to scale deployments automatically",
		docs.Summary("This configuration will automatically set up and associate the "+
			"current deployment with a horizontal pod autoscaler in Kuberentes. Note that "+
			"for this to work, you must also define resource limits and requests for a deployment "+
			"otherwise the metrics-server will not be able to properly determine a deployments "+
			"target CPU utilization"),
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"min_replicas",
				"the minimum amount of pods to have for a deployment",
			)

			doc.SetField(
				"max_replicas",
				"the maximum amount of pods to scale to for a deployment",
			)

			doc.SetField(
				"cpu_percent",
				"the target CPU percent utilization before the horizontal pod autoscaler "+
					"scales up a deployments replicas",
			)
		}),
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
		"image_secret",
		"name of the Kubernetes secrete to use for the image",
		docs.Summary(
			"this references an existing secret, waypoint does not create this secret",
		),
	)

	doc.SetField(
		"kubeconfig",
		"path to the kubeconfig file to use",
		docs.Summary("by default uses from current user's home directory"),
		docs.EnvVar("KUBECONFIG"),
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
		doc.SubFields("probe", func(doc *docs.SubFieldDoc) {
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
				docs.Default("30"),
			)
		}),
	)

	doc.SetField(
		"resources",
		"a map of resource limits and requests to apply to a container on deploy",
		docs.Summary(
			"resource limits and requests for a container. This exists to allow any possible "+
				"resources. For cpu and memory, use those relevant settings instead. "+
				"Keys must start with either `limits_` or `requests_`. Any other options "+
				"will be ignored.",
		),
	)

	doc.SetField(
		"cpu",
		"cpu resource configuration",
		docs.Summary("CPU lets you define resource limits and requests for a container in "+
			"a deployment."),
		doc.SubFields("cpu", func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"request",
				"how much cpu to give the container in cpu cores. Supports m to indicate milli-cores",
			)
			doc.SetField(
				"limit",
				"maximum amount of cpu to give the container. Supports m to indicate milli-cores",
			)
		}),
	)

	doc.SetField(
		"memory",
		"memory resource configuration",
		docs.Summary("Memory lets you define resource limits and requests for a container in "+
			"a deployment."),
		doc.SubFields("memory", func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"request",
				"how much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes",
			)

			doc.SetField(
				"limit",
				"maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes",
			)
		}),
	)

	doc.SetField(
		"scratch_path",
		"a path for the service to store temporary data",
		docs.Summary(
			"a path to a directory that will be created for the service to store temporary data using EmptyDir.",
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
		"service_port",
		"the TCP port that the application is listening on",
		docs.Default(fmt.Sprint(DefaultServicePort)),
		docs.Summary(
			"by default, this config variable is used for exposing a single port for",
			"the container in use. For multi-port configuration, use 'ports' instead.",
		),
	)

	doc.SetField(
		"static_environment",
		"environment variables to control broad modes of the application",
		docs.Summary(
			"environment variables that are meant to configure the container in a static",
			"way. This might be control an image that has multiple modes of operation,",
			"selected via environment variable. Most configuration should use the waypoint",
			"config commands",
		),
	)

	doc.SetField(
		"pod",
		"the configuration for a pod",
		docs.Summary("Pod describes the configuration for a pod when deploying"),
		doc.SubFields("pod", func(doc *docs.SubFieldDoc) {
			commonSubFields["container"](doc)
			doc.SetField(
				"sidecar",
				"a sidecar container within the same pod",
				docs.Summary("Another container to run alongside the app container in the kubernetes pod.",
					"Can be specified multiple times for multiple sidecars.",
				),
				doc.SubFields("sidecar", func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"image",
						"image of the sidecar container",
					)
					commonSubFields["container"](doc)
				}),
			)
			doc.SetField(
				"security_context",
				"holds pod-level security attributes and container settings",
				docs.SubFields(func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"run_as_user",
						"the UID to run the entrypoint of the container process",
					)
					doc.SetField(
						"run_as_non_root",
						"indicates that the container must run as a non-root user",
					)
					doc.SetField(
						"fs_group",
						"a special supplemental group that applies to all containers in a pod",
					)
				}),
			)
		}),
	)

	return doc, nil
}

var mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current deployment is not ready, however your application
might be available or still starting up.
`)

var (
	_ component.Platform         = (*Platform)(nil)
	_ component.PlatformReleaser = (*Platform)(nil)
	_ component.Configurable     = (*Platform)(nil)
	_ component.Documented       = (*Platform)(nil)
	_ component.Destroyer        = (*Platform)(nil)
	_ component.Status           = (*Platform)(nil)
)
