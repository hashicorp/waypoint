package docker

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate sh -c "protoc -I ./ ./*.proto --go_out=plugins=grpc:./"

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}, &Registry{}),
	sdk.WithMappers(PackImageMapper),
}
