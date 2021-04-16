package apply

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_opt=plugins=grpc --go_out=../../../.. waypoint/builtin/k8s/apply/plugin.proto

// Options are the SDK options to use for instantiation for the plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
