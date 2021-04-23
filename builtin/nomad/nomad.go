package nomad

import sdk "github.com/hashicorp/waypoint-plugin-sdk"

//go:generate protoc -I ../../.. ../../../waypoint/builtin/nomad/plugin.proto --go_out=plugins=grpc:../../..

// Options are the SDK options to use for instantiation for
// the Nomad plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
