// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecs

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. -I ../../../thirdparty/proto --go_out=../../../.. --go-grpc_out=../../../.. waypoint/builtin/aws/ecs/plugin.proto

const platformName = "aws-ecs"

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &TaskLauncher{}),
}
