//go:build tools
// +build tools

// To install the following tools at the version used by this repo run:
// $ make tools
// or
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install github.com/kevinburke/go-bindata
import _ "github.com/kevinburke/go-bindata"

//go:generate go install github.com/golang/protobuf/proto
import _ "github.com/golang/protobuf/proto"

//go:generate go install github.com/golang/protobuf/protoc-gen-go
import _ "github.com/golang/protobuf/protoc-gen-go"

//go:generate go install github.com/mitchellh/protoc-gen-go-json
import _ "github.com/mitchellh/protoc-gen-go-json"

// Using a fork of grpc-gateway to fix a bug they have in "nested query param generation"
//go:generate go install github.com/evanphx/grpc-gateway/protoc-gen-swagger
import _ "github.com/evanphx/grpc-gateway/protoc-gen-swagger"

//go:generate go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"

//go:generate go install "-ldflags=-s -w -X github.com/vektra/mockery/mockery.SemVer=1.1.2" github.com/vektra/mockery/cmd/mockery
import _ "github.com/vektra/mockery"
