// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package alb

import "github.com/hashicorp/waypoint-plugin-sdk/component"

func (r *Release) URL() string { return r.Url }

var _ component.Release = (*Release)(nil)
