// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// DefaultPort is the port that a service will forward to the pod(s)
const DefaultPort = 80

// Releaser is the ReleaseManager implementation for Kubernetes.
type Releaser struct {
	config ReleaserConfig
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// DestroyFunc implements component.Destroyer
func (r *Releaser) DestroyFunc() interface{} {
	return r.Destroy
}

// StatusFunc implements component.Status
func (r *Releaser) StatusFunc() interface{} {
	return r.Status
}

func (r *Releaser) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(r.getClientset),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithResource(resource.NewResource(
			resource.WithName("service"),
			resource.WithState(&Resource_Service{}),
			resource.WithCreate(r.resourceServiceCreate),
			resource.WithDestroy(r.resourceServiceDestroy),
			resource.WithStatus(r.resourceServiceStatus),
			resource.WithPlatform(platformName),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("ingress"),
			resource.WithState(&Resource_Ingress{}),
			resource.WithCreate(r.resourceIngressCreate),
			resource.WithDestroy(r.resourceIngressDestroy),
			resource.WithStatus(r.resourceIngressStatus),
			resource.WithPlatform(platformName),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
	)
}

func (r *Releaser) resourceServiceStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Service,
	clientset *clientsetInfo,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Checking status of Kubernetes service resource %q...", state.Name)
	defer s.Abort()

	namespace := r.config.Namespace
	if namespace == "" {
		namespace = clientset.Namespace
	}

	serviceResource := sdk.StatusReport_Resource{
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_ROUTER,
	}
	sr.Resources = append(sr.Resources, &serviceResource)

	serviceResp, err := clientset.Clientset.CoreV1().Services(namespace).Get(ctx, state.Name, metav1.GetOptions{})
	if serviceResp == nil {
		return status.Errorf(codes.FailedPrecondition, "kubernetes service response cannot be empty")
	} else if err != nil {
		if !errors.IsNotFound(err) {
			return err
		} else {
			s.Update("No service resource was found")
			s.Status(terminal.StatusError)
			s.Done()
			s = sg.Add("")

			serviceResource.Name = state.Name
			serviceResource.Health = sdk.StatusReport_MISSING
			serviceResource.HealthMessage = sdk.StatusReport_MISSING.String()

			// Continue on with the rest of our resources
		}
	} else {
		// We found the service, and can use it to populate our resource
		var ipAddress string
		if serviceResp.Spec.LoadBalancerIP != "" {
			ipAddress = serviceResp.Spec.LoadBalancerIP
		} else if serviceResp.Spec.ClusterIP != "" {
			ipAddress = serviceResp.Spec.ClusterIP
		}

		serviceJson, err := json.Marshal(map[string]interface{}{
			"service":   serviceResp,
			"ipAddress": ipAddress,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal k8s service definition to json: %s", err)
		}

		serviceResource.Id = fmt.Sprintf("%s", serviceResp.UID)
		serviceResource.Name = serviceResp.Name
		serviceResource.CreatedTime = timestamppb.New(serviceResp.CreationTimestamp.Time)
		// If we found the service, then it's ready. It doesn't have any other conditions.
		serviceResource.Health = sdk.StatusReport_READY
		serviceResource.HealthMessage = fmt.Sprintf("Service %q exists and is ready", serviceResource.Name)
		serviceResource.StateJson = string(serviceJson)
	}

	s.Update("Finished building report for Kubernetes service resource")
	s.Done()
	return nil
}

func (r *Releaser) resourceServiceCreate(
	ctx context.Context,
	log hclog.Logger,
	target *Deployment,
	result *Release,
	state *Resource_Service,
	csinfo *clientsetInfo,
	sg terminal.StepGroup,
) error {
	step := sg.Add("Initializing Kubernetes client...")
	defer func() { step.Abort() }() // Defer in func in case more steps are added to this func in the future
	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	step.Update("Kubernetes client connected to %s with namespace %s", csinfo.Config.Host, ns)
	step.Done()

	step = sg.Add("Preparing service...")

	clientSet := csinfo.Clientset
	serviceclient := clientSet.CoreV1().Services(ns)

	// Determine if we have a deployment that we manage already
	create := false
	service, err := serviceclient.Get(ctx, result.ServiceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		service = result.newService(result.ServiceName)
		create = true
		err = nil
	}
	if err != nil {
		return err
	}

	// set the service name in state, at this point we've either created it or
	// it already existed
	state.Name = result.ServiceName

	// Update the spec
	service.Spec.Selector = map[string]string{
		"name":  target.Name,
		labelId: target.Id,
	}

	if (r.config.Port != 0 || r.config.NodePort != 0) && r.config.Ports != nil {
		return status.Errorf(codes.FailedPrecondition, "Cannot define both 'ports' and 'port' or 'node_port'."+
			" Use 'ports' for configuring multiple service ports")
	} else if r.config.Ports == nil && (r.config.Port != 0 || r.config.NodePort != 0) {
		r.config.Ports = make([]map[string]string, 1)
		r.config.Ports[0] = map[string]string{
			"port":        strconv.Itoa(int(r.config.Port)),
			"target_port": "http",
			"node_port":   strconv.Itoa(int(r.config.NodePort)),
		}
	} else if r.config.Port == 0 && r.config.NodePort == 0 && r.config.Ports == nil {
		// We don't explicitly set nodeport if Port isn't defined, because
		// k8s will automatically assign a nodeport if unspecified
		r.config.Ports = make([]map[string]string, 1)
		r.config.Ports[0] = map[string]string{
			"target_port": "http",
			"port":        strconv.Itoa(int(DefaultPort)),
		}
	}

	var checkLB bool

	if r.config.LoadBalancer {
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
		checkLB = true
	} else if r.config.Ports[0]["node_port"] != "" && r.config.Ports[0]["node_port"] != "0" {
		service.Spec.Type = corev1.ServiceTypeNodePort
	} else {
		service.Spec.Type = corev1.ServiceTypeClusterIP
	}

	servicePorts := make([]corev1.ServicePort, len(r.config.Ports))
	for i, sp := range r.config.Ports {
		nodePort, _ := strconv.ParseInt(sp["node_port"], 10, 32)
		port, _ := strconv.ParseInt(sp["port"], 10, 32)
		if port == 0 {
			// This likely means port was unset and got parsed to 0
			port = DefaultPort
		}

		var target_port int
		if sp["target_port"] == "" {
			sp["target_port"] = "http"
		} else {
			target_port, err = strconv.Atoi(sp["target_port"])
			if err != nil {
				// it's a string label like 'http', not an integer
				target_port = 0
			}
		}

		servicePorts[i] = corev1.ServicePort{
			Name:     sp["name"],
			Port:     int32(port),
			Protocol: corev1.ProtocolTCP,
			NodePort: int32(nodePort),
		}

		// Because of the type TargetPort is expected to be, we can't pass along
		// an int as a string, it expects the int to actually be an int
		if target_port != 0 {
			servicePorts[i].TargetPort = intstr.FromInt(target_port)
		} else {
			servicePorts[i].TargetPort = intstr.FromString(sp["target_port"])
		}
	}

	service.Spec.Ports = servicePorts

	// Apply Service annotations
	service.Annotations = r.config.Annotations

	// Create/update
	if create {
		step.Update("Creating service...")
		service, err = serviceclient.Create(ctx, service, metav1.CreateOptions{})
	} else {
		step.Update("Updating existing service...")
		service, err = serviceclient.Update(ctx, service, metav1.UpdateOptions{})
	}
	if err != nil {
		return err
	}

	step.Done()
	step = sg.Add("Waiting for service to become ready...")

	// Wait on the IP
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		service, err = serviceclient.Get(ctx, result.ServiceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if checkLB {
			return len(service.Status.LoadBalancer.Ingress) > 0, nil
		} else {
			return service.Spec.ClusterIP != "", nil
		}
	})

	if err != nil {
		return err
	}

	step.Update("Service is ready!")
	step.Done()

	if r.config.LoadBalancer {
		ingress := service.Status.LoadBalancer.Ingress[0]
		result.Url = "http://" + ingress.IP
		if ingress.Hostname != "" {
			result.Url = "http://" + ingress.Hostname
		}

		if service.Spec.Ports[0].Port != 80 {
			result.Url = fmt.Sprintf("%s:%d", result.Url, service.Spec.Ports[0].Port)
		}
	} else if service.Spec.Ports[0].NodePort > 0 {
		nodeclient := clientSet.CoreV1().Nodes()
		nodes, err := nodeclient.List(ctx, metav1.ListOptions{})
		if err != nil {
			// Rather than fail the whole release, report the error and then complete.
			// Print in a standalone step, so the output won't get overwritten if we add more step output later in the future.
			errStep := sg.Add("Cannot determine release URL for nodeport service due to failure to list nodes: %s", err)
			errStep.Status(terminal.StatusError)
			errStep.Done()
		} else {
			nodeIP := nodes.Items[0].Status.Addresses[0].Address
			result.Url = fmt.Sprintf("http://%s:%d", nodeIP, service.Spec.Ports[0].NodePort)
		}
	} else {
		result.Url = fmt.Sprintf("http://%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port)
	}

	return nil
}

func (r *Releaser) resourceServiceDestroy(
	ctx context.Context,
	log hclog.Logger,
	state *Resource_Service,
	sg terminal.StepGroup,
	csinfo *clientsetInfo,
) error {
	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	clientSet := csinfo.Clientset
	serviceclient := clientSet.CoreV1().Services(ns)
	step.Update("Kubernetes client connected to %s with namespace %s", csinfo.Config.Host, ns)
	step.Done()

	step = sg.Add("Deleting service...")
	if err := serviceclient.Delete(ctx, state.Name, metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		log.Debug("no service found, continuing")
	}

	step.Update("Service deleted")
	step.Done()
	return nil
}

func (r *Releaser) resourceIngressCreate(
	ctx context.Context,
	log hclog.Logger,
	target *Deployment,
	result *Release,
	state *Resource_Ingress,
	serviceState *Resource_Service,
	csinfo *clientsetInfo,
	sg terminal.StepGroup,
) error {
	// Preflight config checks
	if r.config.IngressConfig == nil {
		// No ingress config, we're not going to configure one
		return nil
	}

	if r.config.LoadBalancer {
		return status.Error(codes.FailedPrecondition, "A LoadBalancer service type is not "+
			"compatible with an Ingress config. Please pick one or the other for your release")
	}

	if r.config.NodePort != 0 {
		return status.Error(codes.FailedPrecondition, "A NodePort service type is not "+
			"compatible with an Ingress config. Please pick one or the other for your release")
	}

	if r.config.IngressConfig.ClassName != "http" {
		return status.Error(codes.FailedPrecondition, "An ingress stanza must be "+
			"of type \"http\".")
	}

	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	step := sg.Add("Preparing ingress resource...")
	defer func() { step.Abort() }() // Defer in func in case more steps are added to this func in the future

	clientSet := csinfo.Clientset
	serviceclient := clientSet.CoreV1().Services(ns)
	ingressClient := clientSet.NetworkingV1().Ingresses(ns)

	// Determine if we have a deployment that we manage already
	serviceBackend, err := serviceclient.Get(ctx, serviceState.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// We expect resourceServiceCreate to have created a Service prior to
		// creating an ingress resource. Otherwise there is no service backend
		// the ingress resource can refer to
		return status.Errorf(codes.FailedPrecondition, "A service must exist before "+
			"an ingress resource can be created: %s", err)
	}
	if err != nil {
		return err
	}

	var ingressResource *networkingv1.Ingress
	create := false
	ingressResource, err = ingressClient.Get(ctx, serviceState.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// We haven't created an ingress resource yet...
		log.Debug("no ingress resource found, will create one")

		err = nil
		create = true

		// basic ingress resource
		ingressResource = &networkingv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Ingress",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name: result.ServiceName,
			},
		}
	}
	if err != nil {
		return err
	}

	// Set the ingress resource state name to match the service name
	state.Name = result.ServiceName

	var ingressTls networkingv1.IngressTLS
	if r.config.IngressConfig.TlsConfig != nil {
		ingressTls = networkingv1.IngressTLS{
			Hosts:      r.config.IngressConfig.TlsConfig.Hosts,
			SecretName: r.config.IngressConfig.TlsConfig.SecretName,
		}
	}

	// Apply any annotations to the ingress resource
	if r.config.IngressConfig.Annotations != nil {
		ingressResource.ObjectMeta.Annotations = r.config.IngressConfig.Annotations
	}

	// Define the ingress resources backend service it should route traffic to
	ingressBackend := networkingv1.IngressBackend{
		Service: &networkingv1.IngressServiceBackend{
			Name: serviceState.Name,
			Port: networkingv1.ServiceBackendPort{
				Name:   serviceBackend.Spec.Ports[0].Name,
				Number: serviceBackend.Spec.Ports[0].Port,
			},
		},
	}

	if r.config.IngressConfig.PathType == "" {
		r.config.IngressConfig.PathType = "Prefix"
	}
	pathType := networkingv1.PathType(r.config.IngressConfig.PathType)

	// Set the default path to '/' if not set and path type is Prefix
	if pathType == networkingv1.PathTypePrefix && r.config.IngressConfig.Path == "" {
		r.config.IngressConfig.Path = "/"
	}

	httpIngressPath := networkingv1.HTTPIngressPath{
		Path:     r.config.IngressConfig.Path,
		PathType: &pathType,
		Backend:  ingressBackend,
	}

	ingressRule := networkingv1.IngressRule{
		Host: r.config.IngressConfig.Host,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{httpIngressPath},
			},
		},
	}

	ingressResource.Spec = networkingv1.IngressSpec{
		Rules: []networkingv1.IngressRule{ingressRule},
	}

	if r.config.IngressConfig.TlsConfig != nil {
		ingressResource.Spec.TLS = []networkingv1.IngressTLS{ingressTls}
	}

	if r.config.IngressConfig.DefaultBackend {
		ingressResource.Spec.DefaultBackend = &ingressBackend
	}

	// create or update the ingress resource
	if create {
		step.Update("Creating ingress resource...")
		_, err = ingressClient.Create(ctx, ingressResource, metav1.CreateOptions{})
	} else {
		step.Update("Updating existing ingress resource...")
		_, err = ingressClient.Update(ctx, ingressResource, metav1.UpdateOptions{})
	}
	if err != nil {
		return err
	}

	step.Done()
	step = sg.Add("Waiting for ingress resource to become ready...")

	// Wait on load balancer
	var ingress *networkingv1.Ingress
	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		ingress, err = ingressClient.Get(ctx, state.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		return len(ingress.Status.LoadBalancer.Ingress) > 0, nil
	})
	if err != nil {
		return err
	}

	step.Update("Ingress resource is ready!")
	step.Done()

	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		protocol := "http://"
		if r.config.IngressConfig.TlsConfig != nil {
			protocol = "https://"
		}

		if r.config.IngressConfig.Host != "" {
			// We set the requested hostname from the waypoint.hcl if defined
			result.Url = protocol + r.config.IngressConfig.Host
		} else {
			// set the hostname based on the load balancer configured in k8s
			lbIngress := ingress.Status.LoadBalancer.Ingress[0]
			result.Url = protocol + lbIngress.IP
			if lbIngress.Hostname != "" {
				result.Url = protocol + lbIngress.Hostname
			}
		}

		if serviceBackend.Spec.Ports[0].Port != 80 {
			result.Url = fmt.Sprintf("%s:%d", result.Url, serviceBackend.Spec.Ports[0].Port)
		}

		if r.config.IngressConfig.Path != "/" {
			result.Url = fmt.Sprintf("%s%s", result.Url, r.config.IngressConfig.Path)
		}

	} // else we show the cluster IP URL setup by the service resource

	return nil
}

