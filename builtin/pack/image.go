// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pack

func (i *DockerImage) Labels() map[string]string {
	return i.BuildLabels
}
