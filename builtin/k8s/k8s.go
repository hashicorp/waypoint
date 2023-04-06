// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package k8s contains components for deploying to Kubernetes.
package k8s

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. -I ../../thirdparty/proto --go_out=../../.. --go-grpc_out=../../.. waypoint/builtin/k8s/plugin.proto

const platformName = "kubernetes"

// Options are the SDK options to use for instantiation for
// the Kubernetes plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}, &ConfigSourcer{}, &TaskLauncher{}),
}
