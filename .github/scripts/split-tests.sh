#!/usr/bin/env bash

set -euo pipefail

runner_count="${runner_count:-4}"

all_packages=$(go list -json="ImportPath" ./...)

filtered_packages=$(jq --compact-output \
                       --slurp \
                       'map(
                         .ImportPath | 
                         select(
                           (contains("test-e2e") | not)
                         )
                       )'\
                       <<< "${all_packages}")

# _nwise is an undocumented jq function that splits an array
# into arrays of size N where N is the input. We calculate N based
# off of the number of packages divided by the number of GHA runners
# we wish to use. We also track the array index to pass into the matrix
# so that we can uniquely identify the test results.
matrix=$(jq --argjson count $runner_count \
            --compact-output \
            -r \
            '[_nwise(length / $count | ceil)] | to_entries | map({id: .key, packages: .value})' \
            <<< "${filtered_packages}"
)

#             [keys as $k | {group: $k, packages: .[$k]}]' \
echo "matrix=${matrix}" | tee -a $GITHUB_OUTPUT