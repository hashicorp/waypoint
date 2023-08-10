#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


if ! command -v jq &> /dev/null
then
  echo "Please install jq:"
  echo "https://stedolan.github.io/jq/download/"
  exit 1
fi

if ! command -v grpcurl &> /dev/null
then
  echo "Please install grpcurl:"
  echo "https://github.com/fullstorydev/grpcurl#installation"
  exit 1
fi

method=$1
data=$2

if [ -z "$method" ]
then
  echo "Usage: waypoint-grpc.sh <method> [args]"
  echo
  echo "Examples:"
  echo "    waypoint-grpc.sh GetVersionInfo"
  echo "    waypoint-grpc.sh GetProject '{ \"project\": { \"project\": \"example\" } }'"
  exit
fi

default_context=$(waypoint context inspect -json | jq -r .default_context)
context_json=$(waypoint context inspect -json $default_context)
address=$(echo $context_json | jq -r .address)
token=$(echo $context_json | jq -r .auth_token)

grpcurl \
  -insecure \
  -H "client-api-protocol: 1,1" \
  -H "authorization: $token" \
  -d "$data" \
  $address \
  "hashicorp.waypoint.Waypoint.$method"
