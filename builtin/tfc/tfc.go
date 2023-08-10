// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package tfc contains components for syncing outputs of states from Terraform Cloud
package tfc

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation for this plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
