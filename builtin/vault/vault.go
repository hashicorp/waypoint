// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package vault contains components for syncing secrets with Vault.
package vault

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation for this plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
