// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecr

func (i *Image) Name() string {
	return i.Image + ":" + i.Tag
}
