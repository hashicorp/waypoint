// Package k8s contains components for deploying to Kubernetes.
package k8s

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. ../../../waypoint/builtin/k8s/plugin.proto --go_out=plugins=grpc:../../..

// Options are the SDK options to use for instantiation for
// the Kubernetes plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}, &ConfigSourcer{}),
}
