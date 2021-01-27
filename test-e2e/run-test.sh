#!/bin/bash

#TODO: make this a Makefile probably?

echo "Installing dependencies..."
echo

# install packages for building waypoint and running supported platforms:
# - git, curl, (probably more)
# - golang
# - docker
# - k8s (potentially external Digital Ocean service?)
# - nomad (use the nomad dev mode scripts from waypoint-flightlist)

# Build env vars
export GOOS="$(go env GOOS)"
export GOARCH="$(go env GOARCH)"
export GOEXE="$(go env GOEXE)"
export OUTDIR="build/${GOOS}_${GOARCH}"

echo "Building waypoint binary..."
echo

# build waypoint OR download a package, add a switch for this
#   - add param for installing a certain waypoint server

# make tools
# git submodule update for grpc status from api common
# make

echo "Pulling in waypoint-examples for test..."
echo
# Bring in test apps
# clone waypoint-examples

# Test env vars
export WP_BINARY="waypoint"

echo "Starting Waypoint end-to-end tests..."
echo

go test .
