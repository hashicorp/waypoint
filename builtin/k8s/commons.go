package k8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

func createScratchVolumes(scratchSpace []string) []corev1.Volume {
	var volumes []corev1.Volume
	for idx := range scratchSpace {
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
	return volumes
}

func createVolumeMounts(scratchSpace []string, volumes []corev1.Volume) []corev1.VolumeMount {
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
	return volumeMounts
}
