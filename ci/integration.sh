#!/usr/bin/env bash

set -e -u -o pipefail

set -x

go build -o ./waypoint ./cmd/waypoint

cd ci/sinatra || exit 1

../../waypoint build

../../waypoint push

../../waypoint deploy

../../waypoint release

PORT=$(kubectl get service sinatra -o jsonpath="{.spec.ports[0].nodePort}")

test "$(curl -s "localhost:$PORT")" = "Welcome to Waypoint!"
