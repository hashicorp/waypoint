package apprunner

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../../.. -I ../../../thirdparty/proto --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative waypoint/builtin/aws/apprunner/plugin.proto

// App Runner only supports ECR (private) and ECR_PUBLIC image registries;
// It will not have a Builder Component.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
