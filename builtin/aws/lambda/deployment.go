// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package lambda

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)
