// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newServiceAccount returns the basic structure for a service account.
func newServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}
