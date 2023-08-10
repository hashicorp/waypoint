// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package alb

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. -I ../../../thirdparty/proto --go_out=../../../.. --go-grpc_out=../../../.. waypoint/builtin/aws/alb/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Releaser{}),
	sdk.WithMappers(EC2TGMapper, LambdaTGMapper),
}
