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
		healthResult = sdk.StatusReport_PARTIAL
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

func buildStatusReport(
	podList *corev1.PodList,
) sdk.StatusReport {
	var result sdk.StatusReport
	result.External = true
	resources := make([]*sdk.StatusReport_Resource, len(podList.Items))

	// Report on most recently observed status of a deployments pod
	for _, pod := range podList.Items {
		// Overall Pod Health
		podStatus := pod.Status
		result.HealthMessage = podStatus.Message
		result.Health = podPhaseToHealth(podStatus.Phase)

		// Pod containers health
		for i, containerStatus := range podStatus.ContainerStatuses {
			resources[i] = containerStatusToHealth(containerStatus)
		}
	}

	return result
}
