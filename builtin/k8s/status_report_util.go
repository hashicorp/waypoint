package k8s

import (
	corev1 "k8s.io/api/core/v1"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

// Translate a K8S container status into a Waypoint Health Status
func containerStatusToHealth(
	containerStatus corev1.ContainerStatus,
) *sdk.StatusReport_Resource {
	var resourceHealth sdk.StatusReport_Resource
	resourceHealth.Name = containerStatus.Name

	// ContainerStatus.State is a struct of possible container states. If one
	// is set, that is the current state of the container
	if containerStatus.State.Running != nil || containerStatus.Ready == true {
		resourceHealth.Health = sdk.StatusReport_READY
		resourceHealth.HealthMessage = "container is reporting running" // no message defined by k8s api
	} else if containerStatus.State.Waiting != nil {
		resourceHealth.Health = sdk.StatusReport_PARTIAL
		resourceHealth.HealthMessage = containerStatus.State.Waiting.Message
	} else if containerStatus.State.Terminated != nil {
		resourceHealth.Health = sdk.StatusReport_DOWN
		resourceHealth.HealthMessage = containerStatus.State.Terminated.Message
	} else {
		resourceHealth.Health = sdk.StatusReport_UNKNOWN
		resourceHealth.HealthMessage = "container health could not be determined"
	}

	return &resourceHealth
}

// Translate a Pod Phase into a Waypoint Health Status
func podPhaseToHealth(
	phase corev1.PodPhase,
) sdk.StatusReport_Health {
	var healthResult sdk.StatusReport_Health

	switch phase {
	case corev1.PodPending:
		healthResult = sdk.StatusReport_ALIVE
	case corev1.PodRunning:
		healthResult = sdk.StatusReport_READY
	case corev1.PodSucceeded:
		healthResult = sdk.StatusReport_PARTIAL
	case corev1.PodFailed:
		healthResult = sdk.StatusReport_DOWN
	case corev1.PodUnknown:
		healthResult = sdk.StatusReport_UNKNOWN
	default:
		healthResult = sdk.StatusReport_UNKNOWN
	}

	return healthResult
}

// Take a list of Pods and build a Waypoint Status Report based on their reported health
func buildStatusReport(
	podList *corev1.PodList,
) sdk.StatusReport {
	var result sdk.StatusReport
	result.External = true
	var resources []*sdk.StatusReport_Resource

	// Build health for every possible pod for overall health report
	var ready, alive, down, unknown int

	// Report on most recently observed status of a deployments pod
	// Pod resources and its containers will be in order inside Resources
	for _, pod := range podList.Items {
		// Overall Pod Health
		podStatus := pod.Status

		// Determine overall health report based on all pods
		switch podStatus.Phase {
		case corev1.PodPending:
			alive++
		case corev1.PodRunning:
			// Extra checks on the latest condition to ensure pod is reporting ready and running
			for _, c := range podStatus.Conditions {
				if c.Status == corev1.ConditionTrue && c.Type == corev1.PodReady {
					ready++
					break
				}
			}

			alive++
		case corev1.PodSucceeded:
			alive++
		case corev1.PodFailed:
			down++
		case corev1.PodUnknown:
			unknown++
		default:
			unknown++
		}

		podHealth := podPhaseToHealth(podStatus.Phase)
		resources = append(resources, &sdk.StatusReport_Resource{
			Health:        podHealth,
			HealthMessage: podStatus.Message,
			Name:          pod.ObjectMeta.Name,
		})

		// Pod containers health
		for _, containerStatus := range podStatus.ContainerStatuses {
			resources = append(resources, containerStatusToHealth(containerStatus))
		}
	}

	// Overall health status for report
	if ready == len(podList.Items) {
		result.Health = sdk.StatusReport_READY
		result.HealthMessage = "all pods are reporting ready"
	} else if down == len(podList.Items) {
		result.Health = sdk.StatusReport_DOWN
		result.HealthMessage = "all pods are reporting down"
	} else if unknown == len(podList.Items) {
		result.Health = sdk.StatusReport_UNKNOWN
		result.HealthMessage = "status of all pods cannot be determined"
	} else if alive == len(podList.Items) {
		result.Health = sdk.StatusReport_ALIVE
		result.HealthMessage = "all pods are reporting alive"
	} else {
		result.Health = sdk.StatusReport_PARTIAL
		result.HealthMessage = "all pods are reporting a mixed status"
	}

	return result
}
