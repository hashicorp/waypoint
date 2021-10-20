#!/usr/bin/env bash

set -e -u -o pipefail

set -x

echo "Boot up the registry to use:"

docker run -d -p 5000:5000 --restart=always --name registry registry:2

WP="$(pwd)/waypoint"

test -e "$WP"

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

cd ci/sinatra || exit 1

"$WP" init

"$WP" build

"$WP" deploy

"$WP" release

# Smoke test list methods
"$WP" deployment list
"$WP" deployment list -V
"$WP" deployment list -json

## Let things get going.
sleep 10

# Smoke test
"$WP" exec echo "exec working"

"$WP" logs > /tmp/wp-logs &
logs_pid=$!
if [[ ! -f /tmp/wp-logs || $(cat /tmp/wp-logs | wc -l) -eq 0 ]]; then
  echo "logs did not work"
  exit 1
fi
kill $logs_pid


PORT=$(kubectl get service sinatra -o jsonpath="{.spec.ports[0].nodePort}")

test "$(curl -s "localhost:$PORT")" = "Welcome to Waypoint!"