func (r *Releaser) resourceIngressDestroy(
	ctx context.Context,
	state *Resource_Ingress,
	sg terminal.StepGroup,
	csinfo *clientsetInfo,
) error {
	if state.Name == "" {
		// we didn't create an ingress resource, so we can't delete it either
		return nil
	}

	step := sg.Add("Initializing Kubernetes client...")
	defer step.Abort()

	// Prepare our namespace and override if set.
	ns := csinfo.Namespace
	if r.config.Namespace != "" {
		ns = r.config.Namespace
	}

	clientSet := csinfo.Clientset
	ingressClient := clientSet.NetworkingV1().Ingresses(ns)

	step.Update("Kubernetes client connected to %s with namespace %s", csinfo.Config.Host, ns)
	step.Done()

	// get ingress first
	_, err := ingressClient.Get(ctx, state.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		step = sg.Add("No Ingress resource found to delete")
		step.Done()
		return nil
	}
	if err != nil {
		return err
	}

	step = sg.Add("Deleting ingress resource...")
	if err = ingressClient.Delete(ctx, state.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	step.Update("Ingress resource deleted")
	step.Done()
	return nil
}

func (r *Releaser) resourceIngressStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Ingress,
	clientset *clientsetInfo,
	sr *resource.StatusResponse,
) error {
	if state.Name == "" {
		log.Debug("no state found for ingress resource, cannot build status report")
		return nil
	}

	s := sg.Add("Checking status of Kubernetes ingress resource %q...", state.Name)
	defer s.Abort()

	namespace := r.config.Namespace
	if namespace == "" {
		namespace = clientset.Namespace
	}

	ingressResource := sdk.StatusReport_Resource{
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_ROUTER,
	}
	sr.Resources = append(sr.Resources, &ingressResource)

	ingressClient := clientset.Clientset.NetworkingV1().Ingresses(namespace)
	ingressResp, err := ingressClient.Get(ctx, state.Name, metav1.GetOptions{})
	if ingressResp == nil {
		return status.Errorf(codes.FailedPrecondition,
			"kubernetes ingress resource response returned nothing for %q", state.Name)
	} else if err != nil {
		if !errors.IsNotFound(err) {
			return err
		} else {
			s.Update("No ingress resource named %q was found", state.Name)
			s.Status(terminal.StatusError)
			s.Done()
			s = sg.Add("")

			ingressResource.Name = state.Name
			ingressResource.Health = sdk.StatusReport_MISSING
			ingressResource.HealthMessage = sdk.StatusReport_MISSING.String()
		}
	} else {
		// An ingress resource exists
		s.Update("Building status report for ingress resource %q...", state.Name)

		lbIngress := ingressResp.Status.LoadBalancer.Ingress[0]

		var ipAddress string
		if lbIngress.IP != "" {
			ipAddress = lbIngress.IP
		}
		var hostname string
		if lbIngress.Hostname != "" {
			hostname = lbIngress.Hostname
		}
		// we only configure 1 rule, so grab the first service resource
		serviceBackend := ingressResp.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.Service

		ingressJson, err := json.Marshal(map[string]interface{}{
			"ingress":        ingressResp,
			"ipAddress":      ipAddress,
			"hostname":       hostname,
			"serviceBackend": serviceBackend,
		})
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to marshal ingress resource definition to json: %s", err)
		}

		ingressResource.Id = fmt.Sprintf("%s", ingressResp.UID)
		ingressResource.Name = ingressResp.Name
		ingressResource.CreatedTime = timestamppb.New(ingressResp.CreationTimestamp.Time)
		// If we found the ingress resource, then it's ready. It doesn't have any
		// other conditions.
		ingressResource.Health = sdk.StatusReport_READY
		ingressResource.HealthMessage = fmt.Sprintf("Ingress resource %q exists and is ready", ingressResource.Name)
		ingressResource.StateJson = string(ingressJson)
	}

	s.Update("Finished building report for Kubernetes ingress resource")
	s.Done()
	return nil
}

