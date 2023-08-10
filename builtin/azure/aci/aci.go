// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package aci contains components for deploying to Azure ACI.
package aci

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_out=../../../.. --go-grpc_out=../../../.. waypoint/builtin/azure/aci/plugin.proto

// Options are the SDK options to use for instantiation for
// the Azure ACI plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
