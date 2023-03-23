// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pack

func (i *DockerImage) Labels() map[string]string {
	return i.BuildLabels
}
