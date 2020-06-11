#!/usr/bin/env bash

set -e -u -o pipefail

set -x

WP="/go/bin/waypoint"

test -e "$WP"

cd ci/sinatra || exit 1

"$WP" build

"$WP" push

"$WP" deploy

"$WP" release

PORT=$(kubectl get service sinatra -o jsonpath="{.spec.ports[0].nodePort}")

test "$(curl -s "localhost:$PORT")" = "Welcome to Waypoint!"
