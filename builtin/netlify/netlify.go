// Package netlify contains components for deploying to Netlify.
package netlify

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. ../../../waypoint/builtin/netlify/plugin.proto --go_out=plugins=grpc:../../..

// Options are the SDK options to use for instantiation for
// the Netlfiy plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
