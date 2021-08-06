package plugin

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
)

// ArgNamedAny returns an argmapper.Arg that specifies the Any value
// with the proper subtype.
func ArgNamedAny(n string, v *any.Any) argmapper.Arg {
	if v == nil {
		return nil
	}

	msg, err := ptypes.AnyMessageName(v)
	if err != nil {
		// This should never happen.
		panic(err)
	}

	return argmapper.NamedSubtype(n, v, msg)
}
