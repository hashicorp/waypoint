package function_url

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../../.. -I ../../../../thirdparty/proto --go_out=../../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../../.. --go-grpc_opt=paths=source_relative waypoint/builtin/aws/lambda/function_url/plugin.proto

// Options are the SDK options to use for instantiation for the plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Releaser{}),
}
