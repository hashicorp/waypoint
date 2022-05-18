package pack

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. --go_opt=paths=source_relative --go_out=../../.. waypoint/builtin/pack/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}),
	sdk.WithMappers(PackImageMapper),
}
