#!/usr/bin/env bash

set -e -u -o pipefail

set -x

if [ -z "$CI" ]; then  # We are running locally
  export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
fi

# Confirm k8s is working
echo "Confirm kubernetes is working:"
kubectl cluster-info dump


if [ -z "$GITHUB_ACTION" ]; then
  echo "Boot up the registry to use:"
  docker run -d -p 5000:5000 --restart=always --name registry registry:2
fi

WP="$(pwd)/waypoint"

test -e "$WP"

cd ci/sinatra || exit 1

"$WP" init

timeout 3m "$WP" build

"$WP" deploy

"$WP" release

# Smoke test list methods
"$WP" deployment list
"$WP" deployment list -V
"$WP" deployment list -json

## Let things get going.
sleep 10

PORT=$(kubectl get service sinatra -o jsonpath="{.spec.ports[0].nodePort}")

test "$(curl -s "localhost:$PORT")" = "Welcome to Waypoint!"
