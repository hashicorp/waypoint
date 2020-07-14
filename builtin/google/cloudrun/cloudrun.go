// Package cloudrun contains components for deploying to Google Cloud Run.
package cloudrun

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate protoc -I ../../../.. --go_opt=plugins=grpc --go_out=../../../.. waypoint/builtin/google/cloudrun/plugin.proto

// Options are the SDK options to use for instantiation for
// the Google Cloud Run plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}),
}
