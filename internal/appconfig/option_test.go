// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package appconfig

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/plugin"
)

// testWithConfigSourcer can be used by tests to set a specific
// ConfigSourcer implementation. It may make sense to export this
// in the future but we don't have a need for it yet.
func testWithConfigSourcer(n string, cs component.ConfigSourcer) Option {
	return func(w *Watcher) error {
		if w.plugins == nil {
			w.plugins = map[string]*plugin.Instance{}
		}

		w.plugins[n] = &plugin.Instance{
			Component: cs,
		}

		return nil
	}
}
