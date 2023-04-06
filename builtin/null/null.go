// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package null contains components that do [almost] nothing, primarily aimed
// to ease experimentation and testing with Waypoint. For example, the null
// config sourcer can be used to learn about dynamic configuration without
// the complexity of configuring a real remote system such as Vault. This helps
// learn the Waypoint side of things before diving into a more real-world
// system.
package null

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation for this plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}, &Builder{}, &Platform{}, &Releaser{}),
}
