package keepalive

//go:generate sh -c "protoc -I ../.. ../../pkg/keepalive/proto/keepalive.proto --go_out=../.. --go-grpc_out=../.."
