package server

//go:generate sh -c "protoc -I../../thirdparty/proto/api-common-protos -I ../.. ../../pkg/server/proto/server.proto --go_out=plugins=grpc:../.. --go-json_out=../.."
//go:generate mv ./proto/server.pb.json.go ./gen
//go:generate mockery -all -case underscore -dir ./gen -output ./gen/mocks
