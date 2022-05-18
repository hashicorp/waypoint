package exec

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. --go_opt=paths=source_relative --go_out=../../.. --go-grpc_opt=paths=source_relative --go-grpc_out=../../.. waypoint/builtin/exec/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
	sdk.WithMappers(DockerImageMapper),
}
