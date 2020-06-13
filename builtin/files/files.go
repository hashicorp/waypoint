// Package files contains a component for validating local files.
package files

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate sh -c "protoc -I ./ ./*.proto --go_out=plugins=grpc:./"

// Options are the SDK options to use for instantiation for
// the Files plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}, &Registry{}),
}
