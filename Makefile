# A lot of this Makefile right now is temporary since we have a private
# repo so that we can more sanely create

# bin creates the binaries for Devflow
.PHONY: bin
bin:
	go build -o ./devflow ./cmd/devflow
	go build -o ./devflow-entrypoint ./cmd/devflow-entrypoint

.PHONY: bin/linux
bin/linux: # create Linux binaries
	GOOS=linux GOARCH=amd64 $(MAKE) bin

.PHONY: docker/mitchellh
docker/mitchellh:
	docker build -t gcr.io/mitchellh-test/devflow:latest .
	docker push gcr.io/mitchellh-test/devflow:latest
