// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dockerpull

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}),
}
