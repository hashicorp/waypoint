// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package consul contains components for syncing app configuration with Consul.
package consul

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation for this plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
