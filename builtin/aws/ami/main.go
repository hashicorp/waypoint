// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ami

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. --go_out=../../../.. --go-grpc_out=../../../.. waypoint/builtin/aws/ami/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}),
}