// getClientset is a value provider for our resource manager and provides
// the connection information used by resources to interact with Kubernetes.
func (r *Releaser) getClientset() (*clientsetInfo, error) {
	// Get our client
	clientSet, ns, config, err := Clientset(r.config.KubeconfigPath, r.config.Context)
	if err != nil {
		return nil, err
	}

	return &clientsetInfo{
		Clientset: clientSet,
		Namespace: ns,
		Config:    config,
	}, nil
}

// Release creates a Kubernetes service configured for the deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	ui terminal.UI,
	target *Deployment,
	dcr *component.DeclaredResourcesResp,
) (*Release, error) {
	var result Release

	if job.Workspace == "default" {
		result.ServiceName = src.App
	} else {
		result.ServiceName = src.App + "-" + job.Workspace
	}

	sg := ui.StepGroup()
	defer sg.Wait()

	// Create our resource manager and create
	rm := r.resourceManager(log, dcr)
	if err := rm.CreateAll(
		ctx, log, sg, ui,
		target, &result,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	result.ResourceState = rm.State()

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (r *Releaser) Destroy(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) error {

	sg := ui.StepGroup()
	defer sg.Wait()

	rm := r.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if release.ResourceState == nil {
		rm.Resource("service").SetState(&Resource_Service{
			Name: release.ServiceName,
		})
		rm.Resource("ingress").SetState(&Resource_Service{
			Name: release.ServiceName,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return err
		}
	}

	// Destroy
	return rm.DestroyAll(ctx, log, sg, ui)
}

func (r *Releaser) Status(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := r.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if release.ResourceState == nil {
		rm.Resource("service").SetState(&Resource_Service{
			Name: release.ServiceName,
		})
		rm.Resource("ingress").SetState(&Resource_Service{
			Name: release.ServiceName,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return nil, err
		}
	}

	step := sg.Add("Gathering health report for Kubernetes release...")
	defer step.Abort()

	resources, err := rm.StatusAll(ctx, log, sg, ui)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resource manager failed to generate resource statuses: %s", err)
	}

	if len(resources) == 0 {
		// This shouldn't happen - the status func for the releaser should always return a resource or an error.
		return nil, status.Errorf(codes.Internal, "no resources generated for release - cannot determine status.")
	}

	var serviceResource *sdk.StatusReport_Resource
	for _, r := range resources {
		if r.Type == "service" {
			serviceResource = r
			break
		}
	}
	if serviceResource == nil {
		return nil, status.Errorf(codes.Internal, "no service resource found - cannot determine overall health")
	}

	// Create our status report
	result := sdk.StatusReport{
		External:      true,
		GeneratedTime: timestamppb.Now(),
		Resources:     resources,
		Health:        serviceResource.Health,
		HealthMessage: serviceResource.HealthMessage,
	}

	log.Debug("status report complete")

	// update output based on main health state
	step.Update("Finished building report for Kubernetes platform")
	step.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	// More UI detail for non-ready resources
	for _, resource := range result.Resources {
		if resource.Health != sdk.StatusReport_READY {
			st.Step(terminal.StatusWarn, fmt.Sprintf("Resource %q is reporting %q", resource.Name, resource.Health.String()))
		}
	}

	return &result, nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	// Annotations to be applied to the kube service.
	Annotations map[string]string `hcl:"annotations,optional"`

	// Ingress represents an config for setting up an ingress resource.
	IngressConfig *IngressConfig `hcl:"ingress,block"`

	// KubeconfigPath is the path to the kubeconfig file. If this is
	// blank then we default to the home directory.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// Context specifies the kube context to use.
	Context string `hcl:"context,optional"`

	// Load Balancer sets whether or not the service will be a load
	// balancer type service
	LoadBalancer bool `hcl:"load_balancer,optional"`

	// Port configures the port that is used to access the service.
	// The default is 80.
	// Not valid if `Ports` is already defined
	// If defined, will internally be stored into `Ports`
	Port uint `hcl:"port,optional"`

	// A full resource of options to define ports for a service
	Ports []map[string]string `hcl:"ports,optional"`

	// NodePort configures a port to access the service on whichever node
	// is running service.
	// Not valid if `Ports` is already defined
	// If defined, will internally be stored into `Ports`
	NodePort uint `hcl:"node_port,optional"`

	// Namespace is the Kubernetes namespace to target the deployment to.
	Namespace string `hcl:"namespace,optional"`
}

// IngressConfig holds various options to configure an Ingress resource with
// during a release. It currently only spports 'http' based route rules.
type IngressConfig struct {
	// Annotations to be applied to the ingress service.
	Annotations map[string]string `hcl:"annotations,optional"`

	// Currently Waypoint only supports "HTTP" rule-backed ingress resources.
	// We include this stanza label in the future for when Kubernetes has other
	// kinds of rule types.
	ClassName string `hcl:",label"`

	// If set, this will configure the given ingress resources backend service
	// as the default service that accepts traffic if no route rules match from
	// the inbound request.
	// Defaults to false
	DefaultBackend bool `hcl:"default,optional"`

	// If set, this option will configure the ingress controller to accept
	// traffic from the defined hostname. IPs are not allowed, nor are `:` delimiters.
	// Wildcards are allowed to a certain extent. For more details, check out the
	// k8s go client package.
	Host string `hcl:"host,optional"`

	// Defines the kind of rule the path will be. Possible values are:
	// 'Exact', 'Prefix', and 'ImplementationSpecific'.
	// We believe most users expect a 'Prefix' type, so we default to 'Prefix' if
	// not specified.
	PathType string `hcl:"path_type,optional"`

	// Path represents a rule to route requests to. I.e. if an inbound request
	// matches a route like `/foo`, the ingress controller would see that this
	// resource is configured for that rule, and would route traffic to this
	// ingress resources service backend.
	Path string `hcl:"path,optional"`

	// TlsConfig is an optional config that users can set to enable HTTPS traffic
	TlsConfig *IngressTls `hcl:"tls,block"`
}

// IngressTls holds options required to configure an ingress resource with TLS.
type IngressTls struct {
	// Hosts is a list of hosts included in the TLS certificate
	Hosts []string `hcl:"hosts,optional"`

	// SecretName references the name of the secret created inside Kubernetes
	// associated with the TLS certificate information. If cert-manager is used,
	// this name will refer to the secret cert-manager should create when generating
	// certficiates for the secret.
	SecretName string `hcl:"secret_name,optional"`
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Manipulates the Kubernetes Service activate Deployments")
	doc.Input("k8s.Deployment")
	doc.Output("k8s.Release")

	doc.SetField(
		"annotations",
		"Annotations to be applied to the kube service",
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
		"load_balancer",
		"indicates if the Kubernetes Service should LoadBalancer type",
		docs.Summary(
			"if the Kubernetes Service is not a LoadBalancer and node_port is not",
			"set, then the Service uses ClusterIP",
		),
	)

	doc.SetField(
		"node_port",
		"the TCP port that the Service should consume as a NodePort",
		docs.Summary(
			"if this is set but load_balancer is not, the service will be NodePort type,",
			"but if load_balancer is also set, it will be LoadBalancer",
		),
	)

	doc.SetField(
		"port",
		"the TCP port that the application is listening on",
		docs.Default(fmt.Sprint(DefaultPort)),
	)

	doc.SetField(
		"ports",
		"a map of ports and options that the application is listening on",
		docs.Summary(
			"used to define and configure multiple ports that the application is",
			"listening on. Available keys are 'port', 'node_port', 'name', and 'target_port'.",
			"If 'node_port' is set but 'load_balancer' is not, the service will be",
			" NodePort type. If 'load_balancer' is also set, it will be LoadBalancer.",
			"Ports defined will be TCP protocol.",
			"Note that 'name' is required if defining more than one port.",
		),
	)

	doc.SetField(
		"namespace",
		"namespace to create Service in",
		docs.Summary(
			"namespace is the name of the Kubernetes namespace to create the deployment in",
			"This is useful to create Services in non-default namespaces without creating kubeconfig contexts for each",
		),
	)

	doc.SetField(
		"ingress",
		"Configuration to set up an ingress resource to route traffic to the given "+
			"application from an ingress controller",
		docs.Summary(
			"An ingress resource can be created on release that will route traffic "+
				"to the Kubernetes service. Note that before this happens, the Kubernetes "+
				"cluster must already be configured with an Ingress controller. Otherwise "+
				"there won't be a way for inbound traffic to be routed to the ingress resource.",
		),
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"annotations",
				"Annotations to be applied to the ingress resource",
			)

			doc.SetField(
				"default",
				"sets the ingress resource to be the default backend for any traffic "+
					"that doesn't match existing ingress rule paths",
				docs.Default("false"),
			)

			doc.SetField(
				"host",
				"If set, will configure the ingress resource to have the ingress controller "+
					"route traffic for any inbound requests that match this host. IP addresses "+
					"are not allowed, nor are ':' delimiters. Wildcards are allowed to a "+
					"certain extent. For more details check out the Kubernetes documentation",
			)

			doc.SetField(
				"path_type",
				"defines the kind of rule the path will be for the ingress controller. "+
					"Valid path types are 'Exact', 'Prefix', and 'ImplementationSpecific'.",
				docs.Default("Prefix"),
			)

			doc.SetField(
				"path",
				"The route rule that should be used to route requests to this ingress resource. "+
					"A path must begin with a '/'.",
				docs.Default("/"),
			)

			doc.SetField(
				"tls",
				"A stanza of TLS configuration options for traffic to the ingress resource",
				docs.SubFields(func(doc *docs.SubFieldDoc) {
					doc.SetField(
						"hosts",
						"A list of hosts included in the TLS certificate",
					)

					doc.SetField(
						"secret_name",
						"The Kubernetes secret name that should be used to look up or store TLS configs",
					)
				}),
			)
		}),
	)

	return doc, nil
}

var mixedHealthReleaseWarn = strings.TrimSpace(`
Waypoint detected that the current release is not ready, however your application
might be available or still starting up.
`)

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Destroyer      = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
	_ component.Status         = (*Releaser)(nil)
)
