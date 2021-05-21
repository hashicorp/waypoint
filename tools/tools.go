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

//go:generate go install "-ldflags=-s -w -X github.com/vektra/mockery/mockery.SemVer=1.1.2" github.com/vektra/mockery/cmd/mockery
import _ "github.com/vektra/mockery"
