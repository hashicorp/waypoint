// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package nomad

import sdk "github.com/hashicorp/waypoint-plugin-sdk"

//go:generate protoc -I ../../.. -I ../../thirdparty/proto --go_out=../../.. --go-grpc_out=../../.. waypoint/builtin/nomad/plugin.proto

// Options are the SDK options to use for instantiation for
// the Nomad plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &TaskLauncher{}),
}
