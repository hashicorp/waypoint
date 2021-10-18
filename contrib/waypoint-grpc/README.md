# Waypoint gRPC

If you find yourself wanting to interact directly with Waypoint’s gRPC API, this
script is for you. It uses Waypoint CLI’s current context for server address and
credentials, all you need to specify is the method you want to call and the data
you want to send.

Under the hood it uses the excellent gRPCurl. For advanced uses, you may find it
easier to drop down to using gRPCurl directly.

## Prequisites

1. [grpcurl](https://github.com/fullstorydev/grpcurl#installation)
2. [jq](https://stedolan.github.io/jq/download/)
3. waypoint ;)

## Examples

Query the `GetVersionInfo` method:

```sh
$ waypoint-grpc GetVersionInfo
{
  "info": {
    "api": {
      "current": 1,
      "minimum": 1
    },
```

Get a project:

```sh
$ waypoint-grpc.sh GetProject '{ "project": { "project": "example" } }'
{
  "project": {
    "name": "example",
    "applications": [
      {
        "name": "web",
```
