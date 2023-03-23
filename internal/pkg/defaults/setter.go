// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

// Setter is an interface for setting default values
type Setter interface {
	SetDefaults()
}

func callSetter(v interface{}) {
	if ds, ok := v.(Setter); ok {
		ds.SetDefaults()
	}
}
