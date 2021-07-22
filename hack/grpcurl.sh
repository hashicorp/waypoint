token=$(waypoint user token)

grpcurl \
  -insecure \
  -rpc-header 'client-api-protocol: 1,1' \
  -rpc-header "authorization: $token" \
  -d '{ "application": { "project": "acmecorp", "application": "web" } }' \
  localhost:9701 hashicorp.waypoint.Waypoint/UI_ListDeployments
