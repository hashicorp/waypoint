#!/bin/bash

#TODO: make this a Makefile

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
GOEXE="$(go env GOEXE)"
OUTDIR="build/${GOOS}_${GOARCH}"

# Test env vars
export WP_BINARY="waypoint"

echo "Starting Waypoint end-to-end tests..."
echo

# install packages
# clone waypoint-examples
# build waypoint OR download a package, add a switch for this
#   - add param for installing a certain waypoint server

go test .
