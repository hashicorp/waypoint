# Architecture

This document describes the high-level architecture of Waypoint. The goal
of this document is to give anyone interested in contributing to Waypoint
a high-level overview of how Waypoint works. This tries to stay high level
to guide to the right area of Waypoint, but stops short of explaining in detail
how something might work since its easy for that to get out of sync with the
codebase. The recommended usage for this doc is to use it to hone in on a
specific part of Waypoint, then study the code there for specifics.

## Bird's Eye View

<p align="center"><img src=".github/images/birds-eye.png" width="70%"></p>

At the highest level, Waypoint is made up of _four_ distinct Waypoint-specific components (colored in green above):

1. **Waypoint Client (CLI, UI, etc.)** - This is an API client focused on taking
user input and calling the correct APIs to interact with Waypoint such as
viewing a list of deployments, triggering a new deployment, etc. This communicates
only to the Waypoint server.

2. **Waypoint Server** - The main API server. This responds to all the API requests
and stores data. _This communicates directly to nothing._ The Waypoint clients, runners,
entrypoints all open an outbound connection to the server, _not_ the reverse. Therefore,
the server queues information that it then communicates to other components when they
connect.

3. **Waypoint Runners** - One or more runners which are responsible for executing 
logic for builds, deploys, etc. Importantly, these are the only components in the 
entire architecture that need access to the target platform.

4. **Waypoint Entrypoint** - An optional component a deployment may have that 
connects back to the Waypoint server for features such as application configuration,
logs, and more.

All of these components are in the `hashicorp/waypoint` GitHub repository.
The entrypoint is compiled as a separate binary `waypoint-entrypoint`, 
and the client, server, and runner are compiled as the `waypoint` binary.

### Important Architectural Properties

There are certain properties that we specifically designed for that are 
worth calling out. We'd like to preserve these properties in any improvements
to the project.

**Network access only required for the Waypoint Server.** All components
connect _out_ to the Waypoint server. The Waypoint server is the only 
component that needs to be reachable by other components. This makes it
easy to run runners, entrypoints, clients because they don't need to be
internet-reachable. This helps with security as well.

You can see this in the diagram above by the directions of the arrows.

**3rd party secrets only required by the Waypoint Runner.** Only the runner
needs access to "3rd party secrets" (anything other than a Waypoint token),
such as cloud access credentials, Kubernetes credentials, etc. Only the
runner ever interacts directly with these systems. The client interacts
with the server, and the server interacts only with runners by queueing
jobs. This decision makes the security blast radius more easily understandable.

Note that the Waypoint server can _optionally_ have 3rd party secrets
in the form of [application and runner configuration](https://www.waypointproject.io/docs/app-config),
but users can choose to opt out of this. The runners _need_ secrets
to be practically useful.

You can see this in the diagram above by noting that only runners talk 
to the target platforms.

**The only component with persistent data is the Waypoint Server.** The 
server is the only component that persists data beyond restarts. It is
also the only component that has direct access to the database. All
other data storage and access must go through APIs to the server.
This property is very nice because the client, runner, and entrypoints
are all stateless. They can be safely restarted and they will reconstruct
their state. This makes operating a Waypoint cluster simpler.

## Code Path of a Waypoint Operation (such as `waypoint up`)

The list below shows a high-level overview of the code path of a
full end-to-end `waypoint up` operation. It is similar for other
operations such as `waypoint build` or `waypoint logs` but with
subtle differences.

This isn't meant to be an overview of our internal packages. 
The important internal packages are covered in a lot more
detail in the "Code Map" section. The goal of this instead is
to just give you a sense of how the code "flows".

CLI:

1. `cmd/waypoint`: CLI entry `waypoint up`
2. `internal/cli`: CLI logic
3. `internal/client`: High-level API client to interact with the Waypoint server

Server:

4. `internal/server/singleprocess`: Server-side API request logic
5. `internal/server/singleprocess/state`: Server-side persistant storage logic

Runner: 

6. `internal/runner`: Runner-side accept and execute jobs
7. `internal/core`: Runner-side core Waypoint logic

Entrypoint:

8. `internal/ceb`: Entrypoint logic for deployed applications

The code path of `client => server => runner` is a common pattern 
throughout many operations within Waypoint. The entrypoint can be
considered a special kind of client. Operations only ever flow in this
direction. From a code standpoint, at the highest level, it is
`internal/client`, `internal/server`, `internal/runner`, then
`internal/core`. This is the most common call path in Waypoint.

## Code Map

This section goes through the various paths and packages in the
Waypoint repository and documents what their purpose is and some
invariants and design decisions around these packages that should 
be held true.

TODO
