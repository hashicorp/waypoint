# Telemetry

Waypoint can publish internal telemetry, including server gRPC traces.
This directory contains examples and sample infrastructure for consuming that telemetry.

### Docker-Compose
In this directory is a docker compose file to bring up opencensus infrastructure
for consuming and forwarding metrics (an agent anda collector), a jeager
all-in-one container for viewing trace data, and a datadog agent for
forwarding traced to datadog. If you wish to use the datadog exporter, you'll
need to supply a `DD_API_KEY` environment variable to the docker compose, like so:

```
DD_API_KEY=<redacted> docker-compose up
```

If you do not provide a DD_API_KEY, the datadog agent will fail to start, but the
rest of the infrastructure will function correctly.

### Waypoint Server -> Jeager example


The waypoint server supports multiple exporters, but this example will show exporting
to the opencensus agent, and viewing traces in jeager.

#### Run the docker-compose infrastructure
Run `docker-compose up` within this directroy. You do not need to specify a `DD_API_KEY`

#### Run the waypoint server

Provide the following flags to the waypoint server (either via extra flags passed to `waypoint server install` after a `--`,
or as flags to `waypoint server run`):

- `-telemetry-opencensus-agent-addr=127.0.0.1:55678`
- `-telemetry-opencensus-agent-insecure`
- `-telemetry-opencensus-zpages-addr=localhost:9999` (not required, but useful for troubleshooting)

A full example of a `waypoint server run` command:
```
waypoint server run \
  -advertise-addr=127.0.0.1:9701 -listen-grpc=127.0.0.1:9701 -listen-http=127.0.0.1:9702 \
  -db=/tmp/data.db -accept-tos -advertise-tls-skip-verify -url-enabled -vv \
  -telemetry-opencensus-agent-addr=127.0.0.1:55678 \
  -telemetry-opencensus-zpages-addr=localhost:9999 \
  -telemetry-opencensus-agent-insecure \
```

You may observe in the debug logs that the opencensus agent exporter and zpages server have started.

For detailed server telemetry documentation, see the flags prefixed with `-telemetry`
at https://www.waypointproject.io/commands/server-run#command-options. 

#### Make grpc requests to the waypoint server

[Bootstrap](https://www.waypointproject.io/commands/server-bootstrap) your server if you're running it for the first
time via `server-run`, or run any waypoint CLI commands against your server (i.e. `waypoint context verify`).
When the waypoint server handles any grpc request, it will send a trace to the opencensus agent, which will
forward it to the opencensus collector, which will forward it to jeager.

#### Observe traces in local Jaeger

Wait a moment, then visit http://localhost:16686/. Observe that a service `waypoint` exists. Choose "Find Traces"
with the waypoint service selected, and observe traces!

![jaeger traces](images/jeager-traces.png)

### Troubleshooting

You can troubleshoot the flow of traces by visiting the zPages endpoints:
- on the waypoint server (if enabled with the `-telemetry-opencensus-zpages-addr=localhost:9999` flag) at http://localhost:9999/debug/tracez
- on the opencensus agent container at http://localhost:9998/debug/tracez
- on the opencensus collector container at http://localhost:9997/debug/tracez

Example tracez on from the waypoint server:

![tracez example](images/tracez-example.png)