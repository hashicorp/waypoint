// Package files contains a component for validating local files.
package files

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. ../../../waypoint/builtin/files/plugin.proto --go_out=plugins=grpc:../../..

// Options are the SDK options to use for instantiation for
// the Files plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}, &Registry{}),
}
