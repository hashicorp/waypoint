package ecr

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_opt=plugins=grpc --go_out=../../../.. waypoint/builtin/aws/ecr/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Registry{}),
	sdk.WithMappers(ECRImageMapper),
}
