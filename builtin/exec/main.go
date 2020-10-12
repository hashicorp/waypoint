package exec

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. --go_opt=plugins=grpc --go_out=../../.. waypoint/builtin/exec/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
	sdk.WithMappers(DockerImageMapper),
}
