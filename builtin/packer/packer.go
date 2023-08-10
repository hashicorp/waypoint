// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package packer

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
