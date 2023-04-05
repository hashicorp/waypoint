ASSETFS_PATH?=pkg/server/gen/bindata_ui.go

GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_DIRTY=$$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$$(git describe --tags --always --match "v*")
GIT_IMPORT="github.com/hashicorp/waypoint/internal/version"
GOLDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)"
CRT_GOLDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE) -X $(GIT_IMPORT).Version=$(VERSION) -X $(GIT_IMPORT).Prerelease=$(PRERELEASE)"
GO_CMD?=go
WP_SERVER_PLATFORM?="linux/amd64"

# For changelog generation, default the last release to the last tag on
# any branch, and this release to just be the current branch we're on.
LAST_RELEASE?=$$(git describe --tags $$(git rev-list --tags --max-count=1))
THIS_RELEASE?=$$(git rev-parse --abbrev-ref HEAD)


.PHONY: bin
bin: # Creates the binaries for Waypoint for the current platform
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(GOLDFLAGS) -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(GOLDFLAGS) -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o dist/waypoint ./cmd/waypoint

.PHONY: bin/crt-waypoint
bin/crt-waypoint: # Creates the binaries for Waypoint for the current platform
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(WAYPOINT_GOOS) GOARCH=$(WAYPOINT_GOARCH) go build -ldflags $(CRT_GOLDFLAGS) -tags assetsembedded -o dist/$(CRT_BIN_NAME) ./cmd/waypoint

# bin/cli-only only recompiles waypoint, and doesn't recompile or embed the ceb.
# You can use the binary it produces as a server, runner, or CLI, but it won't contain the CEB, so
# it won't be able to build projects that don't have `disable_entrypoint = true` set in their build hcl.
.PHONY: bin/no-ceb
bin/cli-only: # Builds only the cli with no ceb
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

.PHONY: bin/linux
bin/linux: # Creates the binaries for Waypoint for the linux platform
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	GOOS=linux CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint ./cmd/waypoint

.PHONY: bin/windows
bin/windows: # Create windows binaries
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -tags assetsembedded -o ./waypoint.exe ./cmd/waypoint

.PHONY: bin/crt-assets
bin/crt-assets: # Create assets for caching in CRT
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./internal/assets/ceb/ceb ./cmd/waypoint-entrypoint
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./internal/assets/ceb/ceb-arm64 ./cmd/waypoint-entrypoint
	cd internal/assets && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb && cd ../..

.PHONY: bin/entrypoint
bin/entrypoint: # Create the entrypoint for the current platform
	CGO_ENABLED=0 go build -tags assetsembedded -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: bin/crt-waypoint-entrypoint
bin/crt-waypoint-entrypoint: # Create the entrypoint for the current platform
		CGO_ENABLED=0 go build -ldflags $(CRT_GOLDFLAGS) -tags assetsembedded -o dist/waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: install
install: bin # Build and copy binaries to $GOPATH/bin/waypoint
ifneq ("$(wildcard $(GOPATH)/bin/waypoint)","")
	rm $(GOPATH)/bin/waypoint
endif
	mkdir -p $(GOPATH)/bin
	cp ./waypoint $(GOPATH)/bin/waypoint

.PHONY: format
format: # Format all go code in project
	gofmt -s -w ./

.PHONY: docker/server
docker/server: docker/server-only docker/odr

.PHONY: docker/server-only
docker/server-only: # Builds a Waypoint server docker image
	DOCKER_BUILDKIT=1 docker buildx build \
					--platform $(WP_SERVER_PLATFORM) \
					-t waypoint:dev \
					.

.PHONY: docker/odr
docker/odr: # Builds a Waypoint on-demand runner docker image
	DOCKER_BUILDKIT=1 docker buildx build --target odr \
					--platform $(WP_SERVER_PLATFORM) \
					-t waypoint-odr:dev \
					.

.PHONY: docker/tools
docker/tools: # Creates a docker tools file for generating waypoint server protobuf files
	@echo "Building docker tools image"
	docker build -f tools.Dockerfile -t waypoint-tools:dev .

.PHONY: docker/gen/server
docker/gen/server: docker/tools
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	docker run -v `pwd`:/waypoint -it docker.io/library/waypoint-tools:dev make gen/server

.PHONY: docker/gen/plugins
docker/gen/plugins: docker/tools
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	docker run -v `pwd`:/waypoint -it docker.io/library/waypoint-tools:dev make gen/plugins

# expected to be invoked by make gen/changelog LAST_RELEASE=gitref THIS_RELEASE=gitref
.PHONY: gen/changelog
gen/changelog: # Generates the changelog for Waypoint
	@echo "Generating changelog for $(THIS_RELEASE) from $(LAST_RELEASE)..."
	@echo
	@changelog-build -last-release $(LAST_RELEASE) \
		-entries-dir .changelog/ \
		-changelog-template .changelog/changelog.tmpl \
		-note-template .changelog/note.tmpl \
		-this-release $(THIS_RELEASE)

# generates protos for the plugins inside builtin
.PHONY: gen/plugins
gen/plugins: # Generates plugin protobuf Go files
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	go generate ./builtin/...

.PHONY: gen/server
gen/server: # Generates server protobuf Go files from server.proto
	@test -s "thirdparty/proto/api-common-protos/.git" || { echo "git submodules not initialized, run 'git submodule update --init --recursive' and try again"; exit 1; }
	go generate ./pkg/server

.PHONY: gen/client
gen/client: # Generates grpc-gateway client from server.swagger.json
	go generate ./pkg/client

