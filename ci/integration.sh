#!/usr/bin/env bash

set -e -u -o pipefail

set -x

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Configure Kubernetes"

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Confirm k8s is working
echo "Confirm kubernetes is working:"
kubectl cluster-info


echo "Boot up the registry to use:"
docker run -d -p 5000:5000 --restart=always --name registry.localhost registry:2
#
#echo "Connect the registry network to the k3d network"
#docker network connect k3d-k3s-default registry.localhost

WP="$(pwd)/waypoint"

test -e "$WP"

cd ci/sinatra || exit 1

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Waypoint init"
"$WP" init

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Waypoint build"
"$WP" build

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Waypoint deploy"
# If the registry isn't working and the pods are therefore unable to pull, we get stuck in an infinite wait
timeout 1m "$WP" deploy

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Waypoint release"
"$WP" release

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Waypoint deployment list"
# Smoke test list methods
"$WP" deployment list
"$WP" deployment list -V
"$WP" deployment list -json

## Let things get going.
sleep 10

[[ -n "$GITHUB_ACTION" ]] && echo "::group::Check deployed sinatra service"
kubectl get service sinatra -o 'jsonpath={}' | jq

PORT=$(kubectl get service sinatra -o jsonpath="{.spec.ports[0].nodePort}")

test "$(curl -s "localhost:$PORT")" = "Welcome to Waypoint!"
