// Package netlify contains components for deploying to Netlify.
package netlify

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate sh -c "protoc -I ./ ./*.proto --go_out=plugins=grpc:./"

// Options are the SDK options to use for instantiation for
// the Netlfiy plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
