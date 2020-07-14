// Package k8s contains components for deploying to Kubernetes.
package k8s

import (
	"github.com/hashicorp/waypoint/sdk"
)

//go:generate protoc -I ../../.. --go_opt=plugins=grpc --go_out=../../.. waypoint/builtin/k8s/plugin.proto

// Options are the SDK options to use for instantiation for
// the Google Cloud Run plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}, &Releaser{}),
}
