// Package aci contains components for deploying to Azure ACI.
package aci

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate sh -c "protoc -I ./ ./*.proto --go_out=plugins=grpc:./"

// Options are the SDK options to use for instantiation for
// the Azure ACI plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}),
}
