package server

//go:generate sh -c "protoc -I`go list -m -f \"{{.Dir}}\" github.com/hashicorp/opaqueany` -I../../thirdparty/proto/api-common-protos -I ../.. ../../pkg/server/proto/server.proto --go_out=../.. --go-grpc_out=../.. --go-json_out=../.. --swagger_out=logtostderr=true,fqn_for_swagger_name=true,grpc_api_configuration=./proto/gateway.yml:../.. --grpc-gateway_out ../.. --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true --grpc-gateway_opt grpc_api_configuration=./proto/gateway.yml"
//go:generate sh -c "cat ./proto/server.swagger.json ./proto/swagger.json | jq --slurp 'reduce .[] as ${DOLLAR}item ({}; . * ${DOLLAR}item)' > ./gen/server.swagger.json"
//go:generate rm ./proto/server.swagger.json
//go:generate sh -c "mv ./proto/server.pb.*.go ./gen"
//go:generate mockery --all --case underscore --dir ./gen --output ./gen/mocks
