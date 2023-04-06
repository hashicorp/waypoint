// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function_url

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../../.. -I ../../../../thirdparty/proto --go_out=../../../../.. --go-grpc_out=../../../../.. waypoint/builtin/aws/lambda/function_url/plugin.proto

// Options are the SDK options to use for instantiation for the plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Releaser{}),
}
