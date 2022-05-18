package jobspec

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. -I ../../../thirdparty/proto --go_opt=paths=source_relative --go_out=../../../.. --go-grpc_opt=paths=source_relative --go-grpc_out=../../../.. waypoint/builtin/nomad/jobspec/plugin.proto

// Options are the SDK options to use for instantiation for
// the Nomad plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}),
}
