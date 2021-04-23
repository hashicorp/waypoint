package exec

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. ../../../waypoint/builtin/exec/plugin.proto --go_out=plugins=grpc:../../..

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
	sdk.WithMappers(DockerImageMapper),
}
