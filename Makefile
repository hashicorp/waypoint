# A lot of this Makefile right now is temporary since we have a private
# repo so that we can more sanely create
ASSETFS_PATH?=internal/server/bindata_ui.go

# bin creates the binaries for Waypoint
.PHONY: bin
bin:
	GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	go build -tags assetsembedded -o ./waypoint ./cmd/waypoint
	go build -tags assetsembedded -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: dev
dev:
	GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	cd internal/assets && go generate
	go build -o ./waypoint ./cmd/waypoint
	go build -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: bin/linux
bin/linux: # create Linux binaries
	GOOS=linux GOARCH=amd64 $(MAKE) bin

.PHONY: docker/mitchellh
docker/mitchellh:
	DOCKER_BUILDKIT=1 docker build \
					--ssh default \
					--secret id=ssh.config,src="${HOME}/.ssh/config" \
					--secret id=ssh.key,src="${HOME}/.ssh/config" \
					-t waypoint:latest \
					.

.PHONY: docker/evanphx
docker/evanphx:
	DOCKER_BUILDKIT=1 docker build -f hack/Dockerfile.evanphx \
					--ssh default \
					-t waypoint:latest \
					.

.PHONY: gen/ts
gen/ts:
	@rm -rf ./ui/lib/api-common-protos/google 2> /dev/null
	protoc -I=. \
		-I=./vendor/proto/api-common-protos/ \
		./internal/server/proto/server.proto \
		--js_out=import_style=commonjs:ui/lib/waypoint-pb/ \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:ui/lib/waypoint-client/
	@mv ./ui/lib/waypoint-client/internal/server/proto/* ./ui/lib/waypoint-client/
	@mv ./ui/lib/waypoint-client/server_pb.d.ts ./ui/lib/waypoint-pb/
	@mv ./ui/lib/waypoint-pb/internal/server/proto/* ./ui/lib/waypoint-pb/
	# Hack: fix import of api-common-protos and various JS/TS imports
	# These issues below will help:
	#   https://github.com/protocolbuffers/protobuf/issues/5119
	#   https://github.com/protocolbuffers/protobuf/issues/6341
	find . -type f -wholename './ui/lib/waypoint-pb/*' | xargs sed -i 's/..\/..\/..\/google\/rpc\/status/api-common-protos\/google\/rpc\/status/g' 
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/..\/..\/..\/google\/rpc\/status/api-common-protos\/google\/rpc\/status/g' 
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/.\/server_pb/waypoint-pb/g' 
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/..\/..\/..\/internal\/server\/protwaypoint-pb/waypoint-pb/g' 
	
	protoc \
		-I=./vendor/proto/api-common-protos/ \
		./vendor/proto/api-common-protos/google/**/*.proto \
		--js_out=import_style=commonjs,binary:ui/lib/api-common-protos/ \
		--ts_out=ui/lib/api-common-protos/
	@rm -rf ./ui/lib/waypoint-pb/internal
	@rm -rf ./ui/lib/waypoint-client/internal
	@rm -rf ./ui/vendor/vendor
	@rm -rf ./google

# This currently assumes you have run `ember build` in the ui/ directory
static-assets:
	@go-bindata -pkg server -prefix dist -o $(ASSETFS_PATH) ./ui/dist/...
	@go fmt $(ASSETFS_PATH)

.PHONY: gen/doc
gen/doc:
	@rm -rf ./doc/* 2> /dev/null
	protoc -I=. \
		-I=./vendor/proto/api-common-protos/ \
		--doc_out=./doc --doc_opt=html,index.html \
		./internal/server/proto/server.proto
