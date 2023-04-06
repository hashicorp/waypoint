// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
