//go:build tools
// +build tools

// To install the following tools at the version used by this repo run:
// $ make tools
// or
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install github.com/kevinburke/go-bindata
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install github.com/mitchellh/protoc-gen-go-json

// Using a fork of grpc-gateway to fix a bug they have in "nested query param generation"
//go:generate go install github.com/evanphx/grpc-gateway/protoc-gen-swagger

//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
//go:generate go install "-ldflags=-s -w -X github.com/vektra/mockery/cmd/mockery.SemVer=2.15.0" github.com/vektra/mockery/v2@v2.15.0

import (
	_ "github.com/evanphx/grpc-gateway/protoc-gen-swagger"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/kevinburke/go-bindata"
	_ "github.com/mitchellh/protoc-gen-go-json"
	_ "github.com/vektra/mockery/v2/cmd"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
