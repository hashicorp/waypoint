// Package lambda contains components for deploying to AWS Lambda
package lambda

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. ../../../../waypoint/builtin/aws/lambda/plugin.proto --go_out=plugins=grpc:../../../..

// Options are the SDK options to use for instantiation for
// the Google Cloud Run plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
