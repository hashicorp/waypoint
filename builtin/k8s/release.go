// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

// newService returns the basic structure for a new service.
// This isn't ready to deploy right away.
func (*Release) newService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (r *Release) URL() string { return r.Url }

var (
	_ component.Release = (*Release)(nil)
)
