// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. -I ../../thirdparty/proto/opaqueany --go_out=../../.. --go-grpc_out=../../.. waypoint/builtin/docker/plugin.proto

const platformName = "docker"

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}, &Registry{}, &Platform{}, &TaskLauncher{}),
	// sdk.WithMappers(PackImageMapper),
}
