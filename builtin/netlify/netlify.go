// Package netlify contains components for deploying to Netlify.
package netlify

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate protoc -I ../../.. --go_opt=plugins=grpc --go_out=../../.. waypoint/builtin/netlify/plugin.proto

// Options are the SDK options to use for instantiation for
// the Netlfiy plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
