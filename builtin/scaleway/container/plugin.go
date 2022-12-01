package container

import sdk "github.com/hashicorp/waypoint-plugin-sdk"

//go:generate protoc -I ../../../.. -I ../../../thirdparty/proto/opaqueany --go_opt=plugins=grpc --go_out=../../../.. waypoint/builtin/scaleway/scalewaycontainer/plugin.proto

var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
