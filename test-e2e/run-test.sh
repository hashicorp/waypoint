#!/bin/bash
set -eo pipefail

# Waypoint end to end test runner

# shell spinner: https://www.shellscript.sh/tips/spinner/
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

if [ "${E2E_PLATFORM}" != "Docker" ] && [ "${E2E_PLATFORM}" != "Kubernetes" ] && [ "${E2E_PLATFORM}" != "Ecs" ] && [ "${E2E_PLATFORM}" != "Nomad" ]; then
  echo "Environment variable 'E2E_PLATFORM' must be one of: 'Docker', 'Kubernetes', 'Ecs', 'Nomad'"
  exit 1
fi

# For running script outside of `test-e2e` folder
TESTDIR="${WP_TESTE2E_DIR:-$(pwd)}"

echo "==> Installing dependencies..."
echo

make tools

# TODO: install packages for building waypoint and running supported platforms:
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

# Target working directory for the binary location if not specified
export WP_BINARY="${WP_BINARY:-$TESTDIR/waypoint}"

if [ -z "$WP_EXAMPLES_PATH" ]; then
  echo "WP_EXAMPLES_PATH unset; setting to ${TESTDIR}/waypoint-examples"
  export WP_EXAMPLES_PATH="${TESTDIR}/waypoint-examples"
fi

echo "==> Checking if Waypoint binary is built..."
if [ -f "${WP_BINARY}" ]; then
  "${WP_BINARY}" version
  echo
else
  echo "==> Building waypoint binary..."
  echo
  make
  echo
fi

# TODO: build waypoint OR download a package, add a switch for this
#   - add param for installing a certain waypoint server, allow install from alpha package
#   - export proper vars for binary path and server image later on

# make tools
# git submodule update for grpc status from api common
# make

# Bring in test apps (potentially at a certain sha rather than `main`?)
# git clone --depth 1 git@github.com:hashicorp/waypoint-examples.git
if [ ! -d "$WP_EXAMPLES_PATH" ]; then
  echo "==> Pulling in waypoint-examples for test..."
  echo

  git clone --depth 1 git@github.com:hashicorp/waypoint-examples.git "$TESTDIR/waypoint-examples"
else
  echo "==> Using existing waypoint-examples repo for test..."
  echo
fi

# 

echo
echo "==> Running Waypoint end-to-end tests..."
echo

# TODO: allow for running all platforms, or only certain ones

# only spin for local devs running on machine to show tests aren't frozen
if [ -z "$CI" ]; then
  spin &
  SPIN_PID=$!
  trap 'kill -9 $SPIN_PID' $(seq 0 15)
fi

# Run Docker tests
go test -v "github.com/hashicorp/waypoint/test-e2e" -run "$E2E_PLATFORM"
testResult=$?

# Set up Nomad
# Run Nomad tests

# Set up K8S/K3S
# Run K8S tests

# Set up ECS
# Run ECS tests

if [[ "$testResult" -eq 0 ]]; then
  echo
  echo "==> Cleaning up after finishing tests..."
  echo

  if [[ ! -d WP_EXAMPLES_PATH ]]; then
    # Test clean up
    echo
    echo "* Cleaning up 'waypoint-examples'"
    echo

    rm -rf "$TESTDIR/waypoint-examples"
  fi
fi

# must be at end of script
if [ -z "$CI" ]; then
  kill -9 $SPIN_PID
fi
