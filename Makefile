ASSETFS_PATH?=internal/server/gen/bindata_ui.go

GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_DIRTY=$$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$$(git describe --tags --always --match "v*")
GIT_IMPORT="github.com/hashicorp/waypoint/internal/version"
GOLDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)"
CGO_ENABLED?=0
GO_CMD?=go

# For changelog generation, default the last release to the last tag on
# any branch, and this release to just be the current branch we're on.
LAST_RELEASE?=$$(git describe --tags $$(git rev-list --tags --max-count=1))
THIS_RELEASE?=$$(git rev-parse --abbrev-ref HEAD)

.PHONY: bin
bin: # bin creates the binaries for Waypoint for the current platform
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

.PHONY: bin/windows
bin/windows: # create windows binaries
	GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint.exe ./cmd/waypoint

.PHONY: bin/entrypoint
bin/entrypoint: # create the entrypoint for the current platform
	CGO_ENABLED=0 go build -tags assetsembedded -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: install
install: bin # build and copy binaries to $GOPATH/bin/waypoint
	cp ./waypoint $(GOPATH)/bin/waypoint

.PHONY: test
test: # run tests
	go test ./...

.PHONY: format
format: # format go code
	gofmt -s -w ./

.PHONY: docker/server
docker/server:
	DOCKER_BUILDKIT=1 docker build \
					--ssh default \
					--secret id=ssh.config,src="${HOME}/.ssh/config" \
					--secret id=ssh.key,src="${HOME}/.ssh/config" \
					-t waypoint:dev \
					.

.PHONY: docker/evanphx
docker/evanphx:
	DOCKER_BUILDKIT=1 docker build -f hack/Dockerfile.evanphx \
					--ssh default \
					-t waypoint:latest \
					.

# expected to be invoked by make gen/changelog LAST_RELEASE=gitref THIS_RELEASE=gitref
.PHONY: gen/changelog
gen/changelog:
	@echo "Generating changelog for $(THIS_RELEASE) from $(LAST_RELEASE)..."
	@echo
	@changelog-build -last-release $(LAST_RELEASE) \
		-entries-dir .changelog/ \
		-changelog-template .changelog/changelog.tmpl \
		-note-template .changelog/note.tmpl \
		-this-release $(THIS_RELEASE)

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
	@go-bindata -pkg gen -prefix dist -o $(ASSETFS_PATH) ./ui/dist/...
	@gofmt -s -w $(ASSETFS_PATH)

.PHONY: gen/doc
gen/doc:
	@rm -rf ./doc/* 2> /dev/null
	protoc -I=. \
		-I=./vendor/proto/api-common-protos/ \
		--doc_out=./doc --doc_opt=html,index.html \
		./internal/server/proto/server.proto

.PHONY: tools
tools: # install dependencies and tools required to build
	@echo "Fetching tools..."
	$(GO_CMD) generate -tags tools tools/tools.go
	@echo
	@echo "Done!"