.PHONY: gen/ts
gen/ts: # Generates frontend typescript files
	# Clear existing generated files
	@rm -rf ./ui/lib/api-common-protos/google 2> /dev/null
	@rm -rf ./ui/lib/opaqueany/*.{js,ts} 2> /dev/null
	@rm -rf ./ui/lib/waypoint-client/*.ts 2> /dev/null
	@rm -rf ./ui/lib/waypoint-pb/*.{js,d.ts} 2> /dev/null

	# Generate JS and gRPCWeb libraries from pkg/server/proto/server.proto
	protoc \
		-I=. \
		-I=./thirdparty/proto/api-common-protos/ \
		-I=./thirdparty/proto/opaqueany/ \
		./pkg/server/proto/server.proto \
		--js_out=import_style=commonjs:ui/lib/waypoint-pb/ \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:ui/lib/waypoint-client/

	# Rearrange generated libraries
	@mv ./ui/lib/waypoint-client/pkg/server/proto/* ./ui/lib/waypoint-client/
	@rm -rf ./ui/lib/waypoint-client/pkg
	@mv ./ui/lib/waypoint-client/server_pb.d.ts ./ui/lib/waypoint-pb/
	@mv ./ui/lib/waypoint-pb/pkg/server/proto/* ./ui/lib/waypoint-pb/
	@rm -rf ./ui/lib/waypoint-pb/pkg

	# Hack: fix import of api-common-protos and various JS/TS imports
	# These issues below will help:
	#   https://github.com/protocolbuffers/protobuf/issues/5119
	#   https://github.com/protocolbuffers/protobuf/issues/6341
	find . -type f -wholename './ui/lib/waypoint-pb/*' | xargs sed -i 's/\.\.\/\.\.\/\.\.\/google/api-common-protos\/google/g'
	find . -type f -wholename './ui/lib/waypoint-pb/*' | xargs sed -i 's/\.\.\/\.\.\/\.\.\/any_pb/opaqueany\/any_pb/g'
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/\.\.\/\.\.\/\.\.\/google/api-common-protos\/google/g'
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/\.\/server_pb/waypoint-pb/g'
	find . -type f -wholename './ui/lib/waypoint-client/*' | xargs sed -i 's/\.\.\/\.\.\/\.\.\/pkg\/server\/proto\/server_pb/waypoint-pb/g'

	# Generate JS and TS from thirdparty/proto/api-common-protos
	protoc \
		-I=./thirdparty/proto/api-common-protos/ \
		./thirdparty/proto/api-common-protos/google/**/*.proto \
		--js_out=import_style=commonjs,binary:ui/lib/api-common-protos/ \
		--ts_out=ui/lib/api-common-protos/

	# Generate JS and TS from thirdparty/proto/opaqueany
	protoc \
		-I=./thirdparty/proto/opaqueany/ \
		./thirdparty/proto/opaqueany/*.proto \
		--js_out=import_style=commonjs,binary:ui/lib/opaqueany/ \
		--ts_out=ui/lib/opaqueany/

# This currently assumes you have run `ember build` in the ui/ directory
static-assets: # Generates the UI static assets
	@go-bindata -pkg gen -prefix dist -o $(ASSETFS_PATH) ./ui/dist/...
	@gofmt -s -w $(ASSETFS_PATH)

.PHONY: gen/doc
gen/doc: # generates the server proto docs
	mkdir -p ./doc/
	@rm -rf ./doc/* 2> /dev/null
	protoc -I=. \
		-I=./thirdparty/proto/api-common-protos/ \
		--doc_out=./doc --doc_opt=html,index.html \
		./pkg/server/proto/server.proto

.PHONY: gen/website-docs
gen/website-docs: gen/website-mdx gen/integrations-hcl

.PHONY: gen/integrations-hcl
gen/integrations-hcl: # Generates the HCL docs for integrations
	go run ./cmd/waypoint docs -hcl

.PHONY: gen/website-mdx
gen/website-mdx: # Generates the website markdown files
	go run ./cmd/waypoint docs -website-mdx
	go run ./cmd/waypoint docs -json
	go run ./tools/gendocs
	cd ./website; npx --no-install next-hashicorp format content # only format the content folder in website

.PHONY: tools
tools: # Install dependencies and tools required to build
	@echo "Fetching tools..."
	$(GO_CMD) generate -tags tools tools/tools.go
	@echo
	@echo "Done!"

.PHONY: test
test: # Run tests
	go test ./...

# Run state tests found in pkg/serverstate/statetest/
# To run a specific test, use TESTARGS
# Ex:
#
# >$ TESTARGS="-run TestImpl/runner_ondemand/TestOnDemandRunnerConfig" make test/boltdbstate
.PHONY: test/boltdbstate
test/boltdbstate: # Runs the boltdbstate tests
	@echo "Running state tests..."
	go test -test.v ./internal/server/boltdbstate $(TESTARGS)

.PHONY: test/service
test/service: # Runs the server API function tests. Requires a local postgresql and horizon via docker-compose
	$(warning "Running the full service suite requires `docker-compose up`! Some Tests rely on a local Horizon instance to be running.")
	@echo "Running service API server tests..."
	go test -test.v ./pkg/server/singleprocess/

.PHONY: help
help: # Print valid Make targets
	@echo "Valid targets:"
	@grep --extended-regexp --no-filename '^[a-zA-Z/_-]+:' Makefile | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
