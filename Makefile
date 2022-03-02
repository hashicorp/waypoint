ASSETFS_PATH?=pkg/server/gen/bindata_ui.go

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
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

# bin/cli-only only recompiles waypoint, and doesn't recompile or embed the ceb.
# You can use the binary it produces as a server, runner, or CLI, but it won't contain the CEB, so
# it won't be able to build projects that don't have `disable_entrypoint = true` set in their build hcl.
.PHONY: bin/no-ceb
bin/cli-only:
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

.PHONY: bin/linux
bin/linux: # bin creates the binaries for Waypoint for the linux platform
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	GOOS=linux CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

.PHONY: bin/windows
bin/windows: # create windows binaries
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint.exe ./cmd/waypoint

.PHONY: bin/entrypoint
bin/entrypoint: # create the entrypoint for the current platform
	CGO_ENABLED=0 go build -tags assetsembedded -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: install
install: bin # build and copy binaries to $GOPATH/bin/waypoint
	rm $(GOPATH)/bin/waypoint
	mkdir -p $(GOPATH)/bin
	cp ./waypoint $(GOPATH)/bin/waypoint

.PHONY: test
test: # run tests
	go test ./...

.PHONY: format
format: # format go code
	gofmt -s -w ./

.PHONY: docker/server
docker/server: docker/server-only docker/odr

.PHONY: docker/server-only
docker/server-only:
	DOCKER_BUILDKIT=1 docker build \
					-t waypoint:dev \
					.

.PHONY: docker/odr
docker/odr:
	DOCKER_BUILDKIT=1 docker build --target odr \
					-t waypoint-odr:dev \
					.

.PHONY: docker/tools
docker/tools:
	@echo "Building docker tools image"
	docker build -f tools.Dockerfile -t waypoint-tools:dev .

.PHONY: docker/gen/server
docker/gen/server:
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	docker run -v `pwd`:/waypoint -it docker.io/library/waypoint-tools:dev make gen/server

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

# generates protos for the plugins inside builtin
.PHONY: gen/plugins
gen/plugins:
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	go generate ./builtin/...

.PHONY: gen/server
gen/server:
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	go generate ./pkg/server

.PHONY: gen/ts
gen/ts:
	@rm -rf ./ui/lib/api-common-protos/google 2> /dev/null
	protoc -I=. \
		-I=./thirdparty/proto/api-common-protos/ \
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
		-I=./thirdparty/proto/api-common-protos/ \
		./thirdparty/proto/api-common-protos/google/**/*.proto \
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
	mkdir -p ./doc/
	@rm -rf ./doc/* 2> /dev/null
	protoc -I=. \
		-I=./thirdparty/proto/api-common-protos/ \
		--doc_out=./doc --doc_opt=html,index.html \
		./internal/server/proto/server.proto

.PHONY: gen/website-mdx
gen/website-mdx:
	go run ./cmd/waypoint docs -website-mdx
	go run ./tools/gendocs
	cd ./website; npx --no-install next-hashicorp format

.PHONY: tools
tools: # install dependencies and tools required to build
	@echo "Fetching tools..."
	$(GO_CMD) generate -tags tools tools/tools.go
	@echo
	@echo "Done!"
