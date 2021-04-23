package apply

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. ../../../../waypoint/builtin/k8s/apply/plugin.proto --go_out=plugins=grpc:../../../..

// Options are the SDK options to use for instantiation for the plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
