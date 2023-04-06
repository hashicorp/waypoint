// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/opaqueany"
)

// ArgNamedAny returns an argmapper.Arg that specifies the Any value
// with the proper subtype.
func ArgNamedAny(n string, v *opaqueany.Any) argmapper.Arg {
	if v == nil {
		return nil
	}

	msg := v.MessageName()

	return argmapper.NamedSubtype(n, v, string(msg))
}
