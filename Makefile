# A lot of this Makefile right now is temporary since we have a private
# repo so that we can more sanely create

# bin creates the binaries for Waypoint
.PHONY: bin
bin:
	go build -o ./waypoint ./cmd/waypoint
	go build -o ./waypoint-entrypoint ./cmd/waypoint-entrypoint

.PHONY: bin/linux
bin/linux: # create Linux binaries
	GOOS=linux GOARCH=amd64 $(MAKE) bin

.PHONY: docker/mitchellh
docker/mitchellh: bin/linux
	docker build -t gcr.io/mitchellh-test/waypoint:latest .
	#docker push gcr.io/mitchellh-test/waypoint:latest

.PHONY: k8s/mitchellh
k8s/mitchellh:
	./waypoint install \
		--annotate-service "external-dns.alpha.kubernetes.io/hostname=*.df.gcp.mitchellh.dev.,df.gcp.mitchellh.dev." \
		| kubectl apply -f -


