// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lambda

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)
