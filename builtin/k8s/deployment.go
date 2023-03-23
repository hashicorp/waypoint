// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

// newDeployment returns the basic structure for a new deployment.
// This isn't ready to deploy right away.
func (d *Deployment) newDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		// Note both name and app are included here. 'app' is expected for certain
		// k8s integrations, where as waypoint expects 'name' else where for release
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  name,
					"name": name,
				},
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  name,
						"name": name,
					},

					Annotations: map[string]string{},
				},
			},
		},
	}
}

var _ component.Deployment = (*Deployment)(nil)
