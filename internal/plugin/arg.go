package plugin

import (
	"github.com/evanphx/opaqueany"
	"github.com/hashicorp/go-argmapper"
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
