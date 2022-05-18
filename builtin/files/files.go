// Package files contains a component for validating local files.
package files

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. --go_opt=paths=source_relative --go_out=../../.. --go-grpc_opt=paths=source_relative --go-grpc_out=../../.. waypoint/builtin/files/plugin.proto

// Options are the SDK options to use for instantiation for
// the Files plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}, &Registry{}),
}
