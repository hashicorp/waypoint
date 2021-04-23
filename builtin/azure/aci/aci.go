// Package aci contains components for deploying to Azure ACI.
package aci

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. ../../../../waypoint/builtin/azure/aci/plugin.proto --go_out=plugins=grpc:../../../..

// Options are the SDK options to use for instantiation for
// the Azure ACI plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
