// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

// Name is the full name including the tag.
func (i *Image) Name() string {
	return i.Image + ":" + i.Tag
}
