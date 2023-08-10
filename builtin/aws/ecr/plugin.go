// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ecr

func (i *Image) Name() string {
	return i.Image + ":" + i.Tag
}
