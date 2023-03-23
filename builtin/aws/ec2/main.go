// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_out=../../../.. --go-grpc_out=../../../.. waypoint/builtin/aws/ec2/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
