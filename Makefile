# A lot of this Makefile right now is temporary since we have a private
# repo so that we can more sanely create

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
	DOCKER_BUILDKIT=1 docker build --ssh default -t waypoint:latest .
