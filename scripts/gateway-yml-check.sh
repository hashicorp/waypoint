#!/usr/bin/env bash

set -euo pipefail

echo -e "==> Checking all RPCs are represented in pkg/server/proto/gateway.yml\n"

proto_path=pkg/server/proto/server.proto
gateway_path=pkg/server/proto/gateway.yml
skip=$(cat<<-EOF
  CreateSnapshot
  RestoreSnapshot
  GetDefaultOnDemandRunnerConfig
  DeleteOnDemandRunnerConfig
  UpsertBuild
  UpsertPushedArtifact
  UpsertDeployment
  UpsertRelease
  UpsertStatusReport
  DeleteTrigger
  UpsertPipeline
EOF
)

rpcs=$(
  grep \
    --only-matching \
    --extended-regexp \
    "^\s*rpc \w+" \
    $proto_path \
    | \
  sed \
    -n \
    "s/^ *rpc//p"
)

status=0

for rpc in $rpcs; do
  if [[ $skip = *$rpc* ]]; then
    continue
  fi

  if grep -q $rpc $gateway_path; then
    continue
  fi

  echo -e >&2 "\033[1m$rpc\033[0m is missing"
  status=1
done

if [ $status -eq 0 ]; then
  echo -e "\033[32mSUCCESS: $gateway_path looks good!\033[0m"
else
  echo -e >&2 "\n\033[31mERROR: The RPCs above are missing from $gateway_path\033[0m"
fi

exit $status
