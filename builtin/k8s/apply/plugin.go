package apply

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../../.. --go-grpc_opt=paths=source_relative waypoint/builtin/k8s/apply/plugin.proto

// Options are the SDK options to use for instantiation for the plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
