#!/bin/bash

# Waypoint end to end test runner

spin()
{
  spinner="/|\\-/|\\-"
  while :
  do
    for i in `seq 0 7`
    do
      echo -n "${spinner:$i:1}"
      echo -en "\010"
      sleep 1
    done
  done
}

echo "Beginning Waypoint end-to-end tests..."
echo

echo "==> Installing dependencies..."
echo

echo "Skipping for now"
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

echo "==> Building waypoint binary..."
echo

echo "Skipping for now"
echo "Assuming waypoint is available on the path"
echo

# build waypoint OR download a package, add a switch for this
#   - add param for installing a certain waypoint server, allow install from alpha package
#   - export proper vars for binary path and server image later on

# make tools
# git submodule update for grpc status from api common
# make

# Bring in test apps (potentially at a certain sha rather than `main`?)
# git clone --depth 1 git@github.com:hashicorp/waypoint-examples.git
if [[ ! -v WP_EXAMPLES_PATH ]]; then
  echo "==> Pulling in waypoint-examples for test..."
  echo

  git clone --depth 1 git@github.com:hashicorp/waypoint-examples.git
else
  echo "==> Using existing waypoint-examples repo for test..."
  echo
fi

# Test env vars
export WP_BINARY="waypoint"
export WP_SERVERIMAGE="hashicorp/waypoint:latest"
export WP_SERVERIMAGE_UPGRADE="hashicorp/waypoint:latest"

echo
echo "==> Running Waypoint end-to-end tests..."
echo

# TODO: allow for running all platforms, or only certain ones

# only spin for local devs running on machine to show tests aren't frozen
if [[ ! -v CI_ENV ]]; then
  spin &
  SPIN_PID=$!
  trap "kill -9 $SPIN_PID" `seq 0 15`
fi

go test .
testResult=$?

if [[ ! -v WP_EXAMPLES_PATH ]]; then
  if [[ "$testResult" -eq 0 ]]; then
    echo
    echo "==> Cleaning up 'waypoint-examples'"
    echo

    rm -rf waypoint-examples
  fi
fi

# must be at end of script
if [[ ! -v CI_ENV ]]; then
  kill -9 $SPIN_PID
fi
