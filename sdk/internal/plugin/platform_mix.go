package plugin

import (
	"github.com/mitchellh/devflow/sdk/component"
)

// This file contains the structs that "mix" different platform types
// together so that we can return multiple implements as a single interface{}
// result.
//
// The name of the structs is purposely weird and not very idiomatic, we use:
//
//     platform_X_Y_Z
//
/// Where "X", "Y", "Z" are interfaces that are implemented in alphabetical order.

type platform_Log struct {
	component.ConfigurableNotify
	component.Platform
	component.LogPlatform
}

// This may seem silly but due to the way Go handles multiple embedded
// interfaces with overlapping types (i.e. Configurable and ConfigurableNotify),
// we ran into a bug where it didn't implement EITHER by including BOTH. So
// to be extra sure we have this check that it implements it although it seems
// redundant.
var (
	_ component.Configurable       = (*platform_Log)(nil)
	_ component.ConfigurableNotify = (*platform_Log)(nil)
	_ component.Platform           = (*platform_Log)(nil)
	_ component.LogPlatform        = (*platform_Log)(nil)
)
