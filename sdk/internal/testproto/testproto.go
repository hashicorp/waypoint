// Package testproto contains some protobuf defintions that are used
// in internal tests.
package testproto

//go:generate sh -c "protoc -I ./ ./*.proto --go_out=plugins=grpc:./"
